package app_cmds

import (
	// "strings"
	// "os"
	// "context"
	_ "bufio"
	// "github.com/google/uuid"
	// "time"
	// "fmt"
	"gator/internal/config"
	"gator/internal/database"
	// "gator/internal/feed"
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
func (c *Commands) Run(s *State, cmd Command) error {
	err := c.Registry[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

// registers handler function for command name
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Registry[name] = f
}

