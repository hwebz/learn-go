## Backend Development with Go
Youtube Playlist: https://www.youtube.com/watch?v=TtCfDXfSw_0&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE

### Requirements
1. sqlc
2. Docker
3. golang-migrate ([Guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
4. sqlc ([Guide](https://docs.sqlc.dev/en/latest/overview/install.html))

### Init sqlc
> sqlc init

### Generate code using sqlc
> make sqlc

### How to migrate database
```bash
# 1. Install golang-migrate first
brew install golang-migrate

# 2. Run make up to start pg database
make up

# 3. To run a migration
make migrateup

# 4. To reverse a migration
make migratedown
```