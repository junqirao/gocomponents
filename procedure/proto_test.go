package procedure

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/os/gfile"
)

func TestNewProtoFromYaml(t *testing.T) {
	data := gfile.GetBytes("./test_data/test_http.yaml")
	proto, err := NewProtoFromYaml(data)
	if err != nil {
		t.Fatal(err)
		return
	}
	encode, err := gyaml.Encode(proto)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(string(encode))
}
