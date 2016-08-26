package authorization

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/facebook"
    "golang.org/x/oauth2/google"
)

// GetEmail retrieve the email from the third party (facebook or google)
func GetEmail(c *gin.Context, auth *oauth2.Config, apiToken *oauth2.Token) (email string, err error) {
    var profile Profile
    profile, err = GetProfile(c, auth, apiToken)
    email = profile.Email
    if len(profile.Email) == 0 {
        err = errors.New("Empty Email")
    }
    return
}

// GetProfile retrieve the profile from the third party (facebook or google)
func GetProfile(c *gin.Context, auth *oauth2.Config, apiToken *oauth2.Token) (profile Profile, err error) {
    client := auth.Client(c, apiToken)
    uri := ""
    switch auth.Endpoint {
    case facebook.Endpoint:
        uri = "https://graph.facebook.com/v2.2/me?fields=id,name,email,picture,first_name,last_name"
    case google.Endpoint:
        uri = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
    default:
        uri = ""
    }

    resp, err := client.Get(uri)
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    defer resp.Body.Close()
    contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    switch auth.Endpoint {
    case facebook.Endpoint:
        var p ProfileFacebook
        err = json.Unmarshal(contents, &p)
        if err != nil {
            profile = Profile{}
        } else {
            profile = Profile{
                ID:        p.ID,
                Email:     p.Email,
                FirstName: p.FirstName,
                LastName:  p.LastName,
                Hd:        p.Hd,
                Locale:    p.Locale,
                Name:      p.Name,
            }
        }
    case google.Endpoint:
        var p ProfileGoogle
        err = json.Unmarshal(contents, &p)
        if err != nil {
            profile = Profile{}
        } else {
            profile = Profile{
                ID:        p.ID,
                Email:     p.Email,
                FirstName: p.GivenName,
                LastName:  p.FamilyName,
                Hd:        p.Hd,
                Locale:    p.Locale,
                Name:      p.Name,
            }
        }
    default:
        profile = Profile{}
    }

    if len(profile.Email) == 0 {
        err = errors.New("Empty Email")
    }

    return
}
