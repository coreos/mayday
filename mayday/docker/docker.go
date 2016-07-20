package docker

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"github.com/coreos/mayday/mayday"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	dockerDir = "/var/lib/docker/containers"
)

type DockerContainer struct {
	containerId string        // container id
	file        io.Reader     // config file -- /var/lib/docker/containers/{uuid}/config.v2.json
	content     *bytes.Buffer // a Buffer containing the contents of the file
	link        string        // a link to make in the root of the tarball
}

func New(f io.Reader, uuid string) DockerContainer {
	dc := DockerContainer{containerId: uuid, file: f}
	return dc
}

func (d *DockerContainer) Content() io.Reader {
	if d.content != nil {
		return d.content
	}

	// unmarshal docker config into interface
	var config interface{}

	fileContent, err := ioutil.ReadAll(d.file)
	if err != nil {
		log.Printf("error reading docker container configuration: %s", err)
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		log.Printf("error reading docker container configuration: %s", err)
	}
	configInterface := config.(map[string]interface{})

	if !viper.GetBool("danger") {
		configMap := configInterface["Config"].(map[string]interface{})
		env := configMap["Env"].([]interface{})
		var newEnv []string
		for _, e := range env {
			varSplit := strings.SplitAfterN(e.(string), "=", 2)
			newEnv = append(newEnv, varSplit[0]+"scrubbed by mayday")
		}
		configMap["Env"] = newEnv
	}

	byteContent, err := json.MarshalIndent(configInterface, "", "  ")
	if err != nil {
		log.Print("error saving docker container configuration!")
	}

	d.content = bytes.NewBuffer(byteContent)
	return d.content
}

func (d *DockerContainer) Header() *tar.Header {
	if d.content == nil {
		d.Content()
	}

	var header tar.Header
	header.Name = "/docker/" + d.containerId
	header.Size = int64(d.content.Len())
	header.Mode = 0666
	header.ModTime = time.Now()

	return &header
}

func (d *DockerContainer) Name() string {
	return d.containerId
}

func (d *DockerContainer) Link() string {
	return d.link
}

func getLogs(containers []*DockerContainer) []*mayday.Command {
	var logs []*mayday.Command
	if viper.GetBool("danger") {
		log.Println("Danger mode activated. Dump will include docker container logs and environment variables, which may contain sensitive information.")
		if len(containers) != 0 {
			for _, c := range containers {
				logcmd := []string{"docker", "logs", c.Name()}
				cmd := mayday.NewCommand(logcmd, "")
				cmd.Output = "/docker/" + c.Name() + ".log"
				logs = append(logs, cmd)
			}
		}
	}
	return logs
}

func GetContainers() ([]*DockerContainer, []*mayday.Command, error) {
	var containers []*DockerContainer
	var logs []*mayday.Command

	files, err := ioutil.ReadDir(dockerDir)
	if err != nil {
		return containers, logs, err
	}

	for _, file := range files {
		f, err := os.Open(dockerDir + "/" + file.Name() + "/config.v2.json")
		if err != nil {
			log.Printf("unable to read config for container %s: %s", file.Name(), err)
			continue
		}
		dc := New(f, file.Name())
		containers = append(containers, &dc)
	}

	logs = getLogs(containers)

	return containers, logs, nil
}
