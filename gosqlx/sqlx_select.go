package gosqlx

const (
	selectCount = 1000
)

func select1(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	var id int
	var name string
	var age string
	for i := 1; i <= selectCount; i++ {
		if rows, err := sxd.db.Queryx(selectOneSql, i); err != nil {
			continue
		} else {
			for rows.Next() {
				if err = rows.Scan(&id, &name, &age); err != nil {
					break
				} else {
					succCount++
				}
			}
		}
	}

	return succCount
}

func select2(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	var id int
	var name string
	var age string
	for i := 1; i <= selectCount; i++ {
		if stmt, err := sxd.db.Preparex(selectOneSql); err != nil {
			continue
		} else if rows, err := stmt.Queryx(i); err != nil {
			_ = stmt.Close()
			continue
		} else {
			for rows.Next() {
				if err = rows.Scan(&id, &name, &age); err != nil {
					break
				} else {
					succCount++
				}
			}
		}
	}
	return succCount
}

func select3(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	if stmt, err := sxd.db.Preparex(selectOneSql); err != nil {
		return succCount
	} else {
		var id int
		var name string
		var age string
		for i := 1; i <= selectCount; i++ {
			if rows, err := stmt.Queryx(i); err != nil {
				continue
			} else {
				for rows.Next() {
					if err = rows.Scan(&id, &name, &age); err != nil {
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

func select4(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	var id int
	var name string
	var age string
	for i := 1; i <= selectCount; i++ {
		if rows, err := sxd.db.NamedQuery(selectOneNamedSql, map[string]interface{}{
			"id": i,
		}); err != nil {
			continue
		} else {
			for rows.Next() {
				if err = rows.Scan(&id, &name, &age); err != nil {
					break
				} else {
					succCount++
				}
			}
		}
	}

	return succCount
}

func select5(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	var id int
	var name string
	var age string
	for i := 1; i <= selectCount; i++ {
		if stmt, err := sxd.db.PrepareNamed(selectOneNamedSql); err != nil {
			continue
		} else if rows, err := stmt.Queryx(map[string]interface{}{
			"id": i,
		}); err != nil {
			_ = stmt.Close()
			continue
		} else {
			for rows.Next() {
				if err = rows.Scan(&id, &name, &age); err != nil {
					break
				} else {
					succCount++
				}
			}
		}
	}
	return succCount
}

func select6(sxd *SqlxDemo) int64 {
	var succCount int64 = 0
	if stmt, err := sxd.db.PrepareNamed(selectOneNamedSql); err != nil {
		return succCount
	} else {
		var id int
		var name string
		var age string
		for i := 1; i <= selectCount; i++ {
			if rows, err := stmt.Queryx(map[string]interface{}{
				"id": i,
			}); err != nil {
				continue
			} else {
				for rows.Next() {
					if err = rows.Scan(&id, &name, &age); err != nil {
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
