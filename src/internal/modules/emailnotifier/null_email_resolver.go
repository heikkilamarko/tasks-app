package emailnotifier

type NullEmailResolver struct{}

var _ EmailResolver = (*NullEmailResolver)(nil)

func (r *NullEmailResolver) ResolveEmail(userID string) (string, error) {
	return userID, nil
}
