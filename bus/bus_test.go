package bus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/junqirao/gocomponents/kvdb"
)

type testHandler struct {
}

func (t testHandler) Handle(ctx context.Context, msg *Message) error {
	fmt.Printf("receive message: %+v\n", msg)
	msg.Ack(ctx)
	return nil
}

func (t testHandler) After(ctx context.Context, msg *Message) {
	fmt.Printf("finish message: %+v\n", msg)
}

type testErrorAckHandler struct {
}

func (t testErrorAckHandler) Handle(ctx context.Context, msg *Message) error {
	fmt.Printf("receive message: %+v\n", msg)
	msg.Ack(ctx, fmt.Errorf("test error"))
	return nil
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

	fmt.Println("-------")

	RegisterHandler(context.Background(), "test_error_ack", testErrorAckHandler{})
	err = Push(context.Background(), "test_error_ack", "test", 3, true)
	if err != nil {
		t.Fatal(err)
		return
	}
}
