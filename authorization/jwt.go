package authorization

import (
    "fmt"

    "github.com/dgrijalva/jwt-go"
    "github.com/landru29/api-go/model/token"
    "github.com/spf13/viper"
    "gopkg.in/mgo.v2"
)

// MyCustomClaims is the jwt model for the token
type MyCustomClaims struct {
    Token  string `json:"token"`
    Expiry int64  `json:"expiry"`
    jwt.StandardClaims
}

// EncodeToken encode a jwt token
func EncodeToken(savedToken token.Model) (tokenString string, err error) {
    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "token":  savedToken.ID.Hex(),
        "expiry": savedToken.Expiry,
    })
    tokenString, err = jwtToken.SignedString([]byte(viper.GetString("jwt_secret")))
    return
}

// DecodeToken decode a jwt token
func DecodeToken(tokenString string, db *mgo.Database) (dbToken token.Model, err error) {
    tokenJwt, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(viper.GetString("jwt_secret")), nil
    })
    if err != nil {
        fmt.Println("JWT: Not JWT")
        return
    }
    if claims, ok := tokenJwt.Claims.(*MyCustomClaims); ok && tokenJwt.Valid {
        dbToken, err = token.FindToken(db, claims.Token)
        if err != nil {
            fmt.Println("JWT: Not found in db")
            fmt.Println(claims.Token)
            return
        }
    }
    return
}
