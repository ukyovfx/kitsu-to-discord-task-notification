package setup

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

type diagCheck struct {
	Label  string
	Status string // "ok" | "warn" | "fail"
	Detail string
	Fix    string
}

// DiagnosticsHandler runs pre-flight environment checks on demand.
// Pass a refreshCreds func so credentials are always current (read from DB/env at request time).
func DiagnosticsHandler(db *gorm.DB, refreshCreds func() (kitsuHost, botToken, guildID, webhookURL string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		lang := currentLang(r)

		kitsuHost, botToken, guildID, webhookURL := refreshCreds()
		checks := runDiagnostics(kitsuHost, botToken, guildID, webhookURL, db)

		allOK := true
		anyFail := false
		for _, c := range checks {
			if c.Status != "ok" {
				allOK = false
			}
			if c.Status == "fail" {
				anyFail = true
			}
		}

		var summary string
		if allOK {
			summary = `<div class="diag-banner ok">` + esc(t(lang, "すべてのチェックが通過しました。セットアップ準備完了です。", "All checks passed. Ready to run project setup.")) + `</div>`
		} else if anyFail {
			summary = `<div class="diag-banner fail">` + esc(t(lang, "一部のチェックが失敗しました。セットアップ前に修正してください。", "Some checks failed. Please resolve them before running setup.")) + `</div>`
		} else {
			summary = `<div class="diag-banner warn">` + esc(t(lang, "警告があります。確認してから進んでください。", "Warnings found. Review before continuing.")) + `</div>`
		}

		var rows strings.Builder
		for _, c := range checks {
			icon := "✅"
			rowClass := "diag-ok"
			if c.Status == "warn" {
				icon = "⚠️"
				rowClass = "diag-warn"
			} else if c.Status == "fail" {
				icon = "❌"
				rowClass = "diag-fail"
			}
			fix := ""
			if c.Fix != "" {
				fix = `<div class="diag-fix">` + html.EscapeString(t(lang, "対処: ", "Fix: ")+c.Fix) + `</div>`
			}
			rows.WriteString(fmt.Sprintf(`<tr class="%s"><td class="diag-icon">%s</td><td><strong>%s</strong><div class="diag-detail">%s</div>%s</td></tr>`,
				rowClass, icon, html.EscapeString(c.Label), html.EscapeString(c.Detail), fix))
		}

		rerunURL := withLang("/bot/admin/diagnostics", r)
		body := fmt.Sprintf(`
<style>
.diag-banner{padding:14px 20px;border-radius:var(--radius-md);margin-bottom:18px;font-weight:600}
.diag-banner.ok{background:rgba(142,207,139,.18);border:1px solid rgba(142,207,139,.4);color:#8ecf8b}
.diag-banner.warn{background:rgba(255,200,80,.12);border:1px solid rgba(255,200,80,.35);color:#ffc850}
.diag-banner.fail{background:rgba(255,106,80,.14);border:1px solid rgba(255,106,80,.38);color:#ff6a50}
.diag-ok td{color:var(--text)}
.diag-warn td{color:#ffc850}
.diag-fail td{color:#ff6a50}
.diag-icon{font-size:1.3rem;padding-right:14px;white-space:nowrap}
.diag-detail{color:var(--muted);font-size:.875rem;margin-top:3px}
.diag-fix{margin-top:6px;padding:6px 10px;border-radius:8px;background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);font-size:.82rem;color:var(--muted-2)}
.diag-fix::before{content:"→ "}
</style>
%s
<div class="section-card glass">
  <div class="table-wrap">
    <table><tbody>%s</tbody></table>
  </div>
</div>
<div class="button-row">
  <a class="btn" href="%s">%s</a>
  <a class="btn-ghost" href="%s">%s</a>
</div>`,
			summary,
			rows.String(),
			rerunURL, esc(t(lang, "再チェック", "Re-run checks")),
			withLang("/bot/admin", r), esc(t(lang, "管理画面へ", "Back to Admin")),
		)

		fmt.Fprint(w, adminPage(lang, t(lang, "環境診断", "Environment Diagnostics"), r, body))
	}
}

