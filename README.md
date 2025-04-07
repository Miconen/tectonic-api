# Tectonic API

Robust RESTful API designed to serve as a bridge between our frontend applications and our backend database. Our original goal was to provide a more robust, and efficient solution to handle our frontends, with a particular focus on separating business logic from the frontend for improved maintainability and flexibility.

As we are in the process of rewriting our current frontend, Tectonic API will be our new backend. By providing a well-structured and comprehensive set of endpoints, Tectonic API allows the frontend to focus solely on displaying data and user interactions, while the backend handles the complex business logic.

## Key areas

- **RESTful API**: Tectonic API adheres to the principles of Representational State Transfer (REST), a standard architectural style for networked applications. This makes it easy to use and understand, and allows for seamless integration with various frontend technologies.

- **Separation of Concerns**: By encapsulating the business logic within the API, we ensure that the frontend remains clean and focused on user interactions. This separation of concerns improves maintainability, scalability, and flexibility.

- **Expandability**: Tectonic API is designed to handle multiple frontends, making it a scalable solution for our growing projects.

## Installation

```
git clone https://github.com/yourusername/tectonic-api.git
cd tectonic-api
docker compose --profile dev up --build
```

> [!NOTE]
> Only Docker is required to run the project locally.

## Endpoints

The API provides several endpoints that can be viewed through Swagger. `{baseUrl}/swagger/v1/`

## Testing

### Unit tests

> [!TIP]
> Our unit tests are a part of our automatic CI/CD pipeline!

I've placed a priority on learning and utilizing testing for speeding up my own development. As we move closer towards desired functionality, we hope to achieve broader test coverage.
We're helped on this by this being a rewrite of an existing application, thus making the specs clear to us.

### Integration tests info

> [!TIP]
> Our integration tests are also a part of our automatic CI/CD pipeline!

The API has a lot of interplay and dependency on existing data.
For this we use Golangs native integration testing. We also have a basic Postman pipeline for running through our endpoints which we are phasing out currently.

### Go integration tests
#### Before running Go integration tests you will need...

- Include a `DATABASE_URL` variable in your environment/shell. (`DATABASE_URL=postgresql://postgres:postgres@localhost:5432/tectonic`)
- Make sure your database is running

#### Running the Go integration tests

Run `go test ./routes`

### Postman integration tests (Deprecated)
#### Before running Postman you will need...

- The project running, either locally or hosted.
- Import the `Integrations.postman_collection.json` file.
- Specify a global `baseUrl` variable in Postman. For running locally through Docker: `localhost:8080/api`.
- Specify a global `Authorization` variable in Postman. For running locally through Docker: `123`.

##### Running the Postman tests

To run this test, import it into Postman and select `Run Collection`.

## Frontends

You can find the frontend for this project at [Tectonic Bot](https://github.com/Miconen/tectonic-bot).

## Tech stack

- Go ([mux](https://github.com/gorilla/mux), [pgx](https://github.com/jackc/pgx/v5), [Squirrel](https://github.com/Masterminds/squirrel), [swaggo](https://github.com/swaggo/swag))
- PostgreSQL
- Docker

## Contributing

Currently, the project is maintained by a team of three developers. While we are not actively seeking contributions, we welcome anyone interested in contributing to review our project board and issues.
