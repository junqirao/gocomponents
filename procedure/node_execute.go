package procedure

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/glog"
)

func Execute(ctx context.Context, root *Node, nlc NodeLifeCycle, input map[string]any, async ...bool) (output map[string]any, err error) {
	if root == nil {
		return
	}

	results := gmap.New(true)

	var (
		st    = []*Node{root}
		level []*Node
		exec  = execute
	)

	if len(async) > 0 && async[0] {
		exec = executeLevelAsync
	}

	// level traversal
	for len(st) > 0 {
		level = st
		st = nil

		in := gmap.NewStrAnyMap(true)
		in.Sets(input)
		// execute level
		if err = exec(ctx, level, nlc, in, results); err != nil {
			return
		}

		// collect next level
		for _, node := range level {
			for _, child := range node.Children {
				child.inherit(node)
				st = append(st, child)
			}
		}
	}

	output = results.MapStrAny()
	return
}

func executeLevelAsync(ctx context.Context, nodes []*Node, nlc NodeLifeCycle, input *gmap.StrAnyMap, results *gmap.Map) (err error) {
	if len(nodes) == 1 {
		return execute(ctx, nodes, nlc, input, results)
	}

	var (
		wg   = &sync.WaitGroup{}
		errs = gmap.New(true)
	)

	wg.Add(len(nodes))
	for _, node := range nodes {
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					errs.Set(node, rec.(error))
				}
				wg.Done()
			}()
			var (
				// filter input
				in = node.filterInput(input.MapStrAny())
				// child context
				childCtx = context.WithValue(ctx, CtxKeyNodeName, node.Name)
			)
			// before execute
			if err = nlc.BeforeExecute(childCtx, node, in); err != nil {
				return
			}
			// execute
			res, err := nlc.Execute(childCtx, node, in)
			// after execute
			nlc.AfterExecute(childCtx, node, res, err)
			// check error
			if err != nil {
				errs.Set(node, err)
				return
			}
			// execute script
			handleScript(childCtx, node, nlc, in, input, res, results)
			// set output
			results.Set(node.Name, res)
		}()
	}
	wg.Wait()

	errs.Iterator(func(k any, v any) bool {
		node := k.(*Node)
		if e, ok := v.(error); ok && e != nil {
			if node.Must {
				err = e
				return false
			}
			glog.Infof(ctx, "skip node %s execute error: %v", node.Name, e.Error())
		}
		return true
	})
	return
}

func execute(ctx context.Context, nodes []*Node, nlc NodeLifeCycle, input *gmap.StrAnyMap, results *gmap.Map) (err error) {
	for _, node := range nodes {
		var (
			// filter input
			in = node.filterInput(input.MapStrAny())
			// child context
			childCtx = context.WithValue(ctx, CtxKeyNodeName, node.Name)
		)

		// before execute
		if err = nlc.BeforeExecute(childCtx, node, in); err != nil {
			return
		}
		// execute
		var res any
		res, err = nlc.Execute(childCtx, node, in)
		// after execute
		nlc.AfterExecute(childCtx, node, res, err)
		// check error
		if err != nil {
			if node.Must {
				return
			}
			glog.Infof(childCtx, "skip node %s execute error: %v", node.Name, err.Error())
			continue
		}
		// set output
		results.Set(node.Name, res)
		// execute script
		handleScript(childCtx, node, nlc, in, input, res, results)
	}
	return
}

func handleScript(ctx context.Context, node *Node, nlc NodeLifeCycle,
	nodeIn map[string]any, input *gmap.StrAnyMap,
	output any, results *gmap.Map) {
	if node.Script == "" {
		return
	}

	fm := nlc.BeforeScript(ctx, node, nodeIn, output)
	setDefaultFunctions(ctx, node, fm, input, results)

	var (
		execRes = &ScriptExecuteResult{}
		out     *gvar.Var
		err     error
	)
	defer results.Set(fmt.Sprintf("%s_script", node.Name), execRes)

	out, err = node.Script.Execute(fm, map[string]any{
		"node":   node,
		"output": output,
	})
	execRes.OK = err == nil
	if err != nil {
		execRes.Error = err.Error()
		glog.Warningf(ctx, "node %s script execute error: %v", node.Name, err.Error())
		return
	}

	// must parse map
	execRes.Result = out.Map()
}
