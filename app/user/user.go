package user

import (
	"database/sql"
	"log"
	"web-server/tools/auth"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"pwd"`
	Manager  bool   `json:"manager"`
}

func addUser(db *sql.DB, c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// encode the password
	user.Password = auth.EncodePassword(user.Password)
	insertUserSQL := `INSERT INTO users (name, pwd) VALUES (?, ?)`
	_, err := db.Exec(insertUserSQL, user.Name, user.Password)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("User added successfully"))
}

func userAuth(db *sql.DB, c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// encode the password
	user.Password = auth.EncodePassword(user.Password)
	loginAuthSQL := `SELECT * FROM users WHERE name = ? AND pwd = ?`
	row := db.QueryRow(loginAuthSQL, user.Name, user.Password)
	var id int
	var manager bool
	var name, password string
	err := row.Scan(&id, &name, &password, &manager)
	if err != nil {
		log.Printf("Error scanning user: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		// return c.SendStatus(fiber.StatusBadRequest)
	}
	// assign a token to the user
	return auth.GetToken(c, user.Name, manager)
}

func editUser(db *sql.DB, c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// encode the password
	user.Password = auth.EncodePassword(user.Password)
	editUserSQL := `UPDATE users SET pwd = ? WHERE name = ?`
	res, err := db.Exec(editUserSQL, user.Password, user.Name)
	if err != nil {
		log.Printf("Error updating user: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "User not found"))
	}
	return c.Send([]byte("User updated successfully"))
}

func UserFunc(app *fiber.App, db *sql.DB) {
	app.Post("/loginAuth", func(c *fiber.Ctx) error {
		return userAuth(db, c)
	})
	app.Post("/add-user", func(c *fiber.Ctx) error {
		return addUser(db, c)
	})
	app.Put("/edit-user", func(c *fiber.Ctx) error {
		return editUser(db, c)
	})
}

//forgot password
