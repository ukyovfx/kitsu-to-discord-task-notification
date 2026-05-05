package discord

import (
	"app/src/api/kitsu"
	"app/src/utils/config"
	"app/src/utils/truncate"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type EmbedAuthor struct {
	Name string `json:"name"`
}
type EmbedFooter struct {
	Text string `json:"text"`
}
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type EmbedImage struct {
	URL string `json:"url"`
}
type Embed struct {
	Description string       `json:"description,omitempty"`
	Title       string       `json:"title,omitempty"`
	Color       int          `json:"color,omitempty"`
	Url         string       `json:"url,omitempty"`
	Author      EmbedAuthor  `json:"author,omitempty"`
	Footer      EmbedFooter  `json:"footer,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Image       *EmbedImage  `json:"image,omitempty"`
}
type Payload struct {
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds,omitempty"`
}
type Assignee struct {
	Fullname string
	Email    string
	Phone    string
}
type Template struct {
	ProjectName             string
	GroupName               string
	ParentName              string
	TaskName                string
	TaskType                string
	CurrentStatus           string
	PreviousStatus          string
	CommentContent          string
	CommentAuthor           string
	EntityType              string
	Assignees               []Assignee
	AssigneesStr            string
	TaskURL                 string
	MentionContent          string
	ProcessEmoji            string
	StatusMessage           string
	StatusUpper             string
	StatusEmoji             string
	GoogleDriveURL          string
	PreviewImageURL         string // Kitsu プレビュー画像 URL（空ならスキップ）
	IsCommentOnly           bool   // コメントのみの更新（ステータス変化なし）
	StatusTransitionMessage string // ステータス遷移を説明するメッセージ
}

// Discord APIのメッセージ作成レスポンス
type DiscordMessage struct {
	ID string `json:"id"`
}

// containsIgnoreCase reports whether `target` (case-insensitive) is present in `list`.
func containsIgnoreCase(list []string, target string) bool {
	for _, s := range list {
		if strings.EqualFold(s, target) {
			return true
		}
	}
	return false
}

// statusMessageInfo returns the Japanese message text and emoji for a given status.
// Returns ("", "📌") for statuses we don't have a tailored message for.
func statusMessageInfo(status string) (string, string) {
	switch strings.ToUpper(status) {
	case "WFA":
		return "チェックをお願いします", "👀"
	case "RETAKE":
		return "作業をお願いします！", "🔃"
	case "DONE":
		return "チェックが完了しました。お疲れ様でした！", "✅"
	case "READY":
		return "作業を開始してください", "🟢"
	default:
		return "", "📌"
	}
}

// statusTransitionMessage は前後のステータス遷移に応じた補足メッセージを返す。
// 特に対応しない遷移では空文字列を返す。
func statusTransitionMessage(prev, current string) string {
	if prev == "" || strings.EqualFold(prev, current) {
		return ""
	}
	key := strings.ToUpper(prev) + "→" + strings.ToUpper(current)
	switch key {
	case "RETAKE→DONE", "RETAKE→WFA":
		return "🔁 再修正版がアップされました"
	case "WFA→RETAKE":
		return "🔃 チェック後、リテイクになりました"
	case "READY→WIP":
		return "🚀 作業を開始しました"
	case "WIP→WFA":
		return "📤 チェック依頼が送られました"
	case "WFA→DONE":
		return "🎉 最終承認されました！お疲れ様でした"
	default:
		return ""
	}
}

// 工程ごとのアイコン
func getProcessEmoji(taskType string) string {
	switch strings.ToLower(taskType) {
	case "animation":
		return "🎞️"
	case "fx":
		return "💥"
	case "lighting":
		return "💡"
	case "rendering":
		return "🖼️"
	case "compositing":
		return "✨"
	case "color grading":
		return "🎨"
	case "modeling":
		return "🧊"
	case "texturing":
		return "🖌️"
	case "shading":
		return "🎭"
	case "rigging":
		return "🦴"
	case "lookdev":
		return "👁️"
	default:
		return "🎬"
	}
}

// DeleteMessage は既存のDiscordメッセージを削除する
// 成功時は true、失敗時は false を返す（失敗してもプロセスは継続）
func DeleteMessage(webhookURL, messageID string) bool {
	if webhookURL == "" || messageID == "" {
		return false
	}
	deleteURL := fmt.Sprintf("%s/messages/%s", webhookURL, messageID)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		slog.Warn("DeleteMessage: failed to build request", "err", err, "messageID", messageID)
		return false
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("DeleteMessage: request failed", "err", err, "messageID", messageID)
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	// 204 No Content = 削除成功 / 404 Not Found = すでに削除済み（成功扱い）
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return true
	}
	slog.Warn("DeleteMessage: unexpected status", "status", resp.StatusCode, "messageID", messageID)
	return false
}

