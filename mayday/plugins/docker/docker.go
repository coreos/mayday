package docker

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/mayday/mayday/plugins/command"
	"github.com/spf13/viper"
)

const (
	dockerDir = "/var/lib/docker/containers"
)

// errUnrecognizedFormat is the content that will be returned when a given bit
// of content is malformed in a config v2 map
// In a better world, we would return (bytes, err) for Content, but how tarable
// is designed, it's expected that failures don't happen as best I can tell.
var errUnrecognizedFormat = []byte("unrecognized docker config format")

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

func (d *DockerContainer) Content() *bytes.Buffer {
	if d.content != nil {
		return d.content
	}

	// a docker v2 config is a map of string->???
	var config map[string]json.RawMessage

	fileContent, err := ioutil.ReadAll(d.file)
	if err != nil {
		log.Printf("error reading docker container configuration: %s", err)
		return bytes.NewBuffer(errUnrecognizedFormat)
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		log.Printf("error reading docker container configuration: %s", err)
		return bytes.NewBuffer(errUnrecognizedFormat)
	}
	configData, ok := config["Config"]
	if !ok {
		log.Printf("unrecognized docker config for container %q: no Config key", d.containerId)
		return bytes.NewBuffer(errUnrecognizedFormat)
	}

	if !viper.GetBool("danger") {
		// config.Config is also of type string->???; delay ??? decoding
		var configConfig map[string]json.RawMessage
		if err := json.Unmarshal(configData, &configConfig); err != nil {
			log.Printf("unrecognized docker config.Config for container %q: %v", d.containerId, err)
			return bytes.NewBuffer(errUnrecognizedFormat)
		}
		if configEnv, ok := configConfig["Env"]; ok {
			var envString []string
			if err := json.Unmarshal(configEnv, &envString); err != nil {
				log.Printf("error parsing docker environment variables for %q: %v", d.containerId, err)
				return bytes.NewBuffer([]byte("could not unmarshal Env"))
			}
			var newEnv []string
			for _, e := range envString {
				varSplit := strings.SplitAfterN(e, "=", 2)
				newEnv = append(newEnv, varSplit[0]+"scrubbed by mayday")
			}

			newEnvRaw, err := json.Marshal(newEnv)
			if err != nil {
				log.Print("error marshalling new env: %v", err)
				return bytes.NewBuffer([]byte("could not unmarshal Env"))
			}
			configConfig["Env"] = newEnvRaw
			configConfigRaw, err := json.Marshal(configConfig)
			if err != nil {
				log.Print("error marshalling new config: %v", err)
				return bytes.NewBuffer([]byte("json marshal error"))
			}
			config["Config"] = configConfigRaw
		}
	}

	byteContent, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Print("error saving docker container configuration!")
		return bytes.NewBuffer([]byte("json marshal error :["))
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

func getLogs(containers []*DockerContainer) []*command.Command {
	var logs []*command.Command
	if viper.GetBool("danger") {
		log.Println("Danger mode activated. Dump will include docker container logs and environment variables, which may contain sensitive information.")
		if len(containers) != 0 {
			for _, c := range containers {
				logcmd := []string{"docker", "logs", c.Name()}
				cmd := command.New(logcmd, "")
				cmd.Output = "/docker/" + c.Name() + ".log"
				logs = append(logs, cmd)
			}
		}
	}
	return logs
}

func GetContainers() ([]*DockerContainer, []*command.Command, error) {
	var containers []*DockerContainer
	var logs []*command.Command

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
