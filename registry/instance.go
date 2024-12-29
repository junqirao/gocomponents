package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/junqirao/gocomponents/meta"
)

type (
	// Instance of registry object
	Instance struct {
		Id          string                 `json:"id"`           // uuid
		Host        string                 `json:"host"`         // host
		HostName    string                 `json:"host_name"`    // host name
		Port        int                    `json:"port"`         // port
		ServiceName string                 `json:"service_name"` // service name, usually use it as routing key
		Meta        map[string]interface{} `json:"meta"`         // meta data
	}
	// Service contains instances
	Service struct {
		mu        sync.RWMutex
		Name      string
		instances []*Instance
	}
)

func NewInstance(serviceName ...string) *Instance {
	ins := &Instance{Meta: make(map[string]interface{})}
	if len(serviceName) > 0 && serviceName[0] != "" {
		ins.ServiceName = serviceName[0]
	}
	return ins
}

func (i *Instance) WithAddress(host string, port int) *Instance {
	i.Host = host
	i.Port = port
	return i
}

func (i *Instance) WithMetaData(meta map[string]interface{}) *Instance {
	if i.Meta == nil {
		i.Meta = make(map[string]interface{})
	}
	for k, v := range meta {
		i.Meta[k] = v
	}
	return i
}

// Identity generate identity
func (i *Instance) Identity(separator ...string) string {
	sep := defaultIdentitySeparator
	if len(separator) > 0 {
		sep = separator[0]
	}
	return fmt.Sprintf("%s%s%s@%s", i.ServiceName, sep, i.Id, i.Host)
}

// String of instance
func (i *Instance) String() string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}

func (i *Instance) registryIdentity(prefix ...string) string {
	pfx := defaultRegistryPrefix
	if len(prefix) > 0 && prefix[0] != "" {
		pfx = prefix[0]
	}
	return fmt.Sprintf("%s%s", pfx, i.Identity())
}

func (i *Instance) clone() *Instance {
	m := make(map[string]interface{})
	if i.Meta != nil {
		for k, v := range i.Meta {
			m[k] = v
		}
	}
	return &Instance{
		Id:          i.Id,
		Host:        i.Host,
		HostName:    i.HostName,
		Port:        i.Port,
		ServiceName: i.ServiceName,
		Meta:        m,
	}
}

func (i *Instance) fillInfo() *Instance {
	if i.Id == "" {
		i.Id = uuid.NewV4().String()
	}
	// try fetch host name if not exist
	if i.HostName == "" {
		i.HostName, _ = os.Hostname()
	}
	// set defaultPort if out of range
	if i.Port <= 0 || i.Port > 65535 {
		i.Port = defaultPort
	}
	// try to get ip address it host field not set,
	// if failed to get ipv4 address use hostname as host
	if i.Host == "" {
		if ip := meta.IPV4(); ip != "" {
			i.Host = ip
		} else {
			i.Host = i.HostName
		}
	}
	// overwrite service name if meta is set
	if sn := meta.ServiceName(); sn != "" {
		i.ServiceName = sn
	}
	return i
}

func (s *Service) remove(id string) *Instance {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, instance := range s.instances {
		if instance.Identity() == id {
			s.instances = append(s.instances[:i], s.instances[i+1:]...)
			return instance
		}
	}
	return nil
}

// append instance to instances
func (s *Service) append(instance ...*Instance) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.instances = append(s.instances, instance...)
}

// upsert or insert instance to instances. notice: insertion not in order
func (s *Service) upsert(instance *Instance) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, ori := range s.instances {
		if ori.Identity() == instance.Identity() {
			s.instances[i] = instance
			return
		}
	}
	s.instances = append(s.instances, instance)
}

// Range instances
func (s *Service) Range(h func(instance *Instance) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, instance := range s.instances {
		if !h(instance) {
			break
		}
	}
}

// Len of instance
func (s *Service) Len() int {
	return len(s.instances)
}

// Instances slice copy of this service
func (s *Service) Instances() []*Instance {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ins []*Instance
	for _, instance := range s.instances {
		ins = append(ins, instance.clone())
	}
	return ins
}
