package internal

type FileExporter interface {
	Export(tasks []*Task) ([]byte, error)
}
