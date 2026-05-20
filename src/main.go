package main

import (
	"app/src/api/discord"
	"app/src/api/kitsu"
	"app/src/model"
	"app/src/setup"
	"app/src/utils/basicauth"
	"app/src/utils/config"
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/beefsack/go-rate"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"

	"github.com/pieterclaerhout/go-waitgroup"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MakeKitsuResponse(conf config.Config) []kitsu.MessagePayload {

	tasks := kitsu.GetTasks()
	if conf.Log {
		slog.Info("Got tasks: " + strconv.Itoa(len(tasks.Each)))
	}

	taskStatuses := kitsu.GetTaskStatuses()
	if conf.Log {
		slog.Info("Got taskStatuses: " + strconv.Itoa(len(taskStatuses.Each)))
	}

	entities := kitsu.GetEntities()
	if conf.Log {
		slog.Info("Got entities: " + strconv.Itoa(len(entities.Each)))
	}

	enitityTypes := kitsu.GetEntityTypes()
	if conf.Log {
		slog.Info("Got enitityTypes: " + strconv.Itoa(len(enitityTypes.Each)))
	}

	projects := kitsu.GetProjects()
	if conf.Log {
		slog.Info("Got projects: " + strconv.Itoa(len(projects.Each)))
	}

	taskTypes := kitsu.GetTaskTypes()
	if conf.Log {
		slog.Info("Got taskTypes: " + strconv.Itoa(len(taskTypes.Each)))
	}

	persons := kitsu.GetPersons()
	if conf.Log {
		slog.Info("Got persons: " + strconv.Itoa(len(persons.Each)))
	}

	var comments kitsu.Comments
	if conf.Kitsu.SkipComments == false {
		comments = kitsu.GetComments()
		if conf.Log {
			slog.Info("Got comments: " + strconv.Itoa(len(comments.Each)))
		}
	}

	start := time.Now()

	response := make([]kitsu.MessagePayload, len(tasks.Each))

	wg := waitgroup.NewWaitGroup(conf.Threads)

	for i := 0; i < len(response); i++ {
		wg.BlockAdd()
		go func(i int) {
			defer wg.Done()

			layout := "2006-01-02T15:04:05"
			taskDate, err := time.Parse(layout, tasks.Each[i].UpdatedAt)
			if err != nil {
				slog.Info(err)
			}
			daysCount := int(start.Sub(taskDate).Hours() / 24)

			if conf.IgnoreMessagesDaysOld != 0 && daysCount > conf.IgnoreMessagesDaysOld {
				return
			}

			response[i].Task.Task = tasks.Each[i]

			for _, elem := range taskStatuses.Each {
				if elem.ID == tasks.Each[i].TaskStatusID {
					response[i].TaskStatus.TaskStatus = elem
					break
				}
			}

			for _, elem := range entities.Each {
				if elem.ID == tasks.Each[i].EntityID {
					response[i].Entity.Entity = elem
					break
				}
			}

			for _, elem := range enitityTypes.Each {
				if elem.ID == response[i].Entity.Entity.EntityTypeID {
					response[i].EntityType.EntityType = elem
					break
				}
			}

			for _, elem := range entities.Each {
				if elem.ID == response[i].Entity.Entity.ParentID {
					response[i].Parent.Entity = elem
				}
			}

			for _, elem := range projects.Each {
				if elem.ID == response[i].Entity.Entity.ProjectID {
					response[i].Project.Project = elem
					break
				}
			}

			for _, elem := range taskTypes.Each {
				if elem.ID == tasks.Each[i].TaskTypeID {
					response[i].TaskType.TaskType = elem
					break
				}
			}

			if conf.Kitsu.SkipComments == false {
				var taskComments kitsu.Comments
				for _, elem := range comments.Each {
					if elem.ObjectID == tasks.Each[i].ID {
						taskComments.Each = append(taskComments.Each, elem)
					}
				}

				if len(taskComments.Each) > 0 {
					sort.Slice(taskComments.Each, func(i, j int) bool {
						layout := "2006-01-02T15:04:05"
						a, err := time.Parse(layout, taskComments.Each[i].UpdatedAt)
						if err != nil {
							slog.Info(err)
						}
						b, err := time.Parse(layout, taskComments.Each[j].UpdatedAt)
						if err != nil {
							slog.Info(err)
						}
						return a.Unix() > b.Unix()
					})

					response[i].LatestComment.Comment.Comment = taskComments.Each[0]
				}

				for _, elem := range persons.Each {
					if len(taskComments.Each) > 0 {
						if elem.ID == taskComments.Each[0].PersonID {
							response[i].LatestComment.Author.Person = elem
							break
						}
					}
				}
			}

			if len(tasks.Each[i].Assignees) > 0 {
				for _, assigneeID := range tasks.Each[i].Assignees {
					for _, person := range persons.Each {
						if assigneeID == person.ID {
							response[i].Assignees = append(response[i].Assignees, person)
						}
					}
				}
			}

		}(i)
	}
	wg.Wait()

	if conf.Log {
		log.Printf("Done primary loop in %s", time.Since(start))
	}

	var out []kitsu.MessagePayload
	for _, elem := range response {
		if len(elem.Task.Task.ID) > 0 {
			out = append(out, elem)
		}
	}

	if conf.Log {
		log.Printf("Done secondary loop in %s", time.Since(start))
	}

	return out
}

type notificationRouteStats struct {
	DBRouted         int
	ProjectRouted    int
	TaskTypeRouted   int
	MainFallbackSent int
	Dropped          int
}

func previewTasks(tasks []kitsu.MessagePayload, limit int) []string {
	if limit <= 0 || len(tasks) == 0 {
		return nil
	}
	if limit > len(tasks) {
		limit = len(tasks)
	}
	preview := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		task := tasks[i]
		preview = append(preview, fmt.Sprintf("%s|%s|%s|%s|%s",
			task.Task.ID,
			task.Project.Name,
			task.Entity.Name,
			task.TaskType.TaskType.Name,
			task.TaskStatus.TaskStatus.ShortName,
		))
	}
	return preview
}

