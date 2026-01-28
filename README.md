# gator

`gator` is a local-first RSS blog aggregator CLI written in Go with PostgreSQL for persistence.
It supports multiple users, RSS feed management, following and unfollowing feeds, and browsing aggregated posts directly from the terminal.

Repository:
https://github.com/a-fleming/gator

---

## Requirements

To run `gator`, you will need:

- Go 1.22 or newer
- PostgreSQL (a locally running instance)

This project assumes PostgreSQL is available on `localhost:5432` unless configured otherwise.

**Note:**  
`go run .` is intended for development only. For normal usage, you should install and run the compiled `gator` binary.
Go programs build to statically compiled binaries, so once installed, `gator` can be run without the Go toolchain.

---

## Installation

Install the CLI using `go install`:

```bash
go install github.com/a-fleming/gator@latest
```

Make sure Go’s binary directory is on your PATH:

```bash
export PATH="$PATH:$HOME/go/bin"
```

Verify installation:

```bash
gator
```

---

## PostgreSQL Setup

You will need a PostgreSQL database and credentials.
`gator` connects using a standard PostgreSQL connection URL.

Example:

```text
postgres://postgres:postgres@localhost:5432/gator
```

To confirm PostgreSQL is running and reachable:

```bash
psql "postgres://postgres:postgres@localhost:5432/postgres"
```

If this connects successfully, your local Postgres instance is running.

---

## Configuration

`gator` reads its configuration from the following file:

```text
~/.gatorconfig.json
```

Only `db_url` is required initially.
The other fields are managed by the CLI and may be left blank.
They are automatically set during login and cleared during logout or reset.

Example configuration file contents:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator",
  "current_user_name": "",
  "current_user_id": ""
}
```

---

## Usage

Commands follow this structure:

```bash
gator <command> [arguments...]
```

Running `gator` without a command will result in an error.

---

## Commands

### Register a user

Creates a new user and sets them as the current user.

```bash
gator register <username>
```

Example:

```bash
gator register aaron
```

---

### Login

Logs in as an existing user and updates the config.

```bash
gator login <username>
```

---

### Logout

Clears the current user from the config file.

```bash
gator logout
```

---

### List users

Displays all users.
The currently logged-in user is marked.

```bash
gator users
```

---

### Add a feed (requires login)

Creates a new RSS feed and automatically follows it as the current user.

```bash
gator addfeed <feed-name> <feed-url>
```

Example:

```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

---

### List all feeds

Displays all feeds and the users who added them.

```bash
gator feeds
```

---

### Follow a feed (requires login)

Follows an existing feed by URL.

```bash
gator follow <feed-url>
```

---

### Unfollow a feed (requires login)

Unfollows a feed by URL.

```bash
gator unfollow <feed-url>
```

---

### List followed feeds (requires login)

Displays the feeds followed by the current user.

```bash
gator following
```

---

### Browse posts (requires login)

Displays recent posts for the current user.
If no limit is provided, a default value of 2 is used.

```bash
gator browse
gator browse 10
```

---

### Aggregate feeds continuously

Fetches RSS feeds on a repeating interval and stores new posts.
This command runs until manually stopped.

```bash
gator agg <interval>
```

The interval must be a valid Go duration string, as defined by Go's `time.ParseDuration`.

Supported units include:
- `ms` (milliseconds)
- `s`  (seconds)
- `m`  (minutes)
- `h`  (hours)

Examples:

```bash
gator agg 30s
gator agg 5m
gator agg 1h
gator agg 1m30s
```

Intervals without a unit (for example `30`) are not valid.

---

### Reset (development and testing)

Deletes all users and cascades deletes to related data.
Also clears the current user from the config file.

```bash
gator reset
```

---

## Development

Run in development mode:

```bash
go run .
```

Build a local binary:

```bash
go build
./gator <command> [arguments...]
```

---

## Notes

- Commands that require a logged-in user will fail if no user is set in `~/.gatorconfig.json`.
- RSS requests are made with a custom User-Agent and a request timeout to avoid hanging on slow feeds.
- This project is intended for local use and learning, but supports multiple users within the CLI.

---

## License

MIT
