# OpenOnlineClinic

A microservice-based electronic health record (EHR) management system with REST/gRPC API support.

## Installation

- **Language**: Go 1.21+
- **API**: 
  - REST (Gin Framework)
  - gRPC (protobuf v3)
- **Database**: PostgreSQL + GORM
- **Containerization**: Docker
- **Documentation**: Swagger (planned)
- **Logging**: Zap
- **Validation**: go-playground/validator

## Features 
1. **Patient Management**:
   - CRUD operations
   - Search with pagination
2. **Medical Data**:
   - Allergies (add/remove)
   - Insurance policies
   - Doctor prescriptions
3. **Integrations**:
   - Connection with User service
   - External system support via gRPC

### Requirements:
- Docker 20.10+

### Developer installation
Prepare golang
```bash
go mod tidy
```

Prepare docker
```bash
docker compose build
```

Run docker
```bash
docker compose up
```
```bash
# 1. Clone the repository
git clone https://github.com/Ruletk/OnlineClinic.git
cd OnlineClinic

# 2. Start services
docker-compose up -d --build

# 3. Verify the setup
curl http://localhost:8080/patients

##   Run Unit tests

make test-pkg

make test-service


## Git naming

Use [convential commits](https://www.conventionalcommits.org/ru/v1.0.0/)!

- `feature/feature-name` - New feature
- `fix/jira-number` - Bugfix
- `chore/update` - Very minor changes
- `docs/module` - New documentation for module
- `test/module` - New or update tests for module

Do not add more than 2-3 prefixes in your commits! Split it on several commits.
