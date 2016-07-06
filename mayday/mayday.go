package mayday

import (
	"log"
)

func Run(t Tar, files []File, commands []Command, journals []Journal) error {

	// copy files and run commands
	for i := 0; i < len(files); i++ {
		log.Printf("Collecting file: %q\n", files[i].Name)
		t.Add(files[i].Content, files[i].Header)
	}

	for i := 0; i < len(commands); i++ {
		err := commands[i].Run()
		if err != nil {
			log.Print(err)
		}
		t.Add(commands[i].Content, commands[i].header())
	}

	for i := 0; i < len(journals); i++ {
		err := journals[i].Get()
		if err != nil {
			log.Print(err)
		}
		t.Add(journals[i].Content, journals[i].header())
	}

	// make shortlinks
	for i := 0; i < len(commands); i++ {
		t.MaybeMakeLink(commands[i].Link, "mayday_commands/"+commands[i].outputName())
	}

	for i := 0; i < len(files); i++ {
		t.MaybeMakeLink(files[i].Link, files[i].Name)
	}

	return nil
}
