package server

import(
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "errors"
    "crypto/subtle"
    "crypto/sha256"
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


func CheckUser(username string, password string) (int, error){
  row := db.QueryRow("select uid, salt, saltedhash from users where username = ?;", username)
  var(
    uid int
    salt []byte
    expectedsaltedhash []byte
  )

  if row.Scan(&uid, &salt, &expectedsaltedhash) == sql.ErrNoRows {
    return uid, ErrInvalidCredentials
  } else {
    saltedhash := HashAndSaltPassword([]byte(password), salt)
    if subtle.ConstantTimeCompare(expectedsaltedhash, saltedhash) == 1{
      return uid, nil
    }
  }

  return 0, ErrInvalidCredentials
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
            name varchar(256),
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

func UserExists(uid int) bool {
  row := db.QueryRow("select 1 from users where uid = ?;", uid)

  return row.Scan(nil) != sql.ErrNoRows
}

func UploadFile(filename string, data []byte, owner int) error{
  hasher := sha256.New()
  hasher.Write(data)
  hash := hasher.Sum(nil)

  _, err := db.Exec(`
  INSERT INTO books (DATA, NAME, HASH, OWNER)
  VALUES (?, ?, ?, ?)`, data, filename, hash, owner)
  return err
}

func GetFile(uid int, name string) ([]byte, error){
  row := db.QueryRow("select data from books where owner = ? AND name = ?;", uid, name)

  var data []byte
  err := row.Scan(&data)
  return data, err
}
