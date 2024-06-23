# go-db 数据操作

`基于 https://github.com/go-gorm/gorm 的封装，更多用法到官方`

```go
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
```

```
gorm 事物形式操作需要支持
```