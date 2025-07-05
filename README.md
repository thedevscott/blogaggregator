# Overview
The goal of this project is to integrate a Go application with a PostgreSQL
database using the [sqlc](https://sqlc.dev/) and [goose](https://github.com/pressly/goose) tools for typesafe database communication.

This app can be used as a long running service to continuously fetch posts from
RSS feeds and store them in the database.

## Features
  + Add RSS feeds from across the internet to be collected
  + Store the collected posts in a PostgreSQL database
  + Follow and unfollow RSS feeds that other users have added
  + View summaries of the aggregated posts in the terminal, with a link to the full
  post

# Requirements
Go 1.23+
Postgres 15+
  + Note: Usage on non-mac computers may need more database setup than contained
    here. Unfortunately, I do not have access to those systems at this time to
    provide step by step instructions.
    
# Install
This app can be installed by running:
```bash
go install https://github.com/thedevscott/blogaggregator
```

# Config File
Located in the home directory: ~/.gatorconfig.json
Contents:
```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

internal/config: config internal package used for reading and writing JSON file

# Commands
Before running any commands, ensure the config file is setup.

usage: blogaggregator <command> ...

+ "register"  - Adds a new user to the database ie: blogaggregator register <name>
+ "login"     - Sets the current user in the config ie: blogaggregator login <name>
+ "reset"     - Deletes all users in the database ie: blogaggregator reset
+ "users"     - Lists all the usres in the database ie: blogaggregator users
+ "addfeed"   - Adds a feed to tract ie: blogaggregator addfeed "Hacker News RSS" "https://hnrss.org/newest"
+ "agg"       - Runs the app indefinitely & collects feeds at set interval ie:
blogaggregator agg 60s
+ "feeds"     - List the feeds currently being tracked ie: blogaggregator feeds
+ "follow"    - Have the current user follow a registered feed ie: blogaggregator
follow "https://hnrss.org/newest"
+ "unfollow"  - Stop following a feed ie: blogaggregator unfollow "https://hnrss.org/newest"
+ "following" - List the feeds being followed by the current user ie:
blogaggregator following
+ "browse"    - Browse the posts. A feed must be followed in order to browse it. Takes an optional limit parameter ie:
blogaggregator browse <2>

# Postgres Database
1. [Install Postgres](https://www.postgresql.org/download/)
2. Verify install 
```bash
psql --version
```
3. Start the Postgres server [Docs](https://www.postgresql.org/docs/current/server-start.html)
* Mac: brew services start postgresql@15
* Linux: sudo service postgresql start

4. Connect to the DB
* Mac: psql postgres
* Linux: sudo -u postgres psql

5. Create the DB & connect
```bash
CREATE DATABASE <name>;

\c <name>

# set user password Linux
ALTER USER postgres PASSWORD 'postgres';
```

# Reminder before DB tests
Run goose up/down commands from inside the sql/schema directory
```bash
goose postgres <connection_string> up/down
```
