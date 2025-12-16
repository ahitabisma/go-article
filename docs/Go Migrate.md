## Buat Table Migration (User Service)

```
migrate create -ext sql -dir database/migrations -seq create_users_table
migrate create -ext sql -dir database/migrations -seq create_roles_table
migrate create -ext sql -dir database/migrations -seq create_user_role_table
```

## Migration Up

```
migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/belajar_golang?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations up
```

## Rollback 1 Step

```
migrate -path database/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/belajar_golang?charset=utf8mb4&parseTime=True&loc=Local" down 1
```