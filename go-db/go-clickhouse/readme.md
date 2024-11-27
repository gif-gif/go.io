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

# 批量执行
```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/ClickHouse/clickhouse-go/v2"
)

// 用户数据结构
type User struct {
    ID        int
    Name      string
    Age       int
    CreatedAt time.Time
}

func main() {
    // 建立连接
    conn := clickhouse.OpenDB(&clickhouse.Options{
        Addr: []string{"localhost:9000"},
        Auth: clickhouse.Auth{
            Database: "default",
            Username: "default",
            Password: "",
        },
        MaxOpenConns: 5,
        MaxIdleConns: 5,
    })
    defer conn.Close()

    ctx := context.Background()

    // 方法1：使用批处理预处理语句
    {
        tx, err := conn.Begin()
        if err != nil {
            log.Fatal(err)
        }

        stmt, err := tx.Prepare(`
            INSERT INTO users (id, name, age, created_at)
            VALUES (?, ?, ?, ?)
        `)
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        // 批量插入数据
        users := []User{
            {1, "User1", 25, time.Now()},
            {2, "User2", 30, time.Now()},
            {3, "User3", 35, time.Now()},
        }

        for _, user := range users {
            _, err = stmt.Exec(
                user.ID,
                user.Name,
                user.Age,
                user.CreatedAt,
            )
            if err != nil {
                tx.Rollback()
                log.Fatal(err)
            }
        }

        if err := tx.Commit(); err != nil {
            log.Fatal(err)
        }
    }

    // 方法2：使用批量值语法
    {
        _, err := conn.Exec(ctx, `
            INSERT INTO users (id, name, age, created_at)
            VALUES 
                (?, ?, ?, ?),
                (?, ?, ?, ?),
                (?, ?, ?, ?)
        `,
            4, "User4", 40, time.Now(),
            5, "User5", 45, time.Now(),
            6, "User6", 50, time.Now(),
        )
        if err != nil {
            log.Fatal(err)
        }
    }

    // 方法3：使用原生批量插入语法
    {
        batch, err := conn.PrepareBatch(ctx, `
            INSERT INTO users (id, name, age, created_at)
        `)
        if err != nil {
            log.Fatal(err)
        }

        users := []User{
            {7, "User7", 55, time.Now()},
            {8, "User8", 60, time.Now()},
            {9, "User9", 65, time.Now()},
        }

        for _, user := range users {
            err := batch.Append(
                user.ID,
                user.Name,
                user.Age,
                user.CreatedAt,
            )
            if err != nil {
                log.Fatal(err)
            }
        }

        if err := batch.Send(); err != nil {
            log.Fatal(err)
        }
    }
}
```
