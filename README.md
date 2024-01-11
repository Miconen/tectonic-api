# Tectonic API

Robust RESTful API designed to serve as a bridge between our frontend applications and our backend database. Our original goal was to provide a more robust, and efficient solution to handle our frontends, with a particular focus on separating business logic from the frontend for improved maintainability and flexibility.

As we are in the process of rewriting our current frontend, Tectonic API will be our new backend. By providing a well-structured and comprehensive set of endpoints, Tectonic API allows the frontend to focus solely on displaying data and user interactions, while the backend handles the complex business logic.

## Key areas

* **RESTful API**: Tectonic API adheres to the principles of Representational State Transfer (REST), a standard architectural style for networked applications. This makes it easy to use and understand, and allows for seamless integration with various frontend technologies.

* **Separation of Concerns**: By encapsulating the business logic within the API, we ensure that the frontend remains clean and focused on user interactions. This separation of concerns improves maintainability, scalability, and flexibility.

* **Expandability**: Tectonic API is designed to handle multiple frontends, making it a scalable solution for our growing projects.


## Installation
```
git clone https://github.com/yourusername/tectonic-api.git
cd tectonic-api
docker compose up --build
```

> [!NOTE]
> Only Docker is required to run the project locally.

## Endpoints

The API provides several endpoints:

* `/v1/`: Serves Swagger UI.
* `/v1/user`: Handles user data with GET, POST, and DELETE methods.
* `/v1/users`: Retrieves all users with a GET method.
* `/v1/rsn`: Handles RSN data with GET, POST, and DELETE methods.
* `/v1/guild`: Handles guild data with GET, POST, and DELETE methods.
* `/v1/leaderboard`: Retrieves leaderboard data with a GET method.
* `/v1/time`: Handles time data with POST and DELETE methods.
* `/v1/guild/times`: Updates times channel with a PUT method.
* `/v1/guild/multiplier`: Updates multiplier with a PUT method.

## Frontends

You can find the frontend for this project at [Tectonic Bot](https://github.com/Miconen/tectonic-bot).

## Tech stack

* Go ([mux](https://github.com/gorilla/mux), [pgx](https://github.com/jackc/pgx/v5), [Squirrel](https://github.com/Masterminds/squirrel), [swaggo](https://github.com/swaggo/swag))
* PostgreSQL
* Docker

## Contributing

Currently, the project is maintained by a team of three developers. While we are not actively seeking contributions, we welcome anyone interested in contributing to review our project board and issues.
