package repository

import (
	"bytes"
	"context"
	"database/sql"
	"main/internal/constants"
	"main/internal/microservices/auth/utils/hash"
	"main/internal/microservices/profile"
	proto "main/internal/microservices/profile/proto"
	"main/internal/microservices/profile/utils/images"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"
	"github.com/minio/minio-go/v7"
)

type Storage struct {
	db    *sql.DB
	minio *minio.Client
	redis *redis.Pool
}

func NewStorage(db *sql.DB, minio *minio.Client, redis *redis.Pool) profile.Storage { //nolint:nolintlint,ireturn
	return &Storage{db: db, minio: minio, redis: redis}
}

func (s Storage) GetUserProfile(userID int64) (*proto.ProfileData, error) {
	sqlScript := "SELECT username, email, avatar FROM users WHERE id=$1"

	var name, email, avatar string
	err := s.db.QueryRow(sqlScript, userID).Scan(&name, &email, &avatar)

	if err != nil {
		return nil, err
	}

	avatarURL, err := images.GenerateFileURL(avatar, constants.UserObjectsBucketName)
	if err != nil {
		return nil, err
	}

	return &proto.ProfileData{
		Name:   name,
		Email:  email,
		Avatar: avatarURL,
	}, nil
}

func (s Storage) EditProfile(data *proto.EditProfileData) error {
	sqlScript := "SELECT username, password, salt FROM users WHERE id=$1"

	var oldName, oldPassword, oldSalt string
	err := s.db.QueryRow(sqlScript, data.ID).Scan(&oldName, &oldPassword, &oldSalt)
	if err != nil {
		return err
	}

	notChangedPassword, _ := hash.ComparePasswords(oldPassword, oldSalt, data.Password)

	if !notChangedPassword && len(data.Password) != 0 {
		salt, errUUID := uuid.NewV4()
		if errUUID != nil {
			return errUUID
		}

		hashPassword, errHashAndSalt := hash.SaltHash(data.Password, salt.String())
		if errHashAndSalt != nil {
			return errHashAndSalt
		}

		oldPassword = hashPassword
		oldSalt = salt.String()
	}

	if data.Name != oldName && len(data.Name) != 0 {
		oldName = data.Name
	}

	sqlScript = "UPDATE users SET username = $2, password = $3, salt = $4 WHERE id = $1"

	_, err = s.db.Exec(sqlScript, data.ID, oldName, oldPassword, oldSalt)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) EditAvatar(data *proto.EditAvatarData) (string, error) {
	sqlScript := "SELECT avatar FROM users WHERE id=$1"

	var oldAvatar string
	err := s.db.QueryRow(sqlScript, data.ID).Scan(&oldAvatar)
	if err != nil {
		return "", err
	}

	if len(data.Avatar) != 0 {
		sqlScript := "UPDATE users SET avatar = $2 WHERE id = $1"

		_, err = s.db.Exec(sqlScript, data.ID, data.Avatar)
		if err != nil {
			return "", err
		}

		return oldAvatar, nil
	}

	return "", nil
}

func (s Storage) GetAvatar(userID int64) (string, error) {
	sqlScript := "SELECT avatar FROM users WHERE id=$1"

	var avatar string
	err := s.db.QueryRow(sqlScript, userID).Scan(&avatar)

	if err != nil {
		return "", err
	}

	return avatar, nil
}

func (s Storage) UploadAvatar(data *proto.UploadInputFile) (string, error) {
	imageName := images.GenerateObjectName(data.ID)

	opts := minio.PutObjectOptions{
		UserMetadata:            map[string]string{"x-amz-acl": "public-read"},
		UserTags:                nil,
		Progress:                nil,
		ContentType:             data.ContentType,
		ContentEncoding:         "",
		ContentDisposition:      "",
		ContentLanguage:         "",
		CacheControl:            "",
		Mode:                    "",
		RetainUntilDate:         time.Time{},
		ServerSideEncryption:    nil,
		NumThreads:              0,
		StorageClass:            "",
		WebsiteRedirectLocation: "",
		PartSize:                0,
		LegalHold:               "",
		SendContentMd5:          false,
		DisableContentSha256:    false,
		DisableMultipart:        false,
		ConcurrentStreamParts:   false,
		Internal: minio.AdvancedPutOptions{
			SourceVersionID:    "",
			SourceETag:         "",
			ReplicationStatus:  "",
			SourceMTime:        time.Time{},
			ReplicationRequest: false,
			RetentionTimestamp: time.Time{},
			TaggingTimestamp:   time.Time{},
			LegalholdTimestamp: time.Time{},
		},
	}

	_, err := s.minio.PutObject(
		context.Background(),
		constants.UserObjectsBucketName, // Константа с именем бакета
		imageName,
		bytes.NewReader(data.File),
		data.Size,
		opts,
	)
	if err != nil {
		return "", err
	}

	return imageName, nil
}