func labelsFromSet(set map[string]struct{}) []string {
	if len(set) == 0 {
		return nil
	}
	labels := make([]string, 0, len(set))
	for label := range set {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}

func logRouteDispatch(routeSource string, routeLabels []string, tasks []kitsu.MessagePayload, webhookConfigured bool) {
	logger := slog.Info
	message := "Notification route dispatch"
	if !webhookConfigured {
		logger = slog.Warn
		message = "Notification route has no webhook configured; tasks will not be sent"
	}
	logger(message,
		"routeSource", routeSource,
		"routeLabels", routeLabels,
		"taskCount", len(tasks),
		"taskPreview", previewTasks(tasks, 5))
}

func logDroppedTasks(reason string, tasks []kitsu.MessagePayload) {
	slog.Warn("Notification dropped",
		"reason", reason,
		"taskCount", len(tasks),
		"taskPreview", previewTasks(tasks, 5))
}

func FilterTasks(data []kitsu.MessagePayload, conf config.Config, db *gorm.DB) {
	if len(data) == 0 {
		if conf.Log {
			fmt.Printf("Nothing to do\n")
		}
	}

	var filtered []kitsu.MessagePayload
	for i := 0; i < len(data); i++ {

		dbResult := model.FindTask(db, data[i].Task.ID)

		data[i].PreviousStatusName = dbResult.TaskStatus

		if len(dbResult.TaskID) > 0 {
			statusChanged := dbResult.TaskStatus != data[i].TaskStatus.TaskStatus.ShortName
			timestampChanged := dbResult.TaskUpdatedAt != data[i].Task.Task.UpdatedAt
			commentChanged := dbResult.CommentUpdatedAt != data[i].LatestComment.Comment.UpdatedAt

			if statusChanged || timestampChanged || commentChanged {
				// 繧ｳ繝｡繝ｳ繝医・縺ｿ螟牙喧・医せ繝・・繧ｿ繧ｹ繝ｻ繧ｿ繧､繝繧ｹ繧ｿ繝ｳ繝励・蜷後§・峨°縺ｩ縺・°繧定ｨ倬鹸
				data[i].IsCommentOnly = commentChanged && !statusChanged && !timestampChanged
				model.UpdateTask(db, data[i].Task.Task.ID, data[i].Task.Task.UpdatedAt, data[i].TaskStatus.TaskStatus.ShortName, data[i].LatestComment.Comment.ID, data[i].LatestComment.Comment.UpdatedAt)
			} else {
				continue
			}
		} else {
			model.CreateTask(db, data[i].Task.Task.ID, data[i].Task.Task.UpdatedAt, data[i].TaskStatus.TaskStatus.ShortName, data[i].LatestComment.Comment.ID, data[i].LatestComment.Comment.UpdatedAt)
		}

		if conf.SilentUpdateDB {
			if conf.Log {
				log.Printf("Ignoring message\n")
			}
			continue
		}
		// StatusFilter: only notify on WFA, RETAKE, DONE
		currentStatus := data[i].TaskStatus.TaskStatus.ShortName
		if !isNotifiableStatus(currentStatus) {
			continue
		}

		// Treat "none" status as an assign notification when enabled.
		if strings.EqualFold(currentStatus, "none") {
			if !conf.Notification.NotifyOnAssign {
				continue
			}
			data[i].IsAssignNotification = true
		}
		filtered = append(filtered, data[i])
	}

	// 笏笏笏 DB 繝吶・繧ｹ縺ｮ繝励Ο繧ｸ繧ｧ繧ｯ繝医Ν繝ｼ繝・ぅ繝ｳ繧ｰ・・setup 縺ｧ菴懈・縺輔ｌ縺溘・繝ｭ繧ｸ繧ｧ繧ｯ繝亥━蜈茨ｼ俄楳笏笏
	type dbRoute struct {
		WebhookURL   string
		TasksPayload []kitsu.MessagePayload
		RouteLabels  map[string]struct{}
	}
	stats := notificationRouteStats{}
	dbRoutes := make(map[string]*dbRoute) // webhookURL 竊・tasks
	for f := len(filtered) - 1; f >= 0; f-- {
		projectID := filtered[f].Project.ID
		taskTypeName := filtered[f].TaskType.TaskType.Name
		webhookURL := model.FindWebhookForTask(db, projectID, taskTypeName)
		if webhookURL == "" {
			continue
		}
		if _, ok := dbRoutes[webhookURL]; !ok {
			dbRoutes[webhookURL] = &dbRoute{WebhookURL: webhookURL, RouteLabels: make(map[string]struct{})}
		}
		dbRoutes[webhookURL].RouteLabels[fmt.Sprintf("projectID=%s taskType=%s", projectID, taskTypeName)] = struct{}{}
		dbRoutes[webhookURL].TasksPayload = append(dbRoutes[webhookURL].TasksPayload, filtered[f])
		filtered = append(filtered[:f], filtered[f+1:]...)
	}
	for _, route := range dbRoutes {
		if len(route.TasksPayload) > 0 {
			routeLabels := labelsFromSet(route.RouteLabels)
			stats.DBRouted += len(route.TasksPayload)
			logRouteDispatch("db.project_webhook", routeLabels, route.TasksPayload, route.WebhookURL != "")
			DiscordQueueSend(route.TasksPayload, conf, route.WebhookURL, db, "db.project_webhook", routeLabels)
		}
	}

	// 笏笏笏 譌｢蟄倥・ conf.toml 繝吶・繧ｹ縺ｮ繝ｫ繝ｼ繝・ぅ繝ｳ繧ｰ・亥ｾ梧婿莠呈鋤・俄楳笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏
	type TasksByProject struct {
		ProjectName  string
		TasksPayload []kitsu.MessagePayload
	}
	tasksByProject := make([]TasksByProject, len(conf.Discord.Productions))
	for i := 0; i < len(tasksByProject); i++ {
		for f := len(filtered) - 1; f >= 0; f-- {
			if strings.Contains(strings.ToLower(filtered[f].Project.Name), strings.ToLower(conf.Discord.Productions[i].Production)) {
				tasksByProject[i].ProjectName = filtered[f].Project.Name
				tasksByProject[i].TasksPayload = append(tasksByProject[i].TasksPayload, filtered[f])
				filtered = append(filtered[:f], filtered[f+1:]...)
			}
		}
	}

	if len(tasksByProject) > 0 {
		for i := 0; i < len(tasksByProject); i++ {
			if len(tasksByProject[i].TasksPayload) > 0 {
				routeLabels := []string{conf.Discord.Productions[i].Production}
				stats.ProjectRouted += len(tasksByProject[i].TasksPayload)
				logRouteDispatch("conf.production", routeLabels, tasksByProject[i].TasksPayload, conf.Discord.Productions[i].WebhookURL != "")
				DiscordQueueSend(tasksByProject[i].TasksPayload, conf, conf.Discord.Productions[i].WebhookURL, db, "conf.production", routeLabels)
			}
		}
	}

	// 繧ｿ繧ｹ繧ｯ繧ｿ繧､繝暦ｼ亥ｷ･遞具ｼ牙挨 Webhook 繝ｫ繝ｼ繝・ぅ繝ｳ繧ｰ
	// Route tasks configured in [[discord.taskTypeWebhooks]] before the fallback route.
	type tasksByType struct {
		WebhookURL   string
		TasksPayload []kitsu.MessagePayload
		RouteLabel   string
	}
	ttRoutes := make([]tasksByType, len(conf.Discord.TaskTypeWebhooks))
	for i, tw := range conf.Discord.TaskTypeWebhooks {
		ttRoutes[i].WebhookURL = tw.WebhookURL
		ttRoutes[i].RouteLabel = tw.TaskType
	}
	for f := len(filtered) - 1; f >= 0; f-- {
		for i, tw := range conf.Discord.TaskTypeWebhooks {
			if strings.EqualFold(filtered[f].TaskType.TaskType.Name, tw.TaskType) {
				ttRoutes[i].TasksPayload = append(ttRoutes[i].TasksPayload, filtered[f])
				filtered = append(filtered[:f], filtered[f+1:]...)
				break
			}
		}
	}
	for _, route := range ttRoutes {
		if len(route.TasksPayload) > 0 {
			routeLabels := []string{route.RouteLabel}
			stats.TaskTypeRouted += len(route.TasksPayload)
			logRouteDispatch("conf.task_type", routeLabels, route.TasksPayload, route.WebhookURL != "")
			DiscordQueueSend(route.TasksPayload, conf, route.WebhookURL, db, "conf.task_type", routeLabels)
		}
	}

	if len(filtered) > 0 {
		// 繝｡繧､繝ｳ Webhook: DB 蜆ｪ蜈茨ｼ亥・蝗槭そ繝・ヨ繧｢繝・・譎ゅ↓蜈ｨ菴馴夂衍繝√Ε繝ｳ繝阪Ν縺瑚・蜍慕匳骭ｲ縺輔ｌ繧具ｼ峨√↑縺代ｌ縺ｰ conf.toml
		_, _, mainWebhook := getDiscordSettings(db, conf)
		if mainWebhook != "" {
			stats.MainFallbackSent = len(filtered)
			logRouteDispatch("fallback.main_webhook", []string{"default"}, filtered, true)
			DiscordQueueSend(filtered, conf, mainWebhook, db, "fallback.main_webhook", []string{"default"})
		} else {
			stats.Dropped = len(filtered)
			logDroppedTasks("no main webhook configured for unrouted tasks", filtered)
		}
	}

	if len(data) > 0 && (stats.DBRouted > 0 || stats.ProjectRouted > 0 || stats.TaskTypeRouted > 0 || stats.MainFallbackSent > 0 || stats.Dropped > 0) {
		slog.Info("Notification routing summary",
			"incomingTasks", len(data),
			"dbRouted", stats.DBRouted,
			"projectRouted", stats.ProjectRouted,
			"taskTypeRouted", stats.TaskTypeRouted,
			"mainFallbackSent", stats.MainFallbackSent,
			"dropped", stats.Dropped)
	}
}

func isNotifiableStatus(status string) bool {
	// Only send notifications for these statuses
	notifiableStatuses := []string{"wfa", "retake", "done"}
	lowerStatus := strings.ToLower(status)
	for _, s := range notifiableStatuses {
		if lowerStatus == s {
			return true
		}
	}
	return false
}

func DiscordQueueSend(data []kitsu.MessagePayload, conf config.Config, webhookURL string, db *gorm.DB, routeSource string, routeLabels []string) []kitsu.MessagePayload {
	if webhookURL == "" {
		slog.Warn("Notification send skipped: webhook is empty",
			"routeSource", routeSource,
			"routeLabels", routeLabels,
			"taskCount", len(data),
			"taskPreview", previewTasks(data, 5))
		return data
	}

	rl := rate.New(conf.Discord.RequestsPerMinute, time.Minute)

	// 騾∽ｿ｡蜑阪↓譌｢蟄倥・繝｡繝・そ繝ｼ繧ｸID繝ｻWebhookURL繝ｻ繧ｹ繝ｬ繝・ラID繧奪B縺九ｉ蜿朱寔
	previousMessageIDs := make(map[string]string)
	previousWebhookURLs := make(map[string]string)
	previousThreadIDs := make(map[string]string)
	projectNotificationLanguages := make(map[string]string)
	for _, elem := range data {
		dbResult := model.FindTask(db, elem.Task.ID)
		if dbResult.DiscordMessageID != "" {
			previousMessageIDs[elem.Task.ID] = dbResult.DiscordMessageID
			previousWebhookURLs[elem.Task.ID] = dbResult.WebhookURL
		}
		if dbResult.DiscordThreadID != "" {
			previousThreadIDs[elem.Task.ID] = dbResult.DiscordThreadID
		}
		if _, ok := projectNotificationLanguages[elem.Project.ID]; !ok {
			projectNotificationLanguages[elem.Project.ID] = "ja"
			if project := model.FindProjectByKitsuID(db, elem.Project.ID); project != nil && strings.EqualFold(strings.TrimSpace(project.Language), "en") {
				projectNotificationLanguages[elem.Project.ID] = "en"
			}
		}
	}

	var payload []kitsu.MessagePayload
	sentCount := 0
	failedCount := 0
	for i := 0; i < len(data); i++ {
		payload = append(payload, data[i])

		if (i+1)%conf.Discord.EmbedsPerRequests == 0 || (i+1)%len(data) == 0 {
			if conf.Log {
				log.Printf("Sending bunch of messages: " + strconv.Itoa(len(payload)))
			}

			newResults := discord.SendMessageBunch(conf, payload, webhookURL, previousMessageIDs, previousWebhookURLs, previousThreadIDs, projectNotificationLanguages, db)

			// 譁ｰ縺励＞繝｡繝・そ繝ｼ繧ｸID繝ｻ繧ｹ繝ｬ繝・ラID繧奪B縺ｫ菫晏ｭ・+ 逶｣譟ｻ繝ｭ繧ｰ險倬鹸
			for taskID, res := range newResults {
				var task kitsu.MessagePayload
				for _, p := range payload {
					if p.Task.ID == taskID {
						task = p
						break
					}
				}
				// 逶｣譟ｻ繝ｭ繧ｰ: 謌仙粥繝ｻ螟ｱ謨怜撫繧上★險倬鹸
				projectGuildID := ""
				if projectRow := model.FindProjectByKitsuID(db, task.Project.ID); projectRow != nil {
					projectGuildID = projectRow.DiscordGuildID
				}
				model.WriteAuditLog(db, model.AuditLog{
					TaskID:       taskID,
					ProjectID:    task.Project.ID,
					ProjectName:  task.Project.Name,
					GuildID:      projectGuildID,
					EntityName:   task.Entity.Name,
					TaskType:     task.TaskType.TaskType.Name,
					NewStatus:    task.TaskStatus.TaskStatus.ShortName,
					DiscordMsgID: res.MessageID,
					WebhookURL:   webhookURL,
					Success:      res.MessageID != "",
				})
				if res.MessageID != "" {
					sentCount++
					model.UpdateTaskWithDiscord(
						db,
						taskID,
						task.Task.Task.UpdatedAt,
						task.TaskStatus.TaskStatus.ShortName,
						task.LatestComment.Comment.ID,
						task.LatestComment.Comment.UpdatedAt,
						res.MessageID,
						webhookURL,
						res.ThreadID,
					)
				} else {
					failedCount++
				}
			}

			payload = nil
			rl.Wait()
		}
	}

	if sentCount > 0 || failedCount > 0 {
		slog.Info("Notification send result",
			"routeSource", routeSource,
			"routeLabels", routeLabels,
			"sentCount", sentCount,
			"failedCount", failedCount)
		errSummary := ""
		if failedCount > 0 {
			errSummary = fmt.Sprintf("route=%s failed=%d", routeSource, failedCount)
		}
		setup.Stats.RecordSend(sentCount, failedCount, webhookURL, errSummary)
	}

	return data
}

// 笏笏笏 Kitsu 隱崎ｨｼ諠・ｱ繝倥Ν繝代・ 笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏
// Prefer DB settings, then environment variables, then conf.toml.
func getKitsuCreds(db *gorm.DB, conf config.Config) (hostname, email, password string) {
	hostname = model.GetSetting(db, "kitsu.hostname")
	if hostname == "" {
		hostname = os.Getenv("KITSU_HOSTNAME")
	}
	if hostname == "" {
		hostname = conf.Kitsu.Hostname
	}
	email = model.GetSetting(db, setup.RuntimeKitsuEmailSettingKey)
	if email == "" {
		email = os.Getenv(setup.RuntimeKitsuEmailEnv)
	}
	if email == "" {
		email = conf.Kitsu.Email
	}
	password = os.Getenv(setup.RuntimeKitsuPasswordEnv)
	if password == "" {
		password = conf.Kitsu.Password
	}
	if hostname != "" && !strings.HasSuffix(hostname, "/") {
		hostname += "/"
	}
	return
}

func getDiscordSettings(db *gorm.DB, conf config.Config) (botToken, guildID, webhookURL string) {
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		botToken = conf.Discord.BotToken
	}
	guildID = model.GetSetting(db, "discord.guildID")
	if guildID == "" {
		guildID = os.Getenv("DISCORD_GUILD_ID")
	}
	if guildID == "" {
		guildID = conf.Discord.GuildID
	}
	webhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		webhookURL = conf.Discord.WebhookURL
	}
	return
}

