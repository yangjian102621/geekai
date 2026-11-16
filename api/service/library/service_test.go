package library

import (
	"context"
	"encoding/json"
	"geekai/store/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type memoryRepository struct {
	source   model.LibrarySource
	resource model.LibraryResource
}

func (r *memoryRepository) Sources(context.Context) ([]model.LibrarySource, error) {
	return []model.LibrarySource{r.source}, nil
}
func (r *memoryRepository) Source(context.Context, uint) (model.LibrarySource, error) {
	return r.source, nil
}
func (r *memoryRepository) SaveSource(_ context.Context, item *model.LibrarySource) error {
	item.Id = 1
	r.source = *item
	return nil
}
func (r *memoryRepository) DeleteSource(context.Context, uint) error {
	r.source = model.LibrarySource{}
	return nil
}
func (r *memoryRepository) Resources(context.Context) ([]model.LibraryResource, error) {
	return []model.LibraryResource{r.resource}, nil
}
func (r *memoryRepository) Resource(context.Context, uint) (model.LibraryResource, error) {
	return r.resource, nil
}
func (r *memoryRepository) SaveResource(_ context.Context, item *model.LibraryResource) error {
	item.Id = 2
	r.resource = *item
	return nil
}
func (r *memoryRepository) DeleteResource(context.Context, uint) error {
	r.resource = model.LibraryResource{}
	return nil
}

func TestWeKnoraDirectCredentialLifecycle(t *testing.T) {
	ctx := context.Background()
	sharedKey := "shared-read-key"
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != sharedKey {
			t.Error("request did not use the source key")
		}
		if r.URL.Path == "/api/v1/knowledge-bases" {
			io.WriteString(w, `{"data":[{"id":"remote1","name":"Upstream"}]}`)
			return
		}
		if r.URL.Path != "/api/v1/knowledge-bases/remote1" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		io.WriteString(w, `{"data":{"id":"remote1","name":"Upstream"}}`)
	}))
	defer upstream.Close()
	repo := &memoryRepository{}
	service := NewService(repo)
	sourceInput := SourceInput{
		Name: "WeKnora", BaseURL: upstream.URL,
		APIKey: sharedKey, Enabled: true,
	}
	source, err := service.SaveSource(ctx, sourceInput)
	if err != nil {
		t.Fatal(err)
	}
	if source.Provider != "weknora" {
		t.Fatal("source type was not set to WeKnora")
	}
	collections, err := service.RemoteCollections(ctx, source.Id)
	if err != nil || len(collections) != 1 || collections[0].ID != "remote1" {
		t.Fatalf("could not select remote knowledge base: %+v %v", collections, err)
	}
	input := model.LibraryResource{
		SourceID: source.Id, RemoteID: collections[0].ID, Name: "Gift", Published: true,
	}
	resource, err := service.SaveResource(ctx, input)
	if err != nil {
		t.Fatalf("could not publish with only the source key: %v", err)
	}
	resources, err := service.Resources(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []any{source, resource, resources} {
		data, _ := json.Marshal(value)
		if strings.Contains(string(data), sharedKey) {
			t.Fatal("management response disclosed the source credential")
		}
	}
	catalog, err := service.Catalog(ctx)
	if err != nil || len(catalog) != 1 || catalog[0].Connection == nil {
		t.Fatalf("missing connection: %+v %v", catalog, err)
	}
	connection := catalog[0].Connection
	if connection.APIKey != sharedKey || connection.ServerURL != upstream.URL || connection.KnowledgeBaseID != input.RemoteID {
		t.Fatal("incorrect distributed connection")
	}
	sourceInput.ID, sourceInput.APIKey = source.Id, ""
	if _, err := service.SaveSource(ctx, sourceInput); err != nil || repo.source.Credential != sharedKey {
		t.Fatal("empty edit changed the shared key")
	}
	input.Id, input.Name = resource.Id, "Renamed gift"
	if _, err := service.SaveResource(ctx, input); err != nil {
		t.Fatal(err)
	}
	sharedKey, sourceInput.APIKey = "rotated-read-key", "rotated-read-key"
	if _, err := service.SaveSource(ctx, sourceInput); err != nil {
		t.Fatal(err)
	}
	catalog, _ = service.Catalog(ctx)
	if catalog[0].Connection.APIKey != sharedKey || catalog[0].ID != "2" || catalog[0].Name != input.Name {
		t.Fatal("key rotation was not reflected")
	}
	input.RemoteID = "outside-scope"
	if _, err := service.SaveResource(ctx, input); err == nil || repo.resource.RemoteID != "remote1" {
		t.Fatal("published a knowledge base the shared key cannot access")
	}
	input.RemoteID = "remote1"
	repo.source.Credential = ""
	if _, err := service.SaveResource(ctx, input); err == nil {
		t.Fatal("published without a source key")
	}
	catalog, _ = service.Catalog(ctx)
	if len(catalog) != 0 {
		t.Fatal("distributed a connection without a source key")
	}
	repo.source.Credential = sharedKey
	input.Published = false
	if _, err := service.SaveResource(ctx, input); err != nil {
		t.Fatal(err)
	}
	catalog, _ = service.Catalog(ctx)
	if len(catalog) != 0 {
		t.Fatal("unpublished connection remained visible")
	}
}

func TestWeKnoraSourceCredentialValidation(t *testing.T) {
	for _, key := range []string{"", "bad\r\nkey", strings.Repeat("k", 8193)} {
		repo := &memoryRepository{}
		service := NewService(repo)
		_, err := service.SaveSource(context.Background(), SourceInput{
			Name: "WeKnora", BaseURL: "https://knowledge.example.com", APIKey: key,
		})
		if err == nil || repo.source.Id != 0 {
			t.Fatal("saved an invalid source key")
		}
	}
}
