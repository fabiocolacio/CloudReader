package server

import(
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var(
    db *sql.DB
)

func init() {
    source := fmt.Sprintf("%s:%s@/%s",SqlUser, SqlPass, SqlDb)

    var err error
    db, err = sql.Open("mysql", source)
    if err != nil {
        fmt.Println(err)
    }
}

func InitTables() error {
    query := fmt.Sprintf(
        `create table users(
            uid int primary key auto_increment,
            username varchar(%d),
            salt binary(%d),
            saltedhash binary(%d));`,
        UsernameMaxLength,
        SaltLength,
        KeyHashLength)
    _, err := db.Exec(query)

    query = fmt.Sprintf(
        `create table books(
            owner int,
            path varchar(%d),
            hash binary(%d),
            primary key (owner, path),
            foreign key (owner) references users(uid));`,
        PathMaxLength,
        BookHashLength)
    _, err = db.Exec(query)
    if err != nil {
        fmt.Println(err)
    }

    return err
}

func ResetTables() error {
    _, err := db.Exec(`drop table if exists books`)
    if err != nil {
        fmt.Println(err)
    }

    _, err = db.Exec(`drop table if exists users`)
    if err != nil {
        fmt.Println(err)
    }

    err = InitTables()
    if err != nil {
        fmt.Println(err)
    }

    return err
}

