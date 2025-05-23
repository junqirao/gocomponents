package registry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/junqirao/gocomponents/kvdb"
)

func getConfig() *Config {
	return &Config{
		Name: "test",
	}
}

func TestInitWithoutInstance(t *testing.T) {
	err := InitWithConfig(context.Background(), getConfig(), kvdb.MustGetDatabase(context.Background()))
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestInit(t *testing.T) {
	err := InitWithConfig(context.Background(), getConfig(), kvdb.MustGetDatabase(context.Background()),
		NewInstance("test-service").
			WithAddress("127.0.0.1", 8080).
			WithMetaData(map[string]interface{}{"key": "value"}))
	if err != nil {
		t.Fatal(err)
		return
	}
	r := Registry.(*registry)
	r.cache.Range(func(serviceName, s interface{}) bool {
		service := s.(*Service)
		service.Range(func(instance *Instance) bool {
			t.Log(instance)
			return true
		})
		return true
	})
	t.Log("wait 20 s")
	time.Sleep(time.Second * 60)
}

func TestReRegister(t *testing.T) {
	err := InitWithConfig(context.Background(), getConfig(), kvdb.MustGetDatabase(context.Background()),
		NewInstance("test-service").
			WithAddress("127.0.0.1", 8080).
			WithMetaData(map[string]interface{}{"key": "value"}))
	if err != nil {
		t.Fatal(err)
		return
	}
	go func() {
		for {
			time.Sleep(time.Second)
			t.Logf("currentInstance==nil: %v", currentInstance == nil)
		}
	}()
	time.Sleep(time.Hour)
}

func TestRegistry(t *testing.T) {
	err := InitWithConfig(context.Background(), getConfig(), kvdb.MustGetDatabase(context.Background()),
		NewInstance("test-service").
			WithAddress("127.0.0.1", 8080).
			WithMetaData(map[string]interface{}{"key": "value"}))
	if err != nil {
		t.Fatal(err)
		return
	}

	service, err := Registry.GetService(context.Background(), "test-service")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf("service: %+v\n", service.Instances())
	instance := service.Instances()[0]
	if instance.Id != currentInstance.Id {
		t.Fatal("instance id not equal")
	}

	services, err := Registry.GetServices(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	for serviceName, s := range services {
		fmt.Printf("services[%s]: %+v\n", serviceName, s.Instances())
	}

	Registry.RegisterEventHandler(func(instance *Instance, e EventType) {
		fmt.Printf("event: %s, instance: %+v\n", e, instance)
	})

	err = Registry.Deregister(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
}
