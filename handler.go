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
		claims["group"] = "customer"
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

func (h *handler) sharePub(c echo.Context) error {
	publicKeyByte, _ := ioutil.ReadFile("./keys/demo.rsa.pub")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	key, _ := jwk.New(publicKey)
	jsonbuf, _ := json.MarshalIndent(key, "", "  ")
	type Jwk struct {
		E   string `json:"e"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	}
	var jwk Jwk
	json.Unmarshal([]byte(jsonbuf), &jwk)

	type Jwks struct {
		Keys []Jwk `json:"keys"`
	}

	sliceJwk := []Jwk{jwk}

	jwks := Jwks{sliceJwk}

	return c.JSON(http.StatusOK, jwks)
}