func runDiagnostics(kitsuHost, botToken, guildID, webhookURL string, db *gorm.DB) []diagCheck {
	client := &http.Client{Timeout: 8 * time.Second}
	var checks []diagCheck

	// 1. Kitsu hostname configured
	if kitsuHost == "" {
		checks = append(checks, diagCheck{
			Label:  "Kitsu hostname",
			Status: "fail",
			Detail: "No hostname configured.",
			Fix:    "Set KITSU_HOSTNAME in .env or configure via /bot/admin/bot.",
		})
	} else {
		checks = append(checks, diagCheck{
			Label:  "Kitsu hostname",
			Status: "ok",
			Detail: kitsuHost,
		})
	}

	// 2. Kitsu server reachable
	if kitsuHost != "" {
		pingURL := strings.TrimRight(kitsuHost, "/") + "/api/"
		resp, err := client.Get(pingURL)
		if err != nil {
			checks = append(checks, diagCheck{
				Label:  "Kitsu server reachable",
				Status: "fail",
				Detail: "HTTP request failed: " + err.Error(),
				Fix:    "Ensure Kitsu is running and KITSU_HOSTNAME is correct (include http:// or https://).",
			})
		} else {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				checks = append(checks, diagCheck{
					Label:  "Kitsu server reachable",
					Status: "ok",
					Detail: fmt.Sprintf("HTTP %d — server responded.", resp.StatusCode),
				})
			} else {
				checks = append(checks, diagCheck{
					Label:  "Kitsu server reachable",
					Status: "warn",
					Detail: fmt.Sprintf("HTTP %d — server may be starting up.", resp.StatusCode),
					Fix:    "Check Kitsu container logs.",
				})
			}
		}
	}

	// 3. Kitsu auth valid (JWT token present)
	jwtToken := os.Getenv("KitsuJWTToken")
	if jwtToken == "" {
		checks = append(checks, diagCheck{
			Label:  "Kitsu auth",
			Status: "fail",
			Detail: "No active Kitsu session token. App may have failed to authenticate at startup.",
			Fix:    "Check container logs for 'Initial Kitsu authentication failed'. Verify credentials in /bot/admin/bot.",
		})
	} else if kitsuHost != "" {
		// Verify the token is still valid
		authURL := strings.TrimRight(kitsuHost, "/") + "/api/auth/user"
		req, err := http.NewRequest("GET", authURL, nil)
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			resp, err := client.Do(req)
			if err != nil {
				checks = append(checks, diagCheck{
					Label:  "Kitsu auth",
					Status: "warn",
					Detail: "Could not verify token: " + err.Error(),
					Fix:    "Check Kitsu server availability.",
				})
			} else {
				resp.Body.Close()
				if resp.StatusCode == 200 {
					checks = append(checks, diagCheck{
						Label:  "Kitsu auth",
						Status: "ok",
						Detail: "Session token is valid.",
					})
				} else {
					checks = append(checks, diagCheck{
						Label:  "Kitsu auth",
						Status: "fail",
						Detail: fmt.Sprintf("Token rejected (HTTP %d). Session may have expired.", resp.StatusCode),
						Fix:    "Update credentials in /bot/admin/bot and restart the container.",
					})
				}
			}
		}
	} else {
		checks = append(checks, diagCheck{
			Label:  "Kitsu auth",
			Status: "ok",
			Detail: "Session token present.",
		})
	}

	// 4. Discord bot token configured
	if strings.TrimSpace(botToken) == "" {
		checks = append(checks, diagCheck{
			Label:  "Discord bot token",
			Status: "fail",
			Detail: "No bot token configured.",
			Fix:    "Set DISCORD_BOT_TOKEN in .env or configure via /bot/admin/bot.",
		})
	} else {
		checks = append(checks, diagCheck{
			Label:  "Discord bot token",
			Status: "ok",
			Detail: "Token configured (hidden).",
		})
	}

	// 5. Discord bot token valid + get bot user ID
	botUserID := ""
	if strings.TrimSpace(botToken) != "" {
		body, status, err := botDo("GET", discordAPI+"/users/@me", nil, botToken)
		if err != nil {
			checks = append(checks, diagCheck{
				Label:  "Discord bot valid",
				Status: "fail",
				Detail: "Request failed: " + err.Error(),
				Fix:    "Check network connectivity to discord.com.",
			})
		} else if status == 200 {
			var result struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			}
			if json.Unmarshal(body, &result) == nil {
				botUserID = result.ID
				checks = append(checks, diagCheck{
					Label:  "Discord bot valid",
					Status: "ok",
					Detail: fmt.Sprintf("Bot user: %s (ID: %s)", result.Username, result.ID),
				})
			} else {
				checks = append(checks, diagCheck{
					Label:  "Discord bot valid",
					Status: "ok",
					Detail: "Bot token accepted.",
				})
			}
		} else if status == 401 {
			checks = append(checks, diagCheck{
				Label:  "Discord bot valid",
				Status: "fail",
				Detail: "Token rejected (HTTP 401 Unauthorized).",
				Fix:    "Regenerate the bot token in the Discord Developer Portal and update DISCORD_BOT_TOKEN.",
			})
		} else {
			checks = append(checks, diagCheck{
				Label:  "Discord bot valid",
				Status: "warn",
				Detail: fmt.Sprintf("Unexpected response: HTTP %d.", status),
				Fix:    "Check Discord API status at discordstatus.com.",
			})
		}
	}

	// 6. Discord Guild ID configured
	if strings.TrimSpace(guildID) == "" {
		checks = append(checks, diagCheck{
			Label:  "Discord Guild ID",
			Status: "fail",
			Detail: "No Guild ID configured.",
			Fix:    "Set DISCORD_GUILD_ID in .env. Right-click your Discord server icon → Copy Server ID (requires Developer Mode in Discord settings).",
		})
	} else {
		checks = append(checks, diagCheck{
			Label:  "Discord Guild ID",
			Status: "ok",
			Detail: "Guild ID: " + guildID,
		})
	}

	// 7. Discord guild accessible
	if strings.TrimSpace(botToken) != "" && strings.TrimSpace(guildID) != "" {
		_, status, err := botDo("GET", discordAPI+"/guilds/"+guildID, nil, botToken)
		if err != nil {
			checks = append(checks, diagCheck{
				Label:  "Discord guild accessible",
				Status: "fail",
				Detail: "Request failed: " + err.Error(),
			})
		} else if status == 200 {
			checks = append(checks, diagCheck{
				Label:  "Discord guild accessible",
				Status: "ok",
				Detail: "Bot is a member of the guild.",
			})
		} else if status == 403 {
			checks = append(checks, diagCheck{
				Label:  "Discord guild accessible",
				Status: "fail",
				Detail: "HTTP 403 — bot is not a member of this guild.",
				Fix:    "Invite the bot to your server using the OAuth2 URL from the Discord Developer Portal. Required scopes: bot. Required permissions: Manage Channels, Manage Webhooks.",
			})
		} else if status == 404 {
			checks = append(checks, diagCheck{
				Label:  "Discord guild accessible",
				Status: "fail",
				Detail: "HTTP 404 — Guild ID not found.",
				Fix:    "Verify DISCORD_GUILD_ID is correct (18-digit number). Right-click server → Copy Server ID in Discord.",
			})
		} else {
			checks = append(checks, diagCheck{
				Label:  "Discord guild accessible",
				Status: "warn",
				Detail: fmt.Sprintf("HTTP %d — unexpected response.", status),
			})
		}
	}

	// 8. Discord MANAGE_CHANNELS + MANAGE_WEBHOOKS permissions
	// Proxy check: try to list guild channels (requires VIEW_CHANNELS; 403 = no access at all)
	if strings.TrimSpace(botToken) != "" && strings.TrimSpace(guildID) != "" {
		_, status, err := botDo("GET", discordAPI+"/guilds/"+guildID+"/channels", nil, botToken)
		if err != nil {
			checks = append(checks, diagCheck{
				Label:  "Discord permissions (channels)",
				Status: "warn",
				Detail: "Could not check: " + err.Error(),
			})
		} else if status == 200 {
			checks = append(checks, diagCheck{
				Label:  "Discord permissions (channels)",
				Status: "ok",
				Detail: "Bot can list channels in the guild.",
			})
		} else if status == 403 {
			checks = append(checks, diagCheck{
				Label:  "Discord permissions (channels)",
				Status: "fail",
				Detail: "HTTP 403 — bot cannot list channels.",
				Fix:    "In Discord Developer Portal → OAuth2 → Add permissions: Manage Channels, Manage Webhooks. Re-invite the bot if needed.",
			})
		} else {
			checks = append(checks, diagCheck{
				Label:  "Discord permissions (channels)",
				Status: "warn",
				Detail: fmt.Sprintf("HTTP %d.", status),
			})
		}

		// MANAGE_WEBHOOKS check via member permissions
		if botUserID != "" {
			body, mStatus, mErr := botDo("GET", discordAPI+"/guilds/"+guildID+"/members/"+botUserID, nil, botToken)
			if mErr == nil && mStatus == 200 {
				var member struct {
					Permissions string `json:"permissions"`
				}
				if json.Unmarshal(body, &member) == nil && member.Permissions != "" {
					var perms uint64
					fmt.Sscanf(member.Permissions, "%d", &perms)
					const manageWebhooks = uint64(1 << 29)
					const manageChannels = uint64(1 << 4)
					if perms&manageWebhooks != 0 && perms&manageChannels != 0 {
						checks = append(checks, diagCheck{
							Label:  "Discord permissions (webhooks)",
							Status: "ok",
							Detail: "MANAGE_CHANNELS and MANAGE_WEBHOOKS confirmed.",
						})
					} else {
						missing := []string{}
						if perms&manageChannels == 0 {
							missing = append(missing, "MANAGE_CHANNELS")
						}
						if perms&manageWebhooks == 0 {
							missing = append(missing, "MANAGE_WEBHOOKS")
						}
						checks = append(checks, diagCheck{
							Label:  "Discord permissions (webhooks)",
							Status: "fail",
							Detail: "Missing permissions: " + strings.Join(missing, ", "),
							Fix:    "In Discord Developer Portal → Bot → Grant these permissions and re-invite the bot.",
						})
					}
				} else {
					checks = append(checks, diagCheck{
						Label:  "Discord permissions (webhooks)",
						Status: "warn",
						Detail: "Could not read permission bits from member response.",
						Fix:    "Manually verify the bot has MANAGE_CHANNELS and MANAGE_WEBHOOKS in your Discord server.",
					})
				}
			} else {
				checks = append(checks, diagCheck{
					Label:  "Discord permissions (webhooks)",
					Status: "warn",
					Detail: "Could not retrieve bot member info to verify permissions.",
					Fix:    "Manually verify the bot has MANAGE_CHANNELS and MANAGE_WEBHOOKS.",
				})
			}
		}
	}

	// 9. SQLite responsive
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM sqlite_master").Scan(&count).Error; err != nil {
		checks = append(checks, diagCheck{
			Label:  "Database (SQLite)",
			Status: "fail",
			Detail: "Query failed: " + err.Error(),
			Fix:    "Check that the data/ volume is mounted and not corrupted.",
		})
	} else {
		checks = append(checks, diagCheck{
			Label:  "Database (SQLite)",
			Status: "ok",
			Detail: "Database is responsive.",
		})
	}

	// 10. data/ directory writable
	testPath := "./data/.diag_write_test"
	if err := os.WriteFile(testPath, []byte("ok"), 0600); err != nil {
		checks = append(checks, diagCheck{
			Label:  "data/ directory writable",
			Status: "fail",
			Detail: "Write test failed: " + err.Error(),
			Fix:    "Ensure the data/ directory is mounted with write permissions in docker-compose.yml.",
		})
	} else {
		os.Remove(testPath)
		checks = append(checks, diagCheck{
			Label:  "data/ directory writable",
			Status: "ok",
			Detail: "Write test passed.",
		})
	}

	// 11. Fallback webhook reachable (if configured)
	if strings.TrimSpace(webhookURL) == "" {
		checks = append(checks, diagCheck{
			Label:  "Fallback webhook",
			Status: "warn",
			Detail: "No fallback webhook URL configured.",
			Fix:    "Set DISCORD_WEBHOOK_URL in .env. This is required for unrouted task notifications.",
		})
	} else {
		req, err := http.NewRequest("GET", webhookURL, nil)
		if err != nil {
			checks = append(checks, diagCheck{
				Label:  "Fallback webhook",
				Status: "fail",
				Detail: "Invalid webhook URL: " + err.Error(),
				Fix:    "Recreate the webhook in Discord and update DISCORD_WEBHOOK_URL.",
			})
		} else {
			resp, err := client.Do(req)
			if err != nil {
				checks = append(checks, diagCheck{
					Label:  "Fallback webhook",
					Status: "fail",
					Detail: "Request failed: " + err.Error(),
					Fix:    "Check network connectivity to discord.com.",
				})
			} else {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				if resp.StatusCode == 200 {
					checks = append(checks, diagCheck{
						Label:  "Fallback webhook",
						Status: "ok",
						Detail: "Webhook is alive.",
					})
				} else if resp.StatusCode == 404 {
					checks = append(checks, diagCheck{
						Label:  "Fallback webhook",
						Status: "fail",
						Detail: "HTTP 404 — webhook has been deleted from Discord.",
						Fix:    "Recreate the webhook in the Discord channel and update DISCORD_WEBHOOK_URL in .env.",
					})
				} else {
					checks = append(checks, diagCheck{
						Label:  "Fallback webhook",
						Status: "warn",
						Detail: fmt.Sprintf("HTTP %d — unexpected response.", resp.StatusCode),
					})
				}
			}
		}
	}

	return checks
}
