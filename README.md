# poly

[![Go Reference](https://pkg.go.dev/badge/github.com/ykalchevskiy/poly.svg)](https://pkg.go.dev/github.com/ykalchevskiy/poly)
[![Go Report Card](https://goreportcard.com/badge/github.com/ykalchevskiy/poly)](https://goreportcard.com/report/github.com/ykalchevskiy/poly)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ykalchevskiy/poly)](https://golang.org/dl/)

Poly allows polymorphic JSON serialization in Go.

Example:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/ykalchevskiy/poly"
)

type IsAction interface {
	poly.TypeName // just for convenience
	IsAction()
}

type ActionDismiss struct{}

func (ActionDismiss) IsAction()        {}
func (ActionDismiss) TypeName() string { return "dismiss" }

type ActionDeepLink struct {
	URL string `json:"url"`
}

func (ActionDeepLink) IsAction()        {}
func (ActionDeepLink) TypeName() string { return "deep-link" }

type Action = poly.Poly[IsAction, poly.Types2[ActionDismiss, ActionDeepLink]]

func main() {
	var action Action
	var bytes []byte

	_ = json.Unmarshal([]byte(`{"type":"dismiss"}`), &action)
	fmt.Printf("%T\n", action.Value) // ActionDismiss
	bytes, _ = json.Marshal(action)
	fmt.Println(string(bytes)) // {"type":"dismiss"}

	_ = json.Unmarshal([]byte(`{"type":"deep-link","url":"url1"}`), &action)
	fmt.Printf("%T\n", action.Value) // ActionDeepLink
	bytes, _ = json.Marshal(action)
	fmt.Println(string(bytes)) // {"type":"deep-link","url":"url1"}

	// patch the existing ActionDeepLink
	_ = json.Unmarshal([]byte(`{"url":"url2"}`), &action)
	bytes, _ = json.Marshal(action)
	fmt.Println(string(bytes)) // {"type":"deep-link","url":"url2"}
}
```

See also `polygen` generator for more features: https://github.com/ykalchevskiy/polygen
