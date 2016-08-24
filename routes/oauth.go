package routes

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/model/token"
	"github.com/landru29/api-go/mongo"
	"golang.org/x/oauth2"
)

func getRedirect(c *gin.Context, origin string, tokenObj *oauth2.Token) (uri string, err error) {
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
		Origin:       origin,
	}

	savedToken, err := dbToken.Save(mongo.GetMongoDatabase())
	if err != nil {
		c.Redirect(301, "/error")
		return
	}

	q := redirectURL.Query()
	q.Add("api-token", savedToken.ID.Hex())
	/*q.Add("access-token", token.AccessToken)
	q.Add("origin", origin)
	q.Add("refresh-token", token.RefreshToken)
	q.Add("type", token.TokenType)
	q.Add("expire", token.Expiry.String())*/

	redirectURL.RawQuery = q.Encode()

	uri = redirectURL.String()

	return
}
