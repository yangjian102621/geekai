package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/handler/admin"
	"geekai/service/library"
	"geekai/store/model"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type publicRepository struct {
	library.Repository
	source   model.LibrarySource
	resource model.LibraryResource
	err      error
}

func (r *publicRepository) Sources(context.Context) ([]model.LibrarySource, error) {
	return []model.LibrarySource{r.source}, r.err
}
func (r *publicRepository) Source(context.Context, uint) (model.LibrarySource, error) {
	return r.source, nil
}
func (r *publicRepository) Resources(context.Context) ([]model.LibraryResource, error) {
	return []model.LibraryResource{r.resource}, nil
}
func (r *publicRepository) Resource(context.Context, uint) (model.LibraryResource, error) {
	return r.resource, nil
}

func TestPublicLibraryRoutesAndAdminGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &publicRepository{
		source:   model.LibrarySource{Id: 1, Provider: "weknora", BaseURL: "https://knowledge.example.com", Credential: "read-key", Enabled: true},
		resource: model.LibraryResource{Id: 2, SourceID: 1, RemoteID: "remote1", Name: "Public gift", Published: true},
	}
	service := library.NewService(repo)
	engine := gin.New()
	app := &core.AppServer{Engine: engine, Config: &types.AppConfig{}}
	handler.NewLibraryHandler(app, service).RegisterRoutes()
	admin.NewLibraryHandler(app, service).RegisterRoutes()
	request := func(method, path, body string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	if w := request("GET", "/api/admin/library/sources", ""); w.Code != 401 {
		t.Fatalf("admin route not protected: %d", w.Code)
	}
	if w := request("GET", "/api/v1/library", ""); w.Code != 200 || !strings.Contains(w.Body.String(), "Public gift") || !strings.Contains(w.Body.String(), "read-key") || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("catalog: %s", w.Body.String())
	}
	for _, path := range []string{"", "/documents", "/documents/d1/preview", "/documents/d1/text", "/search"} {
		if w := request("GET", "/api/v1/library/2"+path, ""); w.Code != 404 {
			t.Fatal("legacy content endpoint still registered")
		}
	}
	if w := request("POST", "/api/v1/library/2/search", `{}`); w.Code != 404 {
		t.Fatal("legacy search still registered")
	}
	for _, state := range []string{"unpublished", "disabled"} {
		repo.resource.Published = state != "unpublished"
		repo.source.Enabled = state != "disabled"
		w := request("GET", "/api/v1/library", "")
		var payload struct {
			Data []library.Collection `json:"data"`
		}
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &payload) != nil || len(payload.Data) != 0 {
			t.Fatal("unavailable gift remained visible")
		}
	}
	repo.err = errors.New("private database details")
	if w := request("GET", "/api/v1/library", ""); w.Code != 500 || strings.Contains(w.Body.String(), "private") {
		t.Fatalf("catalog failure was not handled safely: %d %s", w.Code, w.Body.String())
	}
}
