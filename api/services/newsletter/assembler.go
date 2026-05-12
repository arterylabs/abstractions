package newsletter

import (
	"bytes"
	"html/template"
)

// innerTemplate renders the inner body content of a newsletter (heading + sections).
// This HTML is then wrapped by email.BuildEmail() before storage and delivery.
var innerTemplate = template.Must(template.New("newsletter").Parse(`<h1>{{.RepoName}} — Engineering Newsletter</h1>
{{range .Entries}}
<section style="margin:32px 0">
{{.Content}}
</section>
{{end}}`))

type innerData struct {
	RepoName string
	Entries  []assemblerEntry
}

type assemblerEntry struct {
	// Content is trusted HTML produced by the AI layer.
	Content template.HTML
}

// Assembler converts a slice of NewsletterEntry values into inner body HTML.
type Assembler struct{}

// Assemble renders the inner newsletter HTML for the given repo and entries.
// The returned HTML is a fragment (no <html>/<body> wrapper) suitable for
// passing to email.BuildEmail() as BodyHTML.
func (a *Assembler) Assemble(repoName string, entries []NewsletterEntry) (string, error) {
	data := innerData{RepoName: repoName}
	for _, e := range entries {
		data.Entries = append(data.Entries, assemblerEntry{
			Content: template.HTML(e.Content), // #nosec G203 — content is LLM-generated HTML
		})
	}
	var buf bytes.Buffer
	if err := innerTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
