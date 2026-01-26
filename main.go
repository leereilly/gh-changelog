package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const feedURL = "https://github.blog/changelog/feed/"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

func main() {
	pretty := flag.Bool("pretty", false, "Show full content with formatted body")
	flag.Parse()

	items, err := fetchFeed(feedURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching feed: %v\n", err)
		os.Exit(1)
	}

	// Check if "open" command with an ID is provided
	args := flag.Args()
	if len(args) > 0 && args[0] == "open" {
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: gh changelog open #<id>\n")
			os.Exit(1)
		}
		openItem(items, args[1])
		return
	}

	output := formatItems(items, *pretty)
	fmt.Print(output)
}

func fetchFeed(url string) ([]Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseFeed(body)
}

func parseFeed(data []byte) ([]Item, error) {
	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		return nil, err
	}

	items := rss.Channel.Items

	// Sort by date in reverse chronological order
	sort.Slice(items, func(i, j int) bool {
		ti, _ := parseDate(items[i].PubDate)
		tj, _ := parseDate(items[j].PubDate)
		return ti.After(tj)
	})

	return items, nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC1123Z, dateStr)
}

func formatDate(dateStr string) string {
	t, err := parseDate(dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02")
}

func formatRelativeDate(dateStr string) string {
	t, err := parseDate(dateStr)
	if err != nil {
		return dateStr
	}
	
	now := time.Now()
	duration := now.Sub(t)
	
	// Less than 60 minutes
	if duration < 60*time.Minute {
		return "Just now"
	}
	
	// Truncate to start of day for accurate day counting
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	days := int(nowDay.Sub(tDay).Hours() / 24)
	
	if days == 0 {
		return "Today"
	} else if days == 1 {
		return "1 day ago"
	} else {
		return fmt.Sprintf("%d days ago", days)
	}
}

func formatItems(items []Item, pretty bool) string {
	var sb strings.Builder

	if !pretty {
		// Add headers for default format with ANSI underline styling (like gh issue list)
		// Column order: ID (left), TITLE (middle), UPDATED (right)
		underline := "\033[4m"
		reset := "\033[0m"
		sb.WriteString(fmt.Sprintf("%s%-6s%s  %s%-90s%s  %s%s%s\n", underline, "ID", reset, underline, "TITLE", reset, underline, "UPDATED", reset))
	}

	for i, item := range items {
		date := formatDate(item.PubDate)

		if pretty {
			sb.WriteString(fmt.Sprintf("%s - %s\n", date, item.Title))
			sb.WriteString(strings.Repeat("-", 40) + "\n")
			content := item.Content
			if content == "" {
				content = item.Description
			}
			sb.WriteString(htmlToText(content))
			sb.WriteString("\n")
			if i < len(items)-1 {
				sb.WriteString("\n")
			}
		} else {
			relativeDate := formatRelativeDate(item.PubDate)
			// ANSI color code for green (same as gh issue list)
			green := "\033[32m"
			reset := "\033[0m"
			// Format ID without color for width calculation, then add color
			idStr := fmt.Sprintf("#%d", i)
			coloredID := fmt.Sprintf("%s%s%s", green, idStr, reset)
			// Pad the colored ID manually since ANSI codes affect string length
			padding := 6 - len(idStr)
			if padding > 0 {
				coloredID = coloredID + strings.Repeat(" ", padding)
			}
			// Truncate title if it exceeds 90 characters
			title := item.Title
			if len(title) > 90 {
				title = title[:87] + "..."
			}
			// Format: ID (colored, 6 chars), Title (90 chars), Updated
			sb.WriteString(fmt.Sprintf("%s  %-90s  %s\n", coloredID, title, relativeDate))
		}
	}

	return sb.String()
}

func htmlToText(html string) string {
	// Remove DOCTYPE and html/body wrapper
	doctype := regexp.MustCompile(`<!DOCTYPE[^>]*>`)
	html = doctype.ReplaceAllString(html, "")
	htmlTags := regexp.MustCompile(`</?html[^>]*>|</?body[^>]*>`)
	html = htmlTags.ReplaceAllString(html, "")

	// Remove video, img, and other media tags completely
	video := regexp.MustCompile(`<video[^>]*>[\s\S]*?</video>`)
	html = video.ReplaceAllString(html, "")
	img := regexp.MustCompile(`<img[^>]*>`)
	html = img.ReplaceAllString(html, "")

	// Remove "The post ... appeared first on ..." footer
	postFooter := regexp.MustCompile(`<p>The post.*?appeared first on.*?</p>`)
	html = postFooter.ReplaceAllString(html, "")

	// Convert headers to text with newlines
	headers := regexp.MustCompile(`<h[1-6][^>]*>(.*?)</h[1-6]>`)
	html = headers.ReplaceAllString(html, "\n$1\n")

	// Convert list items
	li := regexp.MustCompile(`<li[^>]*>(.*?)</li>`)
	html = li.ReplaceAllString(html, "• $1\n")

	// Remove ul/ol tags
	lists := regexp.MustCompile(`</?[uo]l[^>]*>`)
	html = lists.ReplaceAllString(html, "")

	// Convert paragraphs to double newlines
	pOpen := regexp.MustCompile(`<p[^>]*>`)
	html = pOpen.ReplaceAllString(html, "")
	pClose := regexp.MustCompile(`</p>`)
	html = pClose.ReplaceAllString(html, "\n\n")

	// Convert <br> to newlines
	br := regexp.MustCompile(`<br\s*/?>`)
	html = br.ReplaceAllString(html, "\n")

	// Extract link text from anchors
	links := regexp.MustCompile(`<a[^>]*>([^<]*)</a>`)
	html = links.ReplaceAllString(html, "$1")

	// Remove any remaining HTML tags
	allTags := regexp.MustCompile(`<[^>]+>`)
	html = allTags.ReplaceAllString(html, "")

	// Decode HTML entities
	html = decodeHTMLEntities(html)

	// Clean up whitespace
	multipleNewlines := regexp.MustCompile(`\n{3,}`)
	html = multipleNewlines.ReplaceAllString(html, "\n\n")

	multipleSpaces := regexp.MustCompile(` +`)
	html = multipleSpaces.ReplaceAllString(html, " ")

	// Trim lines
	lines := strings.Split(html, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	html = strings.Join(lines, "\n")

	return strings.TrimSpace(html)
}

func decodeHTMLEntities(s string) string {
	entities := map[string]string{
		"&amp;":    "&",
		"&lt;":     "<",
		"&gt;":     ">",
		"&quot;":   "\"",
		"&#39;":    "'",
		"&apos;":   "'",
		"&nbsp;":   " ",
		"&ndash;":  "-",
		"&mdash;":  "—",
		"&lsquo;":  "'",
		"&rsquo;":  "'",
		"&ldquo;":  "\"",
		"&rdquo;":  "\"",
		"&hellip;": "...",
		"&#8230;":  "...",
		"&#8217;":  "'",
		"&#8220;":  "\"",
		"&#8221;":  "\"",
	}

	for entity, char := range entities {
		s = strings.ReplaceAll(s, entity, char)
	}

	return s
}

func openItem(items []Item, idArg string) {
	// Parse the ID from the argument (e.g., "#0" or "0")
	idArg = strings.TrimPrefix(idArg, "#")
	id, err := strconv.Atoi(idArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", idArg)
		os.Exit(1)
	}

	// Check if ID is within bounds
	if id < 0 || id >= len(items) {
		fmt.Fprintf(os.Stderr, "ID #%d is out of range (0-%d)\n", id, len(items)-1)
		os.Exit(1)
	}

	// Get the item's link
	item := items[id]
	if item.Link == "" {
		fmt.Fprintf(os.Stderr, "No link available for item #%d\n", id)
		os.Exit(1)
	}

	// Open the link in the default browser
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", item.Link)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", item.Link)
	default: // Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", item.Link)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open browser: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Opening: %s\n", item.Title)
}
