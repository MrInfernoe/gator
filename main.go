package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"os"
	"database/sql"
	// "context"
	// "gator/internal/config"
	"gator/internal/app_cmds"
	"gator/internal/database"
	// "gator/internal/feed"
)

func main() {
	configState := app_cmds.NewState()
	// fmt.Println(configState.ConfigPtr)

	commands := app_cmds.NewCommands()
	// fmt.Println(commands)
	commands.Register("login", app_cmds.HandlerLogin)
	commands.Register("register", app_cmds.HandlerRegister)
	commands.Register("reset", app_cmds.HandlerReset)
	commands.Register("users", app_cmds.HandlerGetUsers)
	commands.Register("agg", app_cmds.HandlerAgg)
	commands.Register("addfeed", app_cmds.MiddlewareLoggedIn(app_cmds.HandlerAddFeed))
	commands.Register("feeds", app_cmds.HandlerGetFeeds)
	commands.Register("follow", app_cmds.MiddlewareLoggedIn(app_cmds.HandlerFollow))
	commands.Register("following", app_cmds.MiddlewareLoggedIn(app_cmds.HandlerFollowing))
	commands.Register("unfollow", app_cmds.MiddlewareLoggedIn(app_cmds.HandlerUnfollow))
	commands.Register("browse", app_cmds.MiddlewareLoggedIn(app_cmds.HandlerBrowse))
	// commands.Register("scrapeTest", app_cmds.ScrapeFeeds)
	// fmt.Println(commands)

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

	// fmt.Println("feed test")
	// feed, err := feed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// }
	// fmt.Printf("feed:\n%v\n", feed)
}