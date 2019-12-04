package oauthoidc

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
)

type OidcProvider struct {
}

type OidcUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func init() {
	provider := &OidcProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_OIDC, provider)
}

func (glu *OidcUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *OidcUser) IsValid() bool {
	if glu.Id == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}

func (glu *OidcUser) getAuthData() string {
	return strconv.FormatInt(glu.Id, 10)
}

func (m *OidcProvider) GetUserFromJson(data io.Reader) *model.User {
	oidcUser := OidcUser{Id: 12345, Username: "aina", Login: "aina", Email: "aina@naru.pw", Name: "aina"}
	userId := strconv.FormatInt(12345, 10)
	user := model.User{Username: "aina", FirstName: "aina", LastName: "kusaka", Email: "aina@naru.pw", AuthData: &userId, AuthService: model.USER_AUTH_SERVICE_OIDC}
	if oidcUser.IsValid() {
		return &user
	}

	return &user
}

/*
func userFromGitLabUser(glu *GitLabUser) *model.User {
	user := &model.User{}
	username := glu.Username
	if username == "" {
		username = glu.Login
	}
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	userId := glu.getAuthData()
	user.AuthData = &userId
	user.AuthService = model.USER_AUTH_SERVICE_GITLAB

	return user
}

func gitLabUserFromJson(data io.Reader) *GitLabUser {
	decoder := json.NewDecoder(data)
	var glu GitLabUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}
*/
