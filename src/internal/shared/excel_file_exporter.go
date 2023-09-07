package shared

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/xuri/excelize/v2"
)

type ExcelFileExporter struct {
	Logger *slog.Logger
}

func (e *ExcelFileExporter) ExportTasks(w http.ResponseWriter, tasks []*Task, name string) error {
	f := excelize.NewFile()
	defer f.Close()

	if err := f.SetSheetName("Sheet1", "Tasks"); err != nil {
		return err
	}

	sid, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return err
	}

	if err := f.SetRowStyle("Tasks", 1, 1, sid); err != nil {
		return err
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
		return err
	}

	for i, t := range tasks {
		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			return err
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
			return err
		}
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", name))

	return f.Write(w)
}
