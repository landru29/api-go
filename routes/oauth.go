package routes

import (
	"net/url"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func getRedirect(c *gin.Context, origin string, token *oauth2.Token) (uri string, err error) {
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

	session := sessions.Default(c)

	q := redirectURL.Query()
	session.Set("access-token", token.AccessToken)
	session.Set("origin", origin)
	session.Set("refresh-token", token.RefreshToken)
	session.Set("type", token.TokenType)
	session.Set("expire", token.Expiry.String())
	/*q.Add("access-token", token.AccessToken)
	q.Add("origin", origin)
	q.Add("refresh-token", token.RefreshToken)
	q.Add("type", token.TokenType)
	q.Add("expire", token.Expiry.String())*/

	redirectURL.RawQuery = q.Encode()

	uri = redirectURL.String()

	return
}
