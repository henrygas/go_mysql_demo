
本示例示范了go-sql-driver/mysql对增删查改和事务的支持

## 1. 部署步骤
+ 导入数据库
```sql
mysql -uroot -p < ./sql/ys.sql
```

+ 引入依赖
```go
go get github.com/go-sql-driver/mysql
go get github.com/jmoiron/sqlx
```

+ 运行程序
```go
go build
```

+ 依次体验go-sql增删查改方法和事务的使用[单元测试]

main.go里取消RunSqlDemo()行的注释，增加RunSqlxDemo()行的注释
然后运行 
    
    go build -o ./bin/go_sql_demo.exe .
    cd ./bin/
    go_sql_demo.exe

+ 依次体验go-sqlx增删查改方法和事务的使用

main.go里取消RunSqlxDemo()行的注释，增加RunSqlDemo()行的注释
然后运行 
    
    go build -o ./bin/go_sqlx_demo.exe .
    cd ./bin/
    go_sqlx_demo.exe

+ go-sql性能测试

./gosql/go_sql_demo_test.go
```
TestInsert()        # 增
TestDelete()        # 删
TestUpdate()        # 改
TestSelect()        # 查
```

+ go-sqlx性能测试

./gosqlx/go_sqlx_demo_test.go
```
TestInsert()        # 增
TestDelete()        # 删
TestUpdate()        # 改
TestSelect()        # 查
```

## 2. 注意事项
+ 所有操作记得Close
+ 事务记得rollback或者commit
+ 当同一个查询条件被重复使用，只是参数有所不同时, 可以使用prepare方法, 利用PrepareStatement，可复用SQL语句
+ 事务中不要出现db，因为tx和db使用的不是同一个数据库连接