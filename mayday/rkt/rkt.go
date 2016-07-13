package rkt

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"

	"github.com/coreos/mayday/mayday/rkt/v1alpha"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"log"
	"os/exec"
	"time"
)

const (
	timeout = time.Duration(5 * time.Second)
)

var (
	cmd exec.Cmd
)

type Pod struct {
	*v1alpha.Pod
	content *bytes.Buffer
	link    string
}

func (p *Pod) Content() io.Reader {
	if p.content == nil {
		marshalled, _ := yaml.Marshal(&p.Pod)
		p.content = bytes.NewBuffer(marshalled)
		log.Printf("collecting pod data: %s\n", p.Id)
	}
	return p.content
}

func (p *Pod) Header() *tar.Header {
	if p.content == nil {
		p.Content() // set Pod content to marshalled struct
	}
	var header tar.Header
	header.Name = "rkt/" + p.Id
	header.Size = int64(p.content.Len())
	header.Mode = 0666
	header.ModTime = time.Now()

	return &header
}

func (p *Pod) Name() string {
	return p.Id
}

func (p *Pod) Link() string {
	return p.link
}

var closeApi = func() error {
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Error killing Command: %v", err)
		return err
	}
	return nil
}

var startApi = func() error {
	// start rkt in api mode
	// successful startup is defined as rkt being alive after more than 200
	// milliseconds -- if it doesn't have permission to get the pod listing,
	// it generally closes far before then.
	p, err := exec.LookPath("rkt")
	if err != nil {
		log.Println("could not find rkt in PATH")
		return err
	}
	// Set up the actual Cmd to be run
	cmd = exec.Cmd{
		Path: p,
		Args: []string{"rkt", "api-service"},
	}

	cmd.Start()
	wc := make(chan error, 1)
	go func() {
		wc <- cmd.Wait()
	}()
	select {
	case <-time.After(200 * time.Millisecond):
		// since it's not ended yet, we're probably good to go
		return nil
	case err := <-wc:
		if err != nil {
			return err
		}
		return errors.New("rkt closed too quickly")
	}
}

var podsFromApi = func() ([]*v1alpha.Pod, error) {
	conn, err := grpc.Dial("localhost:15441", grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}

	c := v1alpha.NewPublicAPIClient(conn)
	defer conn.Close()

	podResp, err := c.ListPods(context.Background(), &v1alpha.ListPodsRequest{})
	return podResp.Pods, err
}

func GetPods() ([]*Pod, error) {
	var pods []*Pod

	err := startApi()
	if err != nil {
		return pods, err
	}
	defer closeApi()

	apiPods, err := podsFromApi()
	if err != nil {
		return pods, err
	}

	for _, p := range apiPods {
		pods = append(pods, &Pod{Pod: p})
	}

	return pods, nil
}
