package task

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/junqirao/gocomponents/kvdb"
)

type (
	testHandler struct {
	}
)

func (t testHandler) Created(ctx context.Context, task *Task) {
	fmt.Printf("task created: %+v\n", task.Id)
}

func (t testHandler) Handle(ctx context.Context, task *Task) (err error) {
	fmt.Printf("task execute: %+v\n", task.Id)
	sleep := 2
	if v, ok := task.Meta["sleep"]; ok {
		sleep = gconv.Int(v)
	}
	time.Sleep(time.Second * time.Duration(sleep))
	return nil
}

func (t testHandler) Archive(ctx context.Context, task *Task) (delete bool) {
	fmt.Printf("task archive: %+v\n", task.Id)
	return
}

func TestTask(t *testing.T) {
	ctx := context.Background()
	err := kvdb.Init(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = kvdb.InitStorage(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	Init(ctx)

	const (
		taskType = "test"
		id       = "xxx"
	)

	RegisterHandler(taskType, &testHandler{})
	task, err := Create(ctx, taskType)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("before execute: %+v", task)
	task.Wait(ctx)
	t.Logf("executed: %+v", task)
	if task.FinishReason != "" || task.Status != StatusSuccess {
		t.Fatalf("wants reason='' and status=1 but got %v %v", task.FinishReason, task.Status)
		return
	}

	t.Logf("============================ test timeout ============================")
	task, err = Create(ctx, taskType, WithTimeout(1), WithId(id))
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("before execute: %+v", task)
	task.Wait(ctx)
	t.Logf("executed: %+v", task)
	if task.FinishReason != FinishReasonTimeout || task.Status != StatusFailed {
		t.Fatalf("wants reason=timeout and status=2 but got %v %v", task.FinishReason, task.Status)
		return
	}

	t.Logf("============================ test list ============================")
	tasks, err := ListTasks(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	for k, v := range tasks {
		t.Logf("[%v] %+v", k, v)
	}
	t.Logf("============================ test get one ============================")
	task, err = GetTask(ctx, taskType, id)
	if err != nil {
		return
	}
	t.Logf("%+v", task)
	t.Logf("============================ test unique ============================")
	task, err = Create(ctx, taskType, WithId("unique_test"), WithMeta(map[string]any{
		"sleep": 20,
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", task)
	t.Log("sleep 1s...")
	time.Sleep(time.Second)
	_, err = Create(ctx, taskType, WithId("unique_test"), UniqueTask())
	if !errors.Is(err, ErrTaskAlreadyExists) {
		t.Fatalf("want ErrTaskAlreadyExists but got %v", err)
	}
}
