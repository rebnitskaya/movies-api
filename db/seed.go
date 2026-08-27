package db

import (
	"database/sql"
	"fmt"
)

func SeedDatabase(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// -------------------------------------------------------------------------
	// Genres
	// -------------------------------------------------------------------------

	genres := []string{
		"Action",
		"Adventure",
		"Animation",
		"Comedy",
		"Crime",
		"Drama",
		"Fantasy",
		"Horror",
		"Science Fiction",
		"Thriller",
	}

	for _, name := range genres {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO genres (name)
			VALUES (?)
		`, name)
		if err != nil {
			return fmt.Errorf("insert genre %q: %w", name, err)
		}
	}

	// -------------------------------------------------------------------------
	// Actors
	// -------------------------------------------------------------------------

	actors := []struct {
		name      string
		birthdate string
	}{
		{"Tom Hanks", "1956-07-09"},
		{"Tim Allen", "1953-06-13"},
		{"Robin Wright", "1966-04-08"},
		{"Morgan Freeman", "1937-06-01"},
		{"Gary Sinise", "1955-03-17"},
		{"Harrison Ford", "1942-07-13"},
		{"Mark Hamill", "1951-09-25"},
		{"Carrie Fisher", "1956-10-21"},
		{"Liam Neeson", "1952-06-07"},
		{"Ewan McGregor", "1971-03-31"},
		{"Natalie Portman", "1981-06-09"},
		{"Keanu Reeves", "1964-09-02"},
		{"Laurence Fishburne", "1961-07-30"},
		{"Hugo Weaving", "1960-04-04"},
		{"Leonardo DiCaprio", "1974-11-11"},
		{"Joseph Gordon-Levitt", "1981-02-17"},
		{"Tom Hardy", "1977-09-15"},
		{"Christian Bale", "1974-01-30"},
		{"Heath Ledger", "1979-04-04"},
		{"Michael Caine", "1933-03-14"},
		{"Matt Damon", "1970-10-08"},
		{"Ben Affleck", "1972-08-15"},
		{"Robin Williams", "1951-07-21"},
		{"Dustin Hoffman", "1937-08-08"},
		{"Brad Pitt", "1963-12-18"},
		{"Edward Norton", "1969-08-18"},
		{"Helena Bonham Carter", "1966-05-26"},
		{"Johnny Depp", "1963-06-09"},
		{"Orlando Bloom", "1977-01-13"},
		{"Elijah Wood", "1981-01-28"},
		{"Ian McKellen", "1939-05-25"},
		{"Viggo Mortensen", "1958-10-20"},
		{"Cate Blanchett", "1969-05-14"},
		{"Sean Astin", "1971-02-25"},
		{"Daniel Radcliffe", "1989-07-23"},
		{"Emma Watson", "1990-04-15"},
		{"Rupert Grint", "1988-08-24"},
		{"Alan Rickman", "1946-02-21"},
		{"Robert Downey Jr.", "1965-04-04"},
		{"Chris Evans", "1981-06-13"},
		{"Scarlett Johansson", "1984-11-22"},
		{"Chris Hemsworth", "1983-08-11"},
		{"Mark Ruffalo", "1967-11-22"},
		{"Samuel L. Jackson", "1948-12-21"},
		{"John Travolta", "1954-02-18"},
		{"Uma Thurman", "1970-04-29"},
		{"Bruce Willis", "1955-03-19"},
		{"Milla Jovovich", "1975-12-17"},
		{"Sigourney Weaver", "1949-10-08"},
		{"Bill Murray", "1950-09-21"},
	}

	for _, actor := range actors {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO actors (name, birth_date)
			VALUES (?, ?)
		`, actor.name, actor.birthdate)
		if err != nil {
			return fmt.Errorf("insert actor %q: %w", actor.name, err)
		}
	}

	// -------------------------------------------------------------------------
	// Movies
	// -------------------------------------------------------------------------

	movies := []struct {
		title    string
		year     int
		duration int
		genres   []string
		actors   []string
	}{
		{
			"Forrest Gump", 1994, 142,
			[]string{"Drama", "Comedy"},
			[]string{"Tom Hanks", "Robin Wright", "Gary Sinise"},
		},
		{
			"Toy Story", 1995, 81,
			[]string{"Animation", "Comedy", "Adventure"},
			[]string{"Tom Hanks", "Tim Allen"},
		},
		{
			"Toy Story 2", 1999, 92,
			[]string{"Animation", "Comedy", "Adventure"},
			[]string{"Tom Hanks", "Tim Allen"},
		},
		{
			"Toy Story 3", 2010, 103,
			[]string{"Animation", "Comedy", "Adventure"},
			[]string{"Tom Hanks", "Tim Allen"},
		},
		{
			"Saving Private Ryan", 1998, 169,
			[]string{"Drama", "Action"},
			[]string{"Tom Hanks", "Matt Damon", "Tom Sizemore"},
		},
		{
			"The Green Mile", 1999, 189,
			[]string{"Drama", "Fantasy", "Crime"},
			[]string{"Tom Hanks", "Michael Clarke Duncan", "Sam Rockwell"},
		},
		{
			"Cast Away", 2000, 143,
			[]string{"Drama", "Adventure"},
			[]string{"Tom Hanks", "Helen Hunt"},
		},
		{
			"The Terminal", 2004, 128,
			[]string{"Comedy", "Drama"},
			[]string{"Tom Hanks", "Catherine Zeta-Jones"},
		},
		{
			"Star Wars: Episode IV - A New Hope", 1977, 121,
			[]string{"Science Fiction", "Adventure", "Action"},
			[]string{"Harrison Ford", "Mark Hamill", "Carrie Fisher"},
		},
		{
			"Star Wars: Episode V - The Empire Strikes Back", 1980, 124,
			[]string{"Science Fiction", "Adventure", "Action"},
			[]string{"Harrison Ford", "Mark Hamill", "Carrie Fisher"},
		},
		{
			"Star Wars: Episode VI - Return of the Jedi", 1983, 131,
			[]string{"Science Fiction", "Adventure", "Action"},
			[]string{"Harrison Ford", "Mark Hamill", "Carrie Fisher"},
		},
		{
			"Indiana Jones and the Raiders of the Lost Ark", 1981, 115,
			[]string{"Action", "Adventure"},
			[]string{"Harrison Ford"},
		},
		{
			"Indiana Jones and the Temple of Doom", 1984, 118,
			[]string{"Action", "Adventure"},
			[]string{"Harrison Ford"},
		},
		{
			"The Matrix", 1999, 136,
			[]string{"Science Fiction", "Action", "Thriller"},
			[]string{"Keanu Reeves", "Laurence Fishburne", "Hugo Weaving"},
		},
		{
			"The Matrix Reloaded", 2003, 138,
			[]string{"Science Fiction", "Action"},
			[]string{"Keanu Reeves", "Laurence Fishburne", "Hugo Weaving"},
		},
		{
			"The Matrix Revolutions", 2003, 129,
			[]string{"Science Fiction", "Action"},
			[]string{"Keanu Reeves", "Laurence Fishburne", "Hugo Weaving"},
		},
		{
			"Inception", 2010, 148,
			[]string{"Science Fiction", "Action", "Thriller"},
			[]string{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Tom Hardy"},
		},
		{
			"The Dark Knight", 2008, 152,
			[]string{"Action", "Crime", "Drama", "Thriller"},
			[]string{"Christian Bale", "Heath Ledger", "Michael Caine"},
		},
		{
			"The Dark Knight Rises", 2012, 164,
			[]string{"Action", "Crime", "Drama"},
			[]string{"Christian Bale", "Tom Hardy", "Michael Caine"},
		},
		{
			"Batman Begins", 2005, 140,
			[]string{"Action", "Crime", "Drama"},
			[]string{"Christian Bale", "Michael Caine"},
		},
		{
			"Interstellar", 2014, 169,
			[]string{"Science Fiction", "Drama", "Adventure"},
			[]string{"Matthew McConaughey", "Anne Hathaway", "Michael Caine"},
		},
		{
			"Good Will Hunting", 1997, 126,
			[]string{"Drama"},
			[]string{"Matt Damon", "Ben Affleck", "Robin Williams"},
		},
		{
			"The Martian", 2015, 144,
			[]string{"Science Fiction", "Adventure", "Drama"},
			[]string{"Matt Damon", "Jessica Chastain"},
		},
		{
			"Ocean's Eleven", 2001, 116,
			[]string{"Crime", "Comedy", "Thriller"},
			[]string{"George Clooney", "Brad Pitt", "Matt Damon"},
		},
		{
			"Fight Club", 1999, 139,
			[]string{"Drama", "Thriller"},
			[]string{"Brad Pitt", "Edward Norton", "Helena Bonham Carter"},
		},
		{
			"Se7en", 1995, 127,
			[]string{"Crime", "Drama", "Thriller"},
			[]string{"Brad Pitt", "Morgan Freeman"},
		},
		{
			"The Shawshank Redemption", 1994, 142,
			[]string{"Drama", "Crime"},
			[]string{"Morgan Freeman", "Tim Robbins"},
		},
		{
			"Bruce Almighty", 2003, 101,
			[]string{"Comedy", "Fantasy"},
			[]string{"Morgan Freeman", "Jim Carrey"},
		},
		{
			"Pirates of the Caribbean: The Curse of the Black Pearl", 2003, 143,
			[]string{"Adventure", "Fantasy", "Action"},
			[]string{"Johnny Depp", "Orlando Bloom", "Keira Knightley"},
		},
		{
			"Pirates of the Caribbean: Dead Man's Chest", 2006, 151,
			[]string{"Adventure", "Fantasy", "Action"},
			[]string{"Johnny Depp", "Orlando Bloom", "Keira Knightley"},
		},
		{
			"The Lord of the Rings: The Fellowship of the Ring", 2001, 178,
			[]string{"Fantasy", "Adventure", "Drama"},
			[]string{"Elijah Wood", "Ian McKellen", "Viggo Mortensen", "Sean Astin"},
		},
		{
			"The Lord of the Rings: The Two Towers", 2002, 179,
			[]string{"Fantasy", "Adventure", "Drama"},
			[]string{"Elijah Wood", "Ian McKellen", "Viggo Mortensen", "Sean Astin"},
		},
		{
			"The Lord of the Rings: The Return of the King", 2003, 201,
			[]string{"Fantasy", "Adventure", "Drama"},
			[]string{"Elijah Wood", "Ian McKellen", "Viggo Mortensen", "Sean Astin"},
		},
		{
			"Harry Potter and the Philosopher's Stone", 2001, 152,
			[]string{"Fantasy", "Adventure"},
			[]string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint", "Alan Rickman"},
		},
		{
			"Harry Potter and the Chamber of Secrets", 2002, 161,
			[]string{"Fantasy", "Adventure"},
			[]string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint", "Alan Rickman"},
		},
		{
			"Harry Potter and the Prisoner of Azkaban", 2004, 142,
			[]string{"Fantasy", "Adventure", "Drama"},
			[]string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"},
		},
		{
			"Iron Man", 2008, 126,
			[]string{"Action", "Science Fiction", "Adventure"},
			[]string{"Robert Downey Jr.", "Gwyneth Paltrow", "Jeff Bridges"},
		},
		{
			"The Avengers", 2012, 143,
			[]string{"Action", "Science Fiction", "Adventure"},
			[]string{"Robert Downey Jr.", "Chris Evans", "Scarlett Johansson", "Chris Hemsworth"},
		},
		{
			"Captain America: The First Avenger", 2011, 124,
			[]string{"Action", "Adventure", "Science Fiction"},
			[]string{"Chris Evans", "Samuel L. Jackson"},
		},
		{
			"Thor", 2011, 115,
			[]string{"Action", "Fantasy", "Adventure"},
			[]string{"Chris Hemsworth", "Natalie Portman", "Samuel L. Jackson"},
		},
		{
			"The Avengers: Age of Ultron", 2015, 141,
			[]string{"Action", "Science Fiction", "Adventure"},
			[]string{"Robert Downey Jr.", "Chris Evans", "Chris Hemsworth", "Mark Ruffalo"},
		},
		{
			"Guardians of the Galaxy", 2014, 121,
			[]string{"Action", "Science Fiction", "Adventure", "Comedy"},
			[]string{"Chris Pratt", "Bradley Cooper", "Vin Diesel"},
		},
		{
			"Pulp Fiction", 1994, 154,
			[]string{"Crime", "Drama", "Thriller"},
			[]string{"John Travolta", "Samuel L. Jackson", "Uma Thurman"},
		},
		{
			"Kill Bill: Vol. 1", 2003, 111,
			[]string{"Action", "Crime", "Thriller"},
			[]string{"Uma Thurman", "Lucy Liu", "David Carradine"},
		},
		{
			"Die Hard", 1988, 132,
			[]string{"Action", "Thriller"},
			[]string{"Bruce Willis", "Alan Rickman"},
		},
		{
			"The Fifth Element", 1997, 126,
			[]string{"Science Fiction", "Action", "Adventure"},
			[]string{"Bruce Willis", "Milla Jovovich", "Gary Oldman"},
		},
		{
			"Resident Evil", 2002, 100,
			[]string{"Action", "Horror", "Science Fiction"},
			[]string{"Milla Jovovich", "Michelle Rodriguez"},
		},
		{
			"Alien", 1979, 117,
			[]string{"Science Fiction", "Horror", "Thriller"},
			[]string{"Sigourney Weaver", "Ian Holm"},
		},
		{
			"Aliens", 1986, 137,
			[]string{"Science Fiction", "Horror", "Action"},
			[]string{"Sigourney Weaver", "Michael Biehn"},
		},
		{
			"Groundhog Day", 1993, 101,
			[]string{"Comedy", "Fantasy"},
			[]string{"Bill Murray", "Andie MacDowell"},
		},
		{
			"Ghostbusters", 1984, 105,
			[]string{"Comedy", "Fantasy"},
			[]string{"Bill Murray", "Dan Aykroyd", "Sigourney Weaver"},
		},
		{
			"Jumanji", 1995, 104,
			[]string{"Adventure", "Fantasy", "Comedy"},
			[]string{"Robin Williams", "Kirsten Dunst"},
		},
	}

	// -------------------------------------------------------------------------
	// Insert movies and relationships
	// -------------------------------------------------------------------------

	for _, movie := range movies {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO movies (title, release_year, duration)
			VALUES (?, ?, ?)
		`, movie.title, movie.year, movie.duration)

		if err != nil {
			return fmt.Errorf("insert movie %q: %w", movie.title, err)
		}

		// MovieGenre
		for _, genre := range movie.genres {
			_, err := tx.Exec(`
				INSERT OR IGNORE INTO genres_movies (movie_id, genre_id)
				SELECT m.id, g.id
				FROM movies m, genres g
				WHERE m.title = ?
				  AND m.release_year = ?
				  AND g.name = ?
			`, movie.title, movie.year, genre)

			if err != nil {
				return fmt.Errorf(
					"insert genre relationship for %q: %w",
					movie.title,
					err,
				)
			}
		}

		// MovieActor
		for _, actor := range movie.actors {
			_, err := tx.Exec(`
				INSERT OR IGNORE INTO movie_actors (movie_id, actor_id)
				SELECT m.id, a.id
				FROM movies m, actors a
				WHERE m.title = ?
				  AND m.release_year = ?
				  AND a.name = ?
			`, movie.title, movie.year, actor)

			if err != nil {
				return fmt.Errorf(
					"insert actor relationship for %q: %w",
					movie.title,
					err,
				)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}
