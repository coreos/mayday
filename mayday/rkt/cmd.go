package rkt

import (
	"strings"
)

// TODO make this less horrible looking, and less fragile, if possible
func ProcessRktOutput(output string) []Pod {
	var pods []Pod
	var currpod Pod

	lines := strings.Split(output, "\n")

	for i, l := range lines {
		if i == 0 {
			continue
		}
		if l == "" {
			break
		}
		cols := strings.Split(l, "\t")
		if cols[0] != "" {
			// new pod!
			// if previous pod exists, save it to pods
			if currpod.initialized {
				pods = append(pods, currpod)
			}

			currpod = Pod{
				initialized: true,
				Uuid:        cols[0],
				State:       cols[len(cols)-4],
				Created:     cols[len(cols)-3],
				Started:     cols[len(cols)-2],
				Network:     cols[len(cols)-1]}
			currpod.Apps = []App{{
				Name:      cols[1],
				ImageName: cols[2],
				ImageId:   cols[len(cols)-5]}}
		} else {
			var newApp App
			for j, c := range cols {
				if c != "" {
					newApp = App{
						Name:      cols[j],
						ImageName: cols[j+1],
						ImageId:   cols[j+2],
					}
					break
				}
			}
			currpod.Apps = append(currpod.Apps, newApp)
		}
	}
	// save current pod, if it's been initialized
	if currpod.initialized {
		pods = append(pods, currpod)
	}
	return pods
}
