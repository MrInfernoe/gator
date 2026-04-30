package app_cmds

import (
	"fmt"
	"context"
	"os"
	"strings"
	"github.com/google/uuid"
	"time"
	"gator/internal/database"
	"gator/internal/feed"
)

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
	user, err := s.DbQPtr.GetUser(context.Background(), name)
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
	fmt.Println("database has been cleared")
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

func HandlerAgg(s *State, cmd Command) error {
	// if len(cmd.Args) < 1 {
	// 	return fmt.Errorf("URL required")
	// }
	// if len(cmd.Args) > 1 {
	// 	return fmt.Errorf("too many arguments")
	// }
	ctx := context.Background()
	// feedURL := cmd.Args[0]
	feedURL := "https://www.wagslane.dev/index.xml"
	rssfeed, err := feed.FetchFeed(ctx, feedURL)
	if err != nil {
		return err
	}

	fmt.Println(rssfeed)

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	// get current user and connect feed to user
	if len(cmd.Args) < 2 {
		return fmt.Errorf("feed name and url required")
	}
	
	ctx := context.Background()
	id := uuid.New()
	created := time.Now()
	updated := time.Now()
	feed_name := cmd.Args[0]
	url := cmd.Args[1]
	params := database.CreateFeedParams{id, created, updated, feed_name, url, user.ID}

	createdFeed, err := s.DbQPtr.CreateFeed(ctx, params)
	if err != nil {
		return err
	}
	fmt.Printf("added to feeds:\n%v\n", createdFeed)
	newCmd := Command{Name: "following", Args: []string{url}}
	err = HandlerFollow(s, newCmd, user)
	if err != nil {
		return err
	}

	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("too many arguments")
	}
	// print all db feeds: name, url, user
	ctx := context.Background()
	feedsList, err := s.DbQPtr.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feedInfo := range feedsList {
		fmt.Printf("%v %v %v\n", feedInfo.FeedName, feedInfo.Url, feedInfo.UserName.String)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough arguments")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	url := cmd.Args[0]
	ctx := context.Background()
	feed, err := s.DbQPtr.GetFeed(ctx, url)
	if err != nil {
		return err
	}
	// user, err := s.DbQPtr.GetUser(ctx, s.ConfigPtr.Current_user_name)
	// if err != nil {
	// 	return err
	// }


	id := uuid.New()
	created := time.Now()
	updated := time.Now()
	// user_id := user.ID
	feed_id := feed.ID
	params := database.CreateFeedFollowParams{id, created, updated, user.ID, feed_id}


	// s is address of state, holds address of database.Queries create user acts on address of Queries
	createdFollow, err := s.DbQPtr.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}
	fmt.Printf("%v followed %v\n", createdFollow.UserName, createdFollow.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("too many arguments")
	}
	// user_name := s.ConfigPtr.Current_user_name
	ctx := context.Background()
	followedFeeds, err := s.DbQPtr.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return err
	}
	if len(followedFeeds) == 0 {
		fmt.Printf("user %v is not following any feeds\n", user.Name)
		return nil
	}
	fmt.Printf("user %v is following these feeds:\n", user.Name)
	for _, followedFeed := range followedFeeds {
		fmt.Println(followedFeed.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing url")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	}
	feed_url := cmd.Args[0]
	user_id := user.ID

	ctx := context.Background()
	feed, err := s.DbQPtr.GetFeed(ctx, feed_url)
	if err != nil {
		return err
	}
	feed_id := feed.ID

	deleteParams := database.DeleteFeedFollowForUserFeedParams{user_id, feed_id}
	err = s.DbQPtr.DeleteFeedFollowForUserFeed(ctx, deleteParams)
	if err != nil {
		return err
	}
	return nil
}