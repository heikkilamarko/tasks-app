package ui

import "tasks-app/internal/shared"

type AttachmentsUpdate struct {
	Inserted []string
	Deleted  map[int]string
}

func BuildAttachmentsUpdate(current []*shared.Attachment, updated []string) *AttachmentsUpdate {
	var inserted []string

	deleted := make(map[int]string)
	for _, c := range current {
		deleted[c.ID] = c.FileName
	}

	for _, uName := range updated {
		var id int
		var exists bool
		for cID, cName := range deleted {
			if cName == uName {
				exists = true
				id = cID
				break
			}
		}

		if exists {
			delete(deleted, id)
		} else {
			inserted = append(inserted, uName)
		}
	}

	return &AttachmentsUpdate{inserted, deleted}
}