// 笏笏笏 繝昴・繝ｪ繝ｳ繧ｰ謗剃ｻ門宛蠕｡ 笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏笏
// pollMu prevents overlapping polling cycles in the same process.
var pollMu sync.Mutex

func runOnePoll(conf config.Config, db *gorm.DB) {
	if !pollMu.TryLock() {
		slog.Warn("Previous poll still running; skipping this cycle to prevent duplicate Discord messages")
		return
	}
	defer pollMu.Unlock()

	// 繧｢繧ｯ繝・ぅ繝悶↑繧ｻ繝・す繝ｧ繝ｳ・医Ο繧ｰ繧､繝ｳ荳ｭ縺ｮ邂｡逅・・ｼ峨・ Kitsu JWT 繧貞━蜈医＠縺ｦ菴ｿ逕ｨ縺吶ｋ縲・	// 繝ｭ繧ｰ繧｢繧ｦ繝・or 繧ｻ繝・す繝ｧ繝ｳ譛滄剞蛻・ｌ・・h・牙ｾ後・譌｢蟄倥・迺ｰ蠅・､画焚繝医・繧ｯ繝ｳ縺ｮ縺ｾ縺ｾ繝昴・繝ｪ繝ｳ繧ｰ繧堤ｶ咏ｶ壹・
	kitsuResponse := MakeKitsuResponse(conf)
	if conf.Log {
		slog.Info("Done MakeKitsuResponse")
	}
	taskCount := len(kitsuResponse)
	FilterTasks(kitsuResponse, conf, db)
	if conf.Log {
		slog.Info("Done FilterTasks")
	}
	setup.Stats.RecordPoll(taskCount)
}

