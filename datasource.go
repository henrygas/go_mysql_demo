package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//var datasourceName = "root:1234@tcp(127.0.0.1:3306)/ys?timeout=90s&charset=utf8mb4&collation=utf8mb4_general_ci"

func testOpen() {
	db, err := sql.Open("mysql", "root:1234@/ys")
	if err != nil {
		log.Fatal("fail to open mysql driver, ", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("fail to ping mysql, ", err)
		return
	} else {
		log.Println("succ to ping mysql!")
	}
}

func testInsertAndQuery() {
	db, err := sql.Open("mysql", "root:1234@/ys")
	if err != nil {
		log.Fatal("fail to open mysql", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("fail to ping mysql", err)
		return
	}

	// INSERT
	stmtIn, err := db.Prepare("INSERT INTO square_num(num, square_num) VALUES(?, ?)")
	if err != nil {
		log.Fatal("fail to prepare INSERT statement", err)
		return
	}
	defer stmtIn.Close()
	for i := 26; i <= 26; i++ {
		n, err := stmtIn.Exec(i, (i * i))
		if err != nil {
			log.Fatal(fmt.Sprintf("fail to INSERT num=%d", i), err)
		} else {
			log.Printf("succ to INSERT, num=%d, rowsAffected=%d\n", i, n)
		}
	}

	// SELECT
	stmtQuery, err := db.Prepare("SELECT square_num FROM square_num WHERE num = ?")
	if err != nil {
		log.Fatal("failed to prepare SELECT statement", err)
		return
	}
	defer stmtQuery.Close()
	var squareNum int
	var num int

	num = 13
	err = stmtQuery.QueryRow(num).Scan(&squareNum)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to SELECT square_num BY num=%d", num), err)
	} else {
		log.Printf("succ to SELECT, num=%d, square_num=%d\n", num, squareNum)
	}

	num = 1
	err = stmtQuery.QueryRow(num).Scan(&squareNum)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to SELECT square_num BY num=%d", num), err)
	} else {
		log.Printf("succ to SELECT, num=%d, square_num=%d\n", num, squareNum)
	}
}

func testIgnoreNullValues() {
	db, err := sql.Open("mysql", "root:1234@/ys")
	if err != nil {
		log.Fatal("failed to open mysql", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM square_num")
	if err != nil {
		log.Fatal("failed to SELECT * FROM square_num", err)
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("failed to get Columns()", err)
		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i, value := range values {
		scanArgs[i] = &value
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal("failed to rows.Scan")
			return
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(value)
			}
			fmt.Println(columns[i], ":", value)
		}
		fmt.Println("==================================")
	}

}

func main() {
	//testOpen()
	//testInsertAndQuery()
	testIgnoreNullValues()
}
