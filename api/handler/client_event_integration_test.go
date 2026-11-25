package handler

import (
	"context"
	"geekai/core/types"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// Run the real desktop sender against GeekAI; only the database is mocked.
func TestThinkBuddyClientEventIntegration(t *testing.T) {
	root := os.Getenv("THINKBUDDY_TEST_ROOT")
	if root == "" {
		t.Skip("set THINKBUDDY_TEST_ROOT to run the real desktop event sender")
	}
	python := os.Getenv("THINKBUDDY_TEST_PYTHON")
	if python == "" {
		python = filepath.Join(root, ".venv", "Scripts", "python.exe")
	}
	engine, mock := clientEventRouter(t, types.AppConfig{})
	server := httptest.NewServer(engine)
	defer server.Close()
	for _, name := range []string{"app_start", "login", "error_backend", "app_exit"} {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `client_event`").WithArgs(
			nil, "device-1", "session-1", "1.2.3", "Windows 11", "127.0.0.1", name,
			`{"source": "desktop", "flag": false, "integer": 9007199254740993}`,
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, python, "-B", "-c", `
import os
import sys
from pathlib import Path

sys.path.insert(0, str(Path(sys.argv[1]) / "desktop-shell"))
os.environ["THINKBUDDY_API_BASE"] = sys.argv[2]
from desktop.auth.config import telemetry_url
from desktop.telemetry_reporter import TelemetryReporter

endpoint = telemetry_url()
assert endpoint == sys.argv[2] + "/api/events", endpoint
reporter = TelemetryReporter(endpoint, {
    "machine_id": "device-1", "session_id": "session-1",
    "app_ver": "1.2.3", "os": "Windows 11",
}, access_token=lambda: "")
props = {"source": "desktop", "flag": False, "integer": 9007199254740993}
for event in ("app_start", "login", "error_backend"):
    reporter._post_sync(event, props, "", 1792576514)
    assert reporter._fails == 0, event + " was rejected"
reporter.flush_exit("app_exit", props)
assert reporter._fails == 0, "app_exit was rejected"
reporter.close()
`, root, server.URL)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("ThinkBuddy event sender failed: %v\n%s", err, output)
	}
}
