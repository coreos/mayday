package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/mayday/mayday"
)

const (
	dirPrefix = "mayday"
)

var (
	files = []mayday.File{
		{
			Path: "/proc/vmstat",
		},
		{
			Path: "/proc/meminfo",
			Link: "meminfo",
		},
		{
			Path: "/etc/os-release",
			Link: "os-release",
		},
		{
			Path: "/proc/mounts",
			Link: "mounts",
		},
	}
	commands = []mayday.Command{
		{
			Args: []string{"hostname"},
			Link: "hostname",
		},
		{
			Args: []string{"date"},
			Link: "date",
		},
		{
			Args: []string{"systemd-cgls"},
		},
		{
			Args: []string{"systemd-cgtop", "-n1"},
		},
		{
			Args: []string{"ps", "fauxwww"},
			Link: "ps",
		},
		{
			Args: []string{"lsmod"},
			Link: "lsmod",
		},
		{
			Args: []string{"lspci"},
			Link: "lspci",
		},
		{
			Args: []string{"lsof", "-b", "-M", "-n", "-l"},
			Link: "lsof",
		},
		{
			Args: []string{"blkid"},
		},
		{
			Args: []string{"btrfs", "fi", "show"},
		},
		{
			Args: []string{"df", "-al"},
			Link: "df",
		},
		{
			Args: []string{"df", "-ali"},
		},
		{
			Args: []string{"free", "-m"},
			Link: "free",
		},
	}
)

func main() {
	now := time.Now().Format("200601021504-")
	ws, err := ioutil.TempDir("", dirPrefix+now)
	if err != nil {
		log.Fatalf("error creating output directory: %v", err)
	}
	log.Printf("Using workspace directory %q", ws)

	if err := os.MkdirAll(path.Join(ws, mayday.OutputDir), 0700); err != nil {
		log.Fatalf("error creating command output directory: %v", err)
	}
	// TODO(jonboulle): parallelise
	// TODO(jonboulle): handle errors
	for _, f := range files {
		if err := f.Collect(ws); err != nil {
			fmt.Println(err)
		}
	}
	for _, c := range commands {
		if err := c.Run(ws); err != nil {
			fmt.Println(err)
		}
	}

	// Build the output tar file
	tfn := ws + ".tar.gz"
	f, err := os.OpenFile(tfn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("error opening file for output: %v", err)
	}
	gzw := gzip.NewWriter(f)
	if err := buildTarFile(gzw, ws); err != nil {
		log.Fatal(err.Error())
	}
	if err := gzw.Close(); err != nil {
		log.Fatalf("error closing zipfile: %v", err)
	}
	fmt.Printf("Output saved in %v\n", tfn)
	fmt.Println("All done!")
}

// buildTarFile writes a tar-formatted file containing the recursive contents
// of the given directory to the provided io.Writer
func buildTarFile(w io.Writer, dir string) error {
	tw := tar.NewWriter(w)
	wf := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			// TODO(jonboulle): pass instead of failing?
			return err
		}
		var link string
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err = os.Readlink(p)
			link = strings.TrimPrefix(link, dir+"/")
		}
		if err != nil {
			return err
		}
		h, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		h.Name = strings.TrimPrefix(p, path.Dir(dir)+"/")
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if fi.IsDir() || link != "" {
			return nil
		}
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("could not open file for tar: %v", err)
		}
		if _, err := io.Copy(tw, f); err != nil {
			return fmt.Errorf("could not copy file: %v", err)
		}
		return nil
	}
	if err := filepath.Walk(dir, wf); err != nil {
		return fmt.Errorf("error collating output: %v", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("error closing tarfile: %v", err)
	}
	return nil
}
