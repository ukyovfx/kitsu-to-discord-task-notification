package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID               uint `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	TaskID           string
	TaskUpdatedAt    string
	TaskStatus       string
	CommentID        string
	CommentUpdatedAt string
	DiscordMessageID string // Discord に送信したメッセージID
	WebhookURL       string // 送信先 Webhook URL
	DiscordThreadID  string // Discord スレッドID（UseThreads=true 時のみ）
}

func CreateTask(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	db.Create(&Task{TaskID: taskID, TaskUpdatedAt: taskUpdatedAt, TaskStatus: taskStatus, CommentUpdatedAt: commentUpdatedAt, CommentID: commentID})
}

func UpdateTask(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	_ = db.Transaction(func(tx *gorm.DB) error {
		var rec Task
		if err := tx.Where("task_id = ?", taskID).First(&rec).Error; err != nil {
			return err
		}
		rec.TaskUpdatedAt = taskUpdatedAt
		rec.TaskStatus = taskStatus
		rec.CommentUpdatedAt = commentUpdatedAt
		rec.CommentID = commentID
		return tx.Save(&rec).Error
	})
}

func UpdateTaskWithDiscord(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt, discordMessageID, webhookURL, threadID string) {
	_ = db.Transaction(func(tx *gorm.DB) error {
		var rec Task
		if err := tx.Where("task_id = ?", taskID).First(&rec).Error; err != nil {
			return err
		}
		rec.TaskUpdatedAt = taskUpdatedAt
		rec.TaskStatus = taskStatus
		rec.CommentUpdatedAt = commentUpdatedAt
		rec.CommentID = commentID
		rec.DiscordMessageID = discordMessageID
		rec.WebhookURL = webhookURL
		if threadID != "" {
			rec.DiscordThreadID = threadID
		}
		return tx.Save(&rec).Error
	})
}

// GetStatusCounts は DB に保存されているタスクのステータス別件数を返す（日次サマリー用）
type StatusCount struct {
	TaskStatus string
	Count      int
}

func GetStatusCounts(db *gorm.DB) []StatusCount {
	var results []StatusCount
	db.Model(&Task{}).
		Select("task_status, count(*) as count").
		Where("deleted_at IS NULL").
		Group("task_status").
		Order("count desc").
		Scan(&results)
	return results
}

func FindTask(db *gorm.DB, taskID string) Task {
	var Task Task
	db.First(&Task, "task_id = ?", taskID)
	return Task
}