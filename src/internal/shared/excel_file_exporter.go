package shared

import (
	"log/slog"

	"github.com/xuri/excelize/v2"
)

type ExcelFileExporter struct {
	Logger *slog.Logger
}

func (e *ExcelFileExporter) Export(tasks []*Task) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	f.SetCellValue("Sheet1", "A2", "Hello world.")

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
