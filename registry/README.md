## Registry

providing basic registry function

### supported database

* etcd

### usage

#### registry

basic registry, provide instance and service management

```golang
package main

import (
	"context"
	"fmt"

	"github.com/junqirao/gocomponents/registry"
)

func main() {
	// optional init instance if you want to register to registry
	ins := registry.NewInstance("your_service_name").
		WithAddress("127.0.0.1", 8080). // provide ip address and port for communication
		WithMetaData(map[string]interface{}{"key": "value"}) // metadata
	// config
	cfg := registry.Config{
		Prefix:            "/test-registry/", // prefix store to database
		HeartBeatInterval: 3,                 // healthy check heartbeat interval default 3s
	}

	// Init registry module with config and sync services info from database and build local caches.
	// if *Instance is provided will be register automatically.
	// if context is done, watch loop will stop and local cache won't be updated anymore.
	err := registry.Init(context.Background(), cfg, ins)
	if err != nil {
		// do something
		return
	}

	// get service from thread safe local cache
	service, err := registry.Registry.GetService(context.Background(), "test-service")
	if err != nil {
		// do something
		return
	}
	fmt.Printf("service: %+v\n", service.Instances())

	// get services from thread safe local cache
	services, err := registry.Registry.GetServices(context.Background())
	if err != nil {
		// do something
		return
	}
	for serviceName, s := range services {
		fmt.Printf("services[%s]: %+v\n", serviceName, s.Instances())
	}

	// register event handler, when instance changes will be triggered
	registry.Registry.RegisterEventHandler(func(instance *registry.Instance, e registry.EventType) {
		fmt.Printf("event: %s, instance: %+v\n", e, instance)
	})

	// if you don't deregister and exit application,
	// the registered instance will delete automatically after heart_beat_interval (default 3s) 
	err = registry.Registry.Deregister(context.Background())
	if err != nil {
		panic(err)
		return
	}
}

```