package gosql

const (
	//datasourceName = "root:1234@tcp(127.0.0.1:3306)/ys?timeout=90s&collation=utf8mb4_bin"
	datasourceName = "root:Jingle@100@tcp(10.21.248.251:3306)/ys?timeout=90s&collation=utf8mb4_bin"

	selectByNumSql   = "SELECT square_num FROM square_num WHERE num = ?"
	updateByNumEqSql = "UPDATE square_num SET square_num = square_num + 1 WHERE num = ?"
	updateByNumLtSql = "UPDATE square_num SET square_num = square_num - 1 WHERE num < 3"
	updateByNumGtSql = "UPDATE square_num SET square_num = square_num + 5 WHERE num > 2"

	truncateSql = "truncate table `square_num`"
	insertSql   = "INSERT INTO `square_num`(num, square_num) VALUES (?, ?)"
	deleteSql   = "DELETE FROM `square_num` WHERE `num` = ?"
	updateSql   = "UPDATE `square_num` SET `square_num` = `num` * 2 WHERE `num` = ?"
	selectSql   = "SELECT square_num FROM square_num WHERE num = ?"
)
