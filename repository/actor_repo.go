package repository

import (
	"database/sql"
	"fmt"
	m "movies_api/models"
)

func (r actorRepository) FindAllActors() ([]m.Actor, error) {
	query := `
	SELECT a.id, a.name, a.birth_date, m.id, m.title, m.release_year, m.duration
		FROM actors a
		LEFT JOIN movie_actors ma ON a.id = ma.actor_id
		LEFT JOIN movies m ON ma.movie_id = m.id
		ORDER by a.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}
	defer rows.Close()

	actorsMap := make(map[int]*m.Actor)

	var actorIDs []int

	for rows.Next() {
		var actor m.Actor

		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var releaseYear sql.NullInt64
		var duration sql.NullInt64

		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
			&movieID,
			&movieTitle,
			&releaseYear,
			&duration,
		)

		if err != nil {
			return nil, fmt.Errorf("Something happened during query execution: %w", err)
		}

		if _, exists := actorsMap[actor.Id]; !exists {
			actor.Movies = []m.Movie{}
			actorsMap[actor.Id] = &actor
			actorIDs = append(actorIDs, actor.Id)
		}

		if movieID.Valid {
			movie := m.Movie{
				Id:          int(movieID.Int64),
				Title:       movieTitle.String,
				ReleaseYear: int(releaseYear.Int64),
				Duration:    int(duration.Int64),
			}

			actorsMap[actor.Id].Movies = append(actorsMap[actor.Id].Movies, movie)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Something happened during query execution: %w", err)
	}

	actors := make([]m.Actor, 0, len(actorIDs))

	for _, id := range actorIDs {
		actors = append(actors, *actorsMap[id])
	}

	return actors, nil
}

func (r actorRepository) CreateActor(actor m.Actor) (m.Actor, error) {
	query := `
		INSERT INTO actors (name, birth_date)
		VALUES (?,?)
		RETURNING id,  name, birth_date
	`

	err := r.db.QueryRow(query, actor.Name, actor.BirthDate).Scan(
		&actor.Id,
		&actor.Name,
		&actor.BirthDate,
	)
	if err != nil {
		return m.Actor{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return actor, nil
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
	if rowsAffected == 0 {
		return false, sql.ErrNoRows
	}

	return true, nil
}

func (r actorRepository) FindActorByID(id int) (m.Actor, error) {
	query := `
	SELECT a.id, a.name, a.birth_date, m.id, m.title, m.release_year, m.duration
		FROM actors a
		LEFT JOIN movie_actors ma ON a.id = ma.actor_id
		LEFT JOIN movies m ON ma.movie_id = m.id
		WHERE a.id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return m.Actor{}, err
	}
	defer rows.Close()

	var actor m.Actor
	found := false

	for rows.Next() {
		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var releaseYear sql.NullInt64
		var duration sql.NullInt64

		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
			&movieID,
			&movieTitle,
			&releaseYear,
			&duration,
		)
		if err != nil {
			return m.Actor{}, err
		}
		if !found {
			actor.Movies = []m.Movie{}
			found = true
		}

		if movieID.Valid {
			actor.Movies = append(actor.Movies, m.Movie{
				Id:          int(movieID.Int64),
				Title:       movieTitle.String,
				ReleaseYear: int(releaseYear.Int64),
				Duration:    int(duration.Int64),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return m.Actor{}, err
	}

	if !found {
		return m.Actor{}, sql.ErrNoRows
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

	query += " WHERE id = ? RETURNING id, name, birth_date"
	args = append(args, id)

	var actor m.Actor

	err := r.db.QueryRow(query, args...).Scan(
		&actor.Id,
		&actor.Name,
		&actor.BirthDate,
	)
	if err != nil {
		return m.Actor{}, err
	}

	return actor, nil
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
