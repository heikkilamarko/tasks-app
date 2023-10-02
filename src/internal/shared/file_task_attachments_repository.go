package shared

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

type FileTaskAttachmentsRepository struct {
	Config *Config
}

func (repo *FileTaskAttachmentsRepository) SaveAttachments(ctx context.Context, taskID int, fileHeaders []*multipart.FileHeader) error {
	if err := repo.ensureTaskDir(taskID); err != nil {
		return err
	}

	for _, fileHeader := range fileHeaders {
		srcFile, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(repo.getAttachmentPath(taskID, fileHeader.Filename))
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *FileTaskAttachmentsRepository) DeleteAttachments(ctx context.Context, taskID int, deleted map[int]string) error {
	for _, name := range deleted {
		if err := os.Remove(repo.getAttachmentPath(taskID, name)); err != nil {
			return err
		}
	}
	return nil
}

func (repo *FileTaskAttachmentsRepository) DeleteTask(ctx context.Context, taskID int) error {
	if err := os.RemoveAll(repo.getTaskPath(taskID)); err != nil {
		return err
	}
	return nil
}

func (repo *FileTaskAttachmentsRepository) ensureTaskDir(taskID int) error {
	if err := os.MkdirAll(repo.getTaskPath(taskID), 0755); err != nil {
		return err
	}
	return nil
}

func (repo *FileTaskAttachmentsRepository) getAttachmentPath(taskID int, name string) string {
	return filepath.Join(repo.getTaskPath(taskID), name)
}

func (repo *FileTaskAttachmentsRepository) getTaskPath(taskID int) string {
	return filepath.Join(repo.Config.Shared.AttachmentsPath, strconv.Itoa(taskID))
}
