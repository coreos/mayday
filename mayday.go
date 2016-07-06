package main

import (
	"archive/tar"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/mayday/mayday"
)

const (
	dirPrefix     = "/mayday"
	defaultConfig = "/etc/mayday.conf"
)

type Config struct {
	Files    []mayday.File
	Commands []mayday.Command
}

func openConfig() (string, error) {
	configVar := os.Getenv("MAYDAY_CONFIG_FILE")
	configFile := strings.Split(configVar, "=")[0]

	if configFile == "" {
		configFile = defaultConfig
	}

	log.Printf("Reading configuration from %v\n", configFile)
	cbytes, err := ioutil.ReadFile(configFile)
	cstring := string(cbytes[:])
	return cstring, err
}

func readConfig(dat string) ([]mayday.File, []mayday.Command, error) {
	var c Config

	err := json.Unmarshal([]byte(dat), &c)
	if err != nil {
		log.Fatal(err)
	}
	return c.Files, c.Commands, nil
}

func main() {

	conf, err := openConfig()
	if err != nil {
		log.Fatal(err)
	}

	files, commands, err := readConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(files); i++ {
		files[i].Content, err = os.Open(files[i].Name)
		if err != nil {
			log.Fatal(err)
		}

		fi, err := os.Stat(files[i].Name)
		if err != nil {
			log.Fatal(err)
		}

		header, err := tar.FileInfoHeader(fi, files[i].Name)
		files[i].Header = header
		files[i].Header.Name = files[i].Name
		if err != nil {
			log.Fatal(err)
		}
	}

	journals, err := mayday.GetJournals()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().Format("200601021504.999999999")
	ws := os.TempDir() + dirPrefix + now

	var t mayday.Tar
	outputFile := ws + ".tar.gz"
	tarfile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer tarfile.Close()
	t.Init(tarfile)

	mayday.Run(t, files, commands, journals)
	t.Close()

	log.Printf("Output saved in %v\n", outputFile)
	log.Printf("All done!")

	return
}
