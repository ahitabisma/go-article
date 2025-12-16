## Buat Table Migration (User Service)

```
migrate create -ext sql -dir database/migrations -seq create_users_table
migrate create -ext sql -dir database/migrations -seq create_roles_table
migrate create -ext sql -dir database/migrations -seq create_user_role_table
```

## Migration Up

```
migrate -database "postgres://postgres:postgres@localhost:5432/ms-user-service?sslmode=disable" -path database/migrations up
```
