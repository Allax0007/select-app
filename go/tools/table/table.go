package table

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

func createTable(db *sql.DB, c *fiber.Ctx) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL UNIQUE,
		"pwd" TEXT NOT NULL,
		"manager" NUMERIC DEFAULT 1
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Table created successfully"))
}
func getTablesName(db *sql.DB, c *fiber.Ctx) error {
	//get-tables?like={}
	like := c.Query("like")
	querySQL := "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
	if like != "" {
		querySQL += " AND name LIKE '%" + like + "%'"
	}
	rows, err := db.Query(querySQL)
	if err != nil {
		log.Printf("Error getting tables name: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Printf("Error scanning table: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		tables = append(tables, table)
	}
	return c.JSON(tables)
}
func GetColName(q string, db *sql.DB, c *fiber.Ctx) ([]string, error) {
	rows, err := db.Query("PRAGMA table_info(" + q + ")")
	if err != nil {
		log.Printf("Error getting columns name: %v\n", err)
		return nil, c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var columns []string
	for rows.Next() {
		var column, _type string
		var _int, _pk int
		rows.Scan(&_int, &column, &_type, &_int, &_int, &_pk)
		columns = append(columns, column)
	}
	return columns, nil
}
func showTables(db *sql.DB, c *fiber.Ctx) error {
	q := c.Query("q")
	col := c.Query("col")
	if q == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Query parameter 'q' is required"))
	}
	if col == "" {
		col = "*"
	}
	cname, err := GetColName(q, db, c)
	if err != nil {
		log.Printf("Error getting columns name: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	showTablesSQL := `SELECT ` + col + ` FROM [` + q + `] ORDER BY [` + cname[0] + `];`
	rows, err := db.Query(showTablesSQL)
	if err != nil {
		log.Printf("Error showing tables: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var rowsData []interface{}
	columns, err := GetColName(q, db, c)
	if err != nil {
		log.Printf("Error getting columns name: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	rowsData = append(rowsData, columns)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuesPtr := make([]interface{}, len(columns))
		for i := range columns {
			valuesPtr[i] = &values[i]
		}
		rows.Scan(valuesPtr...)
		rowData := make([]interface{}, len(columns))
		for i := range columns {
			rowData[i] = values[i]
		}
		rowsData = append(rowsData, rowData)
	}
	return c.JSON(rowsData)
}
func addColumn(db *sql.DB, c *fiber.Ctx) error {
	//add-column?t={}&&col={}&&type={}
	t := c.Query("t")
	col := c.Query("col")
	_type := c.Query("type")
	if t == "" || col == "" || _type == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Missing query parameter 't', 'col' or 'type'"))
	}
	addColumnSQL := `ALTER TABLE ` + t + ` ADD COLUMN ` + col + ` ` + _type
	_, err := db.Exec(addColumnSQL)
	if err != nil {
		log.Printf("Error adding column: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Column added successfully"))
}
func dropTable(db *sql.DB, c *fiber.Ctx) error {
	//drop-table?t={}
	t := c.Query("t")
	if t == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Missing query parameter 't'"))
	}
	dropTableSQL := `DROP TABLE [` + t + `];`
	_, err := db.Exec(dropTableSQL)
	if err != nil {
		log.Printf("Error dropping table: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Table dropped successfully"))
}
func TableFunc(app *fiber.App, db *sql.DB) {
	app.Get("/create-table", func(c *fiber.Ctx) error {
		return createTable(db, c)
	})
	app.Get("/get-tables", func(c *fiber.Ctx) error {
		return getTablesName(db, c)
	})
	app.Get("/show-tables", func(c *fiber.Ctx) error {
		return showTables(db, c)
	})
	app.Get("/add-column", func(c *fiber.Ctx) error {
		return addColumn(db, c)
	})
	app.Delete("/drop-table", func(c *fiber.Ctx) error {
		return dropTable(db, c)
	})
}
