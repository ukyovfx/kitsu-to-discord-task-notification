## {{.GroupName}} / {{.ParentName}} - {{.TaskName}}
### {{.StatusEmoji}} {{.StatusUpper}}
> {{if .MentionContent}}{{.MentionContent}} {{end}}{{if .StatusMessage}}**{{.StatusMessage}}**{{end}}

変更者: {{.CommentAuthor}}{{if .CommentContent}}
💬 {{.CommentContent}}{{end}}

[📁 Google Drive]({{.GoogleDriveURL}}) ・ [🦊 KITSU]({{.TaskURL}})