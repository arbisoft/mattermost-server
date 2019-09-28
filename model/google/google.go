package oauthgoogle

import (
       "encoding/json"
       "io"
       "io/ioutil"
       "strings"

       "github.com/mattermost/mattermost-server/einterfaces"
       "github.com/mattermost/mattermost-server/model"
)

type GoogleProvider struct {
}

type GoogleUser struct {
       Id    string  `json:"sub"`
       Email string `json:"email"`
       username string `json:"email"`
       // Login    string `json:"login"`
       Name string `json:"name"`
}

func init() {
       provider := &GoogleProvider{}
       einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_GOOGLE, provider)
}

func (glu *GoogleUser) IsValid() bool {
       if glu.Id == "" {
               return false
       }

       if len(glu.Email) == 0 {
               return false
       }

       return true
}

func userFromGoogleUser(glu *GoogleUser) *model.User {
       user := &model.User{}


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
       user.Username = strings.ToLower(splitName[0] + splitName[1])
       user.Email = glu.Email
       userId := glu.Id
       user.AuthData = &userId
       user.AuthService = model.USER_AUTH_SERVICE_GOOGLE

       return user
}

func googleUserFromJson(data io.Reader) *GoogleUser {
       bodyBytes, _ := ioutil.ReadAll(data)
       bodyString := string(bodyBytes)
       var gu GoogleUser
       err := json.Unmarshal([]byte(bodyString), &gu)
       if err == nil {
               return &gu
       } else {
               return nil
       }
}

func (m *GoogleProvider) GetUserFromJson(data io.Reader) *model.User {
       gu := googleUserFromJson(data)
       if gu.IsValid() {
               return userFromGoogleUser(gu)
       }

       return &model.User{}
}
