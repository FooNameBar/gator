package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	cmds "github.com/FooNameBar/gator/internal/commands"
	"github.com/FooNameBar/gator/internal/config"
	"github.com/FooNameBar/gator/internal/data"
	"github.com/FooNameBar/gator/internal/database"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Printf("Error: sql.Open url: %s, %v\n", c.DbUrl, err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	state := data.State{Config: &c, DB: dbQueries}
	commands := cmds.Commands{List: make(map[string]func(*data.State, cmds.Command) error)}

	commands.Register("login", cmds.HandlerLogin)
	commands.Register("register", cmds.HandlerRegister)
	commands.Register("reset", cmds.HandlerReset)
	commands.Register("users", cmds.GetUsers)
	commands.Register("agg", cmds.GetFeeds)
	commands.Register("addfeed", cmds.AddFeed)
	commands.Register("feeds", cmds.ListFeeds)

	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments. Need a command name and its argument")
		os.Exit(1)
	}

	name, cArgs := os.Args[1], os.Args[2:]
	err = commands.Run(&state, cmds.Command{Name: name, Args: cArgs})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
