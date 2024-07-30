package emailnotifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"tasks-app/internal/shared"
)

type ZitadelEmailResolver struct {
	Config *shared.Config
}

var _ EmailResolver = (*ZitadelEmailResolver)(nil)

func (r *ZitadelEmailResolver) ResolveEmail(userID string) (string, error) {
	url, err := url.JoinPath(r.Config.EmailNotifier.ZitadelURL, "/v2/users/", userID)
	if err != nil {
		return "", err
	}

	token := r.Config.EmailNotifier.ZitadelPAT

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body := struct {
		User struct {
			State string `json:"state"`
			Human struct {
				Email struct {
					Email      string `json:"email"`
					IsVerified bool   `json:"isVerified"`
				} `json:"email"`
			} `json:"human"`
		} `json:"user"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", err
	}

	if body.User.State != "USER_STATE_ACTIVE" {
		return "", errors.New("user is not active")
	}

	return body.User.Human.Email.Email, nil
}
