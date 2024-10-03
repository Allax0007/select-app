package data

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"web-server/tools/table"

	"github.com/gofiber/fiber/v2"
)

func SelectData(db *sql.DB, c *fiber.Ctx) error {
	//select-data?t={}&&col={}
	t := c.Query("t")
	col := c.Query("col")
	if t == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Missing query parameter 't'"))
	}
	if col == "" {
		col = "*"
	}
	rows, err := db.Query("SELECT " + col + " FROM " + t)
	if err != nil {
		log.Printf("Error selecting data: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Printf("Error scanning data: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		log.Printf("value: %v\n", value)
		values = append(values, value)
	}
	return c.JSON(values)
}

func deleteData(db *sql.DB, c *fiber.Ctx) error {
	//delete-data?t={}
	t := c.Query("t")
	if t == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	col, err := table.GetColName(t, db, c)
	if err != nil {
		log.Printf("Error getting columns name: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	values := string(c.Body())
	log.Printf("values: %v\n", values) // for debugging
	valuesSlice := strings.Split(strings.TrimSpace(values[1:len(values)-1]), ",")
	deleteDataSQL := `DELETE FROM ` + t + ` WHERE [`
	for i, v := range valuesSlice {
		if v != "null" {
			deleteDataSQL += col[i] + `] = ` + v + ` AND [`
		}
	}
	deleteDataSQL = deleteDataSQL[:len(deleteDataSQL)-6]
	// log.Printf("deleteDataSQL: %v\n", deleteDataSQL) // for debugging
	_, err = db.Exec(deleteDataSQL)
	if err != nil {
		log.Printf("Error deleting data: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Data deleted successfully"))
}

func insertData(db *sql.DB, c *fiber.Ctx) error {
	//insert-data?t={}
	t := c.Query("t")
	if t == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Missing query parameter 't'"))
	}
	values := string(c.Body())
	log.Printf("table: %v\n", t)       // for debugging
	log.Printf("values: %v\n", values) // for debugging
	// values template: {"id":"3","name":"4","pwd":"5"}
	// split values
	valuesSlice := strings.Split(strings.TrimSpace(values[1:len(values)-1]), ",")
	keys := make([]string, 0, len(valuesSlice))
	vals := make([]string, 0, len(valuesSlice))
	for _, v := range valuesSlice {
		v = strings.TrimSpace(v)
		keyVal := strings.Split(v, ":")
		keys = append(keys, strings.TrimSpace(keyVal[0]))
		vals = append(vals, strings.TrimSpace(keyVal[1]))
	}
	log.Printf("keys: %v, vals: %v\n", keys, vals) // for debugging
	// insert data according to keys and vals
	insertDataSQL := `INSERT INTO ` + t + ` (`
	for _, k := range keys {
		insertDataSQL += k + `,`
	}
	insertDataSQL = insertDataSQL[:len(insertDataSQL)-1] + `) VALUES (`
	for _, v := range vals {
		insertDataSQL += v + `,`
	}
	insertDataSQL = insertDataSQL[:len(insertDataSQL)-1] + `);`
	log.Printf("insertDataSQL: %v\n", insertDataSQL) // for debugging
	_, err := db.Exec(insertDataSQL)
	if err != nil {
		log.Printf("Error inserting data: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Data inserted successfully"))
}

func UpdateData(db *sql.DB, c *fiber.Ctx) error {
	//update-data?t={}&&id={}&&col&&val
	t := c.Query("t")
	// id := c.Query("id")
	adjCol, err := strconv.Atoi(c.Query("col"))
	if err != nil {
		log.Printf("Error converting adjCol to int: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	val := c.Query("val")
	if t == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Missing query parameters"))
	}
	col, err := table.GetColName(t, db, c)
	if err != nil {
		log.Printf("Error getting columns name: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// log.Printf("columns: %v\n", col) // for debugging
	// log.Printf("c.Body: %v\n", string(c.Body()[0])) // for debugging
	values := string(c.Body())
	log.Printf("values: %v\n", values) // for debugging
	// values template: ["W","W-2",null,35]
	// split values
	valuesSlice := strings.Split(strings.TrimSpace(values[1:len(values)-1]), ",")
	// log.Printf("valuesSlice: %v\n", valuesSlice) // for debugging

	UpdateDataSQL := `UPDATE ` + t + ` SET [` + col[adjCol] + `] = "` + val + `" WHERE [`
	for i, v := range valuesSlice {
		if v != "null" {
			UpdateDataSQL += col[i] + `] = ` + v + ` AND [`
		}
	}
	UpdateDataSQL = UpdateDataSQL[:len(UpdateDataSQL)-6]
	// log.Printf("UpdateDataSQL: %v\n", UpdateDataSQL) // for debugging
	_, err = db.Exec(UpdateDataSQL)
	if err != nil {
		log.Printf("Error updating data: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Data updated successfully"))
}

func DataFunc(app *fiber.App, db *sql.DB) {
	app.Get("/select-data", func(c *fiber.Ctx) error {
		return SelectData(db, c)
	})
	app.Delete("/delete-data", func(c *fiber.Ctx) error {
		return deleteData(db, c)
	})
	app.Post("/insert-data", func(c *fiber.Ctx) error {
		return insertData(db, c)
	})
	app.Put("/update-data", func(c *fiber.Ctx) error {
		return UpdateData(db, c)
	})
}
