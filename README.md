# Spy Cat Agency Management System

A comprehensive CRUD application for managing spy cats, missions, and targets. Built with Go, Gin framework, and PostgreSQL.

## Features

### Spy Cats Management
- Create, read, update, and delete spy cats
- Validate cat breeds using TheCatAPI
- Track cat availability and salary information
- Years of experience tracking

### Mission Management
- Create missions with 1-3 targets
- Assign cats to missions
- Complete missions when all targets are finished
- Mission status tracking

### Target Management
- Add/remove targets from missions
- Update target information
- Mark targets as complete
- Add and update spy notes
- Notes are frozen when target/mission is completed

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Setup and Run

Prereqs: Go, Bash, Docker Desktop (Docker Compose v2 included)

1) Setup (installs swag if missing, generates docs, starts DB):
```bash
make setup
```

2) Start API (defaults to port 3030):
```bash
make start
```

Optional: override port
```bash
PORT=8080 make start
```

The API will be available at `http://localhost:3030` (or the port you chose).

### API Documentation

Once the application is running, you can access:
- **Swagger UI**: `http://localhost:3030/swagger/index.html`
- **API Base URL**: `http://localhost:3030/api/v1`

## API Endpoints

### Spy Cats
- `POST /api/v1/cats` - Create a new spy cat
- `GET /api/v1/cats` - List all spy cats
- `GET /api/v1/cats/{id}` - Get a specific spy cat
- `PUT /api/v1/cats/{id}` - Update a spy cat's salary
- `DELETE /api/v1/cats/{id}` - Delete a spy cat

### Missions
- `POST /api/v1/missions` - Create a new mission
- `GET /api/v1/missions` - List all missions
- `GET /api/v1/missions/{id}` - Get a specific mission
- `PUT /api/v1/missions/{id}` - Update a mission
- `DELETE /api/v1/missions/{id}` - Delete a mission
- `PUT /api/v1/missions/{id}/assign` - Assign a cat to a mission
- `PUT /api/v1/missions/{id}/complete` - Complete a mission

### Targets
- `POST /api/v1/missions/{missionId}/targets` - Add a target to a mission
- `PUT /api/v1/missions/{missionId}/targets/{id}` - Update a target
- `DELETE /api/v1/missions/{missionId}/targets/{id}` - Delete a target
- `PUT /api/v1/missions/{missionId}/targets/{id}/complete` - Complete a target
- `PUT /api/v1/missions/{missionId}/targets/{id}/notes` - Update target notes

## Example Usage

### Create a Spy Cat
```bash
curl -X POST http://localhost:3030/api/v1/cats \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Whiskers",
    "years_experience": 5,
    "breed": "Persian",
    "salary": 50000.0
  }'
```

### Create a Mission
```bash
curl -X POST http://localhost:3030/api/v1/missions \
  -H "Content-Type: application/json" \
  -d '{
    "cat_id": 1,
    "targets": [
      {
        "name": "John Doe",
        "country": "USA"
      },
      {
        "name": "Jane Smith",
        "country": "Canada"
      }
    ]
  }'
```

### Update Target Notes
```bash
curl -X PUT http://localhost:3030/api/v1/missions/1/targets/1/notes \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Target was seen at the coffee shop. Wearing a blue jacket."
  }'
```

## Development

### Project Structure
```
├── cmd/                    # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migrations
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models and DTOs
│   ├── repository/        # Data access layer
│   └── services/          # Business logic layer
├── docs/                  # Swagger documentation
├── docker-compose.yml     # Database setup
└── Makefile              # Build and run commands
```

### Available Commands
- `make docker-up` - Start PostgreSQL database
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make docker-down` - Stop PostgreSQL database
- `make swagger` - Generate Swagger documentation

### Environment Variables
- `DATABASE_URL` - PostgreSQL connection string (default: postgres://spycat:password@localhost:5432/spycat_agency?sslmode=disable)
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (default: development)

## Stopping the Application

To stop the database:
```bash
make docker-down
```
