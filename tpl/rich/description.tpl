{{if .IsAssignNotification}}📋 タスクがアサインされました

{{end}}{{.StatusEmoji}} {{.StatusUpper}}{{if .IsCommentOnly}} 💬{{end}}
{{if .MentionContent}}> {{.MentionContent}}
{{end}}{{if .StatusMessage}}{{.StatusMessage}}{{end}}{{if .StatusTransitionMessage}}
{{.StatusTransitionMessage}}{{end}}

変更者: {{.CommentAuthor}}{{if .CommentContent}}
💬 {{.CommentContent}}{{end}}

**[🦊 KITSU]({{.TaskURL}})**{{if .GoogleDriveURL}} ・ [📁 Drive]({{.GoogleDriveURL}}){{end}}
