package main

import (
	"fmt"
	"os"
	"zlang/repl"
)
import "os/user"

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(os.Stdout, "hello %s !\n", user.Name)
	fmt.Fprintf(os.Stdout, "feel free to type in command line\n")
	repl.Start(os.Stdin, os.Stdout)
}
