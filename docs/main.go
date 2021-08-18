package main

import (
	"database/sql"
	"fmt"

	_ "github.com/nlepage/go-sqljs-driver"
)

func main() {
	db, err := sql.Open("sqljs", "")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE hello (a int, b char)"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO hello VALUES (0, 'hello')"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO hello VALUES (1, 'world')"); err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("SELECT * FROM hello WHERE a=:aval AND b=:bval")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		sql.Named("aval", 1),
		sql.Named("bval", "world"),
	)

	var a int
	var b string
	if err := row.Scan(&a, &b); err != nil {
		panic(err)
	}

	fmt.Printf("a=%v, b=%v\n", a, b)

	rows, err := db.Query("SELECT * FROM hello")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&a, &b); err != nil {
			panic(err)
		}

		fmt.Printf("a=%v, b=%v\n", a, b)
	}
}
