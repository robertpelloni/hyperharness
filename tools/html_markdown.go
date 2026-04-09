// Package tools provides HTML-to-Markdown conversion for WebFetch.
// Ported from superai/opencode/packages/opencode/src/tool/webfetch.ts and
// superai/crush/internal/agent/tools/fetch.go
//
// Converts HTML content to clean Markdown suitable for LLM consumption.
// Handles headings, lists, code blocks, links, emphasis, and strips
// script/style/meta/link tags.
package tools

import (
	"fmt"
	"regexp"
	"strings"
)

// HTMLToMarkdown converts HTML content to Markdown format.
// This is a native Go implementation that replaces the TurndownService
// dependency from the TypeScript version.
func HTMLToMarkdown(html string) string {
	// Strip script, style, meta, link tags and their content
	html = stripTagAndContent(html, "script")
	html = stripTagAndContent(html, "style")
	html = stripTagAndContent(html, "meta")
	html = stripTagAndContent(html, "link")
	html = stripTagAndContent(html, "noscript")
	html = stripTagAndContent(html, "head")

	// Convert HTML elements to Markdown
	html = convertHeadings(html)
	html = convertLists(html)
	html = convertCodeBlocks(html)
	html = convertInlineCode(html)
	html = convertStrong(html)
	html = convertEmphasis(html)
	html = convertLinks(html)
	html = convertImages(html)
	html = convertParagraphs(html)
	html = convertLineBreaks(html)
	html = convertHorizontalRules(html)
	html = convertBlockquotes(html)

	// Clean up remaining tags
	html = stripAllTags(html)

	// Clean up excessive whitespace
	html = cleanWhitespace(html)

	return strings.TrimSpace(html)
}

// ExtractTextFromHTML extracts plain text from HTML, stripping all tags.
func ExtractTextFromHTML(html string) string {
	// Strip script and style
	html = stripTagAndContent(html, "script")
	html = stripTagAndContent(html, "style")
	html = stripTagAndContent(html, "noscript")

	// Replace block elements with newlines
	blockTags := []string{"p", "div", "br", "h1", "h2", "h3", "h4", "h5", "h6",
		"li", "tr", "hr", "blockquote", "pre"}
	for _, tag := range blockTags {
		re := regexp.MustCompile(`(?i)</?` + tag + `[^>]*>`)
		html = re.ReplaceAllString(html, "\n")
	}

	// Strip all remaining tags
	html = stripAllTags(html)

	// Decode common HTML entities
	html = decodeEntities(html)

	// Clean whitespace
	html = cleanWhitespace(html)

	return strings.TrimSpace(html)
}

// stripTagAndContent removes an HTML tag and everything between opening and closing tags.
func stripTagAndContent(html, tag string) string {
	re := regexp.MustCompile(`(?is)<` + tag + `[^>]*>.*?</` + tag + `>`)
	return re.ReplaceAllString(html, "")
}

// convertHeadings converts h1-h6 tags to Markdown headings.
func convertHeadings(html string) string {
	for i := 1; i <= 6; i++ {
		prefix := strings.Repeat("#", i) + " "
		re := regexp.MustCompile(fmt.Sprintf(`(?i)<h%d[^>]*>(.*?)</h%d>`, i, i))
		html = re.ReplaceAllString(html, "\n"+prefix+"$1\n")
	}
	return html
}

// convertLists converts ul/ol/li to Markdown lists.
func convertLists(html string) string {
	// Unordered lists
	ulRe := regexp.MustCompile(`(?is)<ul[^>]*>(.*?)</ul>`)
	html = ulRe.ReplaceAllStringFunc(html, func(match string) string {
		liRe := regexp.MustCompile(`(?is)<li[^>]*>(.*?)</li>`)
		items := liRe.FindAllStringSubmatch(match, -1)
		var lines []string
		for _, item := range items {
			lines = append(lines, "- "+strings.TrimSpace(item[1]))
		}
		return "\n" + strings.Join(lines, "\n") + "\n"
	})

	// Ordered lists
	olRe := regexp.MustCompile(`(?is)<ol[^>]*>(.*?)</ol>`)
	html = olRe.ReplaceAllStringFunc(html, func(match string) string {
		liRe := regexp.MustCompile(`(?is)<li[^>]*>(.*?)</li>`)
		items := liRe.FindAllStringSubmatch(match, -1)
		var lines []string
		for i, item := range items {
			lines = append(lines, fmt.Sprintf("%d. %s", i+1, strings.TrimSpace(item[1])))
		}
		return "\n" + strings.Join(lines, "\n") + "\n"
	})

	return html
}

