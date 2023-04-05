package search

import proto "main/internal/microservices/search/proto"

type Storage interface {
	GetTechnologies(data *proto.SearchText) ([]*proto.Technology, error)
	IsPositionExists(data *proto.SearchText) (bool, error)
}
