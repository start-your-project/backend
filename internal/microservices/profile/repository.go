package profile

import (
	proto "main/internal/microservices/profile/proto"
)

type Storage interface {
	GetUserProfile(userID int64) (*proto.ProfileData, error)
	EditProfile(data *proto.EditProfileData) error
	EditAvatar(data *proto.EditAvatarData) (string, error)
	GetAvatar(userID int64) (string, error)
	UploadAvatar(data *proto.UploadInputFile) (string, error)
	DeleteFile(string) error

	AddLike(data *proto.LikeData) error
	RemoveLike(data *proto.LikeData) error
	GetFavorites(userID int64) ([]*proto.Favorite, error)

	Finish(data *proto.LikeData) error
	Cancel(data *proto.LikeData) error
	GetFinished(data *proto.LikeData) ([]string, error)
}
