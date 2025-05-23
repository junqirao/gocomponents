package registry

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/grace"
	"github.com/junqirao/gocomponents/kvdb"
)

// global variable define
var (
	// Registry Global instance of registry, use it after Init
	Registry Interface
	// currentInstance created at Init
	currentInstance *Instance
	Empty           = new(Instance)
	onceLoad        = sync.Once{}
)

// error define
var (
	ErrAlreadyRegistered = errors.New("already registered")
	ErrServiceNotFound   = errors.New("service not found")
)

type (
	// Interface abstracts registry
	Interface interface {
		// register currentInstance
		register(ctx context.Context, ins *Instance) (err error)
		// Deregister deregister currentInstance
		Deregister(ctx context.Context) (err error)
		// GetService by service name
		GetService(ctx context.Context, serviceName ...string) (service *Service, err error)
		// GetServices of all
		GetServices(ctx context.Context) (services map[string]*Service, err error)
		// RegisterEventHandler register event handler
		RegisterEventHandler(handler EventHandler)
	}

	// EventType of instance change, alias of kvdb.EventType
	EventType = kvdb.EventType
	// EventHandler of instance change
	EventHandler func(i *Instance, e EventType)
	// eventWrapper ...
	eventWrapper struct {
		handler EventHandler
		next    *eventWrapper
	}
)

// Current returns copy of current instance,
// if not registered return Empty Instance
func Current() *Instance {
	if currentInstance == nil {
		return Empty.Clone()
	}
	return currentInstance.Clone()
}

// Init registry module with config and sync services info from database and build local caches.
// if *Instance is provided will be register automatically.
// if context is done, watch loop will stop and local cache won't be updated anymore.
func Init(ctx context.Context, db kvdb.Database, ins ...*Instance) (err error) {
	config := &Config{}
	if err = g.Cfg().MustGet(ctx, "registry").Scan(&config); err != nil {
		return
	}
	return InitWithConfig(ctx, config, db, ins...)
}

func InitWithConfig(ctx context.Context, config *Config, db kvdb.Database, ins ...*Instance) (err error) {
	onceLoad.Do(func() {
		config.check()
		// create registry instance
		if Registry, err = newRegistry(ctx, *config, db); err != nil {
			return
		}
		// collect instance info and register
		var instance *Instance
		if len(ins) > 0 && ins[0] != nil {
			instance = ins[0]
		} else {
			instance = config.Instance
		}
		if instance != nil {
			err = Registry.register(ctx, instance.Clone().fillInfo())
		}
	})
	return
}

type registry struct {
	cli             kvdb.Database
	cfg             *Config
	cache           sync.Map // service_name : *Service
	evs             *eventWrapper
	reRegisterCount int
}

func newRegistry(ctx context.Context, cfg Config, db kvdb.Database) (r Interface, err error) {
	reg := &registry{cfg: &cfg, cli: db}
	// build local cache
	reg.buildCache(ctx)
	// watchAndUpdateCache changes and upsert local cache
	// ** notice if context.Done() watchAndUpdateCache loop will stop
	go reg.watchAndUpdateCache(ctx)

	return reg, nil
}

func (r *registry) register(ctx context.Context, ins *Instance) (err error) {
	// check is already registered
	if currentInstance != nil {
		return ErrAlreadyRegistered
	}
	currentInstance = ins

	// get or create service
	service, err := r.getOrCreateService(ctx, currentInstance.ServiceName)
	if err != nil {
		return
	} else {
		// check if already registered in another machine
		for _, instance := range service.instances {
			if instance.Identity() == currentInstance.Identity() {
				if instance.Host != currentInstance.Host ||
					instance.Port != currentInstance.Port ||
					instance.HostName != currentInstance.HostName {
					return ErrAlreadyRegistered
				}
			}
		}
	}

	// register with heartbeat
	// renew a context in case upstream context closed cause heartbeat timeout
	if err = r.cli.Set(context.Background(),
		currentInstance.registryIdentity(r.cfg.getRegistryPrefix()),
		currentInstance.String(),
		kvdb.WithTTL(r.cfg.HeartBeatInterval),
		kvdb.WithKeepAlive(),
		kvdb.WithKeepAliveStoppedHandler(func(err error) {
			// keep alive stopped
			g.Log().Errorf(ctx, "registry heartbeat stopped: err=%v", err)
			// re-register
			r.reRegister(ctx, ins)
		}),
	); err != nil {
		return
	}
	g.Log().Infof(ctx, "registry success: %s", currentInstance.String())
	// reset counter
	r.reRegisterCount = 0
	// rebuild local cache
	r.buildCache(ctx)
	return
}

