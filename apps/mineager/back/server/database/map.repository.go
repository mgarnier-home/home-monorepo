package database

import (
	"database/sql"
	"mgarnier11/mineager/server/objects/bo"
	"strings"
)

type MapRepository struct {
	database *sql.DB
}

func CreateMapRepository() *MapRepository {
	return &MapRepository{
		database: DB,
	}
}

func InitMapTable() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS maps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			version TEXT,
			description TEXT
		);
	`)

	if err != nil {
		panic(err)
	}
}

func (repo *MapRepository) GetMaps() (maps []*bo.MapBo, error error) {
	rows, err := repo.database.Query("SELECT * FROM maps")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var mapRow bo.MapBo

		err := rows.Scan(&mapRow.ID, &mapRow.Name, &mapRow.Version, &mapRow.Description)

		if err != nil {
			return nil, err
		}

		maps = append(maps, &mapRow)
	}

	return maps, nil
}

func (repo *MapRepository) GetMapByName(name string) (mapRow *bo.MapBo, error error) {
	var mapRowResult bo.MapBo

	err := repo.database.QueryRow("SELECT * FROM maps WHERE name = ?", strings.ToLower(name)).Scan(&mapRowResult.ID, &mapRowResult.Name, &mapRowResult.Version, &mapRowResult.Description)

	if err != nil {
		return nil, err
	}

	return &mapRowResult, nil

}

func (repo *MapRepository) CreateMap(name string, version string, description string) (mapRow *bo.MapBo, error error) {
	_, err := repo.database.Exec("INSERT INTO maps (name, version, description) VALUES (?, ?, ?)", strings.ToLower(name), version, description)

	if err != nil {
		return nil, err
	}

	return repo.GetMapByName(name)
}

func (repo *MapRepository) DeleteMapByName(name string) error {
	_, err := repo.database.Exec("DELETE FROM maps WHERE name = ?", strings.ToLower(name))

	return err
}
