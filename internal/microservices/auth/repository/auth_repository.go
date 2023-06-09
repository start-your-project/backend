package repository

import (
	"database/sql"
	"main/internal/constants"
	proto "main/internal/microservices/auth/proto"
	"main/internal/microservices/auth/utils/hash"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"
)

type Storage struct {
	db    *sql.DB
	redis *redis.Pool
}

func NewStorage(db *sql.DB, redis *redis.Pool) *Storage {
	return &Storage{db: db, redis: redis}
}

func (s Storage) IsUserExists(data *proto.LogInData) (int64, error) {
	var userID int64
	sqlScript := "SELECT id, password, salt, email_confirmed FROM users WHERE email=$1"
	rows, err := s.db.Query(sqlScript, data.Email)
	if err != nil {
		return userID, err
	}
	err = rows.Err()
	if err != nil {
		return userID, err
	}
	// убедимся, что всё закроется при выходе из программы
	defer func() {
		rows.Close()
	}()

	// Из базы пришел пустой запрос, значит пользователя в базе данных нет
	if !rows.Next() {
		return userID, constants.ErrWrongData
	}

	var (
		password, salt string
		emailConfirmed bool
	)
	err = rows.Scan(&userID, &password, &salt, &emailConfirmed)

	// выход при ошибке
	if err != nil {
		return userID, err
	}

	if !emailConfirmed {
		return userID, constants.ErrEmailIsNotConfirmed
	}

	_, err = hash.ComparePasswords(password, salt, data.Password)
	if err != nil {
		return userID, constants.ErrWrongData
	}

	return userID, nil
}

func (s Storage) IsUserUnique(email string) (bool, error) {
	sqlScript := "SELECT id FROM users WHERE email=$1"
	rows, err := s.db.Query(sqlScript, email)
	if err != nil {
		return false, err
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	defer func() {
		rows.Close()
	}()

	if rows.Next() { // Пользователь с таким email зарегистрирован
		return false, nil
	}
	return true, nil
}

func (s Storage) CreateUser(data *proto.SignUpData) (*proto.Hash, error) {
	var userID int64

	salt, err := uuid.NewV4()
	if err != nil {
		return &proto.Hash{Hash: ""}, err
	}

	hashPassword, err := hash.SaltHash(data.Password, salt.String())
	if err != nil {
		return &proto.Hash{Hash: ""}, err
	}

	sqlScript := "INSERT INTO users(username, email, password, salt, avatar, email_confirmed) VALUES($1, $2, $3, $4, $5, FALSE) RETURNING id"

	//nolint:execinquery
	if err = s.db.QueryRow(sqlScript, data.Name, data.Email, hashPassword, salt, constants.DefaultImage).Scan(&userID); err != nil {
		return &proto.Hash{Hash: ""}, err
	}

	return &proto.Hash{Hash: hashPassword}, nil
}

func (s Storage) DeleteUser(hash *proto.Hash) error {
	sqlScript := "DELETE FROM users WHERE password=$1"

	_, err := s.db.Exec(sqlScript, hash.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) ConfirmEmail(data *proto.Hash) error {
	sqlScript := "SELECT id, email_confirmed FROM users WHERE password=$1"
	rows, err := s.db.Query(sqlScript, data.Hash)
	if err != nil {
		return err
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	// убедимся, что всё закроется при выходе из программы
	defer func() {
		rows.Close()
	}()

	// Из базы пришел пустой запрос, значит пользователя в базе данных нет
	if !rows.Next() {
		return constants.ErrWrongData
	}

	var (
		userID         int64
		emailConfirmed bool
	)

	err = rows.Scan(&userID, &emailConfirmed)
	if err != nil {
		return err
	}

	if emailConfirmed {
		return constants.ErrEmailAlreadyConfirmed
	}

	sqlScript = "UPDATE users SET email_confirmed=TRUE WHERE id=$1"

	_, err = s.db.Exec(sqlScript, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) GetEmailLink(domen string) (*proto.EmailLink, error) {
	sqlScript := "SELECT link FROM emails WHERE domen=$1"
	rows, err := s.db.Query(sqlScript, domen)
	if err != nil {
		return &proto.EmailLink{Link: ""}, err
	}
	err = rows.Err()
	if err != nil {
		return &proto.EmailLink{Link: ""}, err
	}
	defer func() {
		rows.Close()
	}()

	if !rows.Next() {
		return &proto.EmailLink{Link: ""}, nil
	}

	var (
		link string
	)

	err = rows.Scan(&link)
	if err != nil {
		return &proto.EmailLink{Link: ""}, err
	}

	return &proto.EmailLink{
		Link: link,
	}, nil
}

func (s Storage) StoreSession(userID int64) (string, error) {
	connRedis := s.redis.Get()
	defer connRedis.Close()

	session, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	_, err = connRedis.Do("SET", session, userID, "EX", int64(30*24*time.Hour.Seconds()))

	if err != nil {
		return "", err
	}

	return session.String(), nil
}

func (s Storage) GetUserID(session string) (int64, error) {
	connRedis := s.redis.Get()
	defer connRedis.Close()

	userID, err := redis.Int64(connRedis.Do("GET", session))
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (s Storage) DeleteSession(session string) error {
	connRedis := s.redis.Get()
	defer connRedis.Close()

	_, err := connRedis.Do("DEL", session)
	if err != nil {
		return err
	}

	return nil
}
