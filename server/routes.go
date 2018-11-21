package server

import(
    "net/http"
)

var(
    httpServer *http.Server
    httpsServer *http.Server
    macKey [256]byte
)

func init() {
    httpServer = &http.Server{
        Addr: HttpAddr,
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

    case "/login":
        // TODO: If it is a GET request, send login HTML page
        // If it is a POST request, log them in or tell them if there is
        // na error

    case "/register":
        // TODO: If it is a GET request, send the register HTML page
        //If it is a POST request, add user to database or send back error

    case "/library":
        // TODO: Send HTML for the user's library
        // Users can logout, upload book, or read book

    case "/logout":
        // TODO: Log user out and send back to home page

    case "/upload":
        // TODO: Send HTML to the upload page if it is a GET request.
        // If it is a POST request, add the book to the database, etc.
        // Send user to Library.

    case "/read":
        // TODO: Send HTML to read the book

    default:
        res.Write([]byte("<h1>Hello World</h1>"))
    }
}

func ReadBody(req *http.Request) (body []byte, err error) {
    body = make([]byte, req.ContentLength) 
    read, err := req.Body.Read(body);

    if int64(read) == req.ContentLength {
        err = nil
    }

    return body, err
}

