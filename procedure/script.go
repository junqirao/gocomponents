package procedure

import (
	"bytes"
	"text/template"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
)

type (
	// Script in go template syntax
	Script              string
	ScriptExecuteResult struct {
		OK     bool           `json:"ok"`
		Error  string         `json:"error"`
		Result map[string]any `json:"result"`
	}
)

func (s Script) Execute(fm template.FuncMap, data any) (res *gvar.Var, err error) {
	res = gvar.New(nil)
	if s == "" {
		return
	}
	tmpl, err := template.New(gmd5.MustEncrypt(s)).Funcs(fm).Parse(string(s))
	if err != nil {
		return
	}
	bb := &bytes.Buffer{}
	err = tmpl.Execute(bb, data)
	res.Set(bb.Bytes())
	return
}
