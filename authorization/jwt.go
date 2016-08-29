package authorization

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/landru29/api-go/model/token"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

// TokenClaims is the jwt model for the token
type TokenClaims struct {
	Token     string `json:"token"`
	Expiry    int64  `json:"expiry"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	jwt.StandardClaims
}

func getMapClaim(cl TokenClaims) jwt.MapClaims {
	m := structs.Map(&cl)
	mapping := make(map[string]interface{}, len(m))
	for k, v := range m {
		mapping[k] = v
	}
	return mapping
}

// EncodeToken encode a jwt token
func EncodeToken(savedToken token.Model) (tokenString string, err error) {
	jsonToken := TokenClaims{
		Token:     savedToken.ID.Hex(),
		Expiry:    savedToken.Expiry,
		Email:     savedToken.Email,
		FirstName: savedToken.FirstName,
		LastName:  savedToken.LastName,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, getMapClaim(jsonToken))
	tokenString, err = jwtToken.SignedString([]byte(viper.GetString("jwt_secret")))
	return
}

// DecodeToken decode a jwt token
func DecodeToken(tokenString string, db *mgo.Database) (dbToken token.Model, profile Profile, err error) {
	tokenJwt, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt_secret")), nil
	})
	profile = Profile{}
	if err != nil {
		fmt.Println("JWT: Not JWT")
		return
	}
	if claims, ok := tokenJwt.Claims.(*TokenClaims); ok && tokenJwt.Valid {
		if claims.Expiry > time.Now().Unix() {
			profile.Email = claims.Email
			profile.FirstName = claims.FirstName
			profile.LastName = claims.LastName
		} else {
			err = errors.New("Expired token")
			return
		}
		dbToken, err = token.FindToken(db, claims.Token)
		if err != nil {
			fmt.Println("JWT: Not found in db")
			fmt.Println(claims.Token)
			return
		}
	}

	return
}