// convertCodeBlocks converts pre/code tags to fenced code blocks.
func convertCodeBlocks(html string) string {
	re := regexp.MustCompile(`(?is)<pre[^>]*>(?:<code[^>]*>)?(.*?)(?:</code>)?</pre>`)
	return re.ReplaceAllString(html, "\n```\n$1\n```\n")
}

// convertInlineCode converts inline code tags.
func convertInlineCode(html string) string {
	re := regexp.MustCompile(`(?i)<code[^>]*>(.*?)</code>`)
	return re.ReplaceAllString(html, "`$1`")
}

// convertStrong converts b/strong tags.
func convertStrong(html string) string {
	for _, tag := range []string{"b", "strong"} {
		re := regexp.MustCompile(`(?i)<` + tag + `[^>]*>(.*?)</` + tag + `>`)
		html = re.ReplaceAllString(html, "**$1**")
	}
	return html
}

// convertEmphasis converts i/em tags.
func convertEmphasis(html string) string {
	for _, tag := range []string{"i", "em"} {
		re := regexp.MustCompile(`(?i)<` + tag + `[^>]*>(.*?)</` + tag + `>`)
		html = re.ReplaceAllString(html, "*$1*")
	}
	return html
}

// convertLinks converts a tags to Markdown links.
func convertLinks(html string) string {
	re := regexp.MustCompile(`(?i)<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`)
	return re.ReplaceAllString(html, "[$2]($1)")
}

// convertImages converts img tags to Markdown images.
func convertImages(html string) string {
	re := regexp.MustCompile(`(?i)<img[^>]*src=["']([^"']*)["'][^>]*alt=["']([^"']*)["'][^>]*/?\s*>`)
	html = re.ReplaceAllString(html, "![$2]($1)")
	// Without alt
	re2 := regexp.MustCompile(`(?i)<img[^>]*src=["']([^"']*)["'][^>]*/?\s*>`)
	html = re2.ReplaceAllString(html, "![]($1)")
	return html
}

// convertParagraphs converts p tags.
func convertParagraphs(html string) string {
	re := regexp.MustCompile(`(?is)<p[^>]*>(.*?)</p>`)
	return re.ReplaceAllString(html, "\n$1\n")
}

// convertLineBreaks converts br tags.
func convertLineBreaks(html string) string {
	re := regexp.MustCompile(`(?i)<br\s*/?>`)
	return re.ReplaceAllString(html, "\n")
}

// convertHorizontalRules converts hr tags.
func convertHorizontalRules(html string) string {
	re := regexp.MustCompile(`(?i)<hr\s*/?>`)
	return re.ReplaceAllString(html, "\n---\n")
}

// convertBlockquotes converts blockquote tags.
func convertBlockquotes(html string) string {
	re := regexp.MustCompile(`(?is)<blockquote[^>]*>(.*?)</blockquote>`)
	return re.ReplaceAllStringFunc(html, func(match string) string {
		subRe := regexp.MustCompile(`(?is)<blockquote[^>]*>(.*?)</blockquote>`)
		subs := subRe.FindStringSubmatch(match)
		if len(subs) < 2 {
			return match
		}
		lines := strings.Split(strings.TrimSpace(subs[1]), "\n")
		for i, line := range lines {
			lines[i] = "> " + line
		}
		return "\n" + strings.Join(lines, "\n") + "\n"
	})
}

// stripAllTags removes all remaining HTML tags.
func stripAllTags(html string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(html, "")
}

// decodeEntities decodes common HTML entities.
func decodeEntities(html string) string {
	replacements := map[string]string{
		"&amp;":  "&",
		"&lt;":   "<",
		"&gt;":   ">",
		"&quot;": `"`,
		"&#39;":  "'",
		"&apos;": "'",
		"&nbsp;": " ",
	}
	for entity, replacement := range replacements {
		html = strings.ReplaceAll(html, entity, replacement)
	}
	return html
}

// cleanWhitespace collapses multiple blank lines and trims.
func cleanWhitespace(s string) string {
	// Replace multiple newlines with double newline
	re := regexp.MustCompile(`\n{3,}`)
	s = re.ReplaceAllString(s, "\n\n")
	// Trim trailing spaces from lines
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}
