package repository

import (
	"database/sql"
	"fmt"
	"movies_api/models"
)

type ActorRepository struct {
	db *sql.DB
}

type ActorRepositoryInterface interface {
	GetAllActors() ([]models.Actor, error)
	CreateActor(models.Actor) bool
}

func (r ActorRepository) GetAllActors() ([]models.Actor, error) {
	query := `
		SELECT id, name, birth_date
		FROM actors
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution:", err)
	}
	defer rows.Close()

	var actors []models.Actor

	for rows.Next() {
		actor := models.Actor{}

		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)

		if err != nil {
			return nil, fmt.Errorf("Something happened during query execution:", err)
		}

		actors = append(actors, actor)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution:", err)
	}

	return actors, nil
}

// func (r ActorRepository) CreateActor (actor models.Actor) bool {
// 	query := `
// 		INSERT INTO actors
// 	`
// }
