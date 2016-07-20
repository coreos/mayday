package main

import (
	"archive/tar"
	"log"
	"os"
	"time"

	"github.com/coreos/mayday/mayday"
	"github.com/coreos/mayday/mayday/docker"
	"github.com/coreos/mayday/mayday/rkt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	dirPrefix = "/mayday"
)

type Config struct {
	Files    []File    `mapstructure:"files"`
	Commands []Command `mapstructure:"commands"`
	Danger   bool
}

type File struct {
	Name string `mapstructure:"name"`
	Link string `mapstructure:"link"`
}

type Command struct {
	Args []string `mapstructure:"args"`
	Link string   `mapstructure:"link"`
}

func openFile(f File) (*mayday.MaydayFile, error) {
	content, err := os.Open(f.Name)
	if err != nil {
		return nil, err
	}
	defer content.Close()

	fi, err := os.Stat(f.Name)
	if err != nil {
		return nil, err
	}

	header, err := tar.FileInfoHeader(fi, f.Name)
	header.Name = f.Name
	if err != nil {
		return nil, err
	}

	opened := mayday.NewFile(content, header, f.Name, f.Link)
	return &opened, nil
}

func main() {
	pflag.BoolP("danger", "d", false, "collect potentially sensitive information (ex, container logs)")

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/mayday")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error reading config: %s \n", err)
	}

	// binds cli flag "danger" to viper config danger
	viper.BindPFlag("danger", pflag.Lookup("danger"))
	// cli arg takes precendence over anything in config files
	pflag.Parse()

	var tarables []mayday.Tarable

	var C Config

	// fill C with configuration data
	viper.Unmarshal(&C)

	journals, err := mayday.ListJournals()
	if err != nil {
		log.Fatal(err)
	}

	pods, rktLogs, err := rkt.GetPods()
	if err != nil {
		log.Println("Could not connect to rkt. Verify mayday has permissions to launch the rkt client.")
		log.Printf("Connection error: %s", err)
	}

	containers, dockerLogs, err := docker.GetContainers()
	if err != nil {
		log.Println("Could not connect to docker. Verify mayday has permissions to read /var/lib/docker.")
		log.Printf("Connection error: %s", err)
	}

	for _, f := range C.Files {
		mf, err := openFile(f)
		if err != nil {
			log.Printf("error opening %s: %s\n", f.Name, err)
		} else {
			tarables = append(tarables, mf)
		}
	}

	for _, c := range C.Commands {
		tarables = append(tarables, mayday.NewCommand(c.Args, c.Link))
	}

	for _, j := range journals {
		tarables = append(tarables, j)
	}

	for _, p := range pods {
		tarables = append(tarables, p)
	}

	for _, l := range rktLogs {
		tarables = append(tarables, l)
	}

	for _, c := range containers {
		tarables = append(tarables, c)
	}

	for _, l := range dockerLogs {
		tarables = append(tarables, l)
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
	t.Init(tarfile, now)

	mayday.Run(t, tarables)
	t.Close()

	log.Printf("Output saved in %v\n", outputFile)
	log.Printf("All done!")

	return
}
