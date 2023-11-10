package heligo_test

import (
	"bytes"
	"context"
	"heligo"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParams(t *testing.T) {
	router := heligo.New()

	router.Handle("GET", "/param1/:p1/:p2/:p3/:p4/*p5", func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		// Param
		v := r.Param("p1")
		if v != "1" {
			t.Fail()
		}
		v = r.Param("p2")
		if v != "2" {
			t.Fail()
		}
		v = r.Param("p3")
		if v != "3" {
			t.Fail()
		}
		v = r.Param("p4")
		if v != "4" {
			t.Fail()
		}
		v = r.Param("p5")
		if v != "5/67" {
			t.Fail()
		}
		// ParamByPos
		param := r.ParamByPos(0)
		if param.Name != "p1" || param.Value != "1" {
			t.Fail()
		}
		param = r.ParamByPos(1)
		if param.Name != "p2" || param.Value != "2" {
			t.Fail()
		}
		param = r.ParamByPos(2)
		if param.Name != "p3" || param.Value != "3" {
			t.Fail()
		}
		param = r.ParamByPos(3)
		if param.Name != "p4" || param.Value != "4" {
			t.Fail()
		}
		param = r.ParamByPos(4)
		if param.Name != "p5" || param.Value != "5/67" {
			t.Fail()
		}
		// Params
		params := r.Params()
		if len(params) != 5 ||
			params[0].Name != "p1" || params[0].Value != "1" ||
			params[1].Name != "p2" || params[1].Value != "2" ||
			params[2].Name != "p3" || params[2].Value != "3" ||
			params[3].Name != "p4" || params[3].Value != "4" ||
			params[4].Name != "p5" || params[4].Value != "5/67" {
			t.Fail()
		}

		return 200, nil
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/param1/1/2/3/4/5/67", nil)
	router.ServeHTTP(w, r)
}

func TestReadJSON(t *testing.T) {
	router := heligo.New()

	type payload struct {
		String string
		Number int
		Bool   bool
		Array  []int
	}

	router.Handle("POST", "/read1", func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		var p payload
		r.ReadJSON(&p)
		if !reflect.DeepEqual(p, payload{"value", 42, true, []int{1, 2, 3}}) {
			t.Fail()
		}
		return 200, nil
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/read1", bytes.NewBuffer([]byte(`{"String": "value", "Number": 42, "Bool": true, "Array": [1,2,3]}`)))
	router.ServeHTTP(w, r)
}
