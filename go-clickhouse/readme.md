# go-clickhouse clickhouse/v2版本 
# 尽支持 http/https database/sql interface ，即将支持 native interface
# 初始化 

```
goclickhouse.Init(clickhouse.Config{
    Driver:   "clickhouse",
    Addr:     "192.168.1.100:9000",
    User:     "root",
    Password: "123456",
    Database: "test",
})
```

# 创建数据库

```
CREATE DATABASE IF NOT EXISTS test;
```

# 创建表

```
sqlstr := `
    CREATE TABLE IF NOT EXISTS user
    (
        name String,
        gender String,
        birthday Date
    )
    ENGINE = MergeTree()
    ORDER BY (name, gender)
    PARTITION BY toYYYYMM(birthday)
`
if _, err := goclickhouse.DB().Exec(sqlstr); err != nil {
    log.Fatal(err)
}
```

# 添加数据

```
func insert() {
	var (
		tx, _   = goclickhouse.DB().Begin()
		stmt, _ = tx.Prepare(`INSERT INTO user(name, gender) VALUES(?, ?)`)
	)

    data := []interface{}{"", ""}
    if _, err := stmt.Exec(data...); err != nil {
        golog.Error(err)
        return
    }

    if err := tx.Commit(); err != nil {
        golog.Error(err)
	}
}
```

# 查询数据

```
rows, err := DB().Query("SELECT name, gender FROM user")
if err != nil {
    golog.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var (
        name string
        gender string
    )
    if err := rows.Scan(&name, &gender); err != nil {
        log.Fatal(err)
    }
    fmt.Println(name, gender)
}

if err := rows.Err(); err != nil {
    golog.Fatal(err)
}
```

# 删除表

```
if _, err := DB().Exec("DROP TABLE user"); err != nil {
	golog.Fatal(err)
}
```