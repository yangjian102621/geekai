package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"geekai/core/types"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestModelAPIStreamingFlushAndCancellation(t *testing.T) {
	h, engine, mock := modelAPITestHandler(t)
	h.RegisterRoutes(h.App.Router)
	cancelled := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"content\":\"first\"}}]}\n\n")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
		close(cancelled)
	}))
	defer upstream.Close()
	key := modelTokenPrefix + strings.Repeat("d", 43)
	expectModelAPIAuth(mock, key, true, 0)
	expectAPIModel(mock, 2)
	expectAPIPower(mock, false)
	expectAPIUpstream(mock, upstream.URL)
	api := httptest.NewServer(engine)
	defer api.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/chat/completions", strings.NewReader(`{"model":"test-model","messages":[{"role":"user","content":"hi"}],"stream":true}`))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	line, err := bufio.NewReader(response.Body).ReadString('\n')
	if err != nil || !strings.Contains(line, "first") {
		t.Fatalf("response was not flushed before upstream finished: %q %v", line, err)
	}
	cancel()
	select {
	case <-cancelled:
	case <-time.After(2 * time.Second):
		t.Fatal("client cancellation did not reach upstream")
	}
}

func expectAPIModel(mock sqlmock.Sqlmock, power int) {
	mock.ExpectQuery("SELECT .*geekai_chat_models.*type = \\? AND enabled = \\?.*\\(open = \\? OR id IN \\(\\?\\)\\).*value = \\?").
		WithArgs("chat", true, true, 8, "test-model").
		WillReturnRows(sqlmock.NewRows([]string{"id", "value", "power", "max_tokens", "key_id"}).AddRow(1, "test-model", power, 1024, 9))
}

