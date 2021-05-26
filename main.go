package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDSQ = "root:root@tcp(127.0.0.1:3306)/mysql-prepared-conn-test"

func main() {
	fmt.Println("Connecting to", mysqlDSQ)

	db, err := sql.Open("mysql", mysqlDSQ)
	if err != nil {
		fmt.Println("Unable to connect to mysql:", err)
		os.Exit(1)
	}
	defer db.Close()

	printAfter(db, "", nil)

	printAfter(db, "do not close stmt", func() error {
		stmt, err := db.Prepare("SELECT ?")
		if err != nil {
			return err
		}

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "closing stmt", func() error {
		stmt, err := db.Prepare("SELECT ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "with transaction: do not close stmt nether commit", func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare("SELECT ?")
		if err != nil {
			return err
		}

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "with transaction: closing stmt but not commit", func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare("SELECT ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "with transaction: do not close stmt but commit", func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Commit()

		stmt, err := tx.Prepare("SELECT ?")
		if err != nil {
			return err
		}

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "with transaction: do not close stmt but rollback", func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		stmt, err := tx.Prepare("SELECT ?")
		if err != nil {
			return err
		}

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "with transaction: close stmt and commit", func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Commit()

		stmt, err := tx.Prepare("SELECT ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})

	printAfter(db, "with transaction: close stmt and rollback", func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		stmt, err := tx.Prepare("SELECT ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		var n int
		row := stmt.QueryRow("1")
		return row.Scan(&n)
	})
}

func printAfter(db *sql.DB, name string, fn func() error) {
	if fn != nil {
		fmt.Println()
		fmt.Println()
		fmt.Printf("Executing func %q\n", name)

		if err := fn(); err != nil {
			fmt.Println("Unable to execute func:", err)
			os.Exit(1)
		}
	}

	if err := printNumberPreparedStmt(db); err != nil {
		fmt.Println("Unable to print number of prepared statements:", err)
		os.Exit(1)
	}
}

func printNumberPreparedStmt(db *sql.DB) error {
	fmt.Println("Printing number of prepared statements:")

	rows, err := db.Query(`SHOW GLOBAL STATUS LIKE 'com_%prepare%'`)
	if err != nil {
		return err
	}
	if err := printLines(rows); err != nil {
		return err
	}
	rows.Close()

	rows, err = db.Query(`SHOW GLOBAL STATUS LIKE 'com_stmt_close'`)
	if err != nil {
		return err
	}
	if err := printLines(rows); err != nil {
		return err
	}
	rows.Close()
	return nil
}

func printLines(rows *sql.Rows) error {
	for rows.Next() {
		var (
			varName string
			value   uint32
		)
		if err := rows.Scan(&varName, &value); err != nil {
			return err
		}
		fmt.Printf("\t%-20s %d\n", varName+":", value)
	}
	return rows.Err()
}
