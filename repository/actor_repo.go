package repository

import (
	"fmt"
	m "movies_api/models"
)

func (r actorRepository) FindAllActors() ([]m.Actor, error) {
	query := `
		SELECT id, name, birth_date
		FROM actors
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}
	defer rows.Close()

	var actors []m.Actor

	for rows.Next() {
		var actor m.Actor

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
	query := `
		DELETE FROM actors
		WHERE id = ?
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r actorRepository) FindActorByID(id int) (m.Actor, error) {
	//or return error(?)
	query := `
		SELECT *
		FROM actors
		WHERE id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return m.Actor{}, err
	}
	defer rows.Close()

	var actor m.Actor
	found := false

	for rows.Next() {
		found = true

		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return m.Actor{}, err
		}
	}

	if !found {
		return m.Actor{}, err
	}

	err = rows.Err()
	if err != nil {
		return m.Actor{}, err
	}

	return actor, nil

}

func (r actorRepository) ReplaceFieldsInActor(id int, fields map[string]string) (m.Actor, error) {
	query := "UPDATE actors SET "
	args := []any{}

	first := true
	for field, value := range fields {
		var column string

		switch field {
		case "name":
			column = "name"
		case "birthDate":
			column = "birth_date"
		default:
			return m.Actor{}, fmt.Errorf("Invalid field: %s", field)
		}

		if !first {
			query += ", "
		}

		query += column + " = ?"
		args = append(args, value)
		first = false
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return m.Actor{}, err
	}

	return r.FindActorByID(id)
}

func (r actorRepository) FindActorsByName(name string) ([]m.Actor, error) {
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

	var actors []m.Actor

	for rows.Next() {
		var actor m.Actor

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
