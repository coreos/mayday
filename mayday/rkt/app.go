package rkt

import (
	"archive/tar"
	"bytes"
	"io"
	"time"
)

type AppTarable struct {
	pod_uuid string
	apps     []App
	link     string        // currently never set to anything
	content  *bytes.Buffer // the contents of the log, populated by Run()
}

func NewAppTarable(pod Pod) *AppTarable {
	af := AppTarable{pod_uuid: pod.Uuid, apps: pod.Apps}
	var buffer bytes.Buffer

	buffer.WriteString("APP\tIMAGE NAME\tIMAGE ID\n")
	for _, a := range af.apps {
		buffer.WriteString(a.Name)
		buffer.WriteString("\t")
		buffer.WriteString(a.ImageName)
		buffer.WriteString("\t")
		buffer.WriteString(a.ImageId)
		buffer.WriteString("\n")
	}

	af.content = &buffer
	return &af
}

func (a *AppTarable) Header() *tar.Header {
	var header tar.Header
	header.Name = "/containers/" + a.pod_uuid + "/apps"
	header.Mode = 0666
	header.Size = int64(a.content.Len())
	header.ModTime = time.Now()

	return &header
}

func (a *AppTarable) Content() io.Reader {
	return a.content
}

func (a *AppTarable) Name() string {
	return a.pod_uuid
}

func (a *AppTarable) Link() string {
	return a.link
}
