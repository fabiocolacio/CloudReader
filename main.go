package main

//main module contain the main function of the application.
import (
	"flag"
	"fmt"
	"github.com/fabiocolacio/cloudreader/server"
)

//variable declarations
var (
	flagInit  bool
	flagReset bool
)


//main function
func main() {
	flag.BoolVar(&flagInit, "init", false, "Initializes databse tables.")
	flag.BoolVar(&flagReset, "reset", false, "Resets database tables.")
	flag.Parse()

	//initialize the tables in the database
	if flagInit {
		fmt.Println("Creating tables!")
		server.InitTables()
	}

	//reset the tables in the database
	if flagReset {
		fmt.Println("Reseting tables!")
		server.ResetTables()
	}

	server.ListenAndServe()
}