func (r *registry) reRegister(ctx context.Context, ins *Instance) {
	if r.reRegisterCount >= r.cfg.MaximumRetry {
		g.Log().Errorf(ctx, "registry re-register count exceed maximum: %d", r.cfg.MaximumRetry)
		// grace exit
		grace.ExecAndExit(ctx)
		return
	}
	g.Log().Infof(ctx, "registry try re-register, count=%d", r.reRegisterCount)
	r.reRegisterCount++
	// clear all resources
	currentInstance = nil
	// backoff and re-register
	backoff := r.fibWithLimit(r.reRegisterCount, 3, 30)
	g.Log().Infof(ctx, "registry re-register in %d seconds, count=%d", backoff, r.reRegisterCount)
	time.Sleep(time.Second * time.Duration(backoff))
	err := r.register(ctx, ins)
	if err != nil {
		g.Log().Errorf(ctx, "registry re-register failed: %v", err)
		r.reRegister(ctx, ins)
		return
	}
}

func (r *registry) fibWithLimit(n, min, max int) int {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}

	val := fib(n)
	if val > max {
		val = max
	}
	if val < min {
		val = min
	}
	return val
}

func (r *registry) Deregister(ctx context.Context) (err error) {
	if currentInstance == nil {
		return
	}
	err = r.cli.Delete(ctx, currentInstance.registryIdentity(r.cfg.getRegistryPrefix()))
	return
}

func (r *registry) GetService(_ context.Context, serviceName ...string) (service *Service, err error) {
	name := currentInstance.ServiceName
	if len(serviceName) > 0 {
		name = serviceName[0]
	}
	value, ok := r.cache.Load(name)
	if ok {
		service = value.(*Service)
	} else {
		err = ErrServiceNotFound
	}
	return
}

func (r *registry) GetServices(_ context.Context) (services map[string]*Service, err error) {
	services = make(map[string]*Service)
	r.cache.Range(func(key, value interface{}) bool {
		services[key.(string)] = value.(*Service)
		return true
	})
	return
}

func (r *registry) RegisterEventHandler(handler EventHandler) {
	if r.evs == nil {
		r.evs = &eventWrapper{handler: handler}
		return
	}
	p := r.evs
	for p != nil && p.next != nil {
		p = p.next
	}
	p.next = &eventWrapper{handler: handler}
}

func (r *registry) buildCache(ctx context.Context) {
	response, err := r.cli.Get(ctx, r.cfg.getRegistryPrefix())
	if err != nil {
		g.Log().Errorf(ctx, "registry failed to build etcd cache: %v", err)
		return
	}
	size := 0
	for _, kv := range response {
		instance := new(Instance)
		if err = kv.Value.Struct(&instance); err != nil {
			return
		}

		serviceName := instance.ServiceName
		v, ok := r.cache.Load(serviceName)
		if !ok || v == nil {
			service := new(Service)
			r.cache.Store(serviceName, service)
			service.append(instance)
		} else {
			v.(*Service).upsert(instance)
		}

		size++
	}

	g.Log().Infof(ctx, "registry etcd cache builded, size=%v", size)
}

func (r *registry) watchAndUpdateCache(ctx context.Context) {
	pfx := r.cfg.getRegistryPrefix()
	err := r.cli.Watch(ctx, pfx, func(ctx context.Context, e kvdb.Event) {
		var instance *Instance
		switch e.Type {
		case kvdb.EventTypeDelete:
			g.Log().Infof(ctx, "registry node delete event: %v", e.Key)
			// find and delete instance by e.key=instance.Identity()
			r.cache.Range(func(key, value interface{}) bool {
				var (
					deleted bool
					service = value.(*Service)
				)

				instanceId := strings.TrimPrefix(e.Key, pfx)
				instance = service.remove(instanceId)
				deleted = instance != nil

				// remove empty service
				if len(service.instances) == 0 {
					r.cache.Delete(key)
				}
				return !deleted
			})
		case kvdb.EventTypeCreate, kvdb.EventTypeUpdate:
			g.Log().Infof(ctx, "registry node register event: %v", e.Key)
			instance = new(Instance)
			if err := e.Value.Struct(&instance); err != nil {
				g.Log().Errorf(ctx, "registry failed to upsert on watchAndUpdateCache: %v", err)
				return
			}

			// get or create service
			service, err := r.getOrCreateService(ctx, instance.ServiceName)
			if err != nil {
				g.Log().Errorf(ctx, "registry failed to upsert on watchAndUpdateCache: %v", err)
				return
			}

			// upsert or insert instance to service
			service.upsert(instance)

			// upsert currentInstance
			if currentInstance != nil && instance.Id == currentInstance.Id {
				currentInstance = instance.Clone()
			}
		}

		r.pushEvent(instance, e.Type)
	})
	if err != nil {
		g.Log().Errorf(ctx, "registry failed to watchAndUpdateCache etcd: %v", err)
	}
}

func (r *registry) pushEvent(instance *Instance, e EventType) {
	ins := instance.Clone()
	p := r.evs
	for p != nil {
		go p.handler(ins, e)
		p = p.next
	}
}

func (r *registry) getOrCreateService(ctx context.Context, serviceName string) (service *Service, err error) {
	service, err = r.GetService(ctx, serviceName)
	switch {
	case errors.Is(err, ErrServiceNotFound):
		service = new(Service)
		r.cache.Store(serviceName, service)
		err = nil
	case err == nil:
	default:
	}
	return
}
