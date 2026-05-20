package log

import (
	"os"
	"regexp"
	"strings"
)

// RedactSecrets replaces sensitive information with [REDACTED]
func RedactSecrets(message string) string {
	if message == "" {
		return message
	}

	// Discord Bot Token pattern: MTxxxxxx.Gxxxx...
	botTokenPattern := regexp.MustCompile(`MT[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)
	message = botTokenPattern.ReplaceAllString(message, "[REDACTED-TOKEN]")

	// Discord Webhook URL pattern: https://discord.com/api/webhooks/xxx/xxx
	webhookPattern := regexp.MustCompile(`webhooks/[0-9]+/[A-Za-z0-9_-]+`)
	message = webhookPattern.ReplaceAllString(message, "webhooks/[REDACTED]/[REDACTED]")

	// Environment variable values (after = or : in logs)
	envVars := []string{
		"KITSU_RUNTIME_PASSWORD",
		"KITSU_PASSWORD",
		"DISCORD_BOT_TOKEN",
		"DISCORD_WEBHOOK_URL",
		"FB_PASSWORD",
		"ADMIN_PASSWORD",
	}
	for _, envVar := range envVars {
		val := os.Getenv(envVar)
		if val != "" && len(val) > 3 {
			message = strings.ReplaceAll(message, val, "[REDACTED]")
		}
	}

	return message
}
