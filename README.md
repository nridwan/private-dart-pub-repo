# Private Dart Pub Repo

Implementation of [Dart Pub Repository Spec V2](https://github.com/dart-lang/pub/blob/master/doc/repository-spec-v2.md) written in Golang.

What will this cover:

- able to upload pub library to custom remote server
- able to authenticate before downloading pub library as dependency
- enable user token management

## Requirement

- Golang
- Postgresql
- S3 compatible storage: AWS S3 / Google Cloud Storage / Alicloud OSS / MinIO / etc
- SMTP server (if you need good user management, since I give admin no control over changing password, you can enable it via code though)

## Deployment
*Disclaimer: sorry for lack of documentation, gotta go fast

### Initial setup
- please update modules/user/seeder.go, change the initial admin account as you desire.

### Via docker/kubernetes/etc
- Check `docker-compose.example.yml`
- Understand the flow, adjust your needs (database, storage, migration)
- use the dockerfile as base of main app docker image
- correctly setup environment variables/config maps (check `.env.example` / `docker-compose.example.yml` to see available envs)
- Once running, if you need to seed first admin, open the docker shell, run `/pubserver db:seed`

### Via manual build
- check `Dockerfile` to see the build command chain
- adjust and make your own build script
- prepare environment
- run the server by using `<executablename> fx` 
  - also, if you need to seed first admin, run `<executablename> db:seed`

## API docs
- Open [docs directory](/docs/)
- Import `pub_server.postman_collection.json` and `pub_server.postman_environment.json` in Postman app
- Use and change BASE_URL in the environment depends on the deployment url


### User - Login
Used to get access token & refresh token to authenticate with endpoints

Endpoint: `Users > Login` (`POST` | `{{BASE_URL}}/v1/users/login`)

Body Params:
- email
- password

Steps:
1. Insert registered (or seeded) user email & password
2. `response_output.detail.access_token` field can be used to access other endpoint 
   except `Pub API` and `User > Forgot Password` as `Authorization` header, using format `Bearer <token>`
3. `response_output.detail.refresh_token` can be used in [User - Refresh Token](#user---refresh-token), using format `Bearer <token>`

### User - Refresh Token
Used to get access token & refresh token without the needs to input email & password again, as long as the previous refresh token is not expired

Endpoint : `Users > Refresh Token` (`POST` | `{{BASE_URL}}/v1/users/refresh`)

Headers :
- Authorization : Bearer token (refresh token)

1. Use refresh token as Auth bearer header, then hit the endpoint

### User - Forgot Password
Used to change existing user's password, since admin has no control over it except for the first password

Endpoints:
- `Users > Forgot Password - Forgot Password OTP` (`POST` | `{{BASE_URL}}/v1/users/forgot-password/otp`)
  - Body Params:
    - email - registered user email
  - Steps:
    - Insert valid email, hit endpoint
    - Check email's inbox, there is OTP, use it in `Users > Forgot Password - Create Password`. There is also time limit specified in the email.
- `Users > Forgot Password - Create Password` (`POST` | `{{BASE_URL}}/v1/users/forgot-password/create-password`)
  - Body Params:
    - email - registered user email
    - otp - OTP sent through email
    - password - new password
  - Steps:
    - Insert needed parameters, hit endpoint
    - Password changed, can be used to [Login](#user---login)

### Admin - User's CRUD
Used to manage registered user in the system

Restriction: Logged In User must be an admin

Endpoints:
- `Users > Create` (`POST` | `{{BASE_URL}}/v1/users`)
  - Header: 
    - Authorization: Bearer token
  - Body Params:
    - name - user name
    - email - registered user email
    - password - user initial password
    - is_admin - whether new user can access admin restricted endpoints or not
    - can_write - whether new user can create / update pub token as write access (can publish package)
  - Steps:
    - Insert valid email, hit endpoint
    - Newly created user can be used
- `Users > List` (`GET` | `{{BASE_URL}}/v1/users`)
  - Header: 
    - Authorization: Bearer token
  - Query params:
    - page: starts from 1, required
    - limit: data fetched per page, required
    - search: search by name / email, optional
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return list of Users
- `Users > Detail` (`GET` | `{{BASE_URL}}/v1/users/{id}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - id: registered user id
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return User detail
- `Users > Update` (`PUT` | `{{BASE_URL}}/v1/users/{id}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - id: registered user id
  - Body Params:
    - name - user name, optional
    - email - registered user email, optional
    - is_admin - whether new user can access admin restricted endpoints or not, optional
    - can_write - whether new user can create / update pub token as write access (can publish package), optional
  - Steps:
    - Insert needed parameters, hit endpoint
    - Filled parameter should be updated (can partially update user)
    - Note: when changing can_write to false, all that user's pub token write access will be updated to false.
      when updating can_write to true again, respective user must re-enable their token write access via `Pub Token > Update`, or create new pub token.
- `Users > Delete` (`DELETE` | `{{BASE_URL}}/v1/users/{id}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - id: registered user id
  - Steps:
    - Insert needed parameters, hit endpoint
    - User should be deleted and cannot be used to login

### Pub Token
This feature is needed to manage token. user can only create writable access token (write=true) when their user's can_write flag is true. 

- `Pub Token > Create` (`POST` | `{{BASE_URL}}/v1/pubtoken`)
  - Header: 
    - Authorization: Bearer token
  - Body Params:
    - remarks - what the token will be used for, required. please make sure to give meaningful name to sort it out later when revoking
    - write - is the token has capabilities to publish dependencies. can only be filled true if users has `can_write=true`
    - expired at - final day the token can be used, format: `YYYY-MM-DD`
  - Steps:
    - Insert valid email, hit endpoint
    - Will return newly created token, can be used to pull / publish dependencies
- `Pub Token > List` (`GET` | `{{BASE_URL}}/v1/pubtoken`)
  - Header: 
    - Authorization: Bearer token
  - Query params:
    - page: starts from 1, required
    - limit: data fetched per page, required
    - search: search by remarks, optional
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return list of Pub Token
- `Pub Token > Detail` (`GET` | `{{BASE_URL}}/v1/pubtoken/{id}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - id: pub token id
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return 1 Pub Token data
- `Pub Token > Update` (`PUT` | `{{BASE_URL}}/v1/pubtoken/{id}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - id: pub token id
  - Body Params:
    - write - is the token has capabilities to publish dependencies. can only be filled true if users has `can_write=true`
  - Steps:
    - Insert needed parameters, hit endpoint
    - Filled parameter should be updated (can partially update user)
    - Note: when changing can_write to false, all that user's pub token write access will be updated to false.
      when updating can_write to true again, respective user must re-enable their token write access via `Pub Token > Update`, or create new pub token.
- `Pub Token > Delete` (`DELETE` | `{{BASE_URL}}/v1/pubtoken/{id}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - id: pub token id
  - Steps:
    - Insert needed parameters, hit endpoint
    - User should be deleted and cannot be used to login
    - Deleted token access will be automatically revoked

### Pub > Pub API
- This is manual step to upload package to storage. Will be used by dart tool to manage publishing
- Based on [Pub Repository Spec v2](https://github.com/dart-lang/pub/blob/master/doc/repository-spec-v2.md) and inspired by [Unpub](https://github.com/pd4d10/unpub)
- Every first upload is considered private library, to make it public, admin must update visibility using [Pub > Query](#pub--query)
- Public access can see and use non-private library

Endpoints:
- `Pub > API > Package Version List` (`GET` | `{{BASE_URL}}/v1/pub/packages/:package`)
  - Header: 
    - Authorization: Bearer `<PUBTOKEN>`
      - optional, but when using token, user will be able to see private libraries
      - note: `<PUBTOKEN>` should be changed with token generated using [Pub Token > Create](#pub-token)
  - Path parameter:
    - package: package name
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return Pub library and its versions
- `Pub > API > Package Version Detail` (`GET` | `{{BASE_URL}}/v1/pub/packages/:package/versions/:version`)
  - Header: 
    - Authorization: Bearer `<PUBTOKEN>`
      - optional, but when using token, user will be able to see private libraries
      - note: `<PUBTOKEN>` should be changed with token generated using [Pub Token > Create](#pub-token)
  - Path parameter:
    - package: package name
    - version: version name (semver, example: `1.0.0`)
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return specific version info of library
    - deprecated, newer version of dart tool won't need this, but it is created for backward compatible with old dart version
- `Pub > API > Get Upload URL` (`GET` | `{{BASE_URL}}/v1/pub/packages/versions/new`)
  - Header: 
    - Authorization: Bearer `<PUBTOKEN>`
      - required
      - note: `<PUBTOKEN>` should be changed with token generated using [Pub Token > Create](#pub-token)
  - Restriction:
    - PUBTOKEN should have write access
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return `{{BASE_URL}}/v1/pub/api/packages/versions/newUpload`
- `Pub > API > Upload` (`POST` | `{{BASE_URL}}/v1/pub/packages/versions/newUpload`)
  - Header: 
    - Authorization: Bearer `<PUBTOKEN>`
      - required
      - note: `<PUBTOKEN>` should be changed with token generated using [Pub Token > Create](#pub-token)
  - Restriction:
    - PUBTOKEN should have write access
  - Form data (multipart):
    - file: tar.gz file from dart tool that's supposed to be uploaded here.
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return redirect to `{{BASE_URL}}/v1/pub/packages/versions/newUploadFinish`. if error, will bring error message as query parameter `error`.
    - On redirected endpoint, it will return success/error response depending on upload status
    - If success, package version will be inserted.

### Pub > Query
- Public API to see existing pub library uploaded and its version

Restriction: public can only see public package. to see private package, need to login as user

Endpoints:
- `Pub > Query > Package List` (`GET` | `{{BASE_URL}}/v1/pub/query/packages`)
  - Header: 
    - Authorization: Bearer token
  - Query params:
    - page: starts from 1, required
    - limit: data fetched per page, required
    - search: search by package name, optional
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return list of Pub libraries uploaded to the server
- `Pub > Query > Update Package` (`PUT` | `{{BASE_URL}}/v1/pub/query/packages/{package}`)
  - Header: 
    - Authorization: Bearer token
  - Restriction:
    - Only user with admin access can use this feature
  - Path parameter:
    - package: package name (field name)
  - Body Params:
    - private - is the package can only be seen by logged in user / with token access, optional
  - Steps:
    - Insert needed parameters, hit endpoint
    - package visibility should be updated
- `Pub > Query > Version List` (`GET` | `{{BASE_URL}}/v1/pub/query/packages`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - package: package name (field name)
  - Query params:
    - page: starts from 1, required
    - limit: data fetched per page, required
    - search: search by version name, optional
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return list of Pub Versions
- `Pub > Query > Version Detail` (`GET` | `{{BASE_URL}}/v1/pub/query/packages/{package}/versions/{version}`)
  - Header: 
    - Authorization: Bearer token
  - Path parameter:
    - package: package name (field name)
    - version: version name (semver, example: `1.0.0`)
  - Steps:
    - Insert needed parameters, hit endpoint
    - Will return Detail of Pub Version (including changelog & readme)

## User Guides
After successfully run the service we can use the APIs for multiple scenario.

### Admin - Manage User
Generic flow:
1. Login using [Login](#user---login) endpoint.
2. Admin can read / create / update / delete user by using [Admin - User's CRUD](#admin---users-crud)

Deleting user: 
1. use [Users > Delete](#admin---users-crud) endpoint to delete user alongside with their tokens. 
all of the user's pub tokens will be disabled

Revoke write access: 
1. use [Users > Update](#admin---users-crud) endpoint, set `can_write: true`. 
2. all of the user's pub tokens will be automatically set to `write: false`

### Admin - Manage Package Visibility
1. Login using [Login](#user---login) endpoint.
2. Use [Pub > Query > Update Package](#pub--query) endpoint to update package visibility

### User - Setup pub token
Using this method, user will be able to pull / publish libraries

1. Login using [Login](#user---login) endpoint.
2. Create token using [Pub Token > Create](#pub-token), set write depending on usage and permission 
   (do you have access to publish? do you need the token to publish package?)
3. copy the token
4. Follow register token step from [pub.dev documentation](https://dart.dev/tools/pub/cmd/pub-token#add-a-credential-for-the-current-session).
   - Steps:
     - run in terminal/cmd: `dart pub token add {{BASE_URL}}/v1/pub`
     - enter the token copied in step 3
   - Notes:
     - `{{BASE_URL}}` need to be changed to where this app is hosted, for example:`http://localhost:4000`, then it becomes `http://localhost:4000/v1/pub`
     - on windows, probably need to edit manually in `%APPDATA%/dart/pub-credentials.json` due to terminal character limit
     - On postman, this will be applied for Pub API endpoints automatically after creation

### User - Revoke token access
1. Login using [Login](#user---login) endpoint.
2. Delete token using [Pub Token > Delete](#pub-token)

## Development

### How to run

1. install go
2. install toolset:

   - air-verse/air: `go install github.com/air-verse/air@latest`
   - python3: use any method at disposal, but make sure `python` command is linked to python3, since air_build.py need it

3. setup `.env`, see `.env.example`, self explanatory enough to be copied to `.env`
4. run `go mod download`
5. run `air fx` / `air manual`

### Migration

This repository support both auto migration and manual migration. Both are useful.
on development, you can use auto migration, by setting `DB_AUTOMIGRATION=true` in `.env`.
But please make sure to set it to false in production, since it's sometime problematic.

#### Manual Migration - Installation

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

#### Manual Migration - Sync migration file with AutoMigration Model

- Add your automigration model to `loader/main.go`
- run `atlas migrate diff --env gorm`
- your updated sql will be ready in migrations folder

#### Manual Migration - Add your sql

- Add your sql file in migrations directory in format `yyyymmddHHiiss.sql`
- Insert the migration query. this tool is quite different to be honest, there is no down query, just add what you want to add
- run `atlas migrate hash`

#### Manual Migration - Run migration

- run `atlas migrate apply --url "yourdatabaseurl"`
  - example: `atlas migrate apply --url "postgres://postgres@127.0.0.1:5432/golang?sslmode=disable"`

### Opentelemetry

auto integrated for gofiber endpoint and database performance. need to use `monitorService.StartTrace` to add more nexted context for monitoring clarity.
can be disabled depending on environment.

#### How to collect and see the traces

Known choices:

- [jaeger](https://www.jaegertracing.io/docs/1.60/getting-started)
- [tempo+grafana](https://github.com/grafana/tempo/tree/main/example/docker-compose/local)
