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

	if err := f.SetSheetName("Sheet1", "Tasks"); err != nil {
		return nil, err
	}

	sid, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return nil, err
	}

	if err := f.SetRowStyle("Tasks", 1, 1, sid); err != nil {
		return nil, err
	}

	if err := f.SetSheetRow("Tasks", "A1", &[]any{
		"ID",
		"Name",
		"Expires At",
		"Expiring Info At",
		"Expired Info At",
		"Created At",
		"Updated At",
		"Completed At",
	}); err != nil {
		return nil, err
	}

	for i, t := range tasks {
		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			return nil, err
		}

		if err := f.SetSheetRow("Tasks", cell, &[]any{
			t.ID,
			t.Name,
			t.ExpiresAt,
			t.ExpiringInfoAt,
			t.ExpiredInfoAt,
			t.CreatedAt,
			t.UpdatedAt,
			t.CompletedAt,
		}); err != nil {
			return nil, err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
