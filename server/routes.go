package server

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"net/http"
	"strconv"
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
	case "/": HomeRoute(res, req)
	case "/test": TestRoute(res, req)
	case "/login": LoginRoute(res, req)
	case "/register": RegisterRoute(res, req)
    case "/logout": LogoutRoute(res, req)
	case "/library": LibraryRoute(res, req)
	case "/upload": UploadRoute(res, req)
	case "/read": ReadRoute(res, req)
	default: NotFoundRoute(res, req)
	}
}

func HomeRoute(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte(`<h1>CloudReader</h1>
		<body>
		<a href="/login">Login</a>
		<a href="/register">Register</a>
		</body>
		`))
}

func NotFoundRoute(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("Oopsie woopsie this doesn't exist."))
}

func TestRoute(res http.ResponseWriter, req *http.Request) {
    uid := VerifyUser(req)
    if uid == 0 {
        res.Write([]byte(`Wrong Login Information`))
    } else {
        res.Write([]byte(`You have logged in`))
    }
}

func LoginRoute(res http.ResponseWriter, req *http.Request) {
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
            http.SetCookie(res, &session)
            Redirect("/library", res)
        } else {
            res.Write([]byte(err.Error()))
        }
    }
}

func RegisterRoute(res http.ResponseWriter, req *http.Request) {
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

        salt := make([]byte, SaltLength)
        rand.Read(salt)
        saltedhash := pbkdf2.Key([]byte(password), salt, KeyHashIterations, KeyHashLength, KeyHashAlgo)
        err := AddUsers(username, salt, saltedhash)

        if err != nil {
            res.Write([]byte(err.Error()))
        }

        Redirect("/login", res)
    }
}

func LogoutRoute(res http.ResponseWriter, req *http.Request) {
    session := http.Cookie{
        Name:  "session",
        Value: strconv.Itoa(-1),

        //MaxAge: 10 * 60,
        Secure:   false,
        HttpOnly: true,
        SameSite: 1,

        Path: "/",
    }
    http.SetCookie(res, &session)
    Redirect("/login", res)
}

func LibraryRoute(res http.ResponseWriter, req *http.Request) {
    uid := VerifyUser(req)
    if uid != 0 {
        bookString, err := ShowBooks(uid)
        if err == nil{
          res.Write([]byte(`
            <html>
            <head>
            <title> Library </title>
            </head>
            <body>
            <h1> Library </h1>
            <a href="/upload">Upload Book</a>
            <a href="/logout">Logout</a>
            <table>
            `))

            for i := 0; i < len(bookString); i++ {
                fmt.Fprintf(res,`<tr>
                    <td><a href="/read?name=%s">%s</a></td>
                    <td><a href="/read?name=%s" download="%s">download</a></td></tr>`,
                    bookString[i], bookString[i], bookString[i], bookString[i])
            }
            res.Write([]byte(`</table></body></html>`))
        } else {
            res.Write([]byte(`You have no book.`))
        }
    } else {
        res.Write([]byte(`You are not logged in`))
    }
}

func UploadRoute(res http.ResponseWriter, req *http.Request) {
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
                return
            }
            //res.Write([]byte(header.Filename))
            filename := header.Filename
            data, err := ioutil.ReadAll(file)

            err = UploadFile(filename, data, uid)
            if err != nil {
                fmt.Println(err)
            }
						Redirect("/library", res)
        }
    } else {
        res.Write([]byte(`You are not logged in`))
    }
}
func ReadRoute(res http.ResponseWriter, req *http.Request) {
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
}

func ReadBody(req *http.Request) (body []byte, err error) {
	body = make([]byte, req.ContentLength)
	read, err := req.Body.Read(body)

	if int64(read) == req.ContentLength {
		err = nil
	}

	return body, err
}

func Redirect(loc string, res http.ResponseWriter) {
    res.Header().Set("Location", loc)
    res.WriteHeader(302)
}
