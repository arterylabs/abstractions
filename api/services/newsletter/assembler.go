package newsletter

import (
	"bytes"
	"html/template"
)

// assemblerBodyTemplate renders the inner HTML body of the newsletter.
// This is injected into the generic email template wrapper; no outer HTML tags are emitted.
var assemblerBodyTemplate = template.Must(template.New("newsletter-body").Parse(
	`<h2 style="font-size:1.4rem;margin:0 0 1.5rem;color:#333333;">{{.RepoName}} — Engineering Newsletter</h2>` +
		`{{range .Entries}}<div style="margin-bottom:2rem;">{{.Content}}</div>{{end}}`,
))

// assemblerData is the template data.
type assemblerData struct {
	RepoName string
	Entries  []assemblerEntry
}

type assemblerEntry struct {
	// Content is trusted HTML produced by the AI layer.
	Content template.HTML
}

// Assembler converts a slice of NewsletterEntry values into the inner HTML body of a newsletter email.
type Assembler struct{}

// Assemble renders the newsletter inner body HTML for the given repo and entries.
// The returned string is intended to be passed as BodyHTML to email.BuildEmail.
func (a *Assembler) Assemble(repoName string, entries []NewsletterEntry) (string, error) {
	data := assemblerData{RepoName: repoName}
	for _, e := range entries {
		data.Entries = append(data.Entries, assemblerEntry{
			Content: template.HTML(e.Content), // #nosec G203 — content is LLM-generated HTML
		})
	}
	var buf bytes.Buffer
	if err := assemblerBodyTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
