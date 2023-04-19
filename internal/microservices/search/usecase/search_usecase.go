package usecase

import (
	"context"
	"main/internal/microservices/search"
	proto "main/internal/microservices/search/proto"
)

type Service struct {
	storage search.Storage
}

func NewService(storage search.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetTechnologies(ctx context.Context, data *proto.SearchText) (*proto.TechnologiesArr, error) {
	isExists, err := s.storage.IsPositionExists(data)
	if err != nil {
		return &proto.TechnologiesArr{Technology: nil}, err
	}

	if isExists {
		technologies, errGetTechnologies := s.storage.GetTechnologies(data)
		if errGetTechnologies != nil {
			return &proto.TechnologiesArr{Technology: nil}, errGetTechnologies
		}

		return &proto.TechnologiesArr{Technology: technologies}, nil
	}

	return &proto.TechnologiesArr{Technology: nil}, nil
}

func (s *Service) GetTop(ctx context.Context, empty *proto.Empty) (*proto.PositionTop, error) {
	positions, err := s.storage.GetTop()
	if err != nil {
		return &proto.PositionTop{Position: nil}, err
	}

	return &proto.PositionTop{Position: positions}, nil
}

func (s *Service) GetPositions(ctx context.Context, data *proto.GetTechnology) (*proto.NodeInfo, error) {
	positions, err := s.storage.GetPositions(data)
	if err != nil {
		return &proto.NodeInfo{
			Position:    nil,
			TipsToLearn: "",
		}, err
	}

	tipsToLearn, err := s.storage.GetTipsToLearn(data)
	if err != nil {
		return &proto.NodeInfo{
			Position:    nil,
			TipsToLearn: "",
		}, err
	}

	return &proto.NodeInfo{
		Position:    positions,
		TipsToLearn: tipsToLearn,
	}, nil
}

