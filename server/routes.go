package server

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"net/http"
	"strconv"
	//"mime/multipart"
)

var (
	httpServer  *http.Server
	httpsServer *http.Server
	macKey      [256]byte
)

func init() {
	httpServer = &http.Server{
		Addr:    HttpAddr,
		Handler: http.HandlerFunc(routeHandler),
	}
}

func ListenAndServe() error {
	return httpServer.ListenAndServe()
}

func routeHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	switch path {
	case "/":
		// TODO: Show home page

	case "/test":
		// TODO: Check if session cookies work
		uid := VerifyUser(req)
		if uid == 0 {
			res.Write([]byte(`Wrong Login Information`))
		} else {
			res.Write([]byte(`You have logged in`))
		}

	case "/login":
		// TODO: If it is a GET request, send login HTML page
		// If it is a POST request, log them in or tell them if there is
		// na error

		if req.Method == "GET" {
			res.Write([]byte(`
            <html>
            <head>
            <title> Login </title>
            </head>
            <body>
            <h1> Login </h1>
            <form action = "/login" method = "post">
              Username:<br>
                <input type="text" name="Username"><br>
              Password:<br>
                <input type = "password" name = "Password">
                <input type = "submit" value = "Login">
            </form>
            </body>
            </html>
            `))
		} else {
			req.ParseForm()
			username := req.FormValue("Username")
			password := req.FormValue("Password")

			uid, err := CheckUser(username, password)
			if err == nil {

				session := http.Cookie{
					Name:  "session",
					Value: strconv.Itoa(uid),

					//MaxAge: 10 * 60,
					Secure:   false,
					HttpOnly: true,
					SameSite: 1,

					Path: "/",
				}
				fmt.Println(session)
				http.SetCookie(res, &session)

			} else {
				res.Write([]byte(err.Error()))
			}
		}

	case "/register":
		// TODO: If it is a GET request, send the register HTML page
		//If it is a POST request, add user to database or send back error
		if req.Method == "GET" {
			res.Write([]byte(`
            <html>
            <head>
            <title> Register </title>
            </head>
            <body>
            <h1> Register </h1>
            <form action = "/register" method = "post">
  Username:<br>
  <input type="text" name="Username"><br>
  Password:<br>
  <input type = "password" name = "Password">
  <input type = "submit" value = "register">
</form>
</body>
</html>
`))
		} else {
			req.ParseForm()
			username := req.FormValue("Username")
			password := req.FormValue("Password")
			//fmt.Println(username)
			//fmt.Println(password)
			salt := make([]byte, SaltLength)
			rand.Read(salt)
			saltedhash := pbkdf2.Key([]byte(password), salt, KeyHashIterations, KeyHashLength, KeyHashAlgo)
			err := AddUsers(username, salt, saltedhash)

			if err == nil {
				res.WriteHeader(http.StatusOK)
			} else {
				res.Write([]byte(err.Error()))
			}

		}

	case "/library":
		uid := VerifyUser(req)
		if uid != 0 {

		} else {
			res.Write([]byte(`You are not logged in`))
		}
		// TODO: Send HTML for the user's library
		// Users can logout, upload book, or read book

	case "/logout":
		uid := VerifyUser(req)
		if uid != 0 {

		} else {
			res.Write([]byte(`You are not logged in`))
		}

		// TODO: Log user out and send back to home page

	case "/upload":
		uid := VerifyUser(req)
		if uid != 0 {
			if req.Method == "GET" {
				res.Write([]byte(`
          <html>
          <head>
          <title> Upload </title>
          </head>
          <body>
          <h1> Upload </h1>
          <body>
          <form action = "/upload" method = "post" enctype = "multipart/form-data">
          <p>
          Please specify a file, or a set of files:<br>
          <input type="file" name="datafile">
          </p>
          <div>
          <input type="submit" value="Send">
          </div>
          </form>
          </body>
          </html>
          `))
			} else {
				file, header, err := req.FormFile("datafile")
				if err != nil {
					res.Write([]byte("Upload failed"))
					break
				}
				res.Write([]byte(header.Filename))
				filename := header.Filename
				data, err := ioutil.ReadAll(file)
				err = UploadFile(filename, data, uid)
				if err != nil {
					fmt.Println(err)
				}

			}
		} else {
			res.Write([]byte(`You are not logged in`))
		}
		// TODO: Send HTML to the upload page if it is a GET request.
		// If it is a POST request, add the book to the database, etc.
		// Send user to Library.

	case "/read":
		uid := VerifyUser(req)
		if uid != 0 {
			name := req.URL.Query().Get("name")
			//res.Write([]byte(name))
			data, err := GetFile(uid, name)
			if err == nil {
				res.Write(data)
			} else {
				fmt.Println(err)
			}
		} else {
			res.Write([]byte(`You are not logged in`))
		}
		// TODO: Send HTML to read the book

	default:
		res.Write([]byte("<h1>Hello World</h1>"))
	}
}

func ReadBody(req *http.Request) (body []byte, err error) {
	body = make([]byte, req.ContentLength)
	read, err := req.Body.Read(body)

	if int64(read) == req.ContentLength {
		err = nil
	}

	return body, err
}
