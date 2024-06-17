Package jsonapi is simple wrapper for buildin net/http package.
It aims to let developers build json-based web api easier.

[![Go Reference](https://pkg.go.dev/badge/github.com/raohwork/jsonapi.svg)](https://pkg.go.dev/github.com/raohwork/jsonapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/raohwork/jsonapi)](https://goreportcard.com/report/github.com/raohwork/jsonapi)

# Notable changes before v1

### v0.2.x

There're 2 breaking changes:

- logging middleware is removed.
- session middleware is removed.

And 3 tools are deprecated:

- apitool.Client
- apitool.Call
- apitool.ParseResponse

Deprecated components will be removed at v0.3.0

# Usage

Create an api handler is so easy:

```go
// HelloArgs is data structure for arguments passed by POST body.
type HelloArgs struct {
        Name string
        Title string
}

// HelloReply defines data structure this api will return.
type HelloReply struct {
        Message string
}

// HelloHandler greets user with hello
func HelloHandler(q jsonapi.Request) (res interface{}, err error) {
        // Read json objct from request.
        var args HelloArgs
        if err = q.Decode(&args); err != nil {
                // The arguments are not passed in JSON format, returns http
                // status 400 with {"errors": [{"detail": "invalid param"}]}
                err = jsonapi.E400.SetOrigin(err).SetData("invalid param")
                return
        }

        res = HelloReply{fmt.Sprintf("Hello, %s %s", args,Title, args.Name)}
        return
}
```

And this is how we do in main function:

```go
// Suggested usage
apis := []jsonapi.API{
    {"/api/hello", HelloHandler},
}
jsonapi.Register(http.DefaultMux, apis)

// old-school
http.Handle("/api/hello", jsonapi.Handler(HelloHandler))
```

Generated response is a subset of [jsonapi specs](https://jsonapi.org). Refer to
`handler_test.go` for examples.

### Call API with Go

```go
var result MyResult
param := MyParam { ... }
err := callapi.EP(http.MethodPost, uri).Call(ctx, param, &result)
var (
    eClient callapi.EClient
    eFormat callapi.EFormat
)
if err != nil {
    switch {
    case errors.As(err, &eClient):
        log.Fatal("failed to send request:", err)
    case errors.As(err, &eFormat):
        log.Fatal("failed to parse response:", err)
    default:
        log.Fatal("API returns error:", err)
    }
}
```

### Call API with TypeScript

There's a `fetch.ts` providing `grab<T>()` as simple wrapping around `fetch()`.
With following Go code:

```go
type MyStruct struct {
    X int  `json:"x"`
    Y bool `json:"y"`
}

func MyAPI(q jsonapi.Request) (ret interface{}, err error) {
    return []MyStruct{
	    {X: 1, Y: true},
		{X: 2},
	}, nil
}

function main() {
    apis := []jsonapi.API{
	    {"/my-api", MyAPI},
    }
	jsonapi.Register(http.DefaultMux, apis)
	http.ListenAndServe(":80", nil)
}
```

You might write TypeScript code like this:

```ts
export interface MyStruct {
  x?: number;
  y?: boolean;
}

export function getMyApi(): Promise<MyStruct[]> {
  return grab<MyStruct[]>('/my-api');
}

export function postMyApi(): Promise<MyStruct[]> {
  return grab<MyStruct[]>('/my-api', {
    method: 'POST',
	headers: {'Content-Type': 'application/json'},
	body: JSON.stringify('my data')
  });
}
```

# Middleware

```go
func runtimeLog(h jsonapi.Handler) jsonapi.Handler {
    return func(r jsonapi.Request) (data interface{}, err error) {
        log.Printf("entering path %s", r.R().URL.Path)
        begin := time.Now().Unix()
        data, err = h(r)
        log.Printf("processed path %s in %d seconds", r.R().URL.Path, time.Now().Unix()-begin)
        return
    }
}

func main() {
    jsonapi.With(runtimeLog).Register(http.DefaultMux, myapis)
    http.ListenAndServe(":80", nil)
}
```

There're few pre-defined middlewares in package `apitool`, see [godoc](https://pkg.go.dev/github.com/raohwork/jsonapi/apitool).

# License

Mozilla Public License Version 2.0

Copyright 2019- Ronmi Ren <ronmi.ren@gmail.com>
