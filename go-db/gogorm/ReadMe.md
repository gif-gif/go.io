# go-db 数据操作

- https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/
- `基于 https://github.com/go-gorm/gorm 的封装，更多用法到官方`
# 目前已封装
- mysql
- sqlite3
- clickhouse
- postgresql
- sqlserver
- tidb

```go
func testSqlite3() {
    db, err := gogorm.InitSqlite3("./test.db", gogorm.GoDbConfig{
        MaxOpen:      100,
        MaxIdleCount: 10,
    })
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
        return
    }
    err = db.AutoMigrate(&Product{})
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
        return
    }
    
    // Create
    insertProduct := &Product{Code: "D42", Price: 100}
    db.Insert(insertProduct)
    fmt.Println(insertProduct.ID)
    // Read
    var product Product
    tx := db.First(&product, 1) // find product with integer primary key
    if tx.Error != nil {
    fmt.Println("not found first ", tx.Error.Error())
    }
    db.First(&product, "code = ?", "D42")
    // Delete - delete product
    db.Delete(&product, 1)
    
    err = goutils.RemoveFile("./test.db")
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
    }
}
```

```go
func mysqlTest() {
	db, err := gogorm.InitMysql("root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local", gogorm.GoDbConfig{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}
	err = db.AutoMigrate(&Product{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}

	// Create
	insertProduct := &Product{Code: "D42", Price: 100}
	db.Insert(insertProduct)
	fmt.Println(insertProduct.ID)
	// Read
	var product Product
	tx := db.First(&product, 1) // find product with integer primary key
	if tx.Error != nil {
		fmt.Println("not found first ", tx.Error.Error())
	}
	db.First(&product, "code = ?", "D42")
	// Delete - delete product
	db.Delete(&product, 1)

}
```

```go
func testClickhouse() {
	dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
	db, err := gogorm.InitMysql(dsn, gogorm.GoDbConfig{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}
	err = db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&Product{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}
	// Set table options

	// Create
	insertProduct := &Product{Code: "D42", Price: 100}
	db.Insert(insertProduct)
	fmt.Println(insertProduct.ID)
	// Read
	var product Product
	tx := db.First(&product, 1) // find product with integer primary key
	if tx.Error != nil {
		fmt.Println("not found first ", tx.Error.Error())
	}
	db.First(&product, "code = ?", "D42")
	// Delete - delete product
	db.Delete(&product, 1)
}
```

```
mongodb 不是关系性数据库 暂时不支持
```
