package repository

import (
	"database/sql"
	"log"
	proto "main/internal/microservices/search/proto"
	"sort"
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
		var hardSkill sql.NullBool

		if err = rows.Scan(&technology.Name, &technology.Distance, &technology.Professionalism, &hardSkill); err != nil {
			return nil, err
		}
		technology.HardSkill = hardSkill.Bool
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

	var tipsToLearn sql.NullString

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

	return tipsToLearn.String, nil
}

func (s Storage) TechSearch(data *proto.Technologies) ([]*proto.TechSearchPosition, error) {
	profs := make(map[string]float64, 0)

	for _, val := range data.Technology {
		sqlScript := "select name_position, distance from technology_position where name_technology = $1"
		rows, err := s.db.Query(sqlScript, val.Name)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			_ = rows.Err()
			_ = rows.Close()
		}()

		for rows.Next() {
			var profession string
			var distance float64
			if err = rows.Scan(&profession, &distance); err != nil {
				log.Fatal(err)
			}

			// nolint:gosimple
			_, _ = profs[profession]
			profs[profession] += distance
		}
	}
	keys := make([]string, 0, len(profs))

	for key := range profs {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return profs[keys[i]] > profs[keys[j]]
	})

	res := make([]*proto.TechSearchPosition, 0)

	for _, key := range keys {
		res = append(res, &proto.TechSearchPosition{
			Name:    key,
			Percent: float32(profs[key] / float64(len(data.Technology))),
		})
	}

	return res, nil
}
