package cli

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/luizbranco/todo"
)

type CLI struct {
	buf bytes.Buffer
}

func (cli *CLI) TaskTree(tasks []todo.Task) []byte {
	fmt.Fprintln(&cli.buf, "Tasks:\n")
	cli.print(0, tasks)

	return cli.buf.Bytes()
}

func (cli *CLI) print(depth int, tasks []todo.Task) {
	depth++
	for _, t := range tasks {
		prefix := " "
		if t.Done {
			prefix = "\u2713"
		}
		tabs := strings.Repeat(" ", depth-1)
		fmt.Fprintf(&cli.buf, "%s%s[%s] %s\n", prefix, tabs, t.ID, t.Text)
		cli.print(depth, t.Tasks)
	}
}
