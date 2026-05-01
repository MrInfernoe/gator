package app_cmds

import (
	// "strings"
	// "os"
	"database/sql"
	"context"
	_ "bufio"
	"github.com/google/uuid"
	"time"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/feed"
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

func ScrapeFeeds(s *State) error {
	// fmt.Println("DEBUG: starting scrape")
	ctx := context.Background()
	nextFeed, err := s.DbQPtr.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	fetchedParams := database.MarkFeedFetchedParams{nextFeed.Url, time.Now()}
	err = s.DbQPtr.MarkFeedFetched(ctx, fetchedParams)
	if err != nil {
		return err
	}

	fetchedRSSFeed, err := feed.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for _, post := range fetchedRSSFeed.Channel.Items {
	ID := uuid.New()
	CreatedAt := time.Now()
	UpdatedAt := time.Now()
	Title := sql.NullString{post.Title, true}
	Url := post.Link
	Description := sql.NullString{post.Description, true}
	PublishedAt, err := time.Parse(time.RFC1123Z, post.PubDate)
	if err != nil {
		return err
	}
	FeedID := nextFeed.ID
	postParams := database.CreatePostParams{ID, CreatedAt, UpdatedAt, Title, Url, Description, PublishedAt, FeedID}
	createdPost, err := s.DbQPtr.CreatePost(ctx, postParams)
	if fmt.Sprintf("%v", err) == "pq: duplicate key value violates unique constraint \"posts_url_key\" (23505)" {
		continue
	}
	if err != nil {
		return err
	}
	fmt.Printf("From feed \"%v\" added post: \"%v\"\n", nextFeed.Name, createdPost.Title.String)
	}

	// for _, item := range fetchedRSSFeed.Channel.Items {
	// 	fmt.Printf("%v\n", item.Title)
	// }

	return nil
}