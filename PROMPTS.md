# PROMPTS.md

## 1. Golang Migration
- " Buatkan swagger json dari code berikut"
-  "1. Create a new migration:"
```bash
go run cmd/main.go migrate create <migration_name>
```

-  "2. Run pending migrations (up):"

```bash
go run cmd/main.go migrate up
```
-  "3. Rollback the last migration (down):"
```bash
go run cmd/main.go migrate down
```

-  "4. Reset all migrations:"
```bash
go run cmd/main.go migrate reset
```

-  "5. Check current migration version:"
```bash
go run cmd/main.go migrate version
```


## 1. Golang Seed
- "How implement seeding in this project"
- "Seed all tables with sample data"
```bash
go run cmd/main.go seed users
```
- "Seed only users table"
```bash
go run cmd/main.go seed users
```
- "Clear all data from users table"
```bash
go run cmd/main.go seed clear
```
- "Clear and re-seed users table (useful for development)"
```bash
go run cmd/main.go seed refresh
```