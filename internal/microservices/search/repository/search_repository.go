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

	sqlScript = "UPDATE position SET requests_count = requests_count + 1 WHERE name = $1"

	_, err = s.db.Exec(sqlScript, data.Text)
	if err != nil {
		return nil, err
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

func (s Storage) GetTop() ([]*proto.Position, error) {
	sqlScript := "SELECT name FROM position ORDER BY requests_count DESC LIMIT 5"

	positions := make([]*proto.Position, 5)

	rows, err := s.db.Query(sqlScript)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var position proto.Position
		if err = rows.Scan(&position.Name); err != nil {
			return nil, err
		}
		positions = append(positions, &position)
	}

	return positions, nil
}
