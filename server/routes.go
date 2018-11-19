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
    res.Write([]byte("<h1>Hello World</h1>"))
}

