// Package config provides methods for accesing config file in TOML format
package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gookit/slog"
	"github.com/naoina/toml"
)

type Productions struct {
	Production string
	WebhookURL string
}

// UserMapEntry は Kitsu の名前と Discord ID の紐付け
type UserMapEntry struct {
	KitsuName string `toml:"kitsuName"` // Kitsu 上のフルネーム
	DiscordID string `toml:"discordID"` // Discord ユーザー ID
}

// CheckerEntry はタスクタイプ（工程）ごとのチェッカー
type CheckerEntry struct {
	TaskType  string `toml:"taskType"`  // タスクタイプ名（例: "Animation", "FX"）
	DiscordID string `toml:"discordID"` // チェッカーの Discord ID
}

// MentionConfig はメンション機能の全設定
type MentionConfig struct {
	ArtistStatuses  []string       `toml:"artistStatuses"`  // アーティストにメンションするステータス
	CheckerStatuses []string       `toml:"checkerStatuses"` // チェッカーにメンションするステータス
	UserMap         []UserMapEntry `toml:"userMap"`         // Kitsu名 → Discord ID
	Checkers        []CheckerEntry `toml:"checkers"`        // タスクタイプ → チェッカー Discord ID
}

type Config struct {
	TplPreset             string
	IgnoreMessagesDaysOld int
	SilentUpdateDB        bool
	Threads               int
	Debug                 bool
	Log                   bool
	Kitsu                 struct {
		Hostname        string
		Email           string
		Password        string
		SkipComments    bool
		RequestInterval int
	}
	Discord struct {
		EmbedsPerRequests int
		RequestsPerMinute int
		WebhookURL        string
		Productions       []Productions `toml:"productions,omitempty"`
	}
	Mention     MentionConfig `toml:"mention"`
	GoogleDrive struct {
		URL string `toml:"url"`
	} `toml:"googleDrive"`
}

// Read loads conf.toml, expanding ${VAR} placeholders from environment variables.
//
// This lets you keep secrets (Kitsu password, Discord webhook URL) out of the
// committed conf.toml — write `password = "${KITSU_PASSWORD}"` in the file,
// and the value is filled in from the environment at load time.
//
// Only the ${VAR} form is expanded — bare `$VAR` is left alone so that
// passwords containing literal `$` are preserved.
func Read() Config {
	path := "conf.toml"
	if os.Getenv("TEST") == "true" {
		path = os.Getenv("CONF_PATH")
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		slog.Fatal(err)
	}

	expanded := expandEnvBraces(string(raw))

	var config Config
	if err := toml.NewDecoder(strings.NewReader(expanded)).Decode(&config); err != nil {
		slog.Fatal(err)
	}
	return config
}

// expandEnvBraces replaces ${VAR} occurrences with the value of the environment
// variable VAR. Unset variables expand to the empty string. The bare $VAR form
// is intentionally NOT expanded — this avoids mangling values that legitimately
// contain $ (e.g. some auto-generated passwords).
func expandEnvBraces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '$' && s[i+1] == '{' {
			end := strings.IndexByte(s[i+2:], '}')
			if end < 0 {
				b.WriteString(s[i:])
				break
			}
			name := s[i+2 : i+2+end]
			b.WriteString(os.Getenv(name))
			i += 2 + end + 1
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
