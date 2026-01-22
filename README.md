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
```

### Example Output

Default output:
```
2026-01-22  New feature announcement
2026-01-21  Security update released
2026-01-20  API improvements
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
