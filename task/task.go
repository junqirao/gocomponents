package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	uuid "github.com/satori/go.uuid"

	"github.com/junqirao/gocomponents/kvdb"
)

const (
	storageName = "task"
)

const (
	StatusCreated Status = iota
	StatusSuccess
	StatusFailed
)

const (
	FinishReasonTimeout = "timeout"
)

var (
	ErrTaskTypeNotFound  = errors.New("task type not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("task not found")
)

type (
	Task struct {
		Name         string         `json:"name"`
		Id           string         `json:"id"`
		Type         string         `json:"type"`
		Status       Status         `json:"status"`
		CreatedAt    int64          `json:"created_at"`
		FinishedAt   int64          `json:"finished_at"`
		FinishReason string         `json:"finish_reason"`
		Timeout      int64          `json:"timeout"` // second, default 30
		Unique       bool           `json:"unique"`
		Meta         map[string]any `json:"meta"`
		sig          <-chan *Task
	}
	Status  int
	Handler interface {
		Created(ctx context.Context, task *Task)
		Handle(ctx context.Context, task *Task) (err error)
		Archive(ctx context.Context, task *Task) (delete bool)
	}
)

func Init(_ context.Context) {
	// pre load
	_ = kvdb.Storages.GetStorage(storageName)
	// set event handler
	kvdb.Storages.SetEventHandler(storageName, handleTask)
}

type (
	CreateOption func(t *Task)
)

var (
	WithTaskName = func(name string) CreateOption {
		return func(t *Task) {
			t.Name = name
		}
	}
	WithTimeout = func(timeout int64) CreateOption {
		return func(t *Task) {
			t.Timeout = timeout
		}
	}
	WithId = func(id string) CreateOption {
		return func(t *Task) {
			t.Id = id
		}
	}
	UniqueTask = func(unique ...bool) CreateOption {
		return func(t *Task) {
			if len(unique) > 0 {
				t.Unique = unique[0]
			} else {
				t.Unique = true
			}
		}
	}
	WithMeta = func(meta map[string]any) CreateOption {
		return func(t *Task) {
			t.Meta = meta
		}
	}
)

func Create(ctx context.Context, typ string, opts ...CreateOption) (t *Task, err error) {
	t = &Task{
		Type:      typ,
		Status:    StatusCreated,
		CreatedAt: time.Now().UnixMilli(),
	}
	for _, opt := range opts {
		opt(t)
	}
	sto := kvdb.Storages.GetStorage(storageName)
	if t.Unique && t.Id != "" {
		var kvs []*kvdb.KV
		kvs, err = sto.Get(ctx, buildId(t.Type, t.Id))
		if err != nil && !errors.Is(err, kvdb.ErrStorageNotFound) {
			return
		}
		if len(kvs) > 0 {
			tt := new(Task)
			if err = kvs[0].Value.Struct(tt); err != nil {
				return
			}
			if tt.Status == StatusCreated && tt.Type == typ {
				err = ErrTaskAlreadyExists
				return
			}
			return
		}
	}
	if t.Id == "" {
		t.Id = uuid.NewV4().String()
	}
	if t.Name == "" {
		id := ""
		if len(t.Id) < 8 {
			id = t.Id
		} else {
			id = t.Id[:8]
		}
		t.Name = fmt.Sprintf("task_%s", id)
	}
	if t.Timeout == 0 {
		t.Timeout = 30
	}
	if t.Meta == nil {
		t.Meta = make(map[string]any)
	}
	bs, err := json.Marshal(t)
	if err != nil {
		return
	}
	id := buildId(t.Type, t.Id)
	// get notify channel after set storage success
	if err = sto.Set(ctx, id, string(bs)); err == nil {
		t.sig = getNotifyChannel(id)
	}
	return
}

func handleTask(t kvdb.EventType, _ string, value interface{}) {
	ctx := gctx.New()
	// try parse data
	task := new(Task)
	if err := gconv.Struct(value, task); err != nil {
		g.Log().Warningf(ctx, "drop task, failed to unmarshal task, invalid data: %v", value)
		return
	}

	h, err := getHandler(task.Type)
	if err != nil {
		g.Log().Warningf(ctx, "drop task, not handler named of: %s", task.Type)
		return
	}

	if t == kvdb.EventTypeDelete {
		return
	}

	var flow taskFlow
	switch task.Status {
	case StatusCreated:
		flow = &createdTask{h, task}
	case StatusSuccess, StatusFailed:
		flow = &handledTask{h, task}
	default:
		return
	}

	flow.Next(ctx, t)
}

func (t *Task) Wait(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t.Timeout)*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		t.Status = StatusFailed
		t.FinishReason = FinishReasonTimeout
	case newTask := <-t.sig:
		t.Status = newTask.Status
		t.FinishedAt = newTask.FinishedAt
		t.FinishReason = newTask.FinishReason
	}
}

