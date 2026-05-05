## {{.GroupName}} / {{.ParentName}} - {{.TaskName}}
### {{.StatusEmoji}} {{.StatusUpper}}{{if .IsCommentOnly}} 💬{{end}}
> {{if .MentionContent}}{{.MentionContent}} {{end}}{{if .StatusMessage}}**{{.StatusMessage}}**{{end}}{{if .StatusTransitionMessage}}
> {{.StatusTransitionMessage}}{{end}}

変更者: {{.CommentAuthor}}{{if .CommentContent}}
💬 {{.CommentContent}}{{end}}

[📁 Google Drive]({{.GoogleDriveURL}}) ・ [🦊 KITSU]({{.TaskURL}})