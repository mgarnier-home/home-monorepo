package models

import (
	"mgarnier11/mineager/database"
	"strings"
)

type MapBo struct {
	ID          int
	Name        string
	Version     string
	Description string
}

func InitMapTable() {
	_, err := database.DB.Exec(`
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

func GetMaps() (maps []*MapBo, error error) {
	rows, err := database.DB.Query("SELECT * FROM maps")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var mapRow MapBo

		err := rows.Scan(&mapRow.ID, &mapRow.Name, &mapRow.Version, &mapRow.Description)

		if err != nil {
			return nil, err
		}

		maps = append(maps, &mapRow)
	}

	return maps, nil
}

func GetMapByName(name string) (mapRow *MapBo, error error) {
	var mapRowResult MapBo

	err := database.DB.QueryRow("SELECT * FROM maps WHERE name = ?", strings.ToLower(name)).Scan(&mapRowResult.ID, &mapRowResult.Name, &mapRowResult.Version, &mapRowResult.Description)

	if err != nil {
		return nil, err
	}

	return &mapRowResult, nil

}

func CreateMap(name string, version string, description string) (mapRow *MapBo, error error) {
	_, err := database.DB.Exec("INSERT INTO maps (name, version, description) VALUES (?, ?, ?)", strings.ToLower(name), version, description)

	if err != nil {
		return nil, err
	}

	return GetMapByName(name)
}

func DeleteMapByName(name string) error {
	_, err := database.DB.Exec("DELETE FROM maps WHERE name = ?", strings.ToLower(name))

	return err
}
