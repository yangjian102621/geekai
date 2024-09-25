package service

import "geekai/core/types"

type WebsocketService struct {
	Clients *types.LMap[string, *types.WsClient] // clientId => Client
}
