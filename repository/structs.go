package repository

type Genre struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies"`
}

type Actor struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	BirthDate string  `json:"birthDate"`
	Movies    []Movie `json:"movies"`
}

type Movie struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear int     `json:"release_year"`
	Duration    int     `json:"duration"`
	Genres      []Genre `json:"genres"`
	Actors      []Actor `json:"actors"`
}
