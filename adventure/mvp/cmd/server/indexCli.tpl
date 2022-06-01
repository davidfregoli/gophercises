• {{.Title}}
{{range .Paragraphs}}
{{.}}
{{end}}
----------------------
{{range .Options}}
- ({{.Chapter}}) > {{.Text}}
{{else}}
- (intro) > Start over
{{end}}
- (quit) > Exit Adventure
