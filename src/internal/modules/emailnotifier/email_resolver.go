package emailnotifier

type EmailResolver interface {
	ResolveEmail(userID string) (string, error)
}
