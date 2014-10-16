package mayday

import (
	"fmt"
	"log"
	"os"
	"path"
)

func maybeCreateLink(src, dst, workspace string) error {
	if src == "" {
		return nil
	}
	lp := path.Join(workspace, src)
	log.Printf("Creating link (%s -> %s)", lp, dst)
	if err := os.Symlink(dst, lp); err != nil {
		return fmt.Errorf("error creating link: %v", err)
	}
	return nil
}