func expectAPIPower(mock sqlmock.Sqlmock, refund bool) {
	mark, kind, balance := types.PowerSub, types.PowerConsume, 8
	mock.ExpectBegin()
	if refund {
		mark, kind, balance = types.PowerAdd, types.PowerRefund, 10
		mock.ExpectExec("UPDATE `geekai_users` SET `power`=power \\+ \\? WHERE id = \\?").WithArgs(2, 7).WillReturnResult(sqlmock.NewResult(0, 1))
	} else {
		mock.ExpectExec("UPDATE `geekai_users` SET `power`=power - \\? WHERE id = \\? AND power >= \\?").WithArgs(2, 7, 2).WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectQuery("SELECT .*geekai_users").WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", balance))
	mock.ExpectExec("INSERT INTO `geekai_power_logs`").
		WithArgs(7, "alice", kind, 2, balance, "test-model", sqlmock.AnyArg(), mark, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func expectAPIUpstream(mock sqlmock.Sqlmock, endpoint string) {
	mock.ExpectQuery("SELECT .*geekai_api_keys.*enabled").WithArgs(true, "chat", 9).
		WillReturnRows(sqlmock.NewRows([]string{"id", "value", "api_url"}).AddRow(9, "upstream-secret", endpoint))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `geekai_api_keys` SET `last_used_at`").WithArgs(sqlmock.AnyArg(), 9).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

func TestModelAPIChatRelayAndBilling(t *testing.T) {
	normal := `{"id":"chatcmpl-test","object":"chat.completion","choices":[{"index":0,"message":{"role":"assistant","content":"你好","tool_calls":[{"id":"t1","type":"function","function":{"name":"search","arguments":"{}"}}]},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":5,"completion_tokens":3,"total_tokens":8}}`
	stream := "data: {\"id\":\"s1\",\"choices\":[{\"index\":0,\"delta\":{\"reasoning_content\":\"思考\"},\"finish_reason\":null}]}\n\n" +
		"data: {\"id\":\"s1\",\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"t1\",\"function\":{\"name\":\"search\",\"arguments\":\"{}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n" +
		"data: {\"id\":\"s1\",\"choices\":[],\"usage\":{\"total_tokens\":8}}\n\ndata: [DONE]\n\n"
	for _, tc := range []struct {
		name, response, contentType, suffix string
		stream                              bool
		upstreamStatus, status, power       int
		refund, streamError                 bool
	}{
		{name: "json", response: normal, contentType: "application/json", status: 200, power: 2},
		{name: "stream", response: stream, contentType: "text/event-stream", suffix: "/v1", stream: true, status: 200, power: 2},
		{name: "free", response: normal, contentType: "application/json", suffix: "/v1/chat/completions", status: 200},
		{name: "upstream-auth-error", response: `{"error":"upstream-secret"}`, upstreamStatus: 401, status: 502, power: 2, refund: true},
		{name: "malformed-json", response: `{"choices":[null]}`, status: 502, power: 2, refund: true},
		{name: "wrong-stream-type", response: normal, contentType: "application/json", stream: true, status: 502, power: 2, refund: true},
		{name: "empty-stream", response: "data: [DONE]\n\n", contentType: "text/event-stream", stream: true, status: 200, power: 2, refund: true, streamError: true},
		{name: "stream-error-before-output", response: "data: {\"error\":{\"message\":\"upstream-secret\"}}\n\n", contentType: "text/event-stream", stream: true, status: 200, power: 2, refund: true, streamError: true},
		{name: "stream-interruption-after-output", response: "data: {\"choices\":[{\"delta\":{\"content\":\"你好\"}}]}\n\n", contentType: "text/event-stream", stream: true, status: 200, power: 2, streamError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, engine, mock := modelAPITestHandler(t)
			h.RegisterRoutes(h.App.Router)
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer upstream-secret" {
					t.Errorf("incorrect upstream: %s %s", r.URL.Path, r.Header.Get("Authorization"))
				}
				body, _ := io.ReadAll(r.Body)
				var request struct {
					Model     string
					MaxTokens int `json:"max_tokens"`
					Messages  []struct{ Content json.RawMessage }
					Tools     []json.RawMessage
				}
				if json.Unmarshal(body, &request) != nil || request.Model != "test-model" || request.MaxTokens != 1024 || len(request.Tools) != 1 || string(request.Messages[0].Content) != `"  保留空格  "` {
					t.Errorf("request fields changed: %s", body)
				}
				w.Header().Set("Content-Type", tc.contentType)
				if tc.upstreamStatus != 0 {
					w.WriteHeader(tc.upstreamStatus)
				}
				io.WriteString(w, tc.response)
			}))
			defer upstream.Close()
			key := modelTokenPrefix + strings.Repeat("b", 43)
			expectModelAPIAuth(mock, key, true, 0)
			expectAPIModel(mock, tc.power)
			if tc.power > 0 {
				expectAPIPower(mock, false)
			}
			expectAPIUpstream(mock, upstream.URL+tc.suffix)
			if tc.refund {
				expectAPIPower(mock, true)
			}
			body := `{"model":"test-model","user_id":999,"messages":[{"role":"user","content":"  保留空格  "}],"tools":[{"type":"function","function":{"name":"search","parameters":{"type":"object"}}}],"stream":` + map[bool]string{true: "true", false: "false"}[tc.stream] + `}`
			w := modelAPIRequest(engine, "POST", "/v1/chat/completions", "Bearer "+key, body)
			if w.Code != tc.status {
				t.Fatalf("status: %d %s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "upstream-secret") || strings.Contains(w.Body.String(), "event:message") {
				t.Fatalf("leaked secret or custom protocol: %s", w.Body.String())
			}
			if tc.streamError {
				if !strings.Contains(w.Body.String(), `"type":"upstream_error"`) {
					t.Fatalf("missing stream error: %s", w.Body.String())
				}
			} else if tc.status == 200 && w.Body.String() != tc.response {
				t.Fatalf("upstream response was modified: %s", w.Body.String())
			}
			if tc.stream && tc.contentType == "text/event-stream" && w.Header().Get("X-Accel-Buffering") != "no" {
				t.Fatal("stream buffering not disabled")
			}
		})
	}
}

func TestModelAPIChatRejectsBeforeCallingUpstream(t *testing.T) {
	for _, tc := range []struct {
		name, body          string
		model, insufficient bool
		status              int
	}{
		{"missing-model", `{"messages":[{"role":"user","content":"hi"}]}`, false, false, 400},
		{"empty-messages", `{"model":"test-model","messages":[]}`, false, false, 400},
		{"multiple-replies", `{"model":"test-model","messages":[{"role":"user"}],"n":2}`, false, false, 400},
		{"invalid-json", `{`, false, false, 400},
		{"token-limit", `{"model":"test-model","messages":[{"role":"user"}],"max_tokens":2048}`, true, false, 400},
		{"insufficient", `{"model":"test-model","messages":[{"role":"user"}]}`, true, true, 402},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, engine, mock := modelAPITestHandler(t)
			h.RegisterRoutes(h.App.Router)
			key := modelTokenPrefix + strings.Repeat("c", 43)
			expectModelAPIAuth(mock, key, true, 0)
			if tc.model {
				expectAPIModel(mock, 2)
			}
			if tc.insufficient {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `geekai_users`.*power >= \\?").WithArgs(2, 7, 2).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			}
			w := modelAPIRequest(engine, "POST", "/v1/chat/completions", "Bearer "+key, tc.body)
			if w.Code != tc.status {
				t.Fatalf("%d %s", w.Code, w.Body.String())
			}
		})
	}
}
