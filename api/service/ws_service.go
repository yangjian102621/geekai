package service

import "geekai/core/types"

type WebsocketService struct {
	Clients *types.LMap[string, *types.WsClient] // clientId => Client
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		Clients: types.NewLMap[string, *types.WsClient](),
	}
}
