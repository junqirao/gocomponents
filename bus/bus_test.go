package bus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/junqirao/gocomponents/kvdb"
)

var (
	testErr = fmt.Errorf("test error")
)

type testHandler struct {
}

func (t testHandler) Handle(ctx context.Context, msg *Message) {
	fmt.Printf("receive message: %+v\n", msg)
	msg.Ack(ctx)
	return
}

func (t testHandler) After(ctx context.Context, msg *Message) {
	fmt.Printf("finish message: %+v\n", msg)
}

type testErrorAckHandler struct {
}

func (t testErrorAckHandler) Handle(ctx context.Context, msg *Message) {
	fmt.Printf("receive message: %+v\n", msg)
	msg.Ack(ctx, testErr)
	return
}

func (t testErrorAckHandler) After(ctx context.Context, msg *Message) {
	fmt.Printf("finish message: %+v\n", msg)
}

func TestBus(t *testing.T) {
	err := kvdb.Init(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	err = Init(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(time.Second)

	RegisterHandler(context.Background(), "test", testHandler{})
	err = Push(context.Background(), "test", "test", 3, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Second)
	fmt.Println("---------------")
	RegisterHandler(context.Background(), "test_error_ack", testErrorAckHandler{})
	err = Push(context.Background(), "test_error_ack", "test", 3, true)
	if err != nil {
		if err.Error() != testErr.Error() {
			t.Fatalf("expected error: %v, got: %v", testErr, err)
		}
	} else {
		t.Fatalf("expected error: %v, got: %v", testErr, err)
	}
}
