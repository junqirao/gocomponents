package procedure

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"text/template"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/util/gconv"
	uuid "github.com/satori/go.uuid"
)

type testNodeLifeCycle struct {
	ctxKeyId string
}

func (t *testNodeLifeCycle) BeforeExecute(ctx context.Context, node *Node, input map[string]any) (err error) {
	fmt.Println("------------ NodeLifeCycle Start ------------")
	fmt.Printf("[BeforeExecute][uuid:%v][node:%v] input: %+v\n", ctx.Value(t.ctxKeyId), ctx.Value(CtxKeyNodeName), input)
	return
}

func (t *testNodeLifeCycle) Execute(ctx context.Context, node *Node, input map[string]any) (output any, err error) {
	fmt.Printf("[ExecuteNode  ][uuid:%v][node:%v] input: %+v\n", ctx.Value(t.ctxKeyId), ctx.Value(CtxKeyNodeName), input)
	if node.Name == "7" {
		return node.Meta, nil
	}
	return node.Meta["data"], nil
}

func (t *testNodeLifeCycle) AfterExecute(ctx context.Context, node *Node, output any, err error) {
	fmt.Printf("[AfterExecute ][uuid:%v][node:%v] output: %+v, err: %v\n", ctx.Value(t.ctxKeyId), ctx.Value(CtxKeyNodeName), output, err)
}

func (t *testNodeLifeCycle) BeforeScript(ctx context.Context, node *Node, input map[string]any, output any) (fm template.FuncMap) {
	fmt.Printf("[BeforeScript ][uuid:%v][node:%v] input: %+v, output: %+v\n", ctx.Value(t.ctxKeyId), ctx.Value(CtxKeyNodeName), input, output)
	fm = template.FuncMap{
		"conv2str": func(v any) string {
			return gconv.String(v)
		},
		"md5": func(v any) string {
			return gmd5.MustEncrypt(gconv.String(v))
		},
	}
	fmt.Println("------------ NodeLifeCycle End ------------")
	return
}

func TestNode(t *testing.T) {
	root := &Node{
		Name: "1",
		Meta: map[string]any{"data": 1},
		Must: false,
		Children: []*Node{
			{
				Name: "2",
				Meta: map[string]any{"data": 2},
			},
			{
				Name: "3",
				Meta: map[string]any{"data": 3},
				Children: []*Node{
					{
						Name: "4",
						Meta: map[string]any{"data": 4},
					},
					{
						Name: "5",
						Meta: map[string]any{"data": 5},
						Children: []*Node{
							{
								Name: "6",
								Meta: map[string]any{"data": 6},
								Script: `
										{{get_result "1" | set_result "test"}}
										{{.node.Name | set_result "test1"}}
										{{conv2str .output | set_result .node.Name}}`,
							},
						},
					},
					{
						Name: "7",
						Meta: map[string]any{
							"code": 0,
							"data": map[string]any{
								"id":   123,
								"name": "abc",
								"list": []map[string]any{
									{
										"object_id": "0",
										"count":     2,
									},
									{
										"object_id": "1",
										"count":     5,
									},
								},
							},
						},
						Script: `
								{{ $id := .output.data.id | conv2str }}
								{{ set_map_value .output "data.id" $id }}
								{{ .output.data | info "%+v" }}
								{{ .output.data.not_exist.x | info "%+v" }}
								{{ range $i, $v := .output.data.list }}
									{{ if gt $v.count 2}}
										{{ set_map_value $v "test_field" true | info "%+v" }}
									{{ end }}
									{{ md5 $v.object_id | set_map_value $v "object_id" | info "%+v" }}
								{{ end }}`,
					},
				},
			},
		},
	}

	proto := Proto{Node: root}
	if err := proto.Check(context.Background()); err != nil {
		t.Fatal(err)
		return
	}

	bs, _ := gyaml.Encode(root)
	fmt.Println(string(bs))

	nlc := &testNodeLifeCycle{ctxKeyId: "uuid"}
	input := map[string]any{}

	ctx := context.Background()
	t.Log("=============== test sync ===============")
	m, err := ExecuteNode(context.WithValue(ctx, "uuid", uuid.NewV4().String()), root, nlc, input)
	if err != nil {
		t.Fatal(err)
		return
	}
	bs, _ = json.MarshalIndent(m, "", "    ")
	fmt.Println(string(bs))
	doCheck := func(m map[string]any) {
		if m["1"] != 1 {
			t.Fatal("expect 1, but got", m["1"])
		}
		if m["2"] != 2 {
			t.Fatal("expect 2, but got", m["2"])
		}
		if m["3"] != 3 {
			t.Fatal("expect 3, but got", m["3"])
		}
		if m["4"] != 4 {
			t.Fatal("expect 4, but got", m["4"])
		}
		if m["5"] != 5 {
			t.Fatal("expect 5, but got", m["5"])
		}
		if m["6"] != "6" {
			t.Fatal("expect 6, but got", m["6"])
		}
	}
	doCheck(m)
	t.Log("=============== test async ===============")
	m, err = ExecuteNode(context.WithValue(ctx, "uuid", uuid.NewV4().String()), root, nlc, input, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	bs, _ = json.MarshalIndent(m, "", "    ")
	fmt.Println(string(bs))
	doCheck(m)
}
