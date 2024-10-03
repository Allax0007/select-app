package main

import (
	"database/sql"
	"log"
	"web-server/app/data"
	"web-server/app/product"
	"web-server/app/user"
	"web-server/tools/auth"
	"web-server/tools/table"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
}

func main() {
	initDB()
	defer db.Close()

	app := fiber.New()
	app.Static("/", "../dist")

	auth.AuthFunc(app, db)
	user.UserFunc(app, db)
	data.DataFunc(app, db)
	table.TableFunc(app, db)
	product.ProductFunc(app, db)

	log.Fatal(app.Listen(":3000"))
}
