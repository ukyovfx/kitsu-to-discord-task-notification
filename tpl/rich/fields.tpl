[
  {
    "name": "📁 カテゴリ",
    "value": "{{.ProjectName}}",
    "inline": true
  },
  {
    "name": "📊 ステータス",
    "value": "{{if .PreviousStatus}}{{.PreviousStatus}} ➡️ {{end}}{{.CurrentStatus}}",
    "inline": true
  },
  {
    "name": "👤 担当",
    "value": "{{.AssigneesStr}}",
    "inline": true
  }
]