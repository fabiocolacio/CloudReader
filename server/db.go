package server

//db module contain all the essential functions used for
//application database
import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//variable and error object declarations
var (
	db                    *sql.DB
	ErrInvalidCredentials error = errors.New("Invalid username or password.")
	ErrUsernameTaken      error = errors.New("Username already taken.")
	ErrRegistrationFailed error = errors.New("Failed to register user.")
	ErrNoBook error = errors.New("User has no book.")
)

//initialize the database
func init() {
	source := fmt.Sprintf("%s:%s@/%s", SqlUser, SqlPass, SqlDb)

	var err error
	db, err = sql.Open("mysql", source)
	if err != nil {
		fmt.Println(err)
	}
}

//Add new Users object to the Users table in application database.
//username is the string represent the Username
//salt is the byte array used to hash with user Password
//saledthash is byte array of the hashed password concatenated with salt
//return an error if the username is already taken
func AddUsers(username string, salt []byte, saltedhash []byte) error {
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

//CheckUser check wether the user is in the database or not
//username is the string represent the user
//password is the string for the password of the user
//return user id if the user exist; return an error otherwise
func CheckUser(username string, password string) (int, error) {
	row := db.QueryRow("select uid, salt, saltedhash from users where username = ?;", username)
	var (
		uid                int
		salt               []byte
		expectedsaltedhash []byte
	)

	if row.Scan(&uid, &salt, &expectedsaltedhash) == sql.ErrNoRows {
		return uid, ErrInvalidCredentials
	} else {
		saltedhash := HashAndSaltPassword([]byte(password), salt)
		if subtle.ConstantTimeCompare(expectedsaltedhash, saltedhash) == 1 {
			return uid, nil
		}
	}

	return 0, ErrInvalidCredentials
}


//Create the Users and Books tables in the database.
//Return an error if the query does not execute successfully.
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

//Drop the tables in the database and recreate them.
//Return an error if the query does not execute successfully.
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

//Check if the user id is in the database.
//uid is an int represent the user ID in the database.
func UserExists(uid int) bool {
	row := db.QueryRow("select 1 from users where uid = ?;", uid)

	return row.Scan(nil) != sql.ErrNoRows
}

//UploadFile insert the book being uploaded by the user to the database.
//filename is the title of the book
//data is the book itself
//owner is the user who is uploading the book
//Return an error if the query does not execute successfully.
func UploadFile(filename string, data []byte, owner int) error {
	hasher := sha256.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)

	_, err := db.Exec(`
  INSERT INTO books (DATA, NAME, HASH, OWNER)
  VALUES (?, ?, ?, ?)`, data, filename, hash, owner)
	return err
}

//GetFile retrieve the book requested by the user from the database.
//uid is the user ID
//name is the title of the book being retrieved
//return the data of the book if exist; return an error otherwise
func GetFile(uid int, name string) ([]byte, error) {
	row := db.QueryRow("select data from books where owner = ? AND name = ?;", uid, name)

	var data []byte
	err := row.Scan(&data)
	return data, err
}

//ShowBooks show the books owned by a user
//uid is the user ID of the user
//return an array of string containing all the titles of the books
//return an error if there is no book owned by the user
func ShowBooks(uid int) ([]string, error) {
	var books []string
	rows, err := db.Query("select name from books where owner = ?;", uid)
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book string
		if err = rows.Scan(&book); err != nil {
			return books, err
		}
		books = append(books,book)
	}
	if len(books) == 0 {
		err = ErrNoBook
	}
	return books,err
}
