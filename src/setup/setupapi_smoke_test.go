package setup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/src/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&model.Project{}, &model.ProjectWebhook{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func TestSetupStatusHandlerSmoke(t *testing.T) {
	db := newTestDB(t)
	Stats = &RuntimeStats{
		StartTime:       time.Now(),
		webhookFailures: make(map[string]int64),
		webhookLastErr:  make(map[string]string),
		webhookLastSend: make(map[string]time.Time),
	}
	Stats.RecordPoll(7)

	handler := SetupStatusHandler(db, 60, func() (string, string, string, string) {
		return "", "", "", ""
	})

	req := httptest.NewRequest(http.MethodGet, "/api/setup/status", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp SetupStatusResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.SetupComplete {
		t.Fatal("expected setup_complete=false for empty credentials")
	}
	if !resp.Poller.Running {
		t.Fatal("expected poller status to reflect recorded poll activity")
	}
	if resp.Kitsu.Error == nil || resp.Discord.Error == nil {
		t.Fatalf("expected missing-config errors, got %+v", resp)
	}
}

func TestBuildProjectStatusAggregatesChannelsAndWebhooks(t *testing.T) {
	db := newTestDB(t)

	if err := model.CreateProject(db, "p1", "Project One", "cg", "cat-1", "ja"); err != nil {
		t.Fatalf("create project: %v", err)
	}
	if err := model.CreateProject(db, "p2", "Project Two", "cg", "cat-2", "en"); err != nil {
		t.Fatalf("create project: %v", err)
	}
	for _, row := range []struct {
		projectID string
		channel   string
		taskType  string
	}{
		{"p1", "anim", "Animation"},
		{"p1", "general", "*"},
		{"p2", "comp", "Compositing"},
	} {
		if err := model.CreateProjectWebhook(db, row.projectID, row.channel, row.taskType, "https://discord.invalid/"+row.channel, row.channel+"-id"); err != nil {
			t.Fatalf("create webhook: %v", err)
		}
	}

	status := buildProjectStatus(db)
	if !status.Selected {
		t.Fatal("expected selected=true when projects exist")
	}
	if status.ChannelCount != 3 {
		t.Fatalf("expected 3 channels, got %d", status.ChannelCount)
	}
	if status.WebhookCount != 3 {
		t.Fatalf("expected 3 webhooks, got %d", status.WebhookCount)
	}
}