func (s Storage) DeleteFile(name string) error {
	opts := minio.RemoveObjectOptions{
		ForceDelete:      false,
		GovernanceBypass: false,
		VersionID:        "",
		Internal: minio.AdvancedRemoveOptions{
			ReplicationDeleteMarker: false,
			ReplicationStatus:       "",
			ReplicationMTime:        time.Time{},
			ReplicationRequest:      false,
		},
	}

	err := s.minio.RemoveObject(
		context.Background(),
		constants.UserObjectsBucketName,
		name,
		opts,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) AddLike(data *proto.LikeData) error {
	sqlScript := "INSERT INTO user_position(id_user, id_position) VALUES($1, (SELECT id FROM position WHERE name = $2))"

	if _, err := s.db.Exec(sqlScript, data.UserID, data.PositionName); err != nil {
		return err
	}

	return nil
}

func (s Storage) RemoveLike(data *proto.LikeData) error {
	sqlScript := "DELETE FROM user_position WHERE id_user=$1 AND id_position=(SELECT id FROM position WHERE name = $2)"

	if _, err := s.db.Exec(sqlScript, data.UserID, data.PositionName); err != nil {
		return err
	}

	return nil
}

func (s Storage) GetFavorites(userID int64) ([]*proto.Favorite, error) {
	sqlScript := "SELECT position.id, position.name, position.count_technologies, " +
		"COUNT(ut.id_technology) AS count_finished FROM position " +
		"JOIN user_position ON user_position.id_user = $1 AND user_position.id_position = position.id " +
		"JOIN technology_position tp ON position.id = tp.id_position " +
		"LEFT JOIN user_technology ut ON tp.id_technology = ut.id_technology " +
		"GROUP BY position.id;"

	favorites := make([]*proto.Favorite, 0)

	rows, err := s.db.Query(sqlScript, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Err()
		_ = rows.Close()
	}()

	for rows.Next() {
		var favorite proto.Favorite
		if err = rows.Scan(&favorite.PositionID, &favorite.Name, &favorite.CountAll, &favorite.CountFinished); err != nil {
			return nil, err
		}
		favorites = append(favorites, &favorite)
	}

	return favorites, nil
}

func (s Storage) Finish(data *proto.LikeData) error {
	sqlScript := "INSERT INTO user_technology(id_user, id_technology) VALUES($1, (SELECT id FROM technology WHERE name = $2))"

	if _, err := s.db.Exec(sqlScript, data.UserID, data.PositionName); err != nil {
		return err
	}

	return nil
}

func (s Storage) Cancel(data *proto.LikeData) error {
	sqlScript := "DELETE FROM user_technology WHERE id_user=$1 AND id_technology=(SELECT id FROM technology WHERE name = $2)"

	if _, err := s.db.Exec(sqlScript, data.UserID, data.PositionName); err != nil {
		return err
	}

	return nil
}

func (s Storage) GetFinished(data *proto.LikeData) ([]string, error) {
	sqlScript := "SELECT name_technology FROM technology_position JOIN user_technology ut ON technology_position.id_technology = ut.id_technology AND ut.id_user=$1 AND name_position=$2"

	finished := make([]string, 0)

	rows, err := s.db.Query(sqlScript, data.UserID, data.PositionName)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Err()
		_ = rows.Close()
	}()

	for rows.Next() {
		var technology string
		if err = rows.Scan(&technology); err != nil {
			return nil, err
		}
		finished = append(finished, technology)
	}

	return finished, nil
}
