package mayday

import (
	"io/ioutil"
	"log"
)

func AThing() string {
	return "MAYDAY!"
}

type Workspace interface {
	Set(string, []byte)
	Get(string) []byte
}

type InMemWorkspace struct {
	contents map[string][]byte
}

func (w *InMemWorkspace) Get(path string) []byte {
	content := w.contents[path]
	return content
}

func (w *InMemWorkspace) Set(path string, contents []byte) error {
	w.contents[path] = contents
	return nil
}

func GrabFile(w *InMemWorkspace, path string) error {

	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
		return err
	}

	w.Set(path, data)

	return nil
}
