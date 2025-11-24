

### pgsql 访问 （docker）
```bash
docker exec -it light-postgresql psql -U postgres -d light_admin

docker exec -it light-postgresql psql -U postgres -d postgres -c "\dt"
```
