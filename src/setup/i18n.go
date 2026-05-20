package setup

import (
	"net/http"
	"net/url"
	"strings"
)

func currentLang(r *http.Request) string {
	if r == nil {
		return "ja"
	}
	lang := r.URL.Query().Get("lang")
	if lang == "en" {
		return "en"
	}
	return "ja"
}

func withLang(path string, r *http.Request) string {
	return appendLang(path, currentLang(r))
}

func langURL(r *http.Request, lang string) string {
	if r == nil || r.URL == nil {
		return "/?lang=" + lang
	}
	target := r.URL.Path
	if target == "" {
		target = "/"
	}
	if raw := r.URL.RawQuery; raw != "" {
		target += "?" + raw
	}
	return appendLang(target, lang)
}

func nextLang(lang string) string {
	if lang == "en" {
		return "ja"
	}
	return "en"
}

func toggleLangURL(r *http.Request) string {
	return langURL(r, nextLang(currentLang(r)))
}

func t(lang, ja, en string) string {
	if lang == "en" {
		return en
	}
	return ja
}

func appendLang(path, lang string) string {
	if path == "" {
		path = "/"
	}
	u, err := url.Parse(path)
	if err != nil {
		if path == "" {
			path = "/"
		}
		if lang == "" {
			return path
		}
		separator := "?"
		if strings.Contains(path, "?") {
			separator = "&"
		}
		return path + separator + "lang=" + url.QueryEscape(lang)
	}
	if u.Path == "" {
		u.Path = "/"
	}
	values := u.Query()
	values.Set("lang", lang)
	u.RawQuery = values.Encode()
	return u.String()
}
