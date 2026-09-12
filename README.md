# gator is a command line RSS feed tool

## Requirements
- `Postgres` and `go` to run

## Install
- `go install` in the project directory to build and install the binary
- In the root directory, create a config file `.gatorconfig.json` with the following format:
```
{
  "db_url": "postgres://example?sslmode=disable",
  "current_user_name": "username"
}
```

## Commands
- `register` takes an argument and creates a new user from it
- `login` logs in as a specific user
- `users` list all the users that have been registered
- `addfeed` provide a url of an rss feed to scrape later
- `feeds` list the saved rss feeds
- `follow` as a user, follow the given url of an rss feed
- `following` list the feeds for the current user
- `unfollow` remove a followed feed for the current user
- `browse` list the posts scraped from the rss feeds. Defaults to 2 unless limit provided
- `agg` scrape all saved feeds and save the posts for browsing, provide an interval time between scraping in the formate `30s`, `30m`, `1h`
