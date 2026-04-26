package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"os"
	"database/sql"
	// "gator/internal/config"
	"gator/internal/app_cmds"
	"gator/internal/database"
)

func main() {
	configState := app_cmds.NewState()
	// fmt.Println(configState.ConfigPtr)

	commands := app_cmds.NewCommands()
	// fmt.Println(commands)
	commands.Register("login", app_cmds.HandlerLogin)
	// fmt.Println(commands)
	commands.Register("register", app_cmds.HandlerRegister)

	inputArgs := os.Args
	if len(inputArgs) < 2 {
		fmt.Println("error: too few arguments")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", configState.ConfigPtr.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	configState.DbQPtr = dbQueries

	startupCommand := app_cmds.Command{inputArgs[1], inputArgs[2:]}
	err = commands.Run(&configState, startupCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println(configState.ConfigPtr.Db_url)
}