type (
	createdTask struct {
		Handler
		*Task
	}
	handledTask struct {
		Handler
		*Task
	}
	taskFlow interface {
		Next(ctx context.Context, eventType kvdb.EventType)
	}
)

func (c *createdTask) Next(ctx context.Context, eventType kvdb.EventType) {
	id := buildId(c.Type, c.Id)
	mu, err := kvdb.NewMutex(ctx, fmt.Sprintf("task_exec_%s", id))
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	sto := kvdb.Storages.GetStorage(storageName)
	kvs, err := sto.Get(ctx, id)
	if err != nil || len(kvs) == 0 {
		return
	}
	task := new(Task)
	if err = kvs[0].Value.Struct(task); err != nil {
		return
	}
	if task.Status != StatusCreated {
		// already executed by other node
		return
	}

	switch eventType {
	case kvdb.EventTypeCreate:
		c.Created(ctx, c.Task)
		ctx, cancel := context.WithTimeout(ctx, time.Duration(c.Timeout)*time.Second)
		defer cancel()
		ch := make(chan error)
		go func() {
			ch <- c.Handle(ctx, c.Task)
		}()
		select {
		case err := <-ch:
			if err == nil {
				c.Task.Status = StatusSuccess
			} else {
				c.Task.Status = StatusFailed
				c.Task.FinishReason = err.Error()
			}
		case <-ctx.Done():
			// timeout
			c.Task.Status = StatusFailed
			c.FinishReason = FinishReasonTimeout
		}
		c.FinishedAt = time.Now().UnixMilli()
	}

	if err = sto.Set(ctx, id, c); err != nil {
		g.Log().Warningf(ctx, "update task failed: %v", err.Error())
	}
}

func (c *handledTask) Next(ctx context.Context, eventType kvdb.EventType) {
	id := buildId(c.Type, c.Id)
	mu, err := kvdb.NewMutex(ctx, fmt.Sprintf("task_handled_%s", id))
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	sto := kvdb.Storages.GetStorage(storageName)
	kvs, err := sto.Get(ctx, id)
	if err != nil || len(kvs) == 0 {
		// already executed by other node
		return
	}
	task := new(Task)
	if err = kvs[0].Value.Struct(task); err != nil {
		return
	}
	if task.Status != StatusFailed &&
		task.Status != StatusSuccess {
		// invalid status
		return
	}

	switch eventType {
	case kvdb.EventTypeUpdate:
		if c.Archive(ctx, task) {
			_ = sto.Delete(ctx, id)
		}
		notify(task)
	}
}

func GetTask(ctx context.Context, typ, id string) (t *Task, err error) {
	sto := kvdb.Storages.GetStorage(storageName)
	kvs, err := sto.Get(ctx, buildId(typ, id))
	if err != nil {
		return
	}
	if len(kvs) == 0 {
		err = ErrTaskNotFound
		return
	}
	t = new(Task)
	err = kvs[0].Value.Struct(t)
	return
}

func ListTasks(ctx context.Context) (ts []*Task, err error) {
	sto := kvdb.Storages.GetStorage(storageName)
	kvs, err := sto.Get(ctx)
	if err != nil {
		return
	}
	for _, kv := range kvs {
		t := new(Task)
		if err = kv.Value.Struct(t); err != nil {
			return
		}
		ts = append(ts, t)
	}
	return
}

func buildId(typ, id string) string {
	return fmt.Sprintf("%s_%s", typ, id)
}
