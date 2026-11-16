package handler_test

import (
	"context"
	"encoding/json"
	"geekai/core"
	"geekai/handler"
	"geekai/handler/admin"
	"geekai/service/library"
	"geekai/store/model"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

type integrationRepository struct {
	library.Repository
	mu        sync.Mutex
	sources   map[uint]model.LibrarySource
	resources map[uint]model.LibraryResource
}

func (r *integrationRepository) Sources(context.Context) ([]model.LibrarySource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := []model.LibrarySource{}
	for _, item := range r.sources {
		items = append(items, item)
	}
	return items, nil
}
func (r *integrationRepository) Source(_ context.Context, id uint) (model.LibrarySource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.sources[id]
	if !ok {
		return item, library.ErrNotFound
	}
	return item, nil
}
func (r *integrationRepository) SaveSource(_ context.Context, item *model.LibrarySource) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if item.Id == 0 {
		item.Id = uint(len(r.sources) + 1)
	}
	r.sources[item.Id] = *item
	return nil
}
func (r *integrationRepository) Resources(context.Context) ([]model.LibraryResource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := []model.LibraryResource{}
	for _, item := range r.resources {
		items = append(items, item)
	}
	return items, nil
}
func (r *integrationRepository) Resource(_ context.Context, id uint) (model.LibraryResource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.resources[id]
	if !ok {
		return item, library.ErrNotFound
	}
	return item, nil
}
func (r *integrationRepository) SaveResource(_ context.Context, item *model.LibraryResource) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if item.Id == 0 {
		item.Id = uint(len(r.resources) + 1)
	}
	r.resources[item.Id] = *item
	return nil
}

// This opt-in check runs both applications with real HTTP, replacing only the
// upstream service and persistence. Set THINKBUDDY_TEST_ROOT to the app checkout.
func TestPlatformLibraryIntegration(t *testing.T) {
	root := os.Getenv("THINKBUDDY_TEST_ROOT")
	if root == "" {
		t.Skip("set THINKBUDDY_TEST_ROOT to run the two-application integration")
	}
	python := os.Getenv("THINKBUDDY_TEST_PYTHON")
	if python == "" {
		python = filepath.Join(root, ".venv", "Scripts", "python.exe")
	}
	gin.SetMode(gin.TestMode)
	collection := library.Collection{ID: "remote-alpha", Name: "Upstream guide"}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "upstream-read-key" {
			t.Error("upstream credential was not applied")
			w.WriteHeader(401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		write := func(value interface{}) { _ = json.NewEncoder(w).Encode(value) }
		path := r.URL.Path
		switch {
		case path == "/api/v1/knowledge-bases":
			write(gin.H{"data": []library.Collection{collection}})
		case path == "/api/v1/knowledge-bases/remote-alpha":
			write(gin.H{"data": collection})
		case path == "/api/v1/knowledge-bases/remote-alpha/knowledge":
			write(gin.H{"total": 1, "data": []gin.H{{"id": "doc1", "knowledge_base_id": "remote-alpha", "file_name": "guide.pdf", "file_type": "pdf", "file_size": 10}}})
		case strings.HasPrefix(path, "/api/v1/knowledge/") && !strings.HasSuffix(path, "/preview"):
			id := strings.TrimPrefix(path, "/api/v1/knowledge/")
			kb := "remote-alpha"
			if id != "doc1" {
				kb = "other-private-kb"
			}
			write(gin.H{"data": gin.H{"id": id, "knowledge_base_id": kb, "file_name": "guide.pdf"}})
		case path == "/api/v1/chunks/doc1":
			write(gin.H{"total": 1, "data": []gin.H{{"id": "chunk1", "knowledge_id": "doc1", "chunk_index": 0, "content": "Shared learning guide"}}})
		case path == "/api/v1/knowledge/doc1/preview":
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", "inline; filename=guide.pdf")
			if r.Header.Get("Range") != "" {
				w.Header().Set("Content-Range", "bytes 0-3/10")
				w.WriteHeader(http.StatusPartialContent)
				_, _ = w.Write([]byte("%PDF"))
			} else {
				_, _ = w.Write([]byte("%PDF-guide"))
			}
		case path == "/api/v1/knowledge-search":
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["knowledge_base_id"] != "remote-alpha" {
				t.Error("search escaped publication scope")
			}
			write(gin.H{"data": []gin.H{{"id": "chunk1", "content": "Shared answer", "knowledge_id": "doc1", "knowledge_title": "Guide"}}})
		default:
			w.WriteHeader(404)
		}
	}))
	defer upstream.Close()
	repo := &integrationRepository{sources: map[uint]model.LibrarySource{}, resources: map[uint]model.LibraryResource{}}
	service := library.NewService(repo)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Next()
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/library/") && c.Writer.Status() < 400 {
			t.Error("ThinkBuddy sent WeKnora content traffic through GeekAI")
		}
	})
	app := &core.AppServer{Engine: engine}
	handler.NewLibraryHandler(app, service).RegisterRoutes()
	group := engine.Group("/api/admin/library", func(c *gin.Context) {
		if c.GetHeader("Admin-Authorization") != "integration-admin" {
			c.AbortWithStatus(401)
		}
	})
	admin.NewLibraryHandler(app, service).RegisterManagementRoutes(group)
	server := httptest.NewServer(engine)
	defer server.Close()
	cmd := exec.Command(python, "-m", "pytest", "tests/multi_user/test_platform_library_integration.py", "-q")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "THINKBUDDY_TEST_GEEKAI_URL="+server.URL, "THINKBUDDY_TEST_UPSTREAM_URL="+upstream.URL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ThinkBuddy integration failed: %v\n%s", err, output)
	}
	t.Log(string(output))
}
