package product

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"web-server/app/data"
	"web-server/tools/table"

	"github.com/gofiber/fiber/v2"
)

type Product struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Components  []string `json:"comp"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createProduct(db *sql.DB, c *fiber.Ctx) error {
	var product Product
	if err := c.BodyParser(&product); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	log.Printf("product: %v\n", product)

	insertProductSQL := `INSERT INTO Product (name, description) VALUES (?, ?)`
	_, err := db.Exec(insertProductSQL, product.Name, product.Description)
	if err != nil {
		log.Printf("Error inserting product: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// create relation table
	res := createRelTable(db, product.Name)
	if res != nil {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, res.Error()))
	}
	// insert components
	res = addComponents(db, product.Name, product.Components)
	if res != nil {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, res.Error()))
	}
	return c.Send([]byte("Product added successfully"))
}

func deleteProduct(db *sql.DB, c *fiber.Ctx) error {
	q := string(c.Body())
	// log.Printf("q: %v\n", q)
	prod, err := table.GetColName("Product", db, c)
	if err != nil {
		return err
	}
	qSlice := strings.Split(strings.TrimSpace(q[1:len(q)-1]), ",")
	deleteDataSQL := `DELETE FROM Product WHERE `
	for i, v := range qSlice {
		if v != "null" {
			deleteDataSQL += prod[i] + ` = ` + v + ` AND `
		}
	}
	deleteDataSQL = deleteDataSQL[:len(deleteDataSQL)-5]
	_, err = db.Exec(deleteDataSQL)
	if err != nil {
		log.Printf("Error deleting product: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	res := dropRelTable(db, qSlice[0])
	if res != nil {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, res.Error()))
	}
	return c.Send([]byte("Product deleted successfully"))
}

func listProducts(db *sql.DB, c *fiber.Ctx) error {
	return data.SelectData(db, c)
}

func ProductDetails(db *sql.DB, c *fiber.Ctx) error {
	// get-product?name={}
	name := c.Query("name")
	if name == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// get product description
	rows, err := db.Query("SELECT description FROM Product WHERE name='" + name + "';")
	if err != nil {
		log.Printf("Error getting product description: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var desc string
	if rows.Next() {
		err = rows.Scan(&desc)
		if err != nil {
			log.Printf("Error scanning product description: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
	}
	// get components
	rows, err = db.Query("SELECT component FROM Relation" + name + " ORDER BY [order];")
	if err != nil {
		log.Printf("Error getting components: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var components []string
	for rows.Next() {
		var component string
		err := rows.Scan(&component)
		if err != nil {
			log.Printf("Error scanning components: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		components = append(components, component)
	}
	// return as: desc: ["comp1", "comp2", ...]
	product := map[string]interface{}{"desc": desc, "comp": components}
	return c.JSON(product)
}

func createCompRelTable(db *sql.DB, c *fiber.Ctx) error {
	// data: {ComponentA: 'W', ComponentB: 'Y'}
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	log.Printf("data: %v\n", data)
	// check if AtoB exists
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name LIKE '%" + data["ComponentA"] + "%to%" + data["ComponentB"] + "%' OR name LIKE '%" + data["ComponentB"] + "%to%" + data["ComponentA"] + "%';")
	if err != nil {
		log.Printf("Error checking relation table: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	if rows.Next() {
		c.Status(fiber.StatusBadRequest)
		return c.Send([]byte("Relation already exists"))
	}
	// create relation table
	createTableSQL := `CREATE TABLE "` + data["ComponentA"] + `to` + data["ComponentB"] + `" (
						"id"	INTEGER,
						"` + data["ComponentA"] + `id"	INTEGER NOT NULL,
						"` + data["ComponentB"] + `id"	INTEGER NOT NULL,
						"amount"	INTEGER NOT NULL,
						PRIMARY KEY("id"),
						FOREIGN KEY("` + data["ComponentB"] + `id") REFERENCES "` + data["ComponentB"] + `"("id"),
						FOREIGN KEY("` + data["ComponentA"] + `id") REFERENCES "` + data["ComponentA"] + `"("id")
						);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating relation table: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Relation table created successfully"))
}

