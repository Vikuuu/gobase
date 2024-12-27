# gobase

> Automatic schema detection that creates the migrations for your *Go* project.

Write your schema in the Go struct format and migrate it to the database.
This project is to help your database migration and automate the writing of migration file code.

<img align="center" src="assets/img/gobase-img.png" alt="gobase-img" height=30% width=30%>

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Database](#database)
- [Limitation](#limitation)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install this project on your machine you should have *Go* installed in your machine, if not follow this [official step](https://go.dev/doc/install) and install it in your machine.

Now installing Gobase:
```bash
go install github.com/Vikuuu/gobase/cmd/gobase@1.0.1
```

## Usage

To use the gobase you must define a Yaml file in the root of your project directory named `gobase.yml`. The content of the yaml file will be as follows:
```yaml
database: "sqlite3"
schema_data: "./schema/users.go"
migration: "./migrations"
```

- Database: It will tell what database to use.
- Schema Data: This will point to the file in which you have defined your schema in form of struct.
- Migration: The directory in which all your migration files will be stored.

After defining the yaml file and creating a struct, to migrate up to the database use the following command:
```bash
gobase migrate
```
This command will create a migration file in your defined directory of migration, and migrate it to the database.

If your file that defines the struct that you want to convert it to database table, your file will look like this:
```go
package database

type users struct {
    ID   int
    Name string
}
```

Then the corresponding Migration file will look like this:
```sql
-- Up Migration

CREATE TABLE users (
    id INTEGER,
    name TEXT
);


-- Down Migration

DROP TABLE users;
```

Now if you make any changes to your struct, like:
```go
package database

import "time"

type users struct {
    ID        int
    Name      string
    CreatedAt time.Time
}
```
And then again run the `migrate` command like:
```bash
gobase migrate
```

Then the new migration file will look like this:
```sql
-- Up Migration

ALTER TABLE users
ADD COLUMN created_at TIMESTAMP;

-- Down Migration

ALTER TABLE users DROP COLUMN created_at;
```

To migrate down use the following command:
```bash
gobase migrate down
```

## Database

Current the project only supports:
1. Sqlite3

## Limitations

Currently the project has many limitations:
1. Works only on the single table or struct.
2. No table relations.
3. Supports only a single file in the schema directory.
4. Cannot migrate down or up to any desired version, can only migrate down only 1 step.


## Contributing

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-name`
3. Make your changes.
4. Push your branch: `git push origin feature-name`
5. Create a pull request.

## License

This project in licensed under the [MIT License](LICENSE)
