package newsletter

import (
	"bytes"
	"html/template"
)

// assemblerTemplate is the full HTML email template.
// Each NewsletterEntry produces a <section> block.
var assemblerTemplate = template.Must(template.New("newsletter").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.RepoName}} — CodeBrief</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; color: #111; max-width: 680px; margin: 0 auto; padding: 24px; }
  h1   { font-size: 1.6rem; border-bottom: 2px solid #e0e0e0; padding-bottom: 8px; }
  section { margin: 32px 0; }
  footer  { margin-top: 48px; font-size: 0.85rem; color: #666; }
  a       { color: #0066cc; }
</style>
</head>
<body>
<h1>{{.RepoName}} — Engineering Newsletter</h1>
{{range .Entries}}
<section>
{{.Content}}
</section>
{{end}}
<footer>
<p>You are receiving this because you are subscribed to the <strong>{{.RepoName}}</strong> newsletter.
{{if .UnsubscribeURL}}<a href="{{.UnsubscribeURL}}">Unsubscribe</a>{{end}}</p>
</footer>
</body>
</html>
`))

// assemblerData is the template data.
type assemblerData struct {
	RepoName       string
	Entries        []assemblerEntry
	UnsubscribeURL string
}

type assemblerEntry struct {
	// Content is trusted HTML produced by the AI layer.
	Content template.HTML
}

// Assembler converts a slice of NewsletterEntry values into a full HTML email string.
type Assembler struct{}

// Assemble renders the newsletter HTML for the given repo and entries.
// unsubscribeURL may be empty; if set it is rendered in the footer.
func (a *Assembler) Assemble(repoName string, entries []NewsletterEntry, unsubscribeURL string) (string, error) {
	data := assemblerData{
		RepoName:       repoName,
		UnsubscribeURL: unsubscribeURL,
	}
	for _, e := range entries {
		data.Entries = append(data.Entries, assemblerEntry{
			Content: template.HTML(e.Content), // #nosec G203 — content is LLM-generated HTML
		})
	}
	var buf bytes.Buffer
	if err := assemblerTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
