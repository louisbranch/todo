package todo

type Task struct {
	ID    string
	Text  string
	Done  bool
	Tasks []Task
}

type TaskStorage interface {
	All() []Task
	Create(id, text string) error
	Update(id, text string) error
	Finish(id string) error
	Delete(id string) error
}
