package heligo

import (
	"encoding/json"
	"net/http"
)

const MAXPARAMS = 16

type params struct {
	names    [MAXPARAMS]*string
	valuePos [MAXPARAMS]int
	count    int
}

// Request embeds the standard http.Request and the URL parameters in a compressed format
type Request struct {
	*http.Request
	params params
}

func (r *Request) paramValue(i int) string {
	pos := r.params.valuePos[i] >> 16
	end := int(int16(r.params.valuePos[i]))
	if end == 0 {
		return r.Request.URL.Path[pos:]
	} else {
		return r.Request.URL.Path[pos : pos+end]
	}
}

// Param returns a URL parameter by name.
// It returns an empty string if the requested parameter is not found.
func (r *Request) Param(name string) string {
	for i := 0; i < r.params.count; i++ {
		if *r.params.names[i] == name {
			return r.paramValue(i)
		}
	}
	return ""
}

type Param struct {
	Name  string
	Value string
}

// ParamByPos gets a URL parameter by position in the URL (0-based)
func (r *Request) ParamByPos(i int) Param {
	var param Param
	param.Name = *r.params.names[i]
	param.Value = r.paramValue(i)
	return param
}

// Params gets all the URL parameters for the request
func (r *Request) Params() []Param {
	var params []Param
	for i := 0; i < r.params.count; i++ {
		params = append(params, r.ParamByPos(i))
	}
	return params
}

// ReadJSON decodes the JSON in the body into the value pointed by obj
func (r *Request) ReadJSON(obj any) error {
	decoder := json.NewDecoder(r.Request.Body)
	decoder.UseNumber()
	return decoder.Decode(obj)
}
