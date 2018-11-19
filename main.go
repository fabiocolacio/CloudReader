	
package main

import(
    "fmt"
    "flag"
    "github.com/fabiocolacio/cloudreader/server"
)

var(
    flagInit bool
    flagReset bool
)

// cloudreader --init

func main() {
    flag.BoolVar(&flagInit, "init", false, "Initializes databse tables.")
    flag.BoolVar(&flagReset, "reset", false, "Resets database tables.")
    flag.Parse()

    if flagInit {
        fmt.Println("Creating tables!")
        server.InitTables()
    }

    if flagReset {
        fmt.Println("Reseting tables!")
        server.ResetTables()
    }
}


