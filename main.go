package main

import (
	"go_mysql_demo/gosql"
	"go_mysql_demo/gosqlx"
	"log"
)

func RunSqlxDemo() {
	sxd, err := gosqlx.NewSqlxDemo()
	if err != nil {
		log.Fatal("failed to NewSqlxDemo(), ", err)
	}

	err = sxd.InsertRowDemo()
	if err != nil {
		log.Fatal("failed to InsertRowDemo(), ", err)
	}

	err = sxd.InsertByNamedExecDemo()
	if err != nil {
		log.Fatal("failed to InsertByNamedExecDemo(), ", err)
	}

	err = sxd.DeleteDemo()
	if err != nil {
		log.Fatal("failed to DeleteDemo(), ", err)
	}

	err = sxd.UpdateDemo()
	if err != nil {
		log.Fatal("failed to UpdateDemo(), ", err)
	}

	err = sxd.QueryRowDemo()
	if err != nil {
		log.Fatal("failed to QueryRowDemo(), ", err)
	}

	err = sxd.QueryMultiRowDemo()
	if err != nil {
		log.Fatal("failed to QueryMultiRowDemo(), ", err)
	}

	err = sxd.QueryByNamedQueryDemo()
	if err != nil {
		log.Fatal("failed to QueryByNamedQueryDemo(), ", err)
	}

	err = sxd.TransactionDemo()
	if err != nil {
		log.Fatal("failed to TransactionDemo(), ", err)
	}
}

func RunSqlDemo() {
	sd, err := gosql.NewSqlDemo()
	if err != nil {
		log.Fatal("failed to NewSqlDemo(), err: ", err)
	}

	err = sd.TruncateDemo()
	if err != nil {
		log.Fatal("failed to TruncateDemo(), err: ", err)
	}

	err = sd.InsertDemo()
	if err != nil {
		log.Fatal("failed to InsertDemo(), err: ", err)
	}

	err = sd.DeleteDemo()
	if err != nil {
		log.Fatal("failed to DeleteDemo(), err: ", err)
	}

	err = sd.UpdateDemo()
	if err != nil {
		log.Fatal("failed to UpdateDemo(), err: ", err)
	}

	err = sd.QueryDemo()
	if err != nil {
		log.Fatal("failed to QueryDemo(), err: ", err)
	}

	err = sd.TransactionDemo()
	if err != nil {
		log.Fatal("failed to TransactionDemo(), err: ", err)
	}
}

func main() {
	// ----------------- go-sql-driver/mysql + sql demo ------------------------
	//RunSqlDemo()

	// ----------------- go-sql-driver/mysql + sqlx demo ------------------------
	//RunSqlxDemo()
}
