package authorization

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/token"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
)

// APIBaseURL return the base URL of the API
func APIBaseURL() string {
	port := viper.GetString("api_port")
	return viper.GetString("api_protocol") + "://" + viper.GetString("api_host") +
		map[bool]string{true: ":" + port, false: ""}[len(port) > 0] + "/"
}

func buildRedirect(c *gin.Context, origin string, tokenObj *oauth2.Token, profile Profile, db *mgo.Database) (uri string, err error) {
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
		Email:        profile.Email,
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		Origin:       origin,
		Identifier:   profile.Source + profile.ID,
	}

	savedToken, err := dbToken.Save(db)
	if err != nil {
		c.Redirect(301, "/error")
		return
	}

	jwtToken, err := EncodeToken(savedToken)
	if err != nil {
		c.Redirect(301, "/error")
		return
	}

	q := redirectURL.Query()
	q.Add("api-token", jwtToken)
	redirectURL.RawQuery = q.Encode()

	uri = redirectURL.String()

	return
}
