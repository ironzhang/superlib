package bind

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestReq struct {
	A int
	B string
	C float64
	d int8
	E []string
	F []int
}

func TestBindQuery(t *testing.T) {
	req := TestReq{
		A: 1,
		B: "2",
		C: 3.3,
		d: 4,
		E: []string{"5", "6"},
		F: []int{7, 8, 9},
	}
	target := url.Values{
		"A": []string{"1"},
		"B": []string{"2"},
		"C": []string{"3.3"},
		"E": []string{"5", "6"},
		"F": []string{"7", "8", "9"},
	}

	values, err := BindQuery(req, "")
	if err != nil {
		t.Fatalf("bind query: %v", err)
	}
	if got, want := values, target; !reflect.DeepEqual(got, want) {
		t.Errorf("values: got %v, want %v", got, want)
	} else {
		t.Logf("values: %v", got)
	}

	type Param struct {
		Name  string `json:"name" form:"ddd"`
		Test1 string `json:"test,omitempty"`
		Test2 string `json:",omitempty"`
		Test3 string `json:",,"`
		Addr  string
		Addr2 string `json:"-"`
	}
	p := Param{
		Name:  "test",
		Test1: "11",
		Test2: "22",
		Test3: "33",
		Addr:  "127.0.0.1",
		Addr2: "ddd",
	}

	target2 := url.Values{
		"name":  []string{"test"},
		"test":  []string{"11"},
		"Test2": []string{"22"},
		"Test3": []string{"33"},
		"Addr":  []string{"127.0.0.1"},
	}

	values, err = BindQuery(p, "json")
	assert.NoError(t, err)
	if got, want := values, target2; !reflect.DeepEqual(got, want) {
		t.Errorf("values: got %v, want %v", got, want)
	} else {
		t.Logf("values: %v", got)
	}
}
