package template

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/response"
)

var (
	ins = &Templates{
		mu: sync.RWMutex{},
		ts: make(map[string]*Template),
	}
	emptyTemplate = &Template{
		name:    "empty",
		Content: "",
	}
	ErrTemplateNotFound = response.NewCode(404, "template not found", http.StatusNotFound)
)

type (
	Template struct {
		t       *template.Template
		name    string
		Content string
	}
	Templates struct {
		mu sync.RWMutex
		ts map[string]*Template
	}
)

func T(name string) *Template {
	return ins.Get(name)
}

func (t *Templates) Get(name string) *Template {
	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.ts[name]
	if ok {
		return v
	}
	return emptyTemplate
}

func Create(name, content string) (t *Template, err error) {
	t = &Template{
		name:    name,
		Content: content,
	}

	tmpl, err := template.New(name).Parse(t.Content)
	if err != nil {
		return
	}
	t.t = tmpl

	ins.mu.Lock()
	defer ins.mu.Unlock()
	ins.ts[name] = t
	return
}

func (t Template) Parse(params map[string]string) (s string, err error) {
	if t.t == nil {
		err = ErrTemplateNotFound
		return
	}
	bb := &bytes.Buffer{}
	if err = t.t.Execute(bb, params); err != nil {
		return
	}
	s = bb.String()
	return
}

func LoadFromEmbed(ctx context.Context, efs embed.FS, dirName ...string) (err error) {
	dn := "templates"
	if len(dirName) > 0 && dirName[0] != "" {
		dn = dirName[0]
	}
	dir, err := efs.ReadDir(dn)
	if err != nil {
		return
	}
	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}
		var (
			info    fs.FileInfo
			content []byte
		)

		if info, err = entry.Info(); err != nil {
			return
		}
		name := info.Name()
		if content, err = efs.ReadFile(fmt.Sprintf("%s/%s", dn, name)); err != nil {
			return
		}
		name = name[:strings.Index(name, ".")]
		_, err = Create(name, string(content))
		if err != nil {
			return
		}
		g.Log().Infof(ctx, "template loaded from: %s, %s (%vbytes)", info.Name(), name, len(content))
	}
	return
}
