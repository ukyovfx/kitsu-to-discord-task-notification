package request

import (
	"app/src/utils/debug"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gookit/slog"
)

// Do is a thin HTTP wrapper.
//
// Behavior on failure: log via slog.Error and return an empty string.
// The function never calls slog.Fatal — that would crash the entire daemon
// on any transient Kitsu/Discord hiccup.
//
// The signature is preserved for now (no error return). Callers that
// pass `unmarshal` and rely on it being populated must tolerate the empty
// case (the function leaves `unmarshal` zero-initialized on failure).
func Do(token, method, url string, payload, unmarshal interface{}) string {
	var body io.ReadWriter

	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			slog.Error("request.Do: marshal payload failed", "err", err, "url", url)
			return ""
		}
		body = bytes.NewBuffer(buf)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		slog.Error("request.Do: NewRequest failed", "err", err, "url", url)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request.Do: client.Do failed", "err", err, "method", method, "url", url)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error("request.Do: read body failed", "err", err, "url", url)
		return ""
	}

	if os.Getenv("Debug") == "true" {
		debug.Info(resp, respBody)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("request.Do: non-2xx response",
			"status", resp.StatusCode,
			"method", method,
			"url", url)
		return ""
	}

	if unmarshal != nil {
		if err := json.Unmarshal(respBody, &unmarshal); err != nil {
			slog.Error("request.Do: unmarshal failed", "err", err, "url", url)
			return ""
		}
	}

	return string(respBody)
}
