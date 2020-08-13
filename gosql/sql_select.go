package gosql

const (
	selectCount = 1000
)

func Select1(sd *SqlDemo) int64 {
	var succCount int64 = 0
	var squareNum int
	for i := 1; i <= selectCount; i++ {
		if rows, err := sd.db.Query(selectSql, i); err != nil {
			continue
		} else {
			for rows.Next() {
				if err = rows.Scan(&squareNum); err != nil {
					break
				} else {
					succCount++
				}
			}
		}
	}
	return succCount
}

func Select2(sd *SqlDemo) int64 {
	var succCount int64 = 0
	for i := 1; i <= selectCount; i++ {
		stmt, err := sd.db.Prepare(selectSql)
		if err != nil {
			continue
		}

		rows, err := stmt.Query(i)
		if err != nil {
			_ = stmt.Close()
			continue
		}
		var squareNum int
		for rows.Next() {
			if err = rows.Scan(&squareNum); err != nil {
				break
			} else {
				succCount++
			}
		}
	}
	return succCount
}

func Select3(sd *SqlDemo) int64 {
	var succCount int64 = 0
	if stmt, err := sd.db.Prepare(selectSql); err == nil {
		var squareNum int
		for i := 1; i <= selectCount; i++ {
			if rows, err := stmt.Query(i); err != nil {
				continue
			} else {
				for rows.Next() {
					if err = rows.Scan(&squareNum); err != nil {
						break
					} else {
						succCount++
					}
				}
			}
		}
	}
	return succCount
}
