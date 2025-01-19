package procedure

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

// Proto of procedure
type Proto struct {
	// LifeCycle of proto
	LifeCycle `json:"-" yaml:"-"`
	// Name of proto
	Name string `json:"name"`
	// Meta data
	Meta Meta `json:"meta"`
	// Parameters define
	Parameters []*Parameter `json:"parameters"`
	// Node of procedure
	Node *Node `json:"node"`
	// Async is a flag that indicates whether
	// the node is executed asynchronously
	Async bool `json:"async"`
	// Script in go template syntax
	Script Script `json:"script"`
	// validator
	validator *gvalid.Validator
}

func NewProtoFromYaml(data []byte) (p *Proto, err error) {
	p = &Proto{}
	err = gyaml.DecodeTo(data, p)
	return
}

func (p *Proto) BuildValidator() (v *gvalid.Validator) {
	var (
		message = make(map[string]any)
		rules   = make(map[string]string)
	)
	v = gvalid.New()

	for _, parameter := range p.Parameters {
		if parameter.Validate != nil {
			rules[parameter.Name] = parameter.Validate.Rule
			message[parameter.Name] = parameter.Validate.Message
		}
	}

	if len(rules) > 0 {
		v = v.Rules(rules)
	}
	if len(message) > 0 {
		v = v.Messages(message)
	}
	return
}

func (p *Proto) Check(ctx context.Context) (err error) {
	// check proto
	if p == nil {
		return errors.New("empty proto")
	}
	if p.LifeCycle == nil {
		return fmt.Errorf("life cycle not set: %s", p.Name)
	}
	// check node
	if p.Node == nil {
		return errors.New("no executable node")
	}

	var (
		nodeCount int
		errMsg    []string
		traversal func(n *Node, i *int)
	)

	traversal = func(n *Node, i *int) {
		if n == nil {
			return
		}
		for i, child := range n.Children {
			idx := gconv.Int(i)
			traversal(child, &idx)
		}
		nodeCount++
		if n.Name == "" {
			pos := "root"
			if n.directParent != nil {
				pos = fmt.Sprintf("%s.children", n.directParent.Name)
				if i != nil {
					pos = fmt.Sprintf("%s[%d]", pos, *i)
				}
			}
			errMsg = append(errMsg, fmt.Sprintf("%s.name cannot be empty", pos))
			return
		}
	}
	traversal(p.Node, nil)
	if len(errMsg) > 0 {
		err = errors.New(strings.Join(errMsg, ";"))
		glog.Infof(ctx, "proto check node failed %d/%d. errors: %s", nodeCount-len(errMsg), nodeCount, err.Error())
		return
	}
	glog.Infof(ctx, "proto check node success: %d", nodeCount)

	// check parameters
	for i, parameter := range p.Parameters {
		if parameter.Name == "" {
			err = fmt.Errorf("parameters[%d].name cannot be empty", i)
			return
		}
	}
	return
}
