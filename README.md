
本示例示范了go-sql-driver/mysql对增删查改和事务的支持

## 1. 部署步骤
+ 导入数据库
```sql
mysql -uroot -p < ./sql/ys.sql
```

+ 运行程序
```go
go build datasource.go
```

+ 依次体验增删查改方法和事务的使用
```
testInsert()        # 增
testDelete()        # 删
testQuery()         # 查
testUpdate()        # 改
testTransaction()   # 事务
```

## 2. 注意事项
+ 所有操作记得Close
+ 事务记得rollback或者commit
+ 当同一个查询条件被重复使用，只是参数有所不同时, 可以使用prepare方法, 利用PrepareStatement，可复用SQL语句
+ 事务中不要出现db，因为tx和db使用的不是同一个数据库连接