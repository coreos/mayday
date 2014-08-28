package mayday

import (
	"reflect"
	"testing"
)

func TestAThing(t *testing.T) {
	if AThing() != "MAYDAY!" {
		t.Fatalf("incorrect output")
	}
}

func TestInMemWorkspace(t *testing.T) {

	w := InMemWorkspace{
		contents: make(map[string][]byte),
	}

	resolv := []byte("nameserver 8.8.8.8")

	w.Set("resolv", resolv)

	contents := w.Get("resolv")

	if !reflect.DeepEqual(contents, []byte("nameserver 8.8.8.8")) {
		t.Fatal("workspace set/get failed")
	}
}
