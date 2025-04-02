package kvdb

import (
	"context"
	"fmt"
	"testing"
	"time"
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
	err := Init(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	err = InitBus(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(time.Second)

	RegisterBusHandler(context.Background(), "test", testHandler{})
	err = PushMessage(context.Background(), "test", "test", 3, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Second)
	fmt.Println("---------------")
	RegisterBusHandler(context.Background(), "test_error_ack", testErrorAckHandler{})
	err = PushMessage(context.Background(), "test_error_ack", "test", 3, true)
	if err != nil {
		if err.Error() != testErr.Error() {
			t.Fatalf("expected error: %v, got: %v", testErr, err)
		}
	} else {
		t.Fatalf("expected error: %v, got: %v", testErr, err)
	}
}
