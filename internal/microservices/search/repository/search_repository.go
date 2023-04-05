package repository

import (
	"database/sql"
	"log"
	proto "main/internal/microservices/search/proto"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) GetTechnologies(data *proto.SearchText) ([]*proto.Technology, error) {
	sqlScript := "SELECT name_technology, distance, professionalism FROM technology_position WHERE name_position=$1"

	technologies := make([]*proto.Technology, 0)

	rows, err := s.db.Query(sqlScript, data.Text)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var technology proto.Technology
		if err = rows.Scan(&technology.Name, &technology.Distance, &technology.Professionalism); err != nil {
			return nil, err
		}
		technologies = append(technologies, &technology)
	}

	return technologies, nil
}

func (s Storage) IsPositionExists(data *proto.SearchText) (bool, error) {
	sqlScript := "SELECT id FROM position WHERE name=$1"
	log.Println(data.Text)

	rows, err := s.db.Query(sqlScript, data.Text)
	if err != nil {
		return false, err
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	// убедимся, что всё закроется при выходе из программы
	defer func() {
		rows.Close()
	}()

	// Из базы пришел пустой запрос, значит пользователя в базе данных нет
	if !rows.Next() {
		return false, nil
	}

	return true, nil
}
