package emailnotifier

type NullEmailResolver struct{}

func (r *NullEmailResolver) ResolveEmail(userID string) (string, error) {
	return userID, nil
}
