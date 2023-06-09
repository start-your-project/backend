package search

import proto "main/internal/microservices/search/proto"

type Storage interface {
	GetTechnologies(data *proto.SearchText) ([]*proto.Technology, error)
	IsPositionExists(data *proto.SearchText) (bool, error)
	GetTop() ([]*proto.Position, error)
	GetPositions(data *proto.GetTechnology) ([]*proto.Position, error)
	GetTipsToLearn(data *proto.GetTechnology) (string, error)
	TechSearch(data *proto.Technologies) ([]*proto.TechSearchPosition, error)
}
