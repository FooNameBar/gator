package main

import (
	"fmt"
	"os"

	cmds "github.com/FooNameBar/gator/internal/commands"
	"github.com/FooNameBar/gator/internal/config"
	"github.com/FooNameBar/gator/internal/data"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	state := data.State{Config: &c}
	commands := cmds.Commands{List: make(map[string]func(*data.State, cmds.Command) error)}
	commands.Register("login", cmds.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments. Need a command name and its argument")
		os.Exit(1)
	}

	name, cArgs := os.Args[1], os.Args[2:]
	err = commands.Run(&state, cmds.Command{Name: name, Args: cArgs})
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
