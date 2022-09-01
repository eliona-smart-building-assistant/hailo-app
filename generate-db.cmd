go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

go get github.com/volatiletech/sqlboiler/v4
go get github.com/volatiletech/null/v8

docker run --rm -d ^
    --name "eliona_database_sql_boiler_code_generation" ^
    -e "POSTGRES_PASSWORD=secret" ^
    -p "60001:5432" ^
    debezium/postgres:12

docker run --rm ^
    --name "eliona_database_init_sql_boiler_code_generation" ^
    -e "PGPORT=60001" ^
    -e "PGDATABASE=postgres" ^
    -e "PGHOST=host.docker.internal" ^
    -e "PGUSER=postgres" ^
    -e "PGPASSWORD=secret" ^
    eliona.azurecr.io/core/database:develop

docker exec ^
    "eliona_database_sql_boiler_code_generation" ^
    psql -d "postgres://postgres:secret@localhost:5432/iot" -c "ALTER TABLE alarm drop CONSTRAINT alarm_asset_id_subtype_fkey;"

docker exec ^
    "eliona_database_sql_boiler_code_generation" ^
    psql -d "postgres://postgres:secret@localhost:5432/iot" -c "DROP SCHEMA IF EXISTS api CASCADE;"

docker exec ^
    "eliona_database_sql_boiler_code_generation" ^
    psql -d "postgres://postgres:secret@localhost:5432/iot" -c "ALTER TABLE public.pipeline DROP COLUMN cov, DROP COLUMN lpf, DROP COLUMN hws, DROP COLUMN ala;"

sqlboiler psql ^
    -c sqlboiler-public.toml ^
    --wipe --no-tests

sqlboiler psql ^
    -c sqlboiler-versioning.toml ^
    --wipe --no-tests

docker stop "eliona_database_sql_boiler_code_generation"