package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/token"
	"github.com/landru29/api-go/mongo"
	"golang.org/x/oauth2"
)

func getRedirect(c *gin.Context, origin string, tokenObj *oauth2.Token, p Profile) (uri string, err error) {
	uri = ""
	redirect, err := c.Cookie("redirect")
	if err != nil {
		c.Redirect(301, "/error")
		return
	}
	redirectURL, err := url.Parse(redirect)
	if err != nil {
		c.Redirect(301, "/error")
		return
	}

	dbToken := token.Model{
		TokenType:    tokenObj.TokenType,
		AccessToken:  tokenObj.AccessToken,
		RefreshToken: tokenObj.RefreshToken,
		Expiry:       tokenObj.Expiry.Unix(),
		Email:        p.Email,
		Origin:       origin,
	}

	savedToken, err := dbToken.Save(mongo.GetMongoDatabase())
	if err != nil {
		c.Redirect(301, "/error")
		return
	}

	q := redirectURL.Query()
	q.Add("api-token", savedToken.ID.Hex())
	redirectURL.RawQuery = q.Encode()

	uri = redirectURL.String()

	return
}

func getProfile(c *gin.Context, auth *oauth2.Config, apiToken *oauth2.Token, origin string) (profile Profile, err error) {
	client := auth.Client(c, apiToken)
	uri := ""
	switch origin {
	case Facebook:
		uri = "https://graph.facebook.com/v2.2/me?fields=id,name,email,picture,first_name,last_name"
	case Google:
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

	switch origin {
	case Facebook:
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
	case Google:
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

	return
}
