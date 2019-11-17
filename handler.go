package main

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/lestrrat-go/jwx/jwk"
)

type handler struct{}

func (h *handler) login(c echo.Context) error {

	type User struct {
		Name     string `json:name`
		Password string `json:password`
	}
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		panic("ERROR")
	}

	if user.Name == "jon" && user.Password == "password" {

		var privateKey *rsa.PrivateKey
		privateKeyByte, _ := ioutil.ReadFile("./keys/demo.rsa")
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)

		token := jwt.New(jwt.GetSigningMethod("RS256"))

		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Jon Doe"
		claims["role"] = "admin"
		claims["iss"] = "testing@secure.istio.io"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, _ := token.SignedString(privateKey)

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}

func (h *handler) sharePub(c echo.Context) error{
	publicKeyByte, _ := ioutil.ReadFile("./keys/demo.rsa.pub")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	key, _ := jwk.New(publicKey)
	jsonbuf, _ := json.MarshalIndent(key, "", "  ")
	return c.JSON(http.StatusOK, map[string]string{
		"token": string(jsonbuf),
	})

}
