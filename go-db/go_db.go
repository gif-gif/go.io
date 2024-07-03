package godb

import (
	"gorm.io/gorm"
	"time"
)

type GoDB struct {
	DB *gorm.DB
}

func (s *GoDB) Init(config *GoDbConfig) error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}

	if sqlDB != nil {
		if config.MaxIdleCount == 0 {
			config.MaxIdleCount = 10
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(config.MaxIdleCount)
		if config.MaxOpen == 0 {
			config.MaxOpen = 100
		}
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(config.MaxOpen)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		if config.MaxLifetime == 0 {
			config.MaxLifetime = time.Hour
		}
		sqlDB.SetConnMaxLifetime(config.MaxLifetime)
	}

	return nil
}

func (s *GoDB) Model(value interface{}) *gorm.DB {
	return s.DB.Model(value)
}

// 更新用户并完全更新其所有关联
//
// db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
//
// SQL：完全更新地址、用户、电子邮件表，包括现有的关联记录
func (s *GoDB) Session(config *gorm.Session) *gorm.DB {
	return s.DB.Session(config)
}

// 事物开始
//
//	func CreateAnimals(db *gorm.DB) error {
//	 // 再唠叨一下，事务一旦开始，你就应该使用 tx 处理数据
//	 tx := db.Begin()
//	 defer func() {
//	   if r := recover(); r != nil {
//	     tx.Rollback()
//	   }
//	 }()
//
//	 if err := tx.Error; err != nil {
//	   return err
//	 }
//
//	 if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
//	    tx.Rollback()
//	    return err
//	 }
//
//	 if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
//	    tx.Rollback()
//	    return err
//	 }
//
//	 return tx.Commit().Error
//	}
func (s *GoDB) BeginTransaction() *GoDB {
	godb := &GoDB{
		DB: s.DB.Begin(),
	}
	return godb
}

// 事物回滚
func (s *GoDB) Rollback() {
	s.DB.Rollback()
}

// 保存回滚点
func (s *GoDB) SavePoint(name string) {
	s.DB.SavePoint(name)
}

// 事物回滚到name的点击
//
// tx := db.Begin()
// tx.Create(&user1)
//
// tx.SavePoint("sp1")
// tx.Create(&user2)
// tx.RollbackTo("sp1") // Rollback user2
//
// tx.Commit() // Commit user1
func (s *GoDB) RollbackTo(name string) {
	s.DB.RollbackTo(name)
}

// 提交事物
func (s *GoDB) Commit() error {
	return s.DB.Commit().Error
}

// package main
//
// import (
//
//	"gorm.io/gorm"
//	"gorm.io/driver/sqlite"
//
// )
//
//	type Product struct {
//	 gorm.Model
//	 Code  string
//	 Price uint
//	}
//
//	func main() {
//	 db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//	 if err != nil {
//	   panic("failed to connect database")
//	 }
//
//	 // 迁移 schema
//	 db.AutoMigrate(&Product{})
//
//	 // Create
//	 db.Create(&Product{Code: "D42", Price: 100})
//
//	 // Read
//	 var product Product
//	 db.First(&product, 1) // 根据整型主键查找
//	 db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
//
//	 // Update - 将 product 的 price 更新为 200
//	 db.Model(&product).Update("Price", 200)
//	 // Update - 更新多个字段
//	 db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
//	 db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
//
//	 // Delete - 删除 product
//	 db.Delete(&product, 1)
//	}
//
// db.AutoMigrate(&User{})
//
// db.AutoMigrate(&User{}, &Product{}, &Order{})
//
// // Add table suffix when creating tables
// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
func (s *GoDB) AutoMigrate(values ...interface{}) error {
	return s.DB.AutoMigrate(values...)
}

// clickhouse Set table options after AutoMigrate
//
//	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&User{})
func (s *GoDB) Set(key string, value interface{}) (tx *gorm.DB) {
	return s.DB.Set(key, value)
}

func (s *GoDB) Get(key string) (interface{}, bool) {
	return s.DB.Get(key)
}

