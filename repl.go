package main

import (
	"fmt"
	"github.com/abusizhishen/zlang/repl"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(os.Stdout, "hello %s !\n", u.Name)
	fmt.Fprintf(os.Stdout, "feel free to type in command line\n")
	repl.Start(os.Stdin, os.Stdout)
}
