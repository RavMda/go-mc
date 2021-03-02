package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RavMda/go-mc/yggdrasil"
)

var user = flag.String("user", "", "Can be an email address or player name for unmigrated accounts")
var pswd = flag.String("password", "", "Your password")

func main() {
	flag.Parse()

	resp, err := yggdrasil.Authenticate(*user, *pswd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	id, name := resp.SelectedProfile()
	fmt.Println("user:", name)
	fmt.Println("uuid:", id)
	fmt.Println("astk:", resp.AccessToken())
}
