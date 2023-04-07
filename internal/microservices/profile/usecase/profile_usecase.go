package usecase

import (
	"context"
	"main/internal/constants"
	"main/internal/microservices/profile"
	proto "main/internal/microservices/profile/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	storage profile.Storage
}

func NewService(storage profile.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetUserProfile(ctx context.Context, userID *proto.UserID) (*proto.ProfileData, error) {
	userData, err := s.storage.GetUserProfile(userID.ID)
	if err != nil {
		return &proto.ProfileData{
			Name:   "",
			Email:  "",
			Avatar: "",
		}, status.Error(codes.Internal, err.Error())
	}

	return userData, nil
}

func (s *Service) EditProfile(ctx context.Context, data *proto.EditProfileData) (*proto.Empty, error) {
	err := s.storage.EditProfile(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) EditAvatar(ctx context.Context, data *proto.EditAvatarData) (*proto.Empty, error) {
	oldAvatar, err := s.storage.EditAvatar(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if oldAvatar != constants.DefaultImage {
		err = s.storage.DeleteFile(oldAvatar)
		if err != nil {
			return &proto.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}

	return &proto.Empty{}, nil
}

func (s *Service) UploadAvatar(ctx context.Context, data *proto.UploadInputFile) (*proto.FileName, error) {
	name, err := s.storage.UploadAvatar(data)
	if err != nil {
		return &proto.FileName{Name: ""}, status.Error(codes.Internal, err.Error())
	}

	return &proto.FileName{Name: name}, nil
}

func (s *Service) GetAvatar(ctx context.Context, userID *proto.UserID) (*proto.FileName, error) {
	name, err := s.storage.GetAvatar(userID.ID)
	if err != nil {
		return &proto.FileName{Name: ""}, status.Error(codes.Internal, err.Error())
	}

	return &proto.FileName{Name: name}, nil
}

func (s *Service) AddLike(ctx context.Context, data *proto.LikeData) (*proto.Empty, error) {
	if err := s.storage.AddLike(data); err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) RemoveLike(ctx context.Context, data *proto.LikeData) (*proto.Empty, error) {
	err := s.storage.RemoveLike(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) GetFavorites(ctx context.Context, data *proto.UserID) (*proto.Favorites, error) {
	favorites, err := s.storage.GetFavorites(data.ID)
	if err != nil {
		return &proto.Favorites{Favorite: nil}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Favorites{Favorite: favorites}, nil
}
