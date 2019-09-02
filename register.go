// This file is part of jsonapi
//
// jsonapi is distributed in two licenses: The Mozilla Public License,
// v. 2.0 and the GNU Lesser Public License.
//
// jsonapi is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.
//
// See LICENSE for further information.

package jsonapi

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// HTTPMux abstracts http.ServeHTTPMux, so it will be easier to write tests
//
// Only needed methods are added here.
type HTTPMux interface {
	Handle(pattern string, handler http.Handler)
}

// API denotes how a json api handler registers to a servemux
type API struct {
	Pattern string
	Handler func(Request) (interface{}, error)
}

// Register helps you to register many APIHandlers to a http.ServeHTTPMux
func Register(mux HTTPMux, apis []API) {
	reg := http.Handle
	if mux != nil {
		reg = mux.Handle
	}

	for _, api := range apis {
		reg(api.Pattern, Handler(api.Handler))
	}
}

var reCamelToUL *regexp.Regexp
var reCamelToULExcepts *regexp.Regexp

func init() {
	reCamelToUL = regexp.MustCompile(
		`([^A-Z])([A-Z])|([A-Z0-9]+)([A-Z])`,
	)
	reCamelToULExcepts = regexp.MustCompile(
		`^[A-Z0-9]*$`,
	)
}

func findMatchedMethods(
	prefix string, handlers interface{}, conv func(string) string,
) []API {
	v := reflect.ValueOf(handlers)

	ret := make([]API, 0, v.NumMethod())

	for x, t := 0, v.Type(); x < v.NumMethod(); x++ {
		h, ok := v.Method(x).Interface().(func(Request) (interface{}, error))
		if !ok {
			// incorrect signature, skip
			continue
		}

		name := t.Method(x).Name
		if conv != nil {
			name = conv(name)
		}
		ret = append(ret, API{
			Pattern: prefix + "/" + name,
			Handler: h,
		})
	}

	return ret
}

// RegisterAll helps you to register all handler methods
//
// As using reflection to do the job, only exported methods with correct
// signature are registered.
//
// converter is used to convert from method name to url pattern, see
// CovertCamelToSnake for example.
//
// If converter is nil, name will leave unchanged.
func RegisterAll(
	mux HTTPMux, prefix string, handlers interface{},
	converter func(string) string,
) {
	Register(mux, findMatchedMethods(prefix, handlers, converter))
}

// ConvertCamelToSnake is a helper to convert CamelCase to camel_case
func ConvertCamelToSnake(name string) string {
	if reCamelToULExcepts.MatchString(name) {
		return strings.ToLower(name)
	}

	return strings.ToLower(
		reCamelToUL.ReplaceAllString(
			name,
			"${1}${3}_${2}${4}",
		),
	)
}

// ConvertCamelToSlash is a helper to convert CamelCase to camel/case
func ConvertCamelToSlash(name string) string {
	if reCamelToULExcepts.MatchString(name) {
		return strings.ToLower(name)
	}

	return strings.ToLower(
		reCamelToUL.ReplaceAllString(
			name,
			"${1}${3}/${2}${4}",
		),
	)
}
