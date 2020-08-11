package main

import (
	"go_mysql_demo/gosqlx"
	"log"
)

func main() {
	// ----------------- go-sql-driver/mysql + sql demo ------------------------

	//gosql.testOpen()
	//gosql.testInsert()
	//gosql.testDelete()
	//gosql.TestQuery()
	//gosql.testUpdate()
	//gosql.testTransaction()

	// ----------------- go-sql-driver/mysql + sqlx demo ------------------------

	sxd, err := gosqlx.NewSqlxDemo()
	if err != nil {
		log.Fatal("failed to NewSqlxDemo(), ", err)
	}

	//err = sxd.InsertRowDemo()
	//if err != nil {
	//	log.Fatal("failed to InsertRowDemo(), ", err)
	//}

	//err = sxd.InsertByNamedExecDemo()
	//if err != nil {
	//	log.Fatal("failed to InsertByNamedExecDemo(), ", err)
	//}

	//err = sxd.DeleteDemo()
	//if err != nil {
	//	log.Fatal("failed to DeleteDemo(), ", err)
	//}

	//err = sxd.UpdateDemo()
	//if err != nil {
	//	log.Fatal("failed to UpdateDemo(), ", err)
	//}

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

	//err = sxd.TransactionDemo()
	//if err != nil {
	//	log.Fatal("failed to TransactionDemo(), ", err)
	//}
}
