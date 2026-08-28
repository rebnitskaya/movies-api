package repository

import (
	"database/sql"
	"fmt"
	m "movies_api/models"
)

func (r actorRepository) FindAllActors(limit, offset int) ([]m.Actor, error) {
	query := `
		SELECT
	    	a.id,
	     	a.name,
	      	a.birth_date,
	       	m.id,
	    	m.title,
	        m.release_year,
	        m.duration
		FROM (
	    	SELECT id, name, birth_date
	     	FROM actors
	     	ORDER BY id
	      	LIMIT ? OFFSET ?
	    ) a
		LEFT JOIN movie_actors am ON a.id = am.actor_id
		LEFT JOIN movies m ON am.movie_id = m.id
		ORDER BY a.id
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	actorsMap := make(map[int]*m.Actor)

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
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		if _, exists := actorsMap[actor.Id]; !exists {
			actor.Movies = []m.Movie{}
			actorsMap[actor.Id] = &actor
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
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	actors := make([]m.Actor, 0, len(actorsMap))

	for _, id := range actorsMap {
		actors = append(actors, *id)
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
		return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
		return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
			return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}
	}

	err = rows.Err()
	if err != nil {
		return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return m.Actor{}, m.ErrActorNotFound
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
		return false, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if rowsAffected == 0 {
		return false, m.ErrActorNotFound
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
		return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	var actor m.Actor
	found := false

	for rows.Next() {
		found = true
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
			return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
		return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return m.Actor{}, m.ErrActorNotFound
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
			return m.Actor{}, fmt.Errorf("%w: invalid field: %s", m.ErrInvalidInput, field)
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
		return m.Actor{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return actor, nil
}

func (r actorRepository) FindActorsByName(name string) ([]m.ActorWithoutMoviesDto, error) {
	query := `
		SELECT
	    	a.id,
	     	a.name,
	      	a.birth_date,
	       	m.id,
	    	m.title,
	        m.release_year,
	        m.duration
		FROM actors a
		LEFT JOIN movie_actors am ON a.id = am.actor_id
		LEFT JOIN movies m ON am.movie_id = m.id
		WHERE a.name LIKE ? COLLATE NOCASE
		ORDER BY a.id
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	actorsMap := make(map[int]*m.ActorWithoutMoviesDto)

	for rows.Next() {
		var actorID int
		var actorName string
		var birthDate string

		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var releaseYear sql.NullInt64
		var duration sql.NullInt64

		err := rows.Scan(
			&actorID,
			&actorName,
			&birthDate,
			&movieID,
			&movieTitle,
			&releaseYear,
			&duration,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		actor, exists := actorsMap[actorID]

		if !exists {
			actor = &m.ActorWithoutMoviesDto{
				Id:        actorID,
				Name:      actorName,
				BirthDate: birthDate,
				Movies:    []m.MovieWithoutActorsDto{},
			}

			actorsMap[actorID] = actor
		}

		if movieID.Valid {
			actor.Movies = append(actor.Movies, m.MovieWithoutActorsDto{
				Id:          int(movieID.Int64),
				Title:       movieTitle.String,
				ReleaseYear: int(releaseYear.Int64),
				Duration:    int(duration.Int64),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if len(actorsMap) == 0 {
		return []m.ActorWithoutMoviesDto{}, m.ErrActorNotFound
	}

	actors := make([]m.ActorWithoutMoviesDto, 0, len(actorsMap))

	for _, actor := range actorsMap {
		actors = append(actors, *actor)
	}

	return actors, nil
}

func (r actorRepository) CountActors() (int, error) {
	var count int

	query := `
			SELECT COUNT(*)
			FROM actors
		`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
