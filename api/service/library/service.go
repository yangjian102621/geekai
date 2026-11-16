package library

import (
	"context"
	"errors"
	"geekai/store/model"
	"strconv"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type SourceView struct {
	model.LibrarySource
	HasAPIKey bool `json:"has_api_key"`
}

type SourceInput struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	Enabled bool   `json:"enabled"`
}

func (s *Service) Sources(ctx context.Context) ([]SourceView, error) {
	items, err := s.repo.Sources(ctx)
	result := []SourceView{}
	for _, item := range items {
		result = append(result, SourceView{item, item.Credential != ""})
	}
	return result, err
}

func (s *Service) SaveSource(ctx context.Context, input SourceInput) (SourceView, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.APIKey = strings.TrimSpace(input.APIKey)
	if input.Name == "" || len([]rune(input.Name)) > 120 || len(input.APIKey) > 8192 || strings.ContainsAny(input.APIKey, "\r\n") {
		return SourceView{}, errors.New("请输入有效的来源名称和凭据")
	}
	base, err := ValidateURL(input.BaseURL)
	if err != nil {
		return SourceView{}, err
	}
	item := model.LibrarySource{}
	if input.ID != 0 {
		item, err = s.repo.Source(ctx, input.ID)
		if err != nil {
			return SourceView{}, err
		}
		if item.Provider != "weknora" {
			return SourceView{}, errors.New("已有连接的来源类型不可变更，请新建连接")
		}
	}
	// An empty edit preserves the saved key.
	if input.APIKey != "" {
		item.Credential = input.APIKey
	}
	if item.Credential == "" {
		return SourceView{}, errors.New("WeKnora 连接需要 API Key")
	}
	item.Name, item.Provider, item.BaseURL, item.Enabled = input.Name, "weknora", base, input.Enabled
	if err := s.repo.SaveSource(ctx, &item); err != nil {
		return SourceView{}, errors.New("保存资料源失败")
	}
	return SourceView{item, item.Credential != ""}, nil
}

func (s *Service) client(source model.LibrarySource) (*remote, error) {
	if source.Provider != "weknora" {
		return nil, errors.New("不支持的资料源类型")
	}
	return newRemote(source.BaseURL, source.Credential)
}

func (s *Service) RemoteCollections(ctx context.Context, id uint) ([]Collection, error) {
	source, err := s.repo.Source(ctx, id)
	if err != nil {
		return nil, err
	}
	adapter, err := s.client(source)
	if err != nil {
		return nil, err
	}
	return adapter.Collections(ctx)
}

func (s *Service) DeleteSource(ctx context.Context, id uint) error {
	if _, err := s.repo.Source(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteSource(ctx, id)
}

func (s *Service) Resources(ctx context.Context) ([]model.LibraryResource, error) {
	return s.repo.Resources(ctx)
}

func (s *Service) SaveResource(ctx context.Context, input model.LibraryResource) (model.LibraryResource, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.SourceID == 0 || !ValidID(input.RemoteID) || input.Name == "" ||
		len([]rune(input.Name)) > 120 || len(input.Description) > 16000 {
		return model.LibraryResource{}, errors.New("请输入有效的来源、远端知识库 ID 和展示名称")
	}
	source, err := s.repo.Source(ctx, input.SourceID)
	if err != nil {
		return model.LibraryResource{}, err
	}
	item := model.LibraryResource{}
	if input.Id != 0 {
		item, err = s.repo.Resource(ctx, input.Id)
		if err != nil {
			return model.LibraryResource{}, err
		}
	}
	if input.Published && source.Credential == "" {
		return model.LibraryResource{}, errors.New("请先在资料来源中配置 WeKnora API Key，再上架")
	}
	// Verify the selected collection before publishing or changing its mapping.
	if input.Published || item.SourceID != input.SourceID || item.RemoteID != input.RemoteID {
		adapter, err := s.client(source)
		if err != nil {
			return model.LibraryResource{}, err
		}
		if _, err := adapter.Collection(ctx, input.RemoteID); err != nil {
			return model.LibraryResource{}, err
		}
	}
	item.SourceID, item.RemoteID, item.Name = input.SourceID, input.RemoteID, input.Name
	item.Description, item.SortNum, item.Published = input.Description, input.SortNum, input.Published
	if err := s.repo.SaveResource(ctx, &item); err != nil {
		return model.LibraryResource{}, errors.New("保存发布记录失败，同一来源的知识库请勿重复发布")
	}
	return item, nil
}

func (s *Service) DeleteResource(ctx context.Context, id uint) error {
	if _, err := s.repo.Resource(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteResource(ctx, id)
}

func (s *Service) Catalog(ctx context.Context) ([]Collection, error) {
	sources, err := s.repo.Sources(ctx)
	if err != nil {
		return nil, err
	}
	enabled := map[uint]model.LibrarySource{}
	for _, source := range sources {
		if source.Enabled && source.Provider == "weknora" && source.Credential != "" {
			enabled[source.Id] = source
		}
	}
	items, err := s.repo.Resources(ctx)
	if err != nil {
		return nil, err
	}
	result := []Collection{}
	for _, item := range items {
		if source, ok := enabled[item.SourceID]; ok && item.Published {
			info := Collection{
				ID: strconv.FormatUint(uint64(item.Id), 10), Name: item.Name, Description: item.Description,
				Provider: "weknora", Connection: &WeKnoraConnection{source.BaseURL, source.Credential, item.RemoteID},
			}
			result = append(result, info)
		}
	}
	return result, nil
}