func createRelTable(db *sql.DB, t string) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS Relation` + t + ` (
		"order"	INTEGER NOT NULL,
		"component"	TEXT NOT NULL UNIQUE,
		"related"	TEXT,
		PRIMARY KEY("order" AUTOINCREMENT),
		FOREIGN KEY("component") REFERENCES "Component"("type"),
		FOREIGN KEY("related") REFERENCES "Relation` + t + `"("component"),
		CHECK("component" <> "related")
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v\n", err)
		return err
	}
	return nil
}

func chgRelOrder(db *sql.DB, c *fiber.Ctx) error {
	// chg-order?t={}&&ord={}&&dir={}
	t := c.Query("t")
	ordStr := c.Query("ord")
	dirStr := c.Query("dir")
	if t == "" || ordStr == "" || dirStr == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// convert ord and dir to int
	ord, err := strconv.Atoi(ordStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	dir, err := strconv.Atoi(dirStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// if dir is 1, get max order
	if dir == 1 {
		rows, err := db.Query(`SELECT MAX("order") FROM ` + t + `;`)
		if err != nil {
			log.Printf("Error getting max order: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		defer rows.Close()
		var maxOrd int
		if rows.Next() {
			err = rows.Scan(&maxOrd)
			if err != nil {
				log.Printf("Error getting max order: %v\n", err)
				return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
		}
		// if ord is already max order, return
		if ord == maxOrd {
			c.Status(fiber.StatusBadRequest)
			return c.Send([]byte("Order is already max"))
		}
	}

	// if dir is -1 && ord is 1, return
	if dir == -1 && ord == 1 {
		c.Status(fiber.StatusBadRequest)
		return c.Send([]byte("Order is already min"))
	}
	// Start the transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	// Step 1: Temporarily set order of the first row to -1
	_, err = tx.Exec(`UPDATE `+t+` SET "order" = -1 WHERE "order" = ?;`, ord)
	if err != nil {
		log.Printf("Error updating order: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	// Step 2: Update the second row to the original order of the first row
	_, err = tx.Exec(`UPDATE `+t+` SET "order" = ? WHERE "order" = ?;`, ord, ord+dir)
	if err != nil {
		log.Printf("Error updating order: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	// Step 3: Set the first row (with order -1) to the new order
	_, err = tx.Exec(`UPDATE `+t+` SET "order" = ? WHERE "order" = -1;`, ord+dir)
	if err != nil {
		log.Printf("Error updating order: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Order changed successfully"))
}

func dropRelTable(db *sql.DB, t string) error {
	// rm "" from t
	t = strings.ReplaceAll(t, `"`, "")
	dropTableSQL := `DROP TABLE IF EXISTS Relation` + t
	_, err := db.Exec(dropTableSQL)
	if err != nil {
		log.Printf("Error dropping table: %v\n", err)
		return err
	}
	return nil
}

func lsallComp(db *sql.DB, c *fiber.Ctx) error {
	rows, err := db.Query("SELECT type FROM Component ORDER BY type;")
	if err != nil {
		log.Printf("Error getting components: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Printf("Error scanning components: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		values = append(values, value)
	}
	return c.JSON(values)
}

func addCompType(db *sql.DB, c *fiber.Ctx) error {
	// json: "type": "I"
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	log.Printf("data: %v\n", data)
	t := data["type"]
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	_, err = tx.Exec(`INSERT INTO Component (type) VALUES (?)`, t)
	if err != nil {
		log.Printf("Error inserting component type: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// create table
	createTableSQL := `CREATE TABLE IF NOT EXISTS '` + t + `' (
		"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name"	TEXT NOT NULL UNIQUE,
		"description"	TEXT,
		"price"	INTEGER NOT NULL);`
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating component's table: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Component type added successfully"))
}
func delCompType(db *sql.DB, c *fiber.Ctx) error {
	// json: "c": ["I"]
	var data []string
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	t := data[0]
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	_, err = tx.Exec(`DELETE FROM Component WHERE type = ?`, t)
	if err != nil {
		log.Printf("Error deleting component type: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	_, err = tx.Exec(`DROP TABLE IF EXISTS '` + t + `'`)
	if err != nil {
		log.Printf("Error dropping component's table: %v\n", err)
		tx.Rollback()
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.Send([]byte("Component type deleted successfully"))
}

func getComponents(db *sql.DB, c *fiber.Ctx) error {
	prod := c.Query("prod")
	if prod == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	rows, err := db.Query("SELECT component FROM Relation" + prod + " ORDER BY [order];")
	if err != nil {
		log.Printf("Error getting components: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Printf("Error scanning components: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		values = append(values, value)
	}
	return c.JSON(values)
}

func addComponents(db *sql.DB, n string, c []string) error {
	for _, comp := range c {
		insertCompSQL := `INSERT INTO Relation` + n + ` (component) VALUES (?)`
		_, err := db.Exec(insertCompSQL, comp)
		if err != nil {
			log.Printf("Error inserting component: %v\n", err)
			return err
		}
	}
	return nil
}

func getComponentsName(db *sql.DB, c *fiber.Ctx) error {
	t := c.Query("t")
	// rows, err := db.Query("SELECT name FROM Component WHERE type='" + t + "';")
	// rows, err := db.Query("SELECT name, description, price FROM Component WHERE type='" + t + "';")
	rows, err := db.Query("SELECT name, description, price FROM '" + t + "' ORDER BY name;")
	if err != nil {
		log.Printf("Error getting components data: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer rows.Close()
	var values = make(map[string]interface{})
	for rows.Next() {
		var name, desc, price string
		err := rows.Scan(&name, &desc, &price)
		if err != nil {
			log.Printf("Error scanning components data: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		values[name] = map[string]string{"name": name, "desc": desc, "price": price}
	}
	log.Printf("values: %v\n", values)
	return c.JSON(values)
}

func pairComponents(db *sql.DB, c *fiber.Ctx) error {
	prod := c.Query("prod")
	name := c.Query("name")
	next := c.Query("next")
	if prod == "" || name == "" || next == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Missing query parameter(s)"))
	}
	// split name
	nSlice := strings.Split(name, "-")
	if len(nSlice) != 2 || nSlice[0] == "" || nSlice[1] == "" {
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Invalid query parameter 'name'"))
	}
	r, err := db.Query("SELECT CASE WHEN component = '" + nSlice[0] + "' THEN related ELSE component END AS result FROM Relation" + prod + " WHERE (component = '" + nSlice[0] + "' AND related IS NOT NULL) OR related = '" + nSlice[0] + "';")
	if err != nil {
		log.Printf("Error checking relation: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	defer r.Close()
	var rval []string
	for r.Next() {
		var rvalc string
		err = r.Scan(&rvalc)
		if err != nil {
			log.Printf("Error getting relation: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		if rvalc == "" {
			continue
		}
		rval = append(rval, rvalc)
	}
	log.Printf("rval: %v\n", rval) // debug
	var returnValues = make(map[string][]interface{})
	for i := range rval {
		// log.Printf("rval[i]: %v\n", rval[i]) // debug
		rel, err := db.Query("SELECT name FROM sqlite_master WHERE type = 'table' AND (name like '%" + nSlice[0] + "%%" + rval[i] + "%' OR name like '%" + rval[i] + "%%" + nSlice[0] + "%');")
		if err != nil {
			log.Printf("Error getting relation: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		defer rel.Close()
		var relv string
		if rel.Next() {
			err = rel.Scan(&relv)
			if err != nil {
				log.Printf("Error getting relation: %v\n", err)
				return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
		} else {
			c.Status(fiber.StatusBadRequest)
			return c.Send([]byte("Relation exists but no table found"))
		}
		// log.Printf("relv: %v\n", relv) // debug
		rows, err := db.Query("SELECT " + rval[i] + ".name, " + rval[i] + ".description, " + rval[i] + ".price FROM " + nSlice[0] + " INNER JOIN " + relv + " ON " + nSlice[0] + ".id = " + relv + "." + nSlice[0] + "id INNER JOIN " + rval[i] + " ON " + rval[i] + ".id = " + relv + "." + rval[i] + "id WHERE " + nSlice[0] + ".name='" + name + "' ORDER BY " + rval[i] + ".name;")
		// rows, err := db.Query("SELECT " + rval[i] + ".name FROM " + nSlice[0] + " INNER JOIN " + relv + " ON " + nSlice[0] + ".id = " + relv + "." + nSlice[0] + "id INNER JOIN " + rval[i] + " ON " + rval[i] + ".id = " + relv + "." + rval[i] + "id WHERE " + nSlice[0] + ".name='" + name + "' ORDER BY " + rval[i] + ".name;")
		if err != nil {
			log.Printf("Error getting components name: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		defer rows.Close()
		var values = make(map[string]interface{})
		for rows.Next() {
			var name, desc, price string
			err := rows.Scan(&name, &desc, &price)
			if err != nil {
				log.Printf("Error scanning data: %v\n", err)
				return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
			values[name] = map[string]string{"name": name, "desc": desc, "price": price}
		}
		log.Printf("values: %v\n", values) // debug
		returnValues[rval[i]] = make([]interface{}, 0, len(values))
		for _, v := range values {
			returnValues[rval[i]] = append(returnValues[rval[i]], v)
		}
	}
	// if next is not in rval, get unrelated components
	if !contains(rval, next) {
		// log.Printf("next: %v\n", next)
		unrel, err := db.Query("SELECT name, description, price FROM '" + next + "';")
		if err != nil {
			log.Printf("Error getting unrelated components: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		defer unrel.Close()
		var unrelv = make(map[string]interface{})
		for unrel.Next() {
			var name, desc, price string
			err = unrel.Scan(&name, &desc, &price)
			if err != nil {
				log.Printf("Error getting unrelated component: %v\n", err)
				return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
			unrelv[name] = map[string]string{"name": name, "desc": desc, "price": price}
		}
		log.Printf("unrelv: %v\n", unrelv) // debug
		returnValues[next] = make([]interface{}, 0, len(unrelv))
		for _, v := range unrelv {
			returnValues[next] = append(returnValues[next], v)
		}
	}
	log.Printf("values: %v\n", returnValues) // debug
	return c.JSON(returnValues)
}

func CalcPrice(db *sql.DB, c *fiber.Ctx) error {
	// calc-price
	var comp []string
	if err := c.BodyParser(&comp); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	log.Printf("comp: %v\n", comp)
	var price float64
	for _, name := range comp {
		// name:"Z-2" -> table:"Z"
		nSlice := strings.Split(name, "-")
		if len(nSlice) != 2 || nSlice[0] == "" || nSlice[1] == "" {
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Invalid component name"))
		}
		table := nSlice[0]
		rows, err := db.Query("SELECT price FROM " + table + " WHERE name='" + name + "';")
		if err != nil {
			log.Printf("Error getting component price: %v\n", err)
			return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		defer rows.Close()
		var p float64
		if rows.Next() {
			err = rows.Scan(&p)
			if err != nil {
				log.Printf("Error scanning component price: %v\n", err)
				return c.App().ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
		}
		price += p
	}
	return c.Send([]byte(strconv.FormatFloat(price, 'f', -1, 64)))
}

func ProductFunc(app *fiber.App, db *sql.DB) {
	app.Put("/create-product", func(c *fiber.Ctx) error {
		return createProduct(db, c)
	})
	app.Delete("/delete-product", func(c *fiber.Ctx) error {
		return deleteProduct(db, c)
	})
	app.Get("/list-products", func(c *fiber.Ctx) error {
		return listProducts(db, c)
	})
	app.Get("/get-product", func(c *fiber.Ctx) error {
		return ProductDetails(db, c)
	})
	app.Post("/create-comp-rel", func(c *fiber.Ctx) error {
		return createCompRelTable(db, c)
	})
	app.Get("/all-components", func(c *fiber.Ctx) error {
		return lsallComp(db, c)
	})
	app.Get("/get-components", func(c *fiber.Ctx) error {
		return getComponents(db, c)
	})
	app.Post("/add-comp-type", func(c *fiber.Ctx) error {
		return addCompType(db, c)
	})
	app.Delete("/del-comp-type", func(c *fiber.Ctx) error {
		return delCompType(db, c)
	})
	app.Get("/get-comp-name", func(c *fiber.Ctx) error {
		return getComponentsName(db, c)
	})
	app.Get("/pair-components", func(c *fiber.Ctx) error {
		return pairComponents(db, c)
	})
	app.Post("/chg-order", func(c *fiber.Ctx) error {
		return chgRelOrder(db, c)
	})
	app.Post("/calc-price", func(c *fiber.Ctx) error {
		return CalcPrice(db, c)
	})
}