// SendMessage は1タスク分のメッセージを送信してメッセージIDを返す
// 失敗時は空文字列を返す（呼び出し側で判定して旧メッセージ残置等のフォールバックを行う）
func SendMessage(payload Payload, webhookURL string) string {
	// ?wait=true をつけることでDiscordがメッセージIDを返してくれる
	url := webhookURL + "?wait=true"

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("SendMessage: marshal failed", "err", err)
		return ""
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("SendMessage: HTTP post failed", "err", err)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error("SendMessage: read body failed", "err", err)
		return ""
	}

	// 2xx 以外は失敗扱い。429 のときは Retry-After を見て呼び出し側に空を返す
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("SendMessage: non-2xx response",
			"status", resp.StatusCode,
			"retryAfter", resp.Header.Get("Retry-After"),
			"body", string(respBody))
		return ""
	}

	var msg DiscordMessage
	if err := json.Unmarshal(respBody, &msg); err != nil {
		slog.Error("SendMessage: unmarshal failed", "err", err, "body", string(respBody))
		return ""
	}

	return msg.ID
}

// SendMessageBunch は各タスクを処理してメッセージIDのマップを返す
// キーはタスクID、値はDiscordメッセージID
func SendMessageBunch(conf config.Config, data []kitsu.MessagePayload, webHookURL string, previousMessageIDs map[string]string, previousWebhookURLs map[string]string) map[string]string {
	result := make(map[string]string)

	for _, elem := range data {
		var placeholders Template

		placeholders.ProjectName = elem.Project.Name
		placeholders.GroupName = elem.EntityType.Name
		placeholders.ParentName = elem.Parent.Name
		placeholders.TaskName = elem.Entity.Name
		placeholders.TaskType = elem.TaskType.Name
		placeholders.CurrentStatus = elem.TaskStatus.ShortName
		placeholders.StatusUpper = strings.ToUpper(elem.TaskStatus.ShortName)
		placeholders.PreviousStatus = elem.PreviousStatusName
		placeholders.CommentContent = elem.LatestComment.Comment.Text
		placeholders.CommentAuthor = elem.LatestComment.Author.FullName
		placeholders.EntityType = elem.EntityType.EntityType.Name
		placeholders.ProcessEmoji = getProcessEmoji(elem.TaskType.Name)

		// Assignees を解決しつつ、アーティストの Discord ID も収集。
		// 同じ Kitsu 名が userMap に複数行で書かれていても 1 回しか mention しないよう、
		// 既に追加した DiscordID を seen で重複排除する。
		var assigneeNames []string
		var artistMentions []string
		seenDiscordID := make(map[string]bool)
		placeholders.Assignees = make([]Assignee, len(elem.Assignees))
		for i := 0; i < len(elem.Assignees); i++ {
			fullName := elem.Assignees[i].FullName
			placeholders.Assignees[i].Fullname = fullName
			placeholders.Assignees[i].Email = elem.Assignees[i].Email
			assigneeNames = append(assigneeNames, fullName)

			matched := false
			for _, u := range conf.Mention.UserMap {
				if u.KitsuName == fullName {
					if !seenDiscordID[u.DiscordID] {
						artistMentions = append(artistMentions, "<@"+u.DiscordID+">")
						seenDiscordID[u.DiscordID] = true
					}
					matched = true
					break
				}
			}
			if !matched && len(conf.Mention.UserMap) > 0 {
				slog.Warn("Kitsu user not in mention.userMap; will not be @-mentioned",
					"kitsuName", fullName,
					"taskID", elem.Task.ID,
					"hint", "add this person to [[mention.userMap]] in conf.toml")
			}
		}
		if len(assigneeNames) > 0 {
			placeholders.AssigneesStr = strings.Join(assigneeNames, ", ")
		} else {
			placeholders.AssigneesStr = "未割り当て"
		}

		// タスクタイプからチェッカーの Discord ID を検索
		checkerMention := ""
		for _, c := range conf.Mention.Checkers {
			if strings.EqualFold(c.TaskType, elem.TaskType.Name) {
				checkerMention = "<@" + c.DiscordID + ">"
				break
			}
		}

		// ステータスごとのメンション対象を conf.Mention の設定から決める。
		// CheckerStatuses に含まれるステータスではチェッカーをメンション、
		// ArtistStatuses に含まれるステータスではアーティスト全員をメンションする。
		// HereStatuses に含まれるステータスでは @here を追加（緊急通知）。
		// 複数に含まれるステータスでは全てを併記する。
		currentStatus := strings.ToUpper(elem.TaskStatus.ShortName)
		var mentionParts []string
		if containsIgnoreCase(conf.Mention.CheckerStatuses, currentStatus) {
			if checkerMention != "" {
				mentionParts = append(mentionParts, checkerMention)
			} else if len(conf.Mention.Checkers) > 0 {
				slog.Warn("No checker configured for task type; checker will not be @-mentioned",
					"taskType", elem.TaskType.Name,
					"status", currentStatus,
					"taskID", elem.Task.ID,
					"hint", "add this task type to [[mention.checkers]] in conf.toml")
			}
		}
		if containsIgnoreCase(conf.Mention.ArtistStatuses, currentStatus) {
			if artistJoined := strings.Join(artistMentions, " "); artistJoined != "" {
				mentionParts = append(mentionParts, artistJoined)
			}
		}
		// 緊急ステータスは @here でチャンネル全員に通知
		if containsIgnoreCase(conf.Mention.HereStatuses, currentStatus) {
			mentionParts = append(mentionParts, "@here")
		}
		mentionContent := strings.Join(mentionParts, " ")
		statusMessage, statusEmoji := statusMessageInfo(currentStatus)

		placeholders.MentionContent = mentionContent
		placeholders.StatusMessage = statusMessage
		placeholders.StatusEmoji = statusEmoji
		placeholders.GoogleDriveURL = conf.GoogleDrive.URL
		placeholders.IsCommentOnly = elem.IsCommentOnly
		placeholders.StatusTransitionMessage = statusTransitionMessage(elem.PreviousStatusName, elem.TaskStatus.ShortName)

		// TaskURL を組み立て
		category := "assets"
		lowEntityType := strings.ToLower(placeholders.EntityType)
		if strings.Contains(lowEntityType, "shot") || strings.Contains(lowEntityType, "sequence") || strings.Contains(lowEntityType, "episode") {
			category = "shots"
		}
		host := conf.Kitsu.Hostname
		if !strings.HasSuffix(host, "/") {
			host += "/"
		}
		// ショット一覧 or アセット一覧に飛ぶ
		placeholders.TaskURL = fmt.Sprintf("%sproductions/%s/%s", host, elem.Project.ID, category)

		// Kitsu プレビュー画像 URL を組み立て
		// Entity.PreviewFileID は interface{} なので string にキャストする
		if id, ok := elem.Entity.Entity.PreviewFileID.(string); ok && id != "" {
			placeholders.PreviewImageURL = fmt.Sprintf("%sapi/pictures/preview-files/%s/original", host, id)
		}

		// カラーコードを変換
		hexColor := strings.ReplaceAll(elem.TaskStatus.Color, "#", "")
		intColor, _ := strconv.ParseInt(hexColor, 16, 64)

		// テンプレートを展開
		tplPreset := conf.TplPreset
		author := parseTaskTemplate("tpl/"+tplPreset+"/author.tpl", placeholders)
		title := parseTaskTemplate("tpl/"+tplPreset+"/title.tpl", placeholders)
		description := parseTaskTemplate("tpl/"+tplPreset+"/description.tpl", placeholders)
		footer := parseTaskTemplate("tpl/"+tplPreset+"/footer.tpl", placeholders)

		embed := Embed{}
		embed.Title = truncate.TruncateString(title, 256)
		embed.Description = truncate.TruncateString(description, 4096)
		embed.Color = int(intColor)
		embed.Author.Name = truncate.TruncateString(author, 256)
		embed.Footer.Text = truncate.TruncateString(footer, 2048)
		embed.Url = truncate.TruncateString(placeholders.TaskURL, 2000)

		// プレビュー画像がある場合は embed に添付
		if placeholders.PreviewImageURL != "" {
			embed.Image = &EmbedImage{URL: placeholders.PreviewImageURL}
		}

		fieldsRaw := parseTaskTemplate("tpl/"+tplPreset+"/fields.tpl", placeholders)
		if strings.TrimSpace(fieldsRaw) != "" {
			var parsedFields []EmbedField
			err := json.Unmarshal([]byte(fieldsRaw), &parsedFields)
			if err == nil {
				embed.Fields = parsedFields
			}
		}

		// payload を作成
		payload := Payload{}
		if mentionContent != "" {
			payload.Content = mentionContent
		}
		payload.Embeds = []Embed{embed}

		// まず新規メッセージを送信。送信成功時のみ旧メッセージを削除する
		// （順序が逆だと、新規送信失敗時に Discord 上の履歴が消失する）
		newMessageID := SendMessage(payload, webHookURL)
		if newMessageID == "" {
			// 送信失敗時は旧メッセージを残す。DB の messageID も既存値を維持する
			slog.Warn("SendMessage failed; keeping previous message", "taskID", elem.Task.ID)
			if prevMsgID, ok := previousMessageIDs[elem.Task.ID]; ok && prevMsgID != "" {
				result[elem.Task.ID] = prevMsgID
			}
			continue
		}

		// 送信成功 → 旧メッセージを削除
		if prevMsgID, ok := previousMessageIDs[elem.Task.ID]; ok && prevMsgID != "" {
			prevWebhook := previousWebhookURLs[elem.Task.ID]
			if prevWebhook == "" {
				prevWebhook = webHookURL
			}
			DeleteMessage(prevWebhook, prevMsgID)
		}

		result[elem.Task.ID] = newMessageID
	}

	return result
}

func parseTaskTemplate(tplFilePath string, data Template) string {
	tpl, err := ioutil.ReadFile(tplFilePath)
	if err != nil {
		return ""
	}
	t := template.Must(template.New("template").Parse(string(tpl)))
	output := new(bytes.Buffer)
	t.Execute(output, data)
	return strings.TrimSpace(output.String())
}