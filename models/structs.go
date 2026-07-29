package models

type Genre struct {
	Id   int64  `json:"id"` //SQLite and database return int64
	Name string `json:"name"`
}

type Actor struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

type Movie struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear int     `json:"releaseYear"`
	Duration    int     `json:"duration"`
	Genres      []Genre `json:"genres"`
	Actors      []Actor `json:"actors"`
}
