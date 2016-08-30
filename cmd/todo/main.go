package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/luizbranco/todo"
	"github.com/luizbranco/todo/cli"
	"github.com/luizbranco/todo/gob"
)

const root = ""

var help = ` todo help:

  todo add  [task description]
  todo sub  [task id] [subtask description]
  todo edit [task id] [new task description]
  todo done [task id]
  todo rm   [task id]
`

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	file := filepath.Join(usr.HomeDir, ".todo")

	storage, err := gob.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		printTasks(storage)
	} else if len(os.Args) < 3 {
		fmt.Print(help)
	}

	switch os.Args[1] {
	case "add":
		text := strings.Join(os.Args[2:], " ")
		exit(storage.Create(root, text))
	case "edit":
		id := os.Args[2]
		text := strings.Join(os.Args[3:], " ")
		exit(storage.Update(id, text))
	case "sub":
		id := os.Args[2]
		text := strings.Join(os.Args[3:], " ")
		exit(storage.Create(id, text))
	case "done":
		id := os.Args[2]
		exit(storage.Finish(id))
	case "rm":
		id := os.Args[2]
		exit(storage.Delete(id))
	}
	printTasks(storage)
}

func printTasks(storage todo.TaskStorage) {
	printer := &cli.CLI{}
	all := storage.All()
	output := string(printer.TaskTree(all))
	fmt.Print(output)
	os.Exit(0)
}

func exit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
