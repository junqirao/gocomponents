package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/junqirao/gocomponents/kvdb"
	"github.com/junqirao/gocomponents/meta"
)

type (
	testAdaptor struct {
		m map[string]*Record
	}
)

func (t testAdaptor) Store(_ context.Context, record *Record) (err error) {
	bytes, _ := json.Marshal(record)
	fmt.Println(fmt.Sprintf("store record: %s", string(bytes)))
	t.m[record.Name] = record
	return
}

func (t testAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	res = new(RecordQueryResult)
	for _, record := range t.m {
		res.Records = append(res.Records, record)
	}
	return
}

func TestUpdate2LatestConcurrency(t *testing.T) {
	cfg := `
registry:
  management_node: false
  endpoints:
    - 172.18.28.10:2379
    - 172.18.28.10:2380
    - 172.18.28.10:2381
  username:
  password:
`
	adaptor := testAdaptor{m: make(map[string]*Record)}
	fis := make([]*FuncInfo, 0)
	for i := 0; i < 10; i++ {
		fis = append(fis, NewFunc(fmt.Sprintf("test_%v", i), func(ctx context.Context) (err error) {
			// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
			time.Sleep(time.Second)
			return nil
		}, NewFuncConfig()))
	}

	fis2 := make([]*FuncInfo, 0)
	for i := 0; i < 10; i++ {
		idx := i
		fis2 = append(fis2, NewFunc(fmt.Sprintf("test_%v", i), func(ctx context.Context) (err error) {
			// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
			fmt.Println("fis 2 : " + gconv.String(idx))
			time.Sleep(time.Second)
			return nil
		}, NewFuncConfig()))
	}

	fis3 := make([]*FuncInfo, 0)
	for i := 0; i < 10; i++ {
		idx := i
		fis3 = append(fis3, NewFunc(fmt.Sprintf("test_%v", i), func(ctx context.Context) (err error) {
			// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
			fmt.Println("fis 3 : " + gconv.String(idx))
			time.Sleep(time.Second)
			return nil
		}, NewFuncConfig()))
	}

	ctx := context.Background()
	content, err := gcfg.NewAdapterContent(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	g.Cfg().SetAdapter(content)

	// lock
	mu, err := kvdb.NewMutex(ctx, fmt.Sprintf("updater_exec_%s", meta.ServerName()))
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func(tt *testing.T) {
		tt.Log("try concurrence update 1")
		err = ConcurrencyUpdate2Latest(ctx, adaptor, mu, fis...)
		if err != nil {
			wg.Done()
			tt.Error(err)
			return
		}
		wg.Done()
		tt.Log("done concurrence update 1")
	}(t)
	go func(tt *testing.T) {
		tt.Log("try concurrence update 2")
		err = ConcurrencyUpdate2Latest(ctx, adaptor, mu, fis2...)
		if err != nil {
			wg.Done()
			tt.Error(err)
			return
		}
		wg.Done()
		tt.Log("done concurrence update 2")
	}(t)
	go func(tt *testing.T) {
		tt.Log("try concurrence update 3")
		err = ConcurrencyUpdate2Latest(ctx, adaptor, mu, fis3...)
		if err != nil {
			wg.Done()
			tt.Error(err)
			return
		}
		wg.Done()
		tt.Log("done concurrence update 3")
	}(t)

	wg.Wait()
}

func TestKVDatabaseAdaptor(t *testing.T) {
	cfg := `
registry:
  management_node: false
  endpoints:
    - 172.18.28.10:2379
    - 172.18.28.10:2380
    - 172.18.28.10:2381
  username:
  password:
`
	fis := make([]*FuncInfo, 0)
	for i := 0; i < 10; i++ {
		fis = append(fis, NewFunc(fmt.Sprintf("test_%v", i), func(ctx context.Context) (err error) {
			// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
			time.Sleep(time.Second)
			return nil
		}, NewFuncConfig()))
	}

	fis2 := make([]*FuncInfo, 0)
	for i := 0; i < 10; i++ {
		idx := i
		fis2 = append(fis2, NewFunc(fmt.Sprintf("test_%v", i), func(ctx context.Context) (err error) {
			// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
			fmt.Println("fis 2 : " + gconv.String(idx))
			time.Sleep(time.Second)
			return nil
		}, NewFuncConfig()))
	}

	fis3 := make([]*FuncInfo, 0)
	for i := 0; i < 10; i++ {
		idx := i
		fis3 = append(fis3, NewFunc(fmt.Sprintf("test_%v", i), func(ctx context.Context) (err error) {
			// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
			fmt.Println("fis 3 : " + gconv.String(idx))
			time.Sleep(time.Second)
			return nil
		}, NewFuncConfig()))
	}

	ctx := context.Background()
	content, err := gcfg.NewAdapterContent(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	g.Cfg().SetAdapter(content)
	// lock
	mu, err := kvdb.NewMutex(ctx, fmt.Sprintf("updater_exec_%s", meta.ServerName()))
	if err != nil {
		return
	}

	database := kvdb.Raw
	adaptor := NewKVDatabaseAdaptor(database)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func(tt *testing.T) {
		tt.Log("try concurrence update 1")
		err = ConcurrencyUpdate2Latest(ctx, adaptor, mu, fis...)
		if err != nil {
			wg.Done()
			tt.Error(err)
			return
		}
		wg.Done()
		tt.Log("done concurrence update 1")
	}(t)
	go func(tt *testing.T) {
		tt.Log("try concurrence update 2")
		err = ConcurrencyUpdate2Latest(ctx, adaptor, mu, fis2...)
		if err != nil {
			wg.Done()
			tt.Error(err)
			return
		}
		wg.Done()
		tt.Log("done concurrence update 2")
	}(t)
	go func(tt *testing.T) {
		tt.Log("try concurrence update 3")
		err = ConcurrencyUpdate2Latest(ctx, adaptor, mu, fis3...)
		if err != nil {
			wg.Done()
			tt.Error(err)
			return
		}
		wg.Done()
		tt.Log("done concurrence update 3")
	}(t)

	wg.Wait()
}

func TestRetry(t *testing.T) {
	cfg := `
registry:
  management_node: false
  endpoints:
    - 172.18.28.10:2379
    - 172.18.28.10:2380
    - 172.18.28.10:2381
  username:
  password:
`
	fis := make([]*FuncInfo, 0)
	fis = append(fis, NewFunc("test_retry", func(ctx context.Context) (err error) {
		// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
		time.Sleep(time.Second)
		return errors.New("test retry")
	}, NewFuncConfig().Retry()))

	ctx := context.Background()
	content, err := gcfg.NewAdapterContent(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	g.Cfg().SetAdapter(content)
	database := kvdb.Raw
	adaptor := NewKVDatabaseAdaptor(database)
	t.Log("try update")
	err = Update2Latest(ctx, adaptor, fis...)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("done update")

	t.Log("try update")
	fis = make([]*FuncInfo, 0)
	fis = append(fis, NewFunc("test_retry", func(ctx context.Context) (err error) {
		// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
		time.Sleep(time.Second)
		return nil
	}, NewFuncConfig().Retry()))
	err = Update2Latest(ctx, adaptor, fis...)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("done update")
}

func TestMustRetry(t *testing.T) {
	cfg := `
registry:
  management_node: false
  endpoints:
    - 172.18.28.10:2379
    - 172.18.28.10:2380
    - 172.18.28.10:2381
  username:
  password:
`
	e := errors.New("test retry")
	fis := make([]*FuncInfo, 0)
	fis = append(fis, NewFunc("test_retry", func(ctx context.Context) (err error) {
		// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
		time.Sleep(time.Second)
		return e
	}, NewFuncConfig().Must().Retry()))

	ctx := context.Background()
	content, err := gcfg.NewAdapterContent(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	g.Cfg().SetAdapter(content)
	database := kvdb.Raw
	adaptor := NewKVDatabaseAdaptor(database)
	t.Log("try update")
	err = Update2Latest(ctx, adaptor, fis...)
	if !errors.Is(err, e) {
		t.Fatal("error not match")
		return
	}
	t.Log("done update")
}

func TestTimeout(t *testing.T) {
	cfg := `
registry:
  management_node: false
  endpoints:
    - 172.18.28.10:2379
    - 172.18.28.10:2380
    - 172.18.28.10:2381
  username:
  password:
`
	fis := make([]*FuncInfo, 0)
	f := NewFunc("test_timeout", func(ctx context.Context) (err error) {
		time.Sleep(time.Second * 10)
		return nil
	}, NewFuncConfig())
	f.timeout = time.Second
	fis = append(fis, f)

	ctx := context.Background()
	content, err := gcfg.NewAdapterContent(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	g.Cfg().SetAdapter(content)
	database := kvdb.Raw
	adaptor := NewKVDatabaseAdaptor(database)
	t.Log("try update")
	err = Update2Latest(ctx, adaptor, fis...)
	if !errors.Is(err, ErrExecuteTimeout) {
		t.Fatal("error not match")
		return
	}
	t.Log("done update")
}
