package auth

import (
	"database/sql"
	"encoding/base64"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// encodePassword encodes the password to base64
func EncodePassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func GetToken(c *fiber.Ctx, name string, manager bool) error {
	// assign a token to the user
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Unix()
	claims["usr"] = name
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["adm"] = manager
	returnToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	log.Printf("Token: %v\n", returnToken)
	return c.JSON(returnToken)
}

func validateToken(c *fiber.Ctx) error {
	// get the token from the header
	// remove the "Bearer " prefix
	token := c.Get("Authorization")[7:]
	// validate the token
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// validate the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("secret"), nil
	})
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusUnauthorized, err.Error()))
	}
	// // update the token expiration time
	// claims := parsedToken.Claims.(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	// returnToken, err := parsedToken.SignedString([]byte("secret"))
	// if err != nil {
	// 	log.Printf("Error signing token: %v\n", err)
	// 	return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	// }
	// // return parsedToken with status OK
	// c.Status(fiber.StatusOK)
	// log.Printf("Token: %v\n", returnToken)
	// return c.JSON(returnToken)
	return c.SendStatus(fiber.StatusOK)
}

func AuthFunc(app *fiber.App, db *sql.DB) {
	app.Post("/validate", func(c *fiber.Ctx) error {
		return validateToken(c)
	})
}
