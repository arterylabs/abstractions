package newsletter

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"
)

var assemblerBodyTemplate = template.Must(
	template.New("newsletter-body").
		Funcs(template.FuncMap{
			"join": strings.Join,
			"formatTime": func(t *time.Time) string {
				if t == nil {
					return ""
				}
				return t.Format("January 2, 2006")
			},
		}).
		Parse(
			`{{range $i, $e := .Entries}}{{if $i}}<hr style="border:none;border-top:1px solid #eaeaec;margin:2.5rem 0;" />{{end}}<div>{{$e.Content}}</div>{{end}}` +
				`<div style="margin-top:3rem;padding-top:2rem;border-top:1px solid #eaeaec;font-size:0.875rem;color:#6b7280;">` +
				`<p style="margin:0 0 0.75rem;font-size:0.8rem;font-weight:700;letter-spacing:0.05em;text-transform:uppercase;color:#9ca3af;">About this newsletter</p>` +
				`{{if .RepoURL}}<p style="margin:0 0 0.4rem;"><strong style="color:#374151;">Repository:</strong> <a href="{{.RepoURL}}" style="color:#14b7d3;">{{.RepoName}}</a></p>{{end}}` +
				`{{if .RelatedFiles}}<p style="margin:0 0 0.4rem;"><strong style="color:#374151;">Files covered:</strong> {{range $i, $f := .RelatedFiles}}{{if $i}}, {{end}}<code style="font-family:monospace;font-size:0.8rem;background:#f4f4f7;padding:1px 5px;border-radius:3px;border:1px solid #e2e4e8;color:#333;">{{$f}}</code>{{end}}</p>{{end}}` +
				`{{if .IndexedAt}}<p style="margin:0 0 0.75rem;"><strong style="color:#374151;">Last indexed:</strong> {{formatTime .IndexedAt}}</p>{{end}}` +
				`<p style="margin:0;font-size:0.75rem;color:#9ca3af;font-style:italic;">This content is AI-generated. Suggested improvements may not apply to your codebase — review them with that in mind.</p>` +
				`</div>`,
		),
)

type assemblerData struct {
	RepoName     string
	RepoURL      string
	RelatedFiles []string
	IndexedAt    *time.Time
	Entries      []assemblerEntry
}

type assemblerEntry struct {
	// Content is trusted HTML produced by the AI layer.
	Content template.HTML
}

// Assembler converts a slice of NewsletterEntry values into the inner HTML body of a newsletter email.
type Assembler struct{}

// Assemble renders the newsletter inner body HTML for the given repo and entries.
// The returned string is intended to be passed as BodyHTML to email.BuildEmail.
func (a *Assembler) Assemble(repoName, repoURL string, relatedFiles []string, indexedAt *time.Time, entries []NewsletterEntry) (string, error) {
	data := assemblerData{
		RepoName:     repoName,
		RepoURL:      repoURL,
		RelatedFiles: relatedFiles,
		IndexedAt:    indexedAt,
	}
	for _, e := range entries {
		data.Entries = append(data.Entries, assemblerEntry{
			Content: template.HTML(e.Content), // #nosec G203 — content is LLM-generated HTML
		})
	}
	var buf bytes.Buffer
	if err := assemblerBodyTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("newsletter: assemble: execute template: %w", err)
	}
	return buf.String(), nil
}
