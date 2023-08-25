package internal

type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Delete(id int) error
	GetByID(id int) (*Task, error)
	GetAll() ([]*Task, error)
}
