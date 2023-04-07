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

	//jsonFile, err := os.Open("test1.json")
	//if err != nil {
	//	return &proto.TechnologiesArr{Technology: nil}, err
	//}
	//defer func(jsonFile *os.File) {
	//	err = jsonFile.Close()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}(jsonFile)
	//
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//var technologies models.JSONData
	//err = json.Unmarshal(byteValue, &technologies)
	//if err != nil {
	//	return &proto.TechnologiesArr{Technology: nil}, err
	//}
	//
	//res := proto.TechnologiesArr{Technology: nil}
	//for _, technology := range technologies.Additional {
	//	res.Technology = append(res.Technology, &proto.Technology{
	//		Name:            technology.TechnologyName,
	//		Distance:        technology.Distance,
	//		Professionalism: technology.Professionalism,
	//	})
	//}
	//
	//return &res, nil
}
