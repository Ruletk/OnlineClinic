module github.com/Ruletk/OnlineClinic/pkg/database

go 1.23.8

require (
	github.com/Ruletk/OnlineClinic/pkg/config v0.0.0
	github.com/stretchr/testify v1.10.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.26.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Ruletk/OnlineClinic/pkg/config => ../config
