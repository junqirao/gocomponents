package procedure

import (
	"context"
	"text/template"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/pkg/errors"
)

var (
	ErrInvalidParameter  = errors.New("invalid parameter")
	ErrNodeExecuteFailed = errors.New("node execute failed")
)

func (p *Proto) Execute(ctx context.Context) (out any, err error) {
	// input
	input, err := p.LifeCycle.HandleInput(ctx, p)
	if err != nil {
		err = errors.Wrap(ErrInvalidParameter, err.Error())
		return
	}
	// validate
	if p.validator == nil {
		p.validator = p.BuildValidator()
	}
	err = p.validator.Data(input).Run(ctx)
	if err != nil {
		err = errors.Wrap(ErrInvalidParameter, err.Error())
		return
	}
	// execute node and get raw output
	raw, err := ExecuteNode(ctx, p.Node, p.LifeCycle, input, p.Async)
	if err != nil {
		err = errors.Wrap(ErrNodeExecuteFailed, err.Error())
		return
	}
	// execute script
	if p.Script != "" {
		var v *g.Var
		v, err = p.Script.Execute(p.scriptFuncMap(ctx, raw), map[string]any{
			"output": out,
		})
		if err != nil {
			glog.Warningf(ctx, "proto script execute error: %v", err.Error())
			return
		}
		glog.Infof(ctx, "proto script execute result: %v", v.String())
	}
	// handle output
	out, err = p.LifeCycle.HandleOutput(ctx, p, raw)
	if err != nil {
		return
	}
	return
}

func (p *Proto) scriptFuncMap(ctx context.Context, rawOut map[string]any) (fm template.FuncMap) {
	fm = make(template.FuncMap)
	setDefaultFunctions(ctx, fm)
	fm[FuncNameGetResult] = func(k string) any {
		return rawOut[k]
	}
	fm[FuncNameSetResult] = func(k string, v any) error {
		rawOut[k] = v
		return nil
	}
	return
}
