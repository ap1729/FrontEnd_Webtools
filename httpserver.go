package main2

// Gist pulled from https://gist.github.com/wiless/b97637e1b5625248784d
import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

var servedir string
var portid string

func main() {
	// Simple static webserver:

	portid = ":8888"
	if len(os.Args) == 1 {
		servedir = os.Getenv("PWD")
	} else {
		servedir = os.Args[1]
	}

	if len(os.Args) > 2 {
		portid = ":" + os.Args[2]
	}

	adrs, err := net.InterfaceAddrs()
	fmt.Printf("\n Folder Listed : %s", servedir)
	for _, adr := range adrs {
		fmt.Printf("\n Open http://%v%s", strings.Split(adr.String(), "/")[0], portid)
	}

	err = http.ListenAndServe(portid, http.FileServer(http.Dir(servedir)))

	if err == nil {
		fmt.Printf("\n Folder Listed : %s", servedir)
		for _, adr := range adrs {
			fmt.Printf("\n Open http://%v%s", strings.Split(adr.String(), "/")[0], portid)
		}
	} else {
		fmt.Println("Error Starting Listen ", err)
	}

	fmt.Println()
}
