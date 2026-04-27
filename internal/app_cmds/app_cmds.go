package app_cmds

import (
	"strings"
	"os"
	"context"
	_ "bufio"
	"github.com/google/uuid"
	"time"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

// holds config state
type State struct{
	DbQPtr		*database.Queries
	ConfigPtr	*config.Config
}

// initialise state from file
func NewState() State {
	var newState State
	config := config.ReadConfig()
	newState.ConfigPtr = &config
	newState.DbQPtr = &database.Queries{}
	return newState
}

// holds a command
type Command struct {
	Name	string
	Args	[]string
}

type Commands struct {
	Registry map[string]func(*State, Command) error
}

func NewCommands() Commands {
	newCommands := Commands{}
	newCommands.Registry = map[string]func(*State, Command) error {}
	return newCommands
}

// runs given command with state if exists
func (c *Commands) Run(s *State, cmd Command) error { // state???
	err := c.Registry[cmd.Name](s, cmd)								// here
	if err != nil {
		return err
	}
	return nil
}

// registers handler function for command name
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Registry[name] = f
}

// checks for username in database then sets to current in config
func HandlerLogin(s *State, cmd Command) error {
	// expects 1 command argument: username
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: username is required")
	}

	username := cmd.Args[0]

	if _, err := s.DbQPtr.GetUser(context.Background(), username); err != nil {
		if fmt.Sprintf("%v", err) == "sql: no rows in result set" {
			fmt.Println("username not found")
			os.Exit(1)
		}
		return err
	}

	s.ConfigPtr.SetUser(username)
	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

// adds username to database
func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: no name given")
	}

	name := strings.Join(cmd.Args, " ")
	user, err := s.DbQPtr.GetUser(context.Background(), name)		// here
	if err != nil && fmt.Sprintf("%v", err) != fmt.Sprintf("sql: no rows in result set") {
		return fmt.Errorf(fmt.Sprintf("error from get: %v", err))
	}
	if user.Name == name {
		fmt.Println("error: username already exists")
		os.Exit(1)
	}

	ctx := context.Background()

	id := uuid.New()
	created := time.Now()
	updated := time.Now()
	params := database.CreateUserParams{id, created, updated, name}
	// fmt.Printf("%v, %v\n", ctx, params)


	// s is address of state, holds address of database.Queries create user acts on address of Queries
	createdUser, err := s.DbQPtr.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	s.ConfigPtr.SetUser(createdUser.Name)
	fmt.Printf("registered %s to database\n", createdUser.Name)
	return nil
}

// deletes all data from users database
func HandlerReset(s *State, cmd Command) error {
	
	// fmt.Println("Delete all records from \"users\" database? Y/n ")
	
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// if scanner.Text() != "Y" {
	// 	fmt.Println("Delete cancelled")
	// 	return nil
	// }
	ctx := context.Background()
	err := s.DbQPtr.ResetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("database \"users\" has been cleared ")
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	var users []database.User
	ctx := context.Background()
	users, err := s.DbQPtr.GetUsers(ctx)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		fmt.Println("no users found")
	}
	currentUser := s.ConfigPtr.Current_user_name
	for _, user := range users {
		fmt.Printf(user.Name)
		if user.Name == currentUser {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}
	return nil
}
