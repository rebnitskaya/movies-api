package repository

import (
	"database/sql"
	"fmt"
	"movies_api/models"
)

type ActorRepository struct {
	db *sql.DB
}

func (r ActorRepository) FindAllActors() ([]models.Actor, error) {
	query := `
		SELECT *
		FROM actors
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
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
			return nil, fmt.Errorf("Something happened during query execution: %w", err)
		}

		actors = append(actors, actor)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actors, nil
}

func (r ActorRepository) CreateActor(actor models.Actor) (bool, error) {
	query := `
		INSERT INTO actors (name, birth_date)
		VALUES (?,?)
	`

	_, err := r.db.Exec(query, actor.Name, actor.BirthDate)
	if err != nil {
		return false, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return true, nil
}

func (r ActorRepository) FindActorByNameAndBirthDate(actorData models.Actor) (models.Actor, error) {
	query := `
		SELECT *
		FROM actors
		WHERE name = ? AND birth_date = ?
	`

	rows, err := r.db.Query(query, actorData.Name, actorData.BirthDate)
	if err != nil {
		return models.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	defer rows.Close()

	var actor models.Actor

	for rows.Next() {
		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)

		if err != nil {
			return models.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return models.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actor, nil
}
