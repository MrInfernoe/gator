Postgres Go
install gator CLI
Set up config and run
explain commands


gator

gator is a RSS feed aggregator for CLI.
Manage users and feeds by CLI commands.

Features
    User database for individual feed fetching.
    Post limiting.


Install
go install github.com/mrinfernoe/gator


Usage
gator <command> [<args>]

register    <username>      Record this username in the database
login       <username>      Set the current user to this username

reset                       Delete all users of the database and cascade delete through feeds and posts

addfeed     [<url>]         Record this url's feed in the database
follow      [<feed>]        Record the logged in user following the feed
unfollow    [<feed>]        Delete the record for logged in user following this feed

following                   Show all feed names that logged in user is following
browse      [<limit>]       Show 2 or more most recent posts from logged in user's followed feeds

users                       Show all registered users in the database
feeds                       Show all feeds in the database and the users that follow them

agg         <duration>      Records all posts for oldest fetched feed in a loop with a break of duration 
                            between fetches