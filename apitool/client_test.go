package apitool

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/raohwork/jsonapi"
)

// ParamGreeting represents parameters of Greeting API
type ParamGreeting struct {
	Name    string
	Surname string
}

// RespGreeting represents returned type of Greeting API
type RespGreeting struct {
	Name    string
	Surname string
	Greeted bool
}

// greeting is handler of Greeting API
func Greeting(r jsonapi.Request) (interface{}, error) {
	var p ParamGreeting
	if err := r.Decode(&p); err != nil {
		return nil, jsonapi.APPERR.SetData(
			"parameter format error",
		).SetCode("EParamFormat")
	}

	return RespGreeting{
		Name:    p.Name,
		Surname: p.Surname,
		Greeted: true,
	}, nil
}

// RunAPIServer creates and runs an API server
func RunAPIServer() *httptest.Server {
	http.Handle("/greeting", jsonapi.Handler(Greeting))
	return httptest.NewServer(http.DefaultServeMux)
}

func ExampleClient() {
	// start the API server
	server := RunAPIServer()
	defer server.Close()

	client := Call("POST", server.URL+"/greeting", nil)

	var resp RespGreeting
	err := client.Exec(ParamGreeting{Name: "John", Surname: "Doe"}, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Have we greeted to %s %s? %v",
		resp.Name, resp.Surname, resp.Greeted,
	)

	// output: Have we greeted to John Doe? true
}
