package gosqlx

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type SqlxDemo struct {
	db *sqlx.DB
}

const (
	//datasourceName = "root:1234@tcp(127.0.0.1:3306)/ys?timeout=90s&collation=utf8mb4_bin&parseTime=True"
	datasourceName    = "root:Jingle@100@tcp(10.21.248.251:3306)/ys?timeout=90s&collation=utf8mb4_bin&parseTime=True"
	truncateSql       = "truncate table `user`"
	insertOneSql      = "INSERT INTO `user`(name, age) VALUES (?, ?)"
	insertOneNamedSql = "INSERT INTO `user`(name, age) VALUES (:name, :age)"
	deleteOneSql      = "DELETE FROM `user` WHERE `id` = ?"
	deleteOneNamedSql = "DELETE FROM `user` WHERE `id` = :id"
	updateOneSql      = "UPDATE `user` SET `age` = `age` + 100 WHERE `id` = ?"
	updateOneNamedSql = "UPDATE `user` SET `age` = `age` + 100 WHERE `id` = :id"
	selectOneSql      = "SELECT `id`, `name`, `age` from `user` WHERE `id` = ?"
	selectOneNamedSql = "SELECT `id`, `name`, `age` FROM `user` WHERE `id` = :id"
	selectMultiSql    = "SELECT `id`, `name`, `age` from `user` WHERE `id` > ?"
	selectByNameSql   = "SELECT `id`, `name`, `age` from `user` WHERE `name` = :name"
	transUpdateSql1   = "UPDATE `user` set `age` = 21 WHERE `id` = ?"
	transUpdateSql2   = "UPDATE `user` set `age` = 51 WHERE `id` = ?"
)

type User struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func NewSqlxDemo() (*SqlxDemo, error) {
	db, err := sqlx.Connect("mysql", datasourceName)
	if err != nil {
		log.Println("failed to connect mysql, ", err)
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	sxd := SqlxDemo{
		db: db,
	}

	return &sxd, nil
}

// insert one record
func (sxd *SqlxDemo) InsertRowDemo() error {
	result, err := sxd.db.Exec(insertOneSql, "江南小王子", 19)
	if err != nil {
		log.Println("failed to insert one record, ", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get lastInsertId(), ", err)
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get rowsAffected(), ", err)
		return err
	}

	fmt.Printf("lastInsertId: %d, rowsAffected: %d\n", id, n)

	return nil
}

// insert one record by NamedExec
func (sxd *SqlxDemo) InsertByNamedExecDemo() error {
	result, err := sxd.db.NamedExec(insertOneNamedSql, map[string]interface{}{
		"name": "江南草上飞",
		"age":  30,
	})
	if err != nil {
		log.Println("failed to insert by NamedExecDemo, ", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get lastInsertId(), ", err)
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get rowsAffected(), ", err)
		return err
	}

	fmt.Printf("lastInsertId: %d, rowsAffected: %d\n", id, n)
	return nil
}

// delete records
func (sxd *SqlxDemo) DeleteDemo() error {
	result, err := sxd.db.Exec(deleteOneSql, 3)
	if err != nil {
		log.Println("failed to exec deleteOneSql, ", err)
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get rowsAffected, ", err)
		return err
	}
	fmt.Println("delete finish, rows affected: ", n)
	return nil
}

// update records
func (sxd *SqlxDemo) UpdateDemo() error {
	result, err := sxd.db.Exec(updateOneSql, 0)
	if err != nil {
		log.Println("failed to exec updateOneSql, ", err)
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get rows affected, ", err)
		return err
	}

	fmt.Println("finished update, rows affected: ", n)
	return nil
}

// query one record
func (sxd *SqlxDemo) QueryRowDemo() error {
	var u User
	err := sxd.db.Get(&u, selectOneSql, 1)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("select one result: ", err)
			return nil
		} else {
			log.Println("failed to select one, ", err)
			return err
		}
	}
	fmt.Println("select one result: ", u.ID, u.Name, u.Age)
	return nil
}

// query multi records
func (sxd *SqlxDemo) QueryMultiRowDemo() error {
	var users []User
	err := sxd.db.Select(&users, selectMultiSql, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("select multi result: ", err)
			return nil
		} else {
			log.Println("failed to select multi, ", err)
			return err
		}
	}
	fmt.Printf("select multi result-> users: %#v\n", users)

	return nil
}

// query by NamedQuery
func (sxd *SqlxDemo) QueryByNamedQueryDemo() error {
	u := User{
		Name: "江南草上飞",
	}
	rows, err := sxd.db.NamedQuery(selectByNameSql, u)
	if err != nil {
		log.Println("failed to selectByNameSql, ", err)
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	fmt.Println("QueryByNamedQueryDemo result: ")
	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err != nil {
			log.Printf("failed to StructScan, err:%v \n", err)
			continue
		}
		fmt.Printf("\tuser: %#v\n", u)
	}

	return nil
}

// transaction test demo
func (sxd *SqlxDemo) TransactionDemo() error {
	tx, err := sxd.db.Begin()
	if err != nil {
		fmt.Println("begin transaction failed: ", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
			fmt.Println("commit")
		}
	}()

	rs, err := tx.Exec(transUpdateSql1, 2)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec transUpdateSql1 failed")
	}

	rs, err = tx.Exec(transUpdateSql2, 3)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec transUpdateSql2 failed")
	}

	return nil
}
