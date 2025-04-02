# OpenOnlineClinic

## Installation
...
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


## Git naming

Use [convential commits](https://www.conventionalcommits.org/ru/v1.0.0/)!

- `feature/feature-name` - New feature
- `fix/jira-number` - Bugfix
- `chore/update` - Very minor changes
- `docs/module` - New documentation for module
- `test/module` - New or update tests for module

Do not add more than 2-3 prefixes in your commits! Split it on several commits.