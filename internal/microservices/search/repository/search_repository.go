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
	sqlScript := "SELECT name_technology, distance, professionalism, t.hard_skill " +
		"FROM technology_position JOIN technology t ON " +
		"t.id = technology_position.id_technology AND name_position=$1"

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

		if err = rows.Scan(&technology.Name, &technology.Distance, &technology.Professionalism, &technology.HardSkill); err != nil {
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

	positions := make([]*proto.Position, 0)

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

func (s Storage) GetPositions(data *proto.GetTechnology) ([]*proto.Position, error) {
	sqlScript := "SELECT name_position FROM technology_position WHERE name_technology=$1 ORDER BY distance DESC"

	positions := make([]*proto.Position, 0)

	rows, err := s.db.Query(sqlScript, data.Name)
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

func (s Storage) GetTipsToLearn(data *proto.GetTechnology) (string, error) {
	sqlScript := "SELECT tips_to_learn FROM technology WHERE name=$1"

	var tipsToLearn string

	rows, err := s.db.Query(sqlScript, data.Name)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		if err = rows.Scan(&tipsToLearn); err != nil {
			return "", err
		}
	}

	return tipsToLearn, nil
}
