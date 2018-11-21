package server

import(
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "errors"
)

var(
    db *sql.DB
    ErrInvalidCredentials error = errors.New("Invalid username or password.")
    ErrUsernameTaken error = errors.New("Username already taken.")
    ErrRegistrationFailed error = errors.New("Failed to register user.")
)

func init() {
    source := fmt.Sprintf("%s:%s@/%s",SqlUser, SqlPass, SqlDb)

    var err error
    db, err = sql.Open("mysql", source)
    if err != nil {
        fmt.Println(err)
    }
}

func AddUsers(username string, salt []byte, saltedhash []byte) error{
  row := db.QueryRow("select uid from users where username = ?;", username)
var err error
var uid int
if row.Scan(&uid) == sql.ErrNoRows {

    _, err = db.Exec(
        `insert into users (username, salt, saltedhash) values (?, ?, ?);`,
        username, salt, saltedhash)

    if err != nil {
        err = ErrRegistrationFailed
    }
} else {
    err = ErrUsernameTaken
}

return err
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
            data longblob,
            hash binary(%d),
            primary key (owner, hash),
            foreign key (owner) references users(uid));`,
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
