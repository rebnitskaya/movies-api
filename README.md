# Movies API

A REST API written in Go for managing movies, actors, and genres. The project uses SQLite for persistence, the standard `net/http` package for the HTTP server, and a layered architecture that separates HTTP handling, business logic, and database access.

## Features

- Create, read, update, and delete movies
- Create, read, update, and delete actors
- Create, read, update, and delete genres
- Add actors to movies and remove them
- Add genres to movies and remove them
- Search movies by title
- Filter movies by genre, release year, or actor
- Paginate list responses
- Seed the database with initial data using the `--init` flag
- Serve Swagger documentation
- Centralized error handling
- Request timeout and panic recovery middleware


## Architecture

The application follows a layered structure:

```text
HTTP request
  -> handler
  -> service
  -> repository
  -> database/sql
  -> SQLite
```

Each layer has a separate responsibility:

| Layer | Responsibility |
|---|---|
| `main` | Application entry point |
| `cli` | CLI flag parsing |
| `server` | Server creation, dependency wiring, and route registration |
| `handler` | HTTP request/response handling |
| `service` | Business logic and validation |
| `repository` | Database queries and persistence |
| `db` | SQLite connection, schema creation, and seed data |
| `models` | Domain models, DTOs, validation, and application errors |
| `helper` | Shared helpers such as pagination and error handling |
| `middleware` | HTTP middleware such as timeout and panic recovery |

## Project Flow

When a client sends a request, the flow is:

```text
Client
  -> route in server.RegisterRoutes
  -> handler method
  -> service method
  -> repository method
  -> SQLite database
  -> JSON response
```

For example, creating an actor follows this path:

```text
POST /api/actors
  -> ActorHandler.PostActor
  -> ActorService.CreateActor
  -> ActorRepository.FindActorByNameAndBirthDate
  -> ActorRepository.CreateActor
  -> JSON response with 201 Created
```

## Database Schema

The database contains three main entity tables:

- `actors`
- `genres`
- `movies`

It also contains two join tables for many-to-many relationships:

- `movie_actors`
- `genres_movies`

The relationships use foreign keys with `ON DELETE CASCADE`, so related rows in join tables are removed automatically when a movie, actor, or genre is deleted.

```text
actors <-> movie_actors <-> movies
genres <-> genres_movies <-> movies
```

## Getting Started

### Prerequisites

- Go installed
- SQLite support through `github.com/mattn/go-sqlite3`

### Install dependencies

```bash
go mod tidy
```

### Run the server

```bash
go run .
```

### Run the server with seed data

```bash
go run . --init
```

By default, the server runs at:

```text
http://localhost:8080
```

The API base path is:

```text
/api
```

## Swagger

Swagger UI is available at:

```text
http://localhost:8080/swagger/
```

## Pagination

List endpoints support pagination with query parameters:

| Parameter | Description | Default |
|---|---|---|
| `page` | Page number | `1` |
| `size` | Number of items per page | `5` |

Example:

```http
GET /api/movies?page=2&size=10
```

## Endpoints

### Actors

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/actors` | Get a paginated list of actors |
| `GET` | `/api/actors?name={name}` | Search actors by name |
| `GET` | `/api/actors/{id}` | Get an actor by ID |
| `POST` | `/api/actors` | Create an actor |
| `PATCH` | `/api/actors/{id}` | Update an actor |
| `DELETE` | `/api/actors/{id}` | Delete an actor |
| `DELETE` | `/api/actors/{id}?force=true` | Force delete an actor with movie associations |

Example request body for creating an actor:

```json
{
  "name": "Tom Hanks",
  "birthDate": "1956-07-09"
}
```

### Movies

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/movies` | Get a paginated list of movies |
| `GET` | `/api/movies?genre={id}` | Get movies by genre ID |
| `GET` | `/api/movies?year={year}` | Get movies by release year |
| `GET` | `/api/movies?actor={id}` | Get movies by actor ID |
| `GET` | `/api/movies/search?title={text}` | Search movies by title |
| `GET` | `/api/movies/{id}` | Get a movie by ID |
| `POST` | `/api/movies` | Create a movie |
| `PATCH` | `/api/movies/{id}` | Update a movie |
| `DELETE` | `/api/movies/{id}?force=true` | Delete a movie |
| `GET` | `/api/movies/{id}/actors` | Get all actors in a movie |
| `POST` | `/api/movies/{movieID}/actors/{actorID}` | Add an actor to a movie |
| `DELETE` | `/api/movies/{movieID}/actors/{actorID}` | Remove an actor from a movie |
| `POST` | `/api/movies/{movieID}/genres/{genreID}` | Add a genre to a movie |
| `DELETE` | `/api/movies/{movieID}/genres/{genreID}` | Remove a genre from a movie |

Example request body for creating a movie:

```json
{
  "title": "Inception",
  "releaseYear": 2010,
  "duration": 148,
  "actors": [1, 2, 3],
  "genres": [1, 2]
}
```

### Genres

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/genres` | Get a paginated list of genres |
| `GET` | `/api/genres/{id}` | Get a genre by ID |
| `POST` | `/api/genres` | Create a genre |
| `PATCH` | `/api/genres/{id}` | Update a genre |
| `DELETE` | `/api/genres/{id}` | Delete a genre |
| `DELETE` | `/api/genres/{id}?force=true` | Force delete a genre with movie associations |

Example request body for creating a genre:

```json
{
  "name": "Drama"
}
```

## Error Handling

Handlers return errors instead of writing every error response directly. The `helper.GlobalErrorHandler` wrapper converts application errors into HTTP responses.

| Application error | HTTP status |
|---|---|
| `ErrActorNotFound` | `404 Not Found` |
| `ErrMovieNotFound` | `404 Not Found` |
| `ErrGenreNotFound` | `404 Not Found` |
| `ErrInvalidInput` | `400 Bad Request` |
| `ErrBadRequest` | `400 Bad Request` |
| `ErrMovieHasBeenMadeBefore` | `400 Bad Request` |
| Unknown error | `500 Internal Server Error` |

## Example Workflow

1. Create a genre:

```http
POST /api/genres
```

2. Create an actor:

```http
POST /api/actors
```

3. Create a movie:

```http
POST /api/movies
```

4. Add the actor to the movie:

```http
POST /api/movies/{movieID}/actors/{actorID}
```

5. Add the genre to the movie:

```http
POST /api/movies/{movieID}/genres/{genreID}
```

6. Fetch the movie with its relationships:

```http
GET /api/movies/{id}
```

## Notes

- The server listens on `localhost:8080`.
- The database file is created as `movies_api.db`.
- The `--init` flag fills the database with predefined movies, actors, genres, and relationships.
- The project uses repository interfaces, which makes the service layer easier to test with mock implementations.
