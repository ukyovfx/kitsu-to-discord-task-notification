package setup

import (
	"sync"
	"time"
)

// RuntimeStats tracks live operational metrics updated each polling cycle.
// Resets on app restart — acceptable for small-team self-hosted deployments.
type RuntimeStats struct {
	mu sync.RWMutex

	StartTime time.Time

	PollCount         int64
	LastPollTime      time.Time
	LastPollTaskCount int
	LastPollErr       string

	SendSuccessTotal int64
	SendFailureTotal int64
	LastSendTime     time.Time

	webhookFailures  map[string]int64
	webhookLastErr   map[string]string
	webhookLastSend  map[string]time.Time
}

// Stats is the singleton runtime stats instance used by main.go and the health handler.
var Stats = &RuntimeStats{
	StartTime:        time.Now(),
	webhookFailures:  make(map[string]int64),
	webhookLastErr:   make(map[string]string),
	webhookLastSend:  make(map[string]time.Time),
}

// RecordPoll updates polling stats after one cycle completes successfully.
func (s *RuntimeStats) RecordPoll(taskCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PollCount++
	s.LastPollTime = time.Now()
	s.LastPollTaskCount = taskCount
	s.LastPollErr = ""
}

// RecordPollError records a failed poll cycle without updating the success counters.
func (s *RuntimeStats) RecordPollError(errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastPollErr = errMsg
}

// LastPollError returns the error message from the most recent failed poll cycle.
// Returns "" if the last poll succeeded or no poll has run yet.
func (s *RuntimeStats) LastPollError() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.LastPollErr
}

// RecordSend updates Discord send stats for one webhook batch.
// A successful send to a webhook clears its failure counter.
func (s *RuntimeStats) RecordSend(sentCount, failedCount int, webhookURL, errSummary string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SendSuccessTotal += int64(sentCount)
	s.SendFailureTotal += int64(failedCount)
	if sentCount > 0 || failedCount > 0 {
		s.LastSendTime = time.Now()
	}
	if webhookURL == "" {
		return
	}
	if failedCount > 0 {
		s.webhookFailures[webhookURL] += int64(failedCount)
		if errSummary != "" {
			s.webhookLastErr[webhookURL] = errSummary
		}
	}
	if sentCount > 0 {
		delete(s.webhookFailures, webhookURL)
		delete(s.webhookLastErr, webhookURL)
		s.webhookLastSend[webhookURL] = time.Now()
	}
}

// WebhookInactive returns true if the webhook URL has not sent any message for
// at least threshold duration AND the app has been running at least that long.
func (s *RuntimeStats) WebhookInactive(webhookURL string, threshold time.Duration) bool {
	if time.Since(s.StartTime) < threshold {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	lastSend, ok := s.webhookLastSend[webhookURL]
	if !ok {
		return true
	}
	return time.Since(lastSend) > threshold
}

// WebhookHealthEntry describes failure state for one webhook URL since startup.
type WebhookHealthEntry struct {
	URL          string
	FailureCount int64
	LastError    string
}

// WebhookHealthList returns webhook URLs that have recorded failures since startup.
func (s *RuntimeStats) WebhookHealthList() []WebhookHealthEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []WebhookHealthEntry
	for url, count := range s.webhookFailures {
		if count > 0 {
			list = append(list, WebhookHealthEntry{
				URL:          url,
				FailureCount: count,
				LastError:    s.webhookLastErr[url],
			})
		}
	}
	return list
}

// WebhookLastSend returns the time of the last successful send for a given webhook URL.
// Returns zero time if no successful send has been recorded since startup.
func (s *RuntimeStats) WebhookLastSend(webhookURL string) time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.webhookLastSend[webhookURL]
}

// RuntimeSnapshot is a point-in-time copy of all runtime metrics.
type RuntimeSnapshot struct {
	StartTime         time.Time
	PollCount         int64
	LastPollTime      time.Time
	LastPollTaskCount int
	LastPollErr       string
	SendSuccessTotal  int64
	SendFailureTotal  int64
	LastSendTime      time.Time
}

// Snapshot returns a safe copy of current metrics.
func (s *RuntimeStats) Snapshot() RuntimeSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return RuntimeSnapshot{
		StartTime:         s.StartTime,
		PollCount:         s.PollCount,
		LastPollTime:      s.LastPollTime,
		LastPollTaskCount: s.LastPollTaskCount,
		LastPollErr:       s.LastPollErr,
		SendSuccessTotal:  s.SendSuccessTotal,
		SendFailureTotal:  s.SendFailureTotal,
		LastSendTime:      s.LastSendTime,
	}
}
