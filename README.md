# gh-changelog

A GitHub CLI extension to view the latest [GitHub Changelog](https://github.blog/changelog/) entries directly in your terminal.

Built using [GitHub Copilot CLI](https://github.com/features/copilot/cli) for the [GitHub Copilot CLI Challenge](https://dev.to/leereilly/github-changelog-now-in-your-terminal-built-with-copilot-cli-4fn7).

![Demo GIF](images/demo.gif)

## Installation

```
gh extension install leereilly/gh-changelog
```

## Usage

### View recent changelogs

```
gh changelog
```

Example output:

```shell
ID      TITLE                                                                                       UPDATED
#0      GPT-9000-Codex is now available in Visual Studio, JetBrains IDEs, and Emacs                 Today
#1      GitHub Copilot Chat learned sarcasm (enable with COPILOT_SASS=true)                         2 days ago
#2      CodeQL 13.37 as been released                                                               3 days ago
#3      1 vCPU Linux runner now generally available in GitHub Actions                               4 days ago
#4      GitHub Copilot CLI: Plan before you build, steer as you go                                  5 days ago
#5      Install and Use GitHub Copilot CLI directly from the GitHub CLI                             5 days ago
#6      CodeQL 2.23.9 has been released                                                             6 days ago
#7      Strengthen your supply chain with code-to-cloud traceability and SLSA Build Level 3 sec...  6 days ago
#8      Enterprise-scoped budgets that exclude cost center usage in public preview                  6 days ago
#9      GitHub Copilot now supports OpenCode                                                        10 days ago
```

The `--pretty` flag includes the full content with HTML converted to readable plain text. It will be even prettier when someone tackles [#1](https://github.com/leereilly/gh-changelog/issues/1). Anyone?

### View specific changelog

```
gh changelog view 0
```

Example output:

```shell
GPT-9000-Codex is now available in Visual Studio, JetBrains IDEs, and Emacs
---------------------------------------------------------------------------

Today we're excited to announce that GPT-9000-Codex, our most advanced 
AI coding model, is now available across all major IDEs.

What's new:
• 10x faster code completions with 99.7% accuracy
• Full codebase understanding up to 2 million tokens
• Native support for 150+ programming languages
• Real-time pair programming with voice commands

Get started today by updating your GitHub Copilot extension.
```

### Open specific changelog in browser

```
gh changelog open 0
```

## Development

### Requirements

- Go 1.21+

### Build

```
go build -o gh-changelog
```

### Test

```
go test -v
```

## How It Works

This extension fetches the RSS feed from https://github.blog/changelog/feed/ and displays entries in reverse chronological order (newest first).

## License

MIT