func configureSQLite(db *gorm.DB) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)
	sqlDB.SetConnMaxIdleTime(0)

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA foreign_keys=ON;",
	}
	for _, pragma := range pragmas {
		if err := db.Exec(pragma).Error; err != nil {
			return nil, err
		}
	}

	var journalMode string
	var synchronous int
	var busyTimeout int
	var foreignKeys int
	if err := sqlDB.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		return nil, err
	}
	if err := sqlDB.QueryRow("PRAGMA synchronous;").Scan(&synchronous); err != nil {
		return nil, err
	}
	if err := sqlDB.QueryRow("PRAGMA busy_timeout;").Scan(&busyTimeout); err != nil {
		return nil, err
	}
	if err := sqlDB.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys); err != nil {
		return nil, err
	}

	slog.Info("SQLite pragmas configured",
		"journalMode", journalMode,
		"synchronous", synchronous,
		"busyTimeoutMs", busyTimeout,
		"foreignKeys", foreignKeys,
		"maxOpenConns", 1)

	return sqlDB, nil
}

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		f.SetTemplate("[{{datetime}}] [{{level}}] [{{caller}}]\t{{message}} {{data}} {{extra}}\n")
		f.ColorTheme = slog.ColorTheme
	})

	h1 := handler.MustFileHandler("./logs/all-levels.log", handler.WithLogLevels(slog.AllLevels))
	slog.PushHandler(h1)

	// Set log level based on APP_ENV (production = INFO, development = DEBUG)
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		slog.SetLogLevel(slog.InfoLevel)
		slog.Info("App started", "env", appEnv, "log_level", "INFO")
	} else {
		slog.SetLogLevel(slog.DebugLevel)
		slog.Debug("App started", "env", appEnv, "log_level", "DEBUG")
	}

	start := time.Now()

	conf := config.Read()

	if issues := conf.Validate(); len(issues) > 0 {
		for _, msg := range issues {
			if strings.HasPrefix(msg, "[FATAL]") {
				slog.Error("config validation: " + msg)
			} else {
				slog.Warn("config validation: " + msg)
			}
		}
	}

	if conf.Debug {
		os.Setenv("Debug", "true")
	}

	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		slog.Fatal("failed to connect database")
		os.Exit(1)
	}
	sqlDB, err := configureSQLite(db)
	if err != nil {
		slog.Fatal("failed to configure sqlite", "err", err)
		os.Exit(1)
	}
	// Remove the legacy single-column unique index before the composite migration.
	db.Exec("DROP INDEX IF EXISTS idx_checker_maps_task_type")
	db.AutoMigrate(
		&model.Task{},
		&model.Project{},
		&model.ProjectWebhook{},
		&model.UserMap{},
		&model.CheckerMap{},
		&model.Setting{},
		&model.AuditLog{},
		&model.ProjectUserMap{},
		&model.ProjectCheckerMap{},
		&model.ProjectSetting{},
	)
	model.PurgeLegacySensitiveData(db)

	setup.SeedFromConfig(db, conf)
	_, seedGuildID, _ := getDiscordSettings(db, conf)
	model.SeedProjectGuildFallback(db, seedGuildID)

	discord.UserMapResolver = func(projectID, kitsuName, kitsuEmail string) string {
		return model.GetUserMapForProject(db, projectID, kitsuName, kitsuEmail)
	}
	discord.CheckerResolver = func(projectID, taskType string) []string {
		return model.GetCheckerForProject(db, projectID, taskType)
	}
	discord.GoogleDriveURLResolver = func(projectID string) string {
		return model.GetProjectStorageURL(db, projectID)
	}

	if conf.Log {
		slog.Info("Connected to database in %s", time.Since(start))

		if _, err := os.Stat("./dump"); os.IsNotExist(err) {
			err := os.Mkdir("./dump", os.ModeDir)
			if err != nil {
				slog.Fatal("failed to create dir")
				os.Exit(1)
			}
		}
	}

	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	kitsuHostname, kitsuEmail, kitsuPassword := getKitsuCreds(db, conf)
	os.Setenv("KITSU_HOSTNAME", kitsuHostname)

	token := basicauth.AuthForJWTToken(kitsuHostname+"api/auth/login", kitsuEmail, kitsuPassword)
	if token == "" {
		slog.Fatal("Initial Kitsu authentication failed 窶・check hostname/email/password in conf.toml or /bot/admin/bot")
		os.Exit(1)
	}
	os.Setenv("KitsuJWTToken", token)
	if conf.Log {
		slog.Info("Connected to Kitsu in %s", time.Since(start))
	}

	if len(conf.Discord.Productions) > 0 || len(conf.Discord.TaskTypeWebhooks) > 0 {
		slog.Warn("conf.toml routing (Productions/TaskTypeWebhooks) is deprecated — manage channel assignments via Admin UI (/bot/setup) instead")
	}

	c.AddFunc("@every 1h", func() {
		// 繧｢繧ｯ繝・ぅ繝悶↑繧ｻ繝・す繝ｧ繝ｳ・医Ο繧ｰ繧､繝ｳ荳ｭ縺ｮ邂｡逅・・ｼ峨・ JWT 縺後≠繧句ｴ蜷医・縺昴ｌ繧剃ｽｿ縺・・		// 1h 縺斐→縺ｮ Basic 隱崎ｨｼ繝ｪ繝輔Ξ繝・す繝･縺ｯ繧ｹ繧ｭ繝・・縺吶ｋ・・WT 縺ｯ繧ｻ繝・す繝ｧ繝ｳ譛牙柑譛滄俣荳ｭ縺ｯ譛牙柑・峨・		// Runtime refresh no longer reuses admin session tokens.
		// 繧｢繧ｯ繝・ぅ繝悶↑繧ｻ繝・す繝ｧ繝ｳ縺後↑縺代ｌ縺ｰ DB 縺ｮ隱崎ｨｼ諠・ｱ縺ｧ繝ｪ繝輔Ξ繝・す繝･
		h, e, p := getKitsuCreds(db, conf)
		os.Setenv("KITSU_HOSTNAME", h) // hostname 縺悟､峨ｏ縺｣縺溷ｴ蜷医ｂ蜿肴丐
		newToken := basicauth.AuthForJWTToken(h+"api/auth/login", e, p)
		if newToken == "" {
			slog.Warn("Kitsu token refresh failed 窶・keeping previous token until next cycle")
			return
		}
		os.Setenv("KitsuJWTToken", newToken)
		if conf.Log {
			slog.Info("Got new Kitsu token via stored credentials")
		}
	})

	// Initial polling runs in the background so /health can respond during Discord backoff.
	go func() {
		runOnePoll(conf, db)
		if conf.Log {
			slog.Info("Done initial poll", "duration", time.Since(start).String())
		}
	}()

	c.AddFunc("@every "+strconv.Itoa(conf.Kitsu.RequestInterval)+"m", func() {
		runOnePoll(conf, db)
	})

	c.AddFunc("0 3 * * *", func() {
		deleted := model.PurgeOldAuditLogs(db, 90)
		if deleted > 0 {
			slog.Info("audit log purge", "deleted_rows", deleted)
		}
	})

	// HTTP 繧ｵ繝ｼ繝舌・: 繝倥Ν繧ｹ繝√ぉ繝・け + 繝励Ο繧ｸ繧ｧ繧ｯ繝医そ繝・ヨ繧｢繝・・ UI + 邂｡逅・判髱｢
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	setupHandler := func(w http.ResponseWriter, r *http.Request) {
		kitsuHost, _, _ := getKitsuCreds(db, conf)
		botToken, fallbackGuildID, _ := getDiscordSettings(db, conf)
		setup.Handler(kitsuHost, fallbackGuildID, botToken, db)(w, r)
	}

	setupCredsFunc := func() (string, string, string, string) {
		h, _, _ := getKitsuCreds(db, conf)
		tok, gid, wh := getDiscordSettings(db, conf)
		return h, tok, gid, wh
	}

	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		h, _, _ := getKitsuCreds(db, conf)
		setup.LoginHandler(h)(w, r)
	}
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/bot/login", loginHandler)

	mux.HandleFunc("/logout", setup.LogoutHandler())
	mux.HandleFunc("/bot/logout", setup.LogoutHandler())

	mux.HandleFunc("/setup", setup.RequireSession(setupHandler))
	mux.HandleFunc("/bot/setup", setup.RequireSession(setupHandler))

	mux.HandleFunc("/setup-wizard", setup.RequireSession(setup.WizardHandler(db, setupCredsFunc)))
	mux.HandleFunc("/bot/setup-wizard", setup.RequireSession(setup.WizardHandler(db, setupCredsFunc)))

	// Setup diagnostic JSON API — registered under both root and /bot prefix.
	setupAPIRoutes := func(prefix string) {
		mux.HandleFunc(prefix+"/api/setup/status", setup.RequireSession(setup.SetupStatusHandler(
			db, conf.Kitsu.RequestInterval, setupCredsFunc,
		)))
		mux.HandleFunc(prefix+"/api/setup/projects", setup.RequireSession(setup.ProjectsHandler(db)))
		mux.HandleFunc(prefix+"/api/setup/preview-project", setup.RequireSession(setup.PreviewProjectHandler(db, setupCredsFunc)))
		mux.HandleFunc(prefix+"/api/setup/apply-project", setup.RequireSession(setup.ApplyProjectHandler(db, setupCredsFunc)))
		mux.HandleFunc(prefix+"/api/setup/test-kitsu", setup.RequireSession(setup.TestKitsuHandler(db)))
		mux.HandleFunc(prefix+"/api/setup/test-discord", setup.RequireSession(setup.TestDiscordHandler(db)))
		mux.HandleFunc(prefix+"/api/setup/test-notification", setup.RequireSession(setup.TestNotificationHandler(db, setupCredsFunc)))
		mux.HandleFunc(prefix+"/api/setup/mapping", setup.RequireSession(setup.MappingStateHandler(db)))
		mux.HandleFunc(prefix+"/api/setup/mapping/users", setup.RequireSession(setup.SaveUserMappingHandler(db)))
		mux.HandleFunc(prefix+"/api/setup/mapping/checkers", setup.RequireSession(setup.SaveCheckerMappingHandler(db)))
	}
	setupAPIRoutes("")
	setupAPIRoutes("/bot")

	registerAdminRoutes := func(prefix string) {
		mux.HandleFunc(prefix+"/admin", setup.RequireSession(setup.AdminIndex(db)))
		mux.HandleFunc(prefix+"/admin/users", setup.RequireSession(func(w http.ResponseWriter, r *http.Request) {
			h, _, _ := getKitsuCreds(db, conf)
			setup.UsersHandler(db, h)(w, r)
		}))
		mux.HandleFunc(prefix+"/admin/checkers", setup.RequireSession(func(w http.ResponseWriter, r *http.Request) {
			h, _, _ := getKitsuCreds(db, conf)
			setup.CheckersHandler(db, h)(w, r)
		}))
		mux.HandleFunc(prefix+"/admin/drive", setup.RequireSession(setup.DriveHandler(db)))
		// BotHandler: Kitsu 險ｭ螳壻ｿ晏ｭ伜ｾ後↓蜊ｳ譎ょ・謗･邯壹☆繧九さ繝ｼ繝ｫ繝舌ャ繧ｯ繧呈ｳｨ蜈･
		kitsuReconnect := func() {
			h, e, p := getKitsuCreds(db, conf)
			os.Setenv("KITSU_HOSTNAME", h)
			newToken := basicauth.AuthForJWTToken(h+"api/auth/login", e, p)
			if newToken != "" {
				os.Setenv("KitsuJWTToken", newToken)
				slog.Info("Kitsu reconnected via admin UI", "hostname", h, "email", e)
			} else {
				slog.Warn("Kitsu reconnect failed 窶・check credentials in /bot/admin/bot")
			}
		}
		mux.HandleFunc(prefix+"/admin/bot", setup.RequireSession(setup.BotHandler(db, kitsuReconnect)))
		mux.HandleFunc(prefix+"/admin/setup", setup.RequireSession(setup.SetupManualHandler(db, setupCredsFunc)))
		mux.HandleFunc(prefix+"/admin/projects", setup.RequireSession(func(w http.ResponseWriter, r *http.Request) {
			_, fallbackGuildID, _ := getDiscordSettings(db, conf)
			setup.AdminProjectsHandler(db, fallbackGuildID)(w, r)
		}))
		mux.HandleFunc(prefix+"/admin/audit", setup.RequireSession(setup.AuditLogHandler(db)))
		mux.HandleFunc(prefix+"/admin/health", setup.RequireSession(setup.HealthHandler(db)))
		mux.HandleFunc(prefix+"/admin/diagnostics", setup.RequireSession(setup.DiagnosticsHandler(db, func() (string, string, string, string) {
			h, _, _ := getKitsuCreds(db, conf)
			tok, gid, wh := getDiscordSettings(db, conf)
			return h, tok, gid, wh
		})))
	}
	registerAdminRoutes("")
	registerAdminRoutes("/bot")

	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs.html")
	})
	mux.HandleFunc("/site.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		http.ServeFile(w, r, "./site.jsx")
	})
	mux.Handle("/diagrams/", http.StripPrefix("/diagrams/", http.FileServer(http.Dir("./diagrams/"))))

	server := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}
	go func() {
		slog.Info("HTTP server listening on :8090  (/health, /login, /setup, /admin/*, /docs)")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", "err", err)
		}
	}()

	c.Start()

	sigCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	<-sigCtx.Done()
	slog.Info("Shutdown signal received; stopping HTTP server, cron, and SQLite")

	stopCtx := c.Stop()
	select {
	case <-stopCtx.Done():
		slog.Info("Cron stopped cleanly")
	case <-time.After(10 * time.Second):
		slog.Warn("Cron stop timed out")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Warn("HTTP server shutdown had an error", "err", err)
	} else {
		slog.Info("HTTP server shutdown complete")
	}

	if err := sqlDB.Close(); err != nil {
		slog.Warn("SQLite close failed", "err", err)
	} else {
		slog.Info("SQLite connection closed cleanly")
	}
}
