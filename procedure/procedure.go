package procedure

import (
	"context"
	"text/template"
)

type (
	// LifeCycle interface
	LifeCycle interface {
		NodeLifeCycle
		// HandleInput defines how to handle input
		HandleInput(ctx context.Context, proto *Proto) (input map[string]any, err error)
		// HandleOutput defines how to handle output
		HandleOutput(ctx context.Context, proto *Proto, raw map[string]any) (out any, err error)
	}
	// NodeLifeCycle defines node life cycle
	NodeLifeCycle interface {
		// BeforeExecute defines aspect before execute
		// input is copy of ExecuteNode input and filtered.
		BeforeExecute(ctx context.Context, node *Node, input map[string]any) (err error)
		// Execute defines how to execute a node.
		// input is copy of ExecuteNode input and
		// filtered. after execute, output will be
		// saved to ExecuteNode output in key node.Name
		Execute(ctx context.Context, node *Node, input map[string]any) (output any, err error)
		// AfterExecute defines aspect after execute.
		AfterExecute(ctx context.Context, node *Node, output any, err error)
		// BeforeScript defines parameters to build the
		// template.Template, returns the template.FuncMap
		// and the template execute input.
		// this method called after AfterExecute only if
		// node.Script is not empty.
		BeforeScript(ctx context.Context, node *Node, input map[string]any, output any) (fm template.FuncMap)
	}
)
