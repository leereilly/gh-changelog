package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseFeed(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
<channel>
	<title>GitHub Changelog</title>
	<item>
		<title>First Feature</title>
		<pubDate>Mon, 20 Jan 2026 10:00:00 +0000</pubDate>
		<description>First description</description>
		<content:encoded><![CDATA[<p>First content paragraph.</p>]]></content:encoded>
	</item>
	<item>
		<title>Second Feature</title>
		<pubDate>Tue, 21 Jan 2026 12:00:00 +0000</pubDate>
		<description>Second description</description>
		<content:encoded><![CDATA[<p>Second content paragraph.</p>]]></content:encoded>
	</item>
</channel>
</rss>`

	items, err := parseFeed([]byte(xmlData))
	if err != nil {
		t.Fatalf("Failed to parse feed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// Should be sorted reverse chronologically (newest first)
	if items[0].Title != "Second Feature" {
		t.Errorf("Expected first item to be 'Second Feature', got '%s'", items[0].Title)
	}

	if items[1].Title != "First Feature" {
		t.Errorf("Expected second item to be 'First Feature', got '%s'", items[1].Title)
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Thu, 22 Jan 2026 20:45:09 +0000", "2026-01-22"},
		{"Mon, 01 Dec 2025 08:30:00 +0000", "2025-12-01"},
	}

	for _, tt := range tests {
		result := formatDate(tt.input)
		if result != tt.expected {
			t.Errorf("formatDate(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatRelativeDate(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name        string
		timeAgo     time.Duration
		expected    string
	}{
		{"just now - 30 minutes ago", 30 * time.Minute, "Just now"},
		{"just now - 59 minutes ago", 59 * time.Minute, "Just now"},
		{"today - 2 hours ago", 2 * time.Hour, "Today"},
		{"today - 10 hours ago", 10 * time.Hour, "Today"},
		{"yesterday", 25 * time.Hour, "1 day ago"},
		{"two days ago", 2 * 24 * time.Hour, "2 days ago"},
		{"seven days ago", 7 * 24 * time.Hour, "7 days ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDate := now.Add(-tt.timeAgo)
			dateStr := testDate.Format(time.RFC1123Z)
			
			result := formatRelativeDate(dateStr)
			if result != tt.expected {
				t.Errorf("formatRelativeDate(%q) = %q, want %q", dateStr, result, tt.expected)
			}
		})
	}
}

func TestFormatItemsDefault(t *testing.T) {
	items := []Item{
		{Title: "Feature A", PubDate: "Tue, 21 Jan 2026 12:00:00 +0000"},
		{Title: "Feature B", PubDate: "Mon, 20 Jan 2026 10:00:00 +0000"},
	}

	output := formatItems(items, false)

	// Check for headers
	if !strings.Contains(output, "ID  DATE            TITLE") {
		t.Errorf("Expected output to contain header 'ID  DATE            TITLE', got:\n%s", output)
	}

	if !strings.Contains(output, "--  ----            -----") {
		t.Errorf("Expected output to contain header underline, got:\n%s", output)
	}

	// Check for ID column (0 for first item, 1 for second item)
	if !strings.Contains(output, "0   ") {
		t.Errorf("Expected output to contain ID '0', got:\n%s", output)
	}

	if !strings.Contains(output, "1   ") {
		t.Errorf("Expected output to contain ID '1', got:\n%s", output)
	}

	// Check for feature titles
	if !strings.Contains(output, "Feature A") {
		t.Errorf("Expected output to contain 'Feature A', got:\n%s", output)
	}

	if !strings.Contains(output, "Feature B") {
		t.Errorf("Expected output to contain 'Feature B', got:\n%s", output)
	}
}

func TestFormatItemsPretty(t *testing.T) {
	items := []Item{
		{
			Title:   "Feature A",
			PubDate: "Tue, 21 Jan 2026 12:00:00 +0000",
			Content: "<p>This is the content.</p><p>Second paragraph.</p>",
		},
	}

	output := formatItems(items, true)

	if !strings.Contains(output, "2026-01-21 - Feature A") {
		t.Errorf("Expected output to contain '2026-01-21 - Feature A', got:\n%s", output)
	}

	if !strings.Contains(output, "This is the content.") {
		t.Errorf("Expected output to contain body text, got:\n%s", output)
	}

	if strings.Contains(output, "<p>") {
		t.Errorf("Expected no HTML tags in output, got:\n%s", output)
	}
}

func TestHtmlToText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple paragraph",
			input:    "<p>Hello world</p>",
			expected: "Hello world",
		},
		{
			name:     "multiple paragraphs",
			input:    "<p>First paragraph.</p><p>Second paragraph.</p>",
			expected: "First paragraph.\n\nSecond paragraph.",
		},
		{
			name:     "removes images",
			input:    "<p>Text before</p><img src='test.jpg'/><p>Text after</p>",
			expected: "Text before\n\nText after",
		},
		{
			name:     "removes video tags",
			input:    "<p>Hello</p><video src='test.mp4'></video><p>World</p>",
			expected: "Hello\n\nWorld",
		},
		{
			name:     "converts headers",
			input:    "<h2>Header</h2><p>Content</p>",
			expected: "Header\nContent",
		},
		{
			name:     "converts lists",
			input:    "<ul><li>Item one</li><li>Item two</li></ul>",
			expected: "• Item one\n• Item two",
		},
		{
			name:     "extracts link text",
			input:    "<p>Check out <a href='https://example.com'>this link</a> here.</p>",
			expected: "Check out this link here.",
		},
		{
			name:     "decodes HTML entities",
			input:    "<p>It&rsquo;s a &ldquo;test&rdquo; &amp; more</p>",
			expected: "It's a \"test\" & more",
		},
		{
			name:     "removes post footer",
			input:    "<p>Content</p><p>The post Example appeared first on The GitHub Blog.</p>",
			expected: "Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := htmlToText(tt.input)
			if result != tt.expected {
				t.Errorf("htmlToText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	dateStr := "Thu, 22 Jan 2026 20:45:09 +0000"
	parsed, err := parseDate(dateStr)
	if err != nil {
		t.Fatalf("parseDate failed: %v", err)
	}

	if parsed.Year() != 2026 || parsed.Month() != 1 || parsed.Day() != 22 {
		t.Errorf("parseDate returned wrong date: %v", parsed)
	}
}

func TestReverseChronologicalOrder(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
<channel>
	<item>
		<title>Oldest</title>
		<pubDate>Mon, 19 Jan 2026 10:00:00 +0000</pubDate>
	</item>
	<item>
		<title>Middle</title>
		<pubDate>Tue, 20 Jan 2026 10:00:00 +0000</pubDate>
	</item>
	<item>
		<title>Newest</title>
		<pubDate>Wed, 21 Jan 2026 10:00:00 +0000</pubDate>
	</item>
</channel>
</rss>`

	items, err := parseFeed([]byte(xmlData))
	if err != nil {
		t.Fatalf("Failed to parse feed: %v", err)
	}

	expected := []string{"Newest", "Middle", "Oldest"}
	for i, item := range items {
		if item.Title != expected[i] {
			t.Errorf("Item %d: expected %q, got %q", i, expected[i], item.Title)
		}
	}
}
