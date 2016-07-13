package rkt

import (
	"archive/tar"
	"bytes"
	"github.com/coreos/mayday/mayday"
	"io"
	"time"
)

type App struct {
	Name      string
	ImageName string
	ImageId   string
}

// the struct holding all the data about a single pod
type Pod struct {
	initialized bool
	Uuid        string
	Apps        []App
	State       string
	Created     string // these could be dates, but would just mean processing overhead
	Started     string
	Network     string
}

type PodTarable struct {
	pod_uuid string
	link     string        // currently never set to anything
	content  *bytes.Buffer // the contents of the log, populated by Run()
}

func (p *Pod) GetTarables() []mayday.Tarable {
	var tarables []mayday.Tarable
	tarables = append(tarables, NewPodTarable(*p))
	tarables = append(tarables, NewAppTarable(*p))
	// tarables = append(tarables, NewLogTarable(&p))
	return tarables
}

func NewPodTarable(pod Pod) *PodTarable {
	pf := PodTarable{pod_uuid: pod.Uuid}
	var buffer bytes.Buffer

	buffer.WriteString("STATE\tCREATED\tSTARTED\tNETWORKS\n")
	buffer.WriteString(pod.State)
	buffer.WriteString("\t")
	buffer.WriteString(pod.Created)
	buffer.WriteString("\t")
	buffer.WriteString(pod.Started)
	buffer.WriteString("\t")
	buffer.WriteString(pod.Network)
	buffer.WriteString("\n")

	pf.content = &buffer
	return &pf
}

func (p *PodTarable) Header() *tar.Header {
	var header tar.Header
	header.Name = "/containers/" + p.pod_uuid + "/pod"
	header.Mode = 0666
	header.Size = int64(p.content.Len())
	header.ModTime = time.Now()

	return &header
}

func (p *PodTarable) Content() io.Reader {
	return p.content
}

func (p *PodTarable) Name() string {
	return p.pod_uuid
}

func (p *PodTarable) Link() string {
	return p.link
}
