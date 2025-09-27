# Database Seeders

This directory contains database seeding functionality for the MajooTest-Golang project.

## Overview

Seeders are used to populate the database with sample or test data. They are useful for:
- Development environment setup
- Testing scenarios
- Demo data preparation

## Usage

Seeders are executed through the CLI application:

```bash
# From the project root directory
go run cmd/main.go seed <target>
```

## Available Commands

### Seed All Tables
```bash
go run cmd/main.go seed all
```
Seeds all available tables with sample data.

### Seed Users Table
```bash
go run cmd/main.go seed users
```
Populates the users table with sample user accounts.

### Clear Users Table
```bash
go run cmd/main.go seed clear
```
Removes all data from the users table.

### Refresh Users Table
```bash
go run cmd/main.go seed refresh
```
Clears the users table and then populates it with fresh sample data.

## Sample Data

### Users
The users seeder creates the following sample accounts:
- Alice Johnson (alice@example.com)
- Bob Smith (bob@example.com)
- Charlie Brown (charlie@example.com)
- Diana Prince (diana@example.com)
- John Doe (john@example.com)
- Jane Smith (jane@example.com)

All sample users have the password: `password123`

## Adding New Seeders

To add seeders for new tables:

1. Create a new function in the appropriate file (or create a new file in the `seeders/` directory)
2. Follow the naming convention: `Seed<TableName>(db *sql.DB) error`
3. Add the seeder to the switch statement in `cmd/commands/seed.go`
4. Update this README with the new seeder information

## Notes

- Seeders use `ON CONFLICT DO NOTHING` to prevent duplicate data insertion
- All passwords are hashed using bcrypt before storage
- Seeders require the database migrations to be run first
- Use the `refresh` command during development to reset test data