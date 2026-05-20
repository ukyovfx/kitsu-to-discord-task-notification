package log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestRedactSecrets_RuntimeEnvUpdate verifies that secrets set via os.Setenv() AFTER
// startup are still redacted. This is critical because KitsuSync updates DISCORD_BOT_TOKEN
// and Kitsu credentials at runtime via os.Setenv.
func TestRedactSecrets_RuntimeEnvUpdate(t *testing.T) {
	const testVar = "DISCORD_BOT_TOKEN"
	const secret = "super-secret-runtime-token-xyz"

	// Ensure clean state at start and restore after test
	original := os.Getenv(testVar)
	defer os.Setenv(testVar, original)

	// Secret not yet set — should pass through
	result := RedactSecrets("token=" + secret)
	if strings.Contains(result, "[REDACTED]") {
		t.Errorf("should NOT redact before os.Setenv, got: %q", result)
	}

	// Set the secret at runtime (simulates what setupapi.go does)
	os.Setenv(testVar, secret)

	// Now it must be redacted
	result = RedactSecrets("token=" + secret)
	if strings.Contains(result, secret) {
		t.Errorf("should have redacted runtime secret, got: %q", result)
	}
	if !strings.Contains(result, "[REDACTED]") {
		t.Errorf("expected [REDACTED] in result, got: %q", result)
	}
}

func TestRedactingWriter(t *testing.T) {
	const webhookLog = `[2024-01-01] [INFO] [main.go:42]	sending webhook url=https://discord.com/api/webhooks/987654321/TopSecretToken`

	buf := &bytes.Buffer{}
	w := NewRedactingWriter(buf)

	n, err := w.Write([]byte(webhookLog))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	// Write must report the original byte count, not the redacted length
	if n != len(webhookLog) {
		t.Errorf("expected n=%d (original len), got %d", len(webhookLog), n)
	}

	result := buf.String()
	if strings.Contains(result, "987654321") {
		t.Errorf("webhook ID must be redacted, got: %q", result)
	}
	if strings.Contains(result, "TopSecretToken") {
		t.Errorf("webhook token must be redacted, got: %q", result)
	}
	if !strings.Contains(result, "webhooks/[REDACTED]/[REDACTED]") {
		t.Errorf("expected redaction placeholder, got: %q", result)
	}
}

func TestRedactingWriter_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewRedactingWriter(buf)
	n, err := w.Write([]byte(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 bytes, got %d", n)
	}
}

type chunkWriter struct {
	chunk int
	buf   bytes.Buffer
}

func (w *chunkWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if w.chunk <= 0 || w.chunk >= len(p) {
		return w.buf.Write(p)
	}
	return w.buf.Write(p[:w.chunk])
}

type partialErrorWriter struct {
	chunk int
	err   error
	buf   bytes.Buffer
}

func (w *partialErrorWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, w.err
	}
	if w.chunk <= 0 || w.chunk > len(p) {
		w.chunk = len(p)
	}
	_, _ = w.buf.Write(p[:w.chunk])
	return w.chunk, w.err
}

func TestRedactingWriter_PartialWritesAreCompleted(t *testing.T) {
	const secret = "super-secret-runtime-token-xyz"
	t.Setenv("DISCORD_BOT_TOKEN", secret)

	underlying := &chunkWriter{chunk: 5}
	w := NewRedactingWriter(underlying)

	n, err := w.Write([]byte("token=" + secret))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != len("token="+secret) {
		t.Fatalf("expected n=%d, got %d", len("token="+secret), n)
	}
	if strings.Contains(underlying.buf.String(), secret) {
		t.Fatalf("secret leaked in partial write result: %q", underlying.buf.String())
	}
	if !strings.Contains(underlying.buf.String(), "[REDACTED]") {
		t.Fatalf("expected redacted output, got %q", underlying.buf.String())
	}
}

func TestRedactingWriter_PartialWriteErrorReportsConsumedInput(t *testing.T) {
	const secret = "super-secret-runtime-token-xyz"
	t.Setenv("DISCORD_BOT_TOKEN", secret)
	wantErr := os.ErrClosed
	underlying := &partialErrorWriter{chunk: 4, err: wantErr}
	w := NewRedactingWriter(underlying)

	n, err := w.Write([]byte("token=" + secret))
	if !strings.Contains(underlying.buf.String(), "toke") {
		t.Fatalf("expected partial output to be written, got %q", underlying.buf.String())
	}
	if err != wantErr {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
	if n != len("token="+secret) {
		t.Fatalf("expected original byte count %d, got %d", len("token="+secret), n)
	}
}
