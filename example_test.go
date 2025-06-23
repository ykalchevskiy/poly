package poly_test

import (
	"encoding/json"
	"fmt"

	"github.com/ykalchevskiy/poly"
)

type IsAction interface {
	poly.TypeName
	IsAction()
}

type ActionDismiss struct{}

func (ActionDismiss) IsAction() {}

func (ActionDismiss) TypeName() string {
	return "dismiss"
}

type ActionDeepLink struct {
	URL string `json:"url"`
}

func (ActionDeepLink) IsAction() {}

func (ActionDeepLink) TypeName() string {
	return "deep-link"
}

type Action = poly.Poly[IsAction, poly.Types2[ActionDismiss, ActionDeepLink]]

func Example() {
	var action Action
	var err error
	var bytes []byte

	// create ActionDismiss
	if err = json.Unmarshal([]byte(`{"type": "dismiss"}`), &action); err != nil {
		panic(err)
	}

	bytes, err = json.Marshal(action)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	// create ActionDeepLink
	if err = json.Unmarshal([]byte(`{"type": "deep-link", "url": "url"}`), &action); err != nil {
		panic(err)
	}

	bytes, err = json.Marshal(action)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	// patch the existing ActionDeepLink
	if err = json.Unmarshal([]byte(`{"url": "url-2"}`), &action); err != nil {
		panic(err)
	}

	bytes, err = json.Marshal(action)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	// Output:
	// {"type":"dismiss"}
	// {"type":"deep-link","url":"url"}
	// {"type":"deep-link","url":"url-2"}
}
