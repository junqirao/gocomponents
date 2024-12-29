package audit

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type testAdaptor struct {
}

func (t testAdaptor) Store(ctx context.Context, record *Record) (err error) {
	fmt.Println(fmt.Sprintf("store: %+v\n", record))
	return
}

func (t testAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	fmt.Println(fmt.Sprintf("load: %+v\n", params))
	return
}

func TestLogger(t *testing.T) {
	SetAdaptor(&testAdaptor{})
	Logger.Module("module").Log("test", "12345")
	Logger.Module("module").Log("test", map[string]interface{}{
		"a": 1,
		"b": 2,
	})
	time.Sleep(time.Second)
	_ = Logger.Close()
}

type silenceAdaptor struct {
}

func (t silenceAdaptor) Store(ctx context.Context, record *Record) (err error) {
	return
}

func (t silenceAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	return
}

func BenchmarkLogger(b *testing.B) {
	SetAdaptor(&silenceAdaptor{})
	for i := 0; i < b.N; i++ {
		Logger.Module("module").Log("test", i)
	}
}
