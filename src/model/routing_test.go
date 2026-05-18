package model

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newRoutingTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&ProjectWebhook{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func TestFindWebhookForTaskPrefersExactTaskTypeOverFallback(t *testing.T) {
	db := newRoutingTestDB(t)

	if err := CreateProjectWebhook(db, "project-1", "general", "*", "https://discord.invalid/general", "c1"); err != nil {
		t.Fatalf("create fallback webhook: %v", err)
	}
	if err := CreateProjectWebhook(db, "project-1", "comp", "Compositing", "https://discord.invalid/comp", "c2"); err != nil {
		t.Fatalf("create exact webhook: %v", err)
	}

	got := FindWebhookForTask(db, "project-1", "Compositing")
	if got != "https://discord.invalid/comp" {
		t.Fatalf("expected exact task-type webhook, got %q", got)
	}
}

func TestFindWebhookForTaskFallsBackToProjectWildcard(t *testing.T) {
	db := newRoutingTestDB(t)

	if err := CreateProjectWebhook(db, "project-1", "general", "*", "https://discord.invalid/general", "c1"); err != nil {
		t.Fatalf("create fallback webhook: %v", err)
	}

	got := FindWebhookForTask(db, "project-1", "Lighting")
	if got != "https://discord.invalid/general" {
		t.Fatalf("expected fallback webhook, got %q", got)
	}
}

func TestFindWebhookForTaskReturnsEmptyWhenUnconfigured(t *testing.T) {
	db := newRoutingTestDB(t)

	got := FindWebhookForTask(db, "project-1", "Lighting")
	if got != "" {
		t.Fatalf("expected no webhook, got %q", got)
	}
}
