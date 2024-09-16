# Private Dart Pub Repo

Implementation of [Dart Pub Repository Spec V2](https://github.com/dart-lang/pub/blob/master/doc/repository-spec-v2.md) written in Golang.

What will this cover:

- able to upload pub library to custom remote server
- able to authenticate before downloading pub library as dependency
- enable both server token management and user token management

## Requirement

- Golang
- Postgresql

## How to run

1. install go
2. install toolset:

   - air-verse/air: `go install github.com/air-verse/air@latest`
   - python3: use any method at disposal, but make sure `python` command is linked to python3, since air_build.py need it

3. setup `.env`, see `.env.example`, self explanatory enough to be copied to `.env`
4. run `go mod download`
5. run `air fx` / `air manual`

## Migration

This repository support both auto migration and manual migration. Both are useful.
on development, you can use auto migration, by setting `DB_AUTOMIGRATION=true` in `.env`.
But please make sure to set it to false in production, since it's sometime problematic.

### Manual Migration - Installation

Requirement:

- Docker or alternatives (podman)

```
⚠️ If using podman, you need to add alias to `docker` command.
```

```
⚠️ In windows, powershell alias won't work for podman, create `docker.bat` with content:
@echo off
podman %*

then add it to directory that's registered in PATH
```

On linux/mac, just follow this instruction: https://atlasgo.io/docs

On windows, I suggest this step:

- Download windows binary from Manual Installation tab in https://atlasgo.io/docs
- Put it in a folder, for example C:\tools
- rename it to `atlas.exe`

### Manual Migration - Sync migration file with AutoMigration Model

- Add your automigration model to `loader/main.go`
- run `atlas migrate diff --env gorm`
- your updated sql will be ready in migrations folder

### Manual Migration - Add your sql

- Add your sql file in migrations directory in format `yyyymmddHHiiss.sql`
- Insert the migration query. this tool is quite different to be honest, there is no down query, just add what you want to add
- run `atlas migrate hash`

### Manual Migration - Run migration

- run `atlas migrate apply --url "yourdatabaseurl"`
  - example: `atlas migrate apply --url "postgres://postgres@127.0.0.1:5432/golang?sslmode=disable"`

## Process Monitoring

Open `{{BASE_URL}}/monitor` in browser to see. You can disable this by removing the route in `monitor/route.go`

## Opentelemetry

auto integrated for gofiber endpoint and database performance. need to use `monitorService.StartTrace` to add more nexted context for monitoring clarity.

### How to collect and see the traces

Known choices:

- [jaeger](https://www.jaegertracing.io/docs/1.60/getting-started)
- [tempo+grafana](https://github.com/grafana/tempo/tree/main/example/docker-compose/local)
