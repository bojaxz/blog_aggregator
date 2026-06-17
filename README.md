# gator 

Gator is a CLI tool for aggregating blogs! Keep up to date with your favorite blog site straight from your terminal!

## Prerequisites

 - Have Go Installed
 - Have Postgres Installed

## Installation

Make sure you have Go and Postgres installed first

Then install the CLI tool with
```Bash
go install github.com/bojaxz/blog_aggregator@latest
```

## Configuration

Create a config file in your home directory name .gatorconfig.json

you'll need a db_url and a current_user_name

example config would look something like this:

`{
    "db_url":"<connection_string>?sslmode=disable",
    "current_user_name":"<user_name>"
}`

## Usage
In order to get started run `gator register <user_name>`

Note: Use `go run . <command>` during development
Use `gator <command>` after installing the binary

This will create a user with the name you provide and set the current_user_name in the config to that user

To add a blog, run `gator addfeed <blog_name> <blog_url>`

Then aggregate the blog with `gator agg <duration>` with duration being a number and unit
example: `gator agg 30s` to aggregate every 30 seconds

From there you're all set to browse your favorite blogs right from your terminal's command line with
`gator browse <limit>`
Limit will determine the number of results, with the default being 2.

## Notes
follow your favorite feeds with `gator follow <feed url>`
unfollow a feed with `gator unfollow <feed url>`

Please enjoy! and don't spam!!

