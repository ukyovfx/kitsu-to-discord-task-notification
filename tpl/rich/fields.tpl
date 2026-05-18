[
  {
    "name": "📊 ステータス",
    "value": "{{if .PreviousStatus}}{{.PreviousStatus}} ➡️ {{end}}{{.CurrentStatus}}",
    "inline": true
  },
  {
    "name": "👤 担当",
    "value": "{{.AssigneesStr}}",
    "inline": true
  }{{if .GoogleDriveURL}},
  {
    "name": "📁 ストレージ",
    "value": "[Open Storage]({{.GoogleDriveURL}})",
    "inline": true
  }{{end}}
]