package main

import (
	"fmt"
	"github.com/gif-gif/go.io/go-db/gogorm"
	gofile "github.com/gif-gif/go.io/go-file"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"github.com/gogf/gf/util/gconv"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	goio.Init(goio.DEVELOPMENT)
	//testSqlite3()
	//mysqlTest()
	//testTransaction()
	//postgreSqlTest()
	//testHasMany()

	tarTest()
}

func tarTest() {
	err := gogorm.Init(gogorm.Config{
		Name:       "starrocks",
		DBType:     gogorm.DATABASE_STARROCKS,
		DataSource: "root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local",
	})

	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}
}

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
	db.Create(insertProduct)
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

	err = gofile.RemoveFile("./test.db")
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
	}
}

func mysqlTest() {
	err := gogorm.Init(gogorm.Config{
		DataSource: "root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local",
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
	db.Create(insertProduct)
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

func testClickhouse() {
	err := gogorm.Init(gogorm.Config{
		DataSource: "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20",
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
	db.Create(insertProduct)
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

func testTransaction() {
	err := gogorm.Init(gogorm.Config{
		DataSource: "root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local",
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

	tx := db.Begin()

	// Create
	insertProduct := &Product{Code: "D42", Price: 100}
	txd := tx.Create(insertProduct)
	if txd.Error != nil {
		fmt.Println("Insert error ", txd.Error.Error())
		tx.Rollback()
		return
	}

	newId := insertProduct.ID

	fmt.Println("Inserted ID: " + gconv.String(newId))
	// Read
	var product Product
	txd = tx.First(&product, insertProduct.ID) // find product with integer primary key
	if txd.Error != nil {
		fmt.Println("not found first error ", txd.Error.Error())
		tx.Rollback()
		return
	}

	txd = tx.First(&product, "code = ?", "D42")
	if txd.Error != nil {
		fmt.Println("not found first error1 ", txd.Error.Error())
		tx.Rollback()
		return
	}

	deleteProduct := Product{}
	// Delete - delete product
	deleteProduct.ID = newId
	txd = tx.Delete(&deleteProduct)
	if txd.Error != nil {
		fmt.Println("Delete error ", txd.Error.Error())
		tx.Rollback()
	}

	tx.Commit()
}

func testHasMany() {
	// User 有多张 CreditCard，UserID 是外键

	type CreditCard struct {
		gorm.Model
		Number    string
		UserRefer uint
	}

	type User struct {
		gorm.Model
		CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
	}

	err := gogorm.Init(gogorm.Config{
		DataSource: "root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local",
	})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}
	db := gogorm.Default().DB
	err = db.AutoMigrate(&User{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}

	user := User{
		CreditCards: []CreditCard{
			CreditCard{Number: "jinzhu"},
			CreditCard{Number: "jinzhu"},
		},
	}

	db.Create(&user)
	//db.Save(&user)

	// 检索用户列表并预加载信用卡
	var users []User
	err = db.Model(&User{}).Preload("CreditCard").Find(&users).Error
	if err != nil {
		golog.WithTag("godb").Error("检索用户列表并预加载信用卡:" + err.Error())
	} else {
		fmt.Println(users)
	}
}

func postgreSqlTest() {
	err := gogorm.Init(gogorm.Config{
		DataSource: "host=122.228.113.238 user=postgres password=223238 dbname=passwall port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		DBType:     gogorm.DATABASE_POSTGRESQL,
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
	db.Create(insertProduct)
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

func (u *Product) BeforeDelete(tx *gorm.DB) (err error) {
	//if u.Role == "admin" {
	//	return errors.New("admin user not allowed to delete")
	//}

	fmt.Println("BeforeDelete ")
	return nil
}

func (u *Product) BeforeCreate(tx *gorm.DB) (err error) {
	//if u.Role == "admin" {
	//	return errors.New("admin user not allowed to delete")
	//}

	fmt.Println("BeforeCreate ")
	return nil
}
