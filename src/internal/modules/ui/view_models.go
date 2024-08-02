package ui

import (
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"tasks-app/internal/shared"
	"time"

	"github.com/gorilla/csrf"
)

type LanguageRequest struct {
	Language string
}

type ThemeRequest struct {
	Theme string
}

type TimezoneRequest struct {
	Timezone string
}

type TaskRequest struct {
	ID int
}

type TaskAttachmentRequest struct {
	ID   int
	Name string
}

type NewTaskRequest struct {
	Name        string
	ExpiresAt   *time.Time
	Attachments *AttachmentsRequest
}

type UpdateTaskRequest struct {
	ID          int
	Name        string
	ExpiresAt   *time.Time
	Attachments *AttachmentsRequest
}

type AttachmentsRequest struct {
	Names []string
	Files []*multipart.FileHeader
}

type TasksResponse struct {
	UI            *UIModel
	Tasks         []*shared.Task
	IsCreatingNew bool
}

type TaskResponse struct {
	UI   *UIModel
	Task *shared.Task
}

type UIModel struct {
	Title     string
	Theme     string
	Language  string
	Languages []string
	T         map[string]string
	Location  *time.Location
	Timezones []string
	UserID    string
	UserName  string
	HubURL    string
	CSRFToken string
	CSRFField template.HTML
}

func NewUIModel(r *http.Request) *UIModel {
	var userID, userName string
	user, _ := shared.GetUserContext(r.Context())
	if user != nil {
		userID = user.ID
		userName = user.Name
	}

	language := GetLanguage(r)

	return &UIModel{
		Title:     "",
		Theme:     GetTheme(r),
		Language:  language,
		Languages: SupportedLanguages,
		T:         Translations[language],
		Location:  GetLocation(r),
		Timezones: SupportedTimezones,
		UserID:    userID,
		UserName:  userName,
		HubURL:    os.Getenv("APP_UI_HUB_URL"),
		CSRFToken: csrf.Token(r),
		CSRFField: csrf.TemplateField(r),
	}
}

func NewTasksResponse(r *http.Request, tasks []*shared.Task) *TasksResponse {
	return &TasksResponse{
		UI:    NewUIModel(r),
		Tasks: tasks,
	}
}

func NewTaskResponse(r *http.Request, task *shared.Task) *TaskResponse {
	return &TaskResponse{
		UI:   NewUIModel(r),
		Task: task,
	}
}

func ParseSetLanguageRequest(r *http.Request) (*LanguageRequest, error) {
	var errs []error

	lang, err := ParseLanguage(r.FormValue("language"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &LanguageRequest{lang}, nil
}

func ParseSetThemeRequest(r *http.Request) (*ThemeRequest, error) {
	var errs []error

	theme, err := ParseTheme(r.FormValue("theme"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &ThemeRequest{theme}, nil
}

func ParseSetTimezoneRequest(r *http.Request) (*TimezoneRequest, error) {
	var errs []error

	tz, err := ParseTimezone(r.FormValue("timezone"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &TimezoneRequest{tz}, nil
}

func ParseTaskRequest(r *http.Request) (*TaskRequest, error) {
	var errs []error

	id, err := ParseTaskID(r.PathValue("id"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &TaskRequest{id}, nil
}

func ParseTaskAttachmentRequest(r *http.Request) (*TaskAttachmentRequest, error) {
	var errs []error

	id, err := ParseTaskID(r.PathValue("id"))
	if err != nil {
		errs = append(errs, err)
	}

	name, err := ParseTaskAttachmentName(r.PathValue("name"))
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &TaskAttachmentRequest{id, name}, nil
}

func ParseNewTaskRequest(r *http.Request) (*NewTaskRequest, error) {
	if err := r.ParseMultipartForm(1 << 27); err != nil {
		return nil, errors.New("attachments: payload size exceeds limit")
	}

	var errs []error

	name, err := ParseTaskName(r.FormValue("name"))
	if err != nil {
		errs = append(errs, err)
	}

	expiresAt, err := ParseTaskExpiresAt(r.FormValue("expires_at"), GetLocation(r))
	if err != nil {
		errs = append(errs, err)
	}

	attachments, err := ParseTaskAttachments(r)
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &NewTaskRequest{name, expiresAt, attachments}, nil
}

func ParseUpdateTaskRequest(r *http.Request) (*UpdateTaskRequest, error) {
	if err := r.ParseMultipartForm(1 << 27); err != nil {
		return nil, errors.New("attachments: payload size exceeds limit")
	}

	var errs []error

	id, err := ParseTaskID(r.PathValue("id"))
	if err != nil {
		errs = append(errs, err)
	}

	name, err := ParseTaskName(r.FormValue("name"))
	if err != nil {
		errs = append(errs, err)
	}

	expiresAt, err := ParseTaskExpiresAt(r.FormValue("expires_at"), GetLocation(r))
	if err != nil {
		errs = append(errs, err)
	}

	attachments, err := ParseTaskAttachments(r)
	if err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &UpdateTaskRequest{id, name, expiresAt, attachments}, nil
}

func ParseLanguage(value string) (string, error) {
	if !IsValidLanguage(value) {
		return "", fmt.Errorf("language: required, supported values: %s", strings.Join(SupportedLanguages, ", "))
	}

	return value, nil
}

func ParseTheme(value string) (string, error) {
	if !IsValidTheme(value) {
		return "", fmt.Errorf("theme: required, supported values: %s", strings.Join(SupportedThemes, ", "))
	}

	return value, nil
}

func ParseTimezone(value string) (string, error) {
	if !IsValidTimezone(value) {
		return "", fmt.Errorf("timezone: required, supported values: %s", strings.Join(SupportedTimezones, ", "))
	}

	return value, nil
}

func ParseTaskID(value string) (int, error) {
	v, err := strconv.Atoi(value)
	if err != nil || v < 1 {
		return 0, errors.New("id: required, must be an integer greater than 0")
	}

	return v, nil
}

func ParseTaskAttachmentName(value string) (string, error) {
	l := len(value)
	if l < 1 || 200 < l {
		return "", errors.New("name: required, must be between 1 and 200 characters")
	}

	return value, nil
}

func ParseTaskName(value string) (string, error) {
	l := len(value)
	if l < 1 || 200 < l {
		return "", errors.New("name: required, must be between 1 and 200 characters")
	}

	return value, nil
}

func ParseTaskExpiresAt(value string, l *time.Location) (*time.Time, error) {
	if value == "" || l == nil {
		return nil, nil
	}

	v, err := ParseTime(value, l)
	if err != nil {
		return nil, errors.New("expires_at: must be in ISO 8601 format")
	}

	return v, nil
}

func ParseTaskAttachments(r *http.Request) (*AttachmentsRequest, error) {
	files := r.MultipartForm.File["attachments"]

	names := slices.DeleteFunc(
		r.Form["attachments"],
		func(name string) bool {
			return name == ""
		},
	)

	for _, file := range files {
		names = append(names, file.Filename)
	}

	slices.Sort(names)
	names = slices.Compact(names)

	return &AttachmentsRequest{names, files}, nil
}

func ParseTime(t string, l *time.Location) (*time.Time, error) {
	v, err := time.ParseInLocation(TimeFormatISO, t, l)
	if err != nil {
		return nil, err
	}

	v = v.UTC()

	return &v, nil
}
