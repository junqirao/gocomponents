package procedure

import (
	"strings"
)

const (
	CtxKeyNodeName = "__node_name"
)

type (
	// Node of procedure, defines call flow and behavior.
	// nodes will be executed in level order synchronously
	// or asynchronously by NodeHandleFunc within input.
	Node struct {
		// Name of node, output will be overwritten if
		// using same name.
		Name string `json:"name"`
		// Meta data
		Meta map[string]any `json:"meta"`
		// Must execute success flag.
		// checks NodeLifeCycle.Execute return error
		Must bool `json:"must"`
		// Children nodes
		Children []*Node `json:"children"`
		// InputFilter of node, supported alias usage,
		// e.g. "aaa as bbb" empty means no filter,
		// all input will be passed
		InputFilter []string `json:"input_filter"`
		// Script in go template syntax
		Script Script `json:"script"`
		// directParent node
		directParent *Node
	}
)

func (n *Node) inherit(parent *Node) {
	if parent == nil {
		return
	}
	n.directParent = parent
}

func (n *Node) filterInput(in map[string]any) map[string]any {
	// in already cloned
	if n.InputFilter == nil {
		return in
	}
	m := make(map[string]any)
	for _, s := range n.InputFilter {
		if s == "" {
			continue
		}
		var (
			parts   = strings.Split(s, " as ")
			key, as string
		)

		if len(parts) == 1 {
			key = parts[0]
			as = key
		}

		if len(parts) > 1 {
			as = parts[1]
		}
		m[as] = in[key]
	}
	return m
}
