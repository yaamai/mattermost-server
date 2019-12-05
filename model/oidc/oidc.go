package oauthoidc

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v5/model"
	"io"
)

type OidcProvider struct {
}

type OidcUser struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func (m *OidcProvider) GetUserFromJson(data io.Reader) *model.User {

	decoder := json.NewDecoder(data)
	var oidcUser OidcUser
	err := decoder.Decode(&oidcUser)
	if err != nil {
		return nil
	}

	userId := oidcUser.Sub
	user := model.User{
		Username:    oidcUser.Sub,
		FirstName:   oidcUser.Sub,
		Email:       oidcUser.Email,
		AuthData:    &userId,
		AuthService: model.USER_AUTH_SERVICE_OIDC,
	}

	return &user
}
