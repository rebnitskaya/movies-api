package repository

import (
	"fmt"
	m "movies_api/models"
)

func (r actorRepository) FindAllActors() ([]m.Actor, error) {
	return []m.Actor{}, nil
}

func (r actorRepository) CreateActor(actor m.Actor) (bool, error) {
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

func (r actorRepository) FindActorByNameAndBirthDate(name string, birthDate string) (m.Actor, error) {
	query := `
		SELECT *
		FROM actors
		WHERE name = ? AND birth_date = ?
	`

	rows, err := r.db.Query(query, name, birthDate)
	if err != nil {
		return m.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	defer rows.Close()

	var actor m.Actor

	for rows.Next() {
		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)

		if err != nil {
			return m.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return m.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actor, nil
}

func (r actorRepository) DeleteActorByID(id int) (bool, error) {
	return false, nil
}

func (r actorRepository) FindActorByID(id int) (m.Actor, bool) {
	return m.Actor{}, false
}

func (r actorRepository) ReplaceFieldsInActor(id int, fields map[string]string) (m.Actor, bool) {
	return m.Actor{}, false
}

func (r actorRepository) FindActorsByName(name string) ([]m.Actor, error) {
	return []m.Actor{}, nil
}
