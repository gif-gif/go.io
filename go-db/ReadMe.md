# go-db 数据操作

`基于 https://github.com/go-gorm/gorm 的封装，更多用法到官方`

```go
func testSqlite3() {
    db, err := godb.InitSqlite3("./test.db", godb.GoDbConfig{
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
	db, err := godb.InitMysql("root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local", godb.GoDbConfig{})
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
	db, err := godb.InitMysql(dsn, godb.GoDbConfig{})
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
gorm 事物形式操作需要支持
```