package kvdb

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	ErrStorageNotFound = errors.New("not found in storage")
)

type (
	cachedStorage struct {
		db   *storage
		root *storageNode
	}

	storageNode struct {
		name   string
		next   sync.Map     // key: string, value *storageNode
		mu     sync.RWMutex // protect values
		values []*KV
	}
)

func newCachedStorage(ctx context.Context, sto *storage) *cachedStorage {
	s := &cachedStorage{
		db: sto,
	}
	s.buildCache(ctx)
	return s
}

func (c *cachedStorage) Get(_ context.Context, key ...string) (vs []*KV, err error) {
	k := ""
	if len(key) > 0 {
		k = key[0]
	}

	node := c.root
	var dfs func(node *storageNode)
	dfs = func(node *storageNode) {
		if node == nil {
			return
		}

		node.mu.Lock()
		defer node.mu.Unlock()
		if len(node.values) > 0 {
			vs = append(vs, node.values...)
		}

		node.next.Range(func(_, value any) bool {
			dfs(value.(*storageNode))
			return true
		})
	}
	if k == "" {
		dfs(c.root)
		return
	}

	pos := strings.Split(k, c.db.cfg.Separator)
	node = c.root
	for _, po := range pos {
		if v, ok := node.next.Load(po); !ok {
			err = ErrStorageNotFound
			return
		} else {
			node = v.(*storageNode)
		}
	}
	dfs(node)
	return
}

func (c *cachedStorage) Set(ctx context.Context, key string, value interface{}) (err error) {
	err = c.db.Set(ctx, key, value)
	if err != nil {
		return
	}

	c.setCache(c.db.buildStorageKey(key), value)
	return
}

func (c *cachedStorage) SetTTL(ctx context.Context, key string, value interface{}, ttl int64, keepalive ...bool) (err error) {
	err = c.db.SetTTL(ctx, key, value, ttl, keepalive...)
	if err != nil {
		return
	}

	c.setCache(c.db.buildStorageKey(key), value)
	return
}

func (c *cachedStorage) setCache(key string, value interface{}) {
	pos := strings.Split(strings.TrimPrefix(key, c.db.buildStorageKey()), c.db.cfg.Separator)
	node := c.root
	for _, po := range pos {
		if po == "" {
			continue
		}
		var n *storageNode
		if v, ok := node.next.Load(po); !ok {
			n = &storageNode{
				name: po,
			}
			node.next.Store(po, n)
		} else {
			n = v.(*storageNode)
		}
		node = n
	}

	node.updateOrInsertValue(key, value)
}

func (c *cachedStorage) Delete(ctx context.Context, key string) (err error) {
	err = c.db.Delete(ctx, key)
	if err != nil {
		return
	}

	c.remove(key)
	return
}

func (c *cachedStorage) remove(key string) {
	key = strings.TrimPrefix(key, c.db.buildStorageKey())
	pos := strings.Split(key, c.db.cfg.Separator)
	node := c.root

	for i := 0; i < len(pos)-1; i++ {
		po := pos[i]
		if po == "" {
			continue
		}
		if v, ok := node.next.Load(po); ok {
			node = v.(*storageNode)
		} else {
			return
		}
	}

	// in next layer
	if target, has := node.next.Load(pos[len(pos)-1]); has {
		if strings.HasSuffix(key, c.db.cfg.Separator) {
			node.next.Delete(pos[len(pos)-1])
		} else {
			target.(*storageNode).removeValue(c.db.buildStorageKey(key))
		}
	} else {
		// in current layer
		node.removeValue(c.db.buildStorageKey(key))
	}
}

func (c *cachedStorage) buildCache(ctx context.Context) {
	pfx := c.db.buildStorageKey()
	kvs, err := c.db.GetPrefix(ctx, pfx)
	if err != nil {
		g.Log().Errorf(ctx, "failed to build cache: %s", err.Error())
		return
	}

	root := &storageNode{
		name: strings.ReplaceAll(c.db.name, c.db.cfg.Separator, ""),
	}
	for _, kv := range kvs {
		if !strings.HasPrefix(kv.Key, pfx) {
			continue
		}
		pos := strings.Split(strings.TrimPrefix(kv.Key, pfx), c.db.cfg.Separator)
		node := root
		for _, po := range pos {
			var n *storageNode
			if v, ok := node.next.Load(po); !ok {
				n = &storageNode{
					name: po,
				}
				node.next.Store(po, n)
			} else {
				n = v.(*storageNode)
			}

			node = n
		}

		node.appendValues(kv)
	}
	c.root = root
}

func (c *cachedStorage) handleEvent(t EventType, key string, value interface{}) {
	key = c.db.buildStorageKey(key)
	switch t {
	case EventTypeUpdate, EventTypeCreate:
		c.setCache(key, value)
	case EventTypeDelete:
		c.remove(key)
	}
}

func (n *storageNode) appendValues(vs ...*KV) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.values = append(n.values, vs...)
}

func (n *storageNode) updateOrInsertValue(key string, value interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for i, kv := range n.values {
		if kv.Key == key {
			n.values[i].Value = g.NewVar(value)
			return
		}
	}

	n.values = append(n.values, &KV{Key: key, Value: g.NewVar(value)})
}

func (n *storageNode) removeValue(key string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for i, value := range n.values {
		if value.Key == key {
			n.values = append(n.values[:i], n.values[i+1:]...)
		}
	}
}
