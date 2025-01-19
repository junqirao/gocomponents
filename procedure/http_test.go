package procedure

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
)

type testLifeCycle struct {
	testNodeLifeCycle
	GoFrameHTTPInputHandler
}

func (t testLifeCycle) HandleInput(ctx context.Context, proto *Proto) (input map[string]any, err error) {
	fmt.Println("------------ ProtoLifeCycle Start ------------")
	fmt.Printf("[HandleInput  ][proto:%v] \n", proto.Name)
	return t.GoFrameHTTPInputHandler.HandleInput(ctx, proto)
}

func (t testLifeCycle) HandleOutput(ctx context.Context, proto *Proto, raw map[string]any) (out any, err error) {
	var (
		outNode   = proto.Node
		traversal func(n *Node)
	)
	traversal = func(n *Node) {
		if n == nil {
			return
		}
		for _, child := range n.Children {
			traversal(child)
		}
		if n.Meta.Get("final_result", false).Bool() {
			fmt.Println("get final result node: ", n.Name)
			outNode = n
			return
		}
	}
	traversal(outNode)
	out = raw[outNode.Name]
	fmt.Printf("[HandleOutput ][proto:%v] output: %+v\n", proto.Name, out)
	return
}

func (t testLifeCycle) Execute(ctx context.Context, node *Node, input map[string]any) (output any, err error) {
	fmt.Printf("[ExecuteNode  ][uuid:%v][node:%v] input: %+v\n", ctx.Value(t.ctxKeyId), ctx.Value(CtxKeyNodeName), input)
	return map[string]any{
		"data": gconv.Int(input["test"]) * node.Meta.Get("times", 1).Int(),
	}, nil
}

func TestHTTPProto(t *testing.T) {
	data := gfile.GetBytes("./test_data/test_http.yaml")
	proto, err := NewProtoFromYaml(data)
	if err != nil {
		t.Fatal(err)
		return
	}

	s := g.Server()
	s.SetPort(8000)
	lc := &testLifeCycle{
		testNodeLifeCycle:       testNodeLifeCycle{},
		GoFrameHTTPInputHandler: GoFrameHTTPInputHandler{},
	}
	s.Group("/", func(group *ghttp.RouterGroup) {
		err = BindGHTTPRouter(context.Background(), group, lc, func(r *ghttp.Request, res any, err error) {
			fmt.Println("------------ ProtoLifeCycle End ------------")
			if err != nil {
				r.Response.WriteHeader(http.StatusInternalServerError)
				r.Response.WriteJsonExit(g.Map{
					"message": err.Error(),
				})
				return
			}
			r.Response.WriteJsonExit(g.Map{
				"data": res,
			})
		}, proto)
		if err != nil {
			t.Fatal(err)
			return
		}
	})
	s.Run()
}
