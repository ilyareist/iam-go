package main

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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

		publicKeyByte, _ := ioutil.ReadFile("./keys/demo.rsa.pub")
		publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)

		token := jwt.New(jwt.GetSigningMethod("RS256"))

		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Jon Doe"
		claims["role"] = "admin"
		claims["iss"] = "testing@secure.istio.io"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		key, _ := jwk.New(publicKey)

		jsonbuf, _ := json.MarshalIndent(key, "", "  ")


		os.Stdout.Write(jsonbuf)

		t, _ := token.SignedString(privateKey)

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}

func (h *handler) sharePub(c echo.Context) {
	//pubKey:="-----BEGINPUBLICKEY-----MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2jHRPcMcMFyj5a9dbDYtSdDmlviEsXvj9bC8FCaNCtHq44hXvw3uFuAZ9hUlYA+yo6i+5IXbFO0RNmMVLATmrr0R2XpPkXry4k4x5b8lh2HJkLzwYbioWu4ijAo92C7uxVMNq99y6YprcrwsRyKApMv4C9WewUOObgoe/6QfYN5Bhen5kWrEgLmyt8cPkPSTK54A4Ki3+i58qnlH4h8GtD9b03VYeeV4cbcqyuQPaiZZF5q7VRI/XkzJcr/IdzTg9Pt0bEKthyvao0NfEPxN8++u8dq8Uz1W/uJL7dpV3r2nCK/dc1hahMIrjqH96VtL6LYGJogJ5ykrSWm8/KHz+QIDAQAB-----ENDPUBLICKEY-----"
	//return c.JSON(http.StatusOK, map[string])
}
