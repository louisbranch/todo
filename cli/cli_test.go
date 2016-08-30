package cli

import (
	"bytes"
	"testing"

	"github.com/luizbranco/todo"
)

func TestPrint(t *testing.T) {
	tasks := []todo.Task{
		todo.Task{
			ID:   "1",
			Text: "First",
			Tasks: []todo.Task{
				todo.Task{
					ID:   "3",
					Text: "Third",
					Tasks: []todo.Task{
						todo.Task{ID: "5", Text: "Fifth"},
					},
				},
				todo.Task{
					ID:   "4",
					Text: "Forth",
					Done: true,
				},
			},
		},
		todo.Task{ID: "2", Text: "Second"},
	}
	cli := &CLI{}

	text := `Tasks:

 [1] First
  [3] Third
   [5] Fifth
✓ [4] Forth
 [2] Second
`

	got := cli.TaskTree(tasks)
	want := []byte(text)

	if !bytes.Equal(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}

}
