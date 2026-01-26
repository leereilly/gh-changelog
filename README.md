# gh-changelog

A GitHub CLI extension to view the latest GitHub Changelog entries directly in your terminal.

## Installation

```sh
gh extension install leereilly/gh-changelog
```

## Usage

```sh
# List recent changelog entries (date and title only)
gh changelog

# Show full content with formatted body
gh changelog --pretty

# Open a specific changelog entry in your browser
gh changelog open #0
gh changelog open 0
```

### Example Output

Default output:
```
ID      TITLE                                                                                       UPDATED
#0      New feature announcement                                                                    Today
#1      Security update released                                                                    1 day ago
#2      API improvements                                                                            2 days ago
```

With `--pretty`:
```
2026-01-22 - New feature announcement
----------------------------------------
Full description of the feature with
formatted content and bullet points.

2026-01-21 - Security update released
----------------------------------------
Details about the security update...
```

Opening a changelog entry:
```sh
$ gh changelog open #0
Opening: New feature announcement
# Opens https://github.blog/changelog/... in your default browser
```

## Development

### Requirements

- Go 1.21+

### Build

```sh
go build -o gh-changelog
```

### Test

```sh
go test -v
```

## How It Works

This extension fetches the RSS feed from https://github.blog/changelog/feed/ and displays entries in reverse chronological order (newest first). The `--pretty` flag includes the full content with HTML converted to readable plain text.

## License

MIT
