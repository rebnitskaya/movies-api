package repository

import (
	"fmt"
)

func (r actorRepository) FindAllActors() ([]Actor, error) {
	query := `
		SELECT id, name, birth_date
		FROM actors
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}
	defer rows.Close()

	var actors []Actor

	for rows.Next() {
		var actor Actor

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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actors, nil
}

func (r actorRepository) CreateActor(actor Actor) (bool, error) {
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

func (r actorRepository) FindActorByNameAndBirthDate(name string, birthDate string) (Actor, error) {
	query := `
		SELECT *
		FROM actors
		WHERE name = ? AND birth_date = ?
	`

	rows, err := r.db.Query(query, name, birthDate)
	if err != nil {
		return Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	defer rows.Close()

	var actor Actor

	for rows.Next() {
		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)

		if err != nil {
			return Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actor, nil
}

func (r actorRepository) DeleteActorByID(id int) (bool, error) {
	return false, nil
}

func (r actorRepository) FindActorByID(id int) (Actor, bool) {
	//or return error(?)
	query := `
		SELECT *
		FROM actors
		WHERE id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return Actor{}, false
	}
	defer rows.Close()

	var actor Actor
	found := false

	for rows.Next() {
		found = true

		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return Actor{}, false
		}
	}

	if !found {
		return Actor{}, false
	}

	err = rows.Err()
	if err != nil {
		return Actor{}, false
	}

	return actor, true
}

func (r actorRepository) ReplaceFieldsInActor(id int, fields map[string]string) (Actor, bool) {
	return Actor{}, false
}

func (r actorRepository) FindActorsByName(name string) ([]Actor, error) {
	query := `
		SELECT id, name, birth_date
		FROM actors
		WHERE name = ?
	`

	rows, err := r.db.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}
	defer rows.Close()

	var actors []Actor

	for rows.Next() {
		var actor Actor

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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actors, nil
}
