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
    err := gogorm.Init(gogorm.Config{
        DataSource: "./test.db",
        DBType:     gogorm.DATABASE_SQLITE,
    })
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
        return
    }
    db := gogorm.Default().DB
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
    err := gogorm.Init(gogorm.Config{
        DataSource: "root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Localb",
    })
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
        return
    }
    db := gogorm.Default().DB
	
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
    err := gogorm.Init(gogorm.Config{
        DataSource: dsn,
		DBType:     gogorm.DATABASE_CLICKHOUSE,
    })
	
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
        return
    }
    db := gogorm.Default().DB
	
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
