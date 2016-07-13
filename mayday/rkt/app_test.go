package rkt

import (
	"bytes"
	"github.com/coreos/mayday/mayday"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppTarable(t *testing.T) {
	assert.Implements(t, (*mayday.Tarable)(nil), new(AppTarable))

	var pod Pod
	var app1 App
	var app2 App

	app1.Name = "app1"
	app1.ImageName = "image1"
	app1.ImageId = "111"

	app2.Name = "app2"
	app2.ImageName = "image2"
	app2.ImageId = "222"

	pod.Uuid = "abc123"
	pod.Apps = []App{app1, app2}

	af := NewAppTarable(pod)

	content := new(bytes.Buffer)
	content.ReadFrom(af.Content())

	res := "APP\tIMAGE NAME\tIMAGE ID\n" +
		"app1\timage1\t111\n" +
		"app2\timage2\t222\n"

	assert.Equal(t, content.String(), res)
}

func TestAppHeader(t *testing.T) {

	pod := Pod{
		Uuid: "abc123",
		Apps: []App{{Name: "app1", ImageName: "image1", ImageId: "111"}},
	}

	af := NewAppTarable(pod)

	assert.Equal(t, af.Name(), "abc123")

	hdr := af.Header()
	assert.Equal(t, hdr.Name, "/containers/abc123/apps")
	assert.Equal(t, hdr.Size, int64(af.content.Len()))
}