//	type Result struct {
//	 ID   int
//	 Name string
//	 Age  int
//	}
//
// var result Result
// db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)
//
// db.Raw("SELECT id, name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)
//
// var age int
// db.Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Scan(&age)
//
// var users []User
// db.Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Scan(&users)
func (s *GoDB) Raw(sql string, dest interface{}, values ...interface{}) (tx *gorm.DB) {
	return s.DB.Raw(sql, values...).Scan(dest)
}

// db.Exec("DROP TABLE users")
// db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})
//
// // Exec with SQL Expression
// db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
func (s *GoDB) Exec(sql string, values ...interface{}) (tx *gorm.DB) {
	return s.DB.Exec(sql, values...)
}

func (s *GoDB) Exe1c(table string, dest interface{}, where string, args ...interface{}) (tx *gorm.DB) {
	return s.DB.Table(table).Where(where, args).Scan(dest)
}

// Insert or Batch Insert
//
// user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
//
//	users := []*User{
//		 {Name: "Jinzhu", Age: 18, Birthday: time.Now()},
//		 {Name: "Jackson", Age: 19, Birthday: time.Now()},
//		}
//
// result := db.Create(&user)  pass pointer of data to Create
//
// user.ID              returns inserted data's primary key
//
// result.Error         returns error
//
// result.RowsAffected  returns inserted records count
//
// Create Hooks
//
//	func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
//	 u.UUID = uuid.New()
//
//	 if u.Role == "admin" {
//	   return errors.New("invalid role")
//	 }
//	 return
//	}
//
// Default Values
//
//	type User struct {
//	 ID   int64
//	 Name string `gorm:"default:galeone"`
//	 Age  int64  `gorm:"default:18"`
//	}
func (s *GoDB) Create(value interface{}) (tx *gorm.DB) {
	return s.DB.Create(value)
}

// 特殊字段处理 Insert or Batch Insert
func (s *GoDB) InsertSpecified(fields []string, exclude bool, value interface{}) (tx *gorm.DB) {
	if exclude {
		return s.DB.Omit(fields...).Create(&value)
	} else {
		fs := []interface{}{}
		for _, field := range fields {
			fs = append(fs, field)
		}
		return s.DB.Select(fs[0], fs[1:]...).Create(&value)
	}
}

func (s *GoDB) Select(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.Select(value, conds...)
}

// Get the first record ordered by primary key
// SELECT * FROM users ORDER BY id LIMIT 1;
//
// var product Product
// db.First(&product, 1) // find product with integer primary key
// db.First(&product, "code = ?", "D42") // find product with code D42
func (s *GoDB) First(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.First(value, conds...)
}

// Get last record, ordered by primary key desc
//
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
func (s *GoDB) Last(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.Last(value, conds...)
}

// Get one record, no specified order
// SELECT * FROM users LIMIT 1;
//
//	check error ErrRecordNotFound
//
// errors.Is(result.Error, gorm.ErrRecordNotFound)
func (s *GoDB) FindOne(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.Take(value, conds...)
}

//	func (u *User) AfterFind(tx *gorm.DB) (err error) {
//	 // Custom logic after finding a user
//	 if u.Role == "" {
//	   u.Role = "user" // Set default role if not specified
//	 }
//	 return
//	}
func (s *GoDB) Find(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.Find(value, conds...)
}

func (s *GoDB) FirstOrCreate(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.FirstOrCreate(value, conds...)
}

func (s *GoDB) Save(value interface{}) (tx *gorm.DB) {
	return s.DB.Save(value)
}

func (s *GoDB) Update(model interface{}, column string, value interface{}) (tx *gorm.DB) {
	return s.DB.Model(model).Update(column, value)
}

// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
func (s *GoDB) Updates(model interface{}, value interface{}) (tx *gorm.DB) {
	return s.DB.Model(model).Updates(value)
}

// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
func (s *GoDB) UpdatesByMap(model interface{}, value map[string]interface{}) (tx *gorm.DB) {
	return s.DB.Model(model).Updates(value)
}

// 如果你的模型包含了 gorm.DeletedAt字段（该字段也被包含在gorm.Model中），那么该模型将会自动获得软删除的能力！
//
// 当调用Delete时，GORM并不会从数据库中删除该记录，而是将该记录的DeleteAt设置为当前时间，而后的一般查询方法将无法查找到此条记录。
func (s *GoDB) Delete(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	return s.DB.Delete(value, conds...)
}
