# Interview-Focused Rebuild Guide

## Goal

Rebuild the project as an `admin-only` backend that is easy to explain in interviews.

Focus on:

- configuration and boot flow
- database setup
- JWT authentication
- RBAC with roles, menus, departments, and API permissions
- clean HTTP handlers and services
- a small but real test suite

Do not start by cloning every module in this repository. The first version should be small, coherent, and demonstrably correct.

## What "Done" Looks Like

By the end of this guide, you should be able to demo:

- `go-admin -h`
- `go-admin config -c config/settings.yml`
- `go-admin server -c config/settings.yml`
- `POST /api/v1/login`
- authenticated access to core admin endpoints
- permission allowed and denied behavior for different roles
- passing unit and API tests

## Weekly Schedule

## Week 1: Project Skeleton and CLI

### Focus

Get the project compiling. Keep the CLI simple and working.

### Build Order

1. `go.mod`
2. `config/settings.yml`
3. `config/extend.go`
4. `common/global/adm.go`
5. `common/global/topic.go`
6. `common/global/logo.go`
7. `common/global/casbin.go`
8. `main.go`
9. `cmd/cobra.go`
10. `cmd/version/server.go`
11. `cmd/config/server.go`

### Milestone Demo

- `go build` succeeds
- `./go-admin -h` prints help
- `./go-admin version` prints a version
- `./go-admin config -c config/settings.yml` loads and prints config

### Test Targets

- smoke test for CLI boot
- config load test for valid and invalid config paths

## Week 2: Shared Models and Database Setup

### Focus

Create the shared data types and database initialization path used by the rest of the app.

### Build Order

1. `common/models/by.go`
2. `common/models/type.go`
3. `common/models/user.go`
4. `common/models/response.go`
5. `common/models/menu.go`
6. `common/models/migrate.go`
7. `common/database/open.go`
8. `common/database/open_sqlite3.go`
9. `common/database/initialize.go`
10. `common/storage/initialize.go`
11. `common/service/service.go`
12. `common/apis/api.go`
13. `common/ip.go`
14. `common/response/binding.go`

### Milestone Demo

- server code can load config and establish a DB connection
- base response and model types compile cleanly

### Test Targets

- DB factory test for supported drivers
- service base test for error aggregation
- binding error formatting test

## Week 3: DTOs and CRUD Primitives

### Focus

Build the reusable query and CRUD helpers before touching business modules.

### Build Order

1. `common/dto/type.go`
2. `common/dto/pagination.go`
3. `common/dto/search.go`
4. `common/dto/order.go`
5. `common/dto/generate.go`
6. `common/dto/auto_form.go`
7. `common/actions/type.go`
8. `common/actions/permission.go`
9. `common/actions/index.go`
10. `common/actions/view.go`
11. `common/actions/create.go`
12. `common/actions/update.go`
13. `common/actions/delete.go`

### Milestone Demo

- a sample GORM query can paginate, filter, and order correctly
- permission scope logic can be applied as a GORM scope

### Test Targets

- pagination scope test
- search-tag parser test
- order builder test
- permission scope test for each data-scope mode

## Week 4: Middleware and Authentication

### Focus

Build the request pipeline and make login work end to end.

### Build Order

1. `common/middleware/request_id.go`
2. `common/middleware/logger.go`
3. `common/middleware/db.go`
4. `common/middleware/customerror.go`
5. `common/middleware/header.go`
6. `common/middleware/demo.go`
7. `common/middleware/handler/user.go`
8. `common/middleware/handler/role.go`
9. `common/middleware/handler/login.go`
10. `common/middleware/handler/auth.go`
11. `common/middleware/handler/ping.go`
12. `common/middleware/auth.go`
13. `common/middleware/permission.go`
14. `common/middleware/init.go`

### Milestone Demo

- request ID appears in logs
- `/info` responds
- `POST /api/v1/login` returns a token
- protected routes reject missing or invalid tokens

### Test Targets

- login handler test for good credentials
- login handler test for bad credentials
- JWT middleware test
- role enforcement test
- middleware chain smoke test with `httptest`

## Week 5: Server Boot and Admin Core Foundation

### Focus

Boot a real server and define the admin core models.

### Build Order

1. `cmd/api/server.go`
2. `app/admin/models/casbin_rule.go`
3. `app/admin/models/datascope.go`
4. `app/admin/models/sys_user.go`
5. `app/admin/models/sys_role.go`
6. `app/admin/models/sys_menu.go`
7. `app/admin/models/sys_dept.go`
8. `app/admin/models/sys_api.go`

### Milestone Demo

- `./go-admin server -c config/settings.yml` starts cleanly
- server responds on `/info`
- database tables for core admin models can be migrated or created

### Test Targets

- model hook test for password hashing
- model serialization test for user and role payloads
- boot smoke test for server startup and shutdown

## Week 6: Users and Roles

### Focus

Finish the two most important interview modules first: users and roles.

### Build Order

1. `app/admin/service/dto/sys_user.go`
2. `app/admin/service/sys_user.go`
3. `app/admin/apis/sys_user.go`
4. `app/admin/router/sys_user.go`
5. `app/admin/service/dto/sys_role.go`
6. `app/admin/service/sys_role.go`
7. `app/admin/service/sys_role_menu.go`
8. `app/admin/apis/sys_role.go`
9. `app/admin/router/sys_role.go`

### Milestone Demo

- create, list, update, and delete users
- create, list, update, and delete roles
- assign menus to roles
- user password reset and profile endpoints work

### Test Targets

- `sys_user` service tests for CRUD and password changes
- `sys_role` service tests for CRUD and role-menu assignment
- API tests for `sys_user` handlers with `httptest`
- permission regression test for forbidden updates

## Week 7: Menus, Departments, and API Permissions

### Focus

Complete the RBAC story so you can explain access control clearly.

### Build Order

1. `app/admin/service/dto/sys_menu.go`
2. `app/admin/service/sys_menu.go`
3. `app/admin/apis/sys_menu.go`
4. `app/admin/router/sys_menu.go`
5. `app/admin/service/dto/sys_dept.go`
6. `app/admin/service/sys_dept.go`
7. `app/admin/apis/sys_dept.go`
8. `app/admin/router/sys_dept.go`
9. `app/admin/service/dto/sys_api.go`
10. `app/admin/service/sys_api.go`
11. `app/admin/apis/sys_api.go`
12. `app/admin/router/sys_api.go`

### Milestone Demo

- menu tree endpoints work
- department tree endpoints work
- API permission registry works
- one role can access an endpoint that another role cannot

### Test Targets

- menu tree test
- department tree test
- API registration test
- end-to-end permission test with two roles

## Week 8: Router Wiring, Captcha, and Interview Polish

### Focus

Wire the admin module together, add minimal tests that prove the system works, and prepare the demo story.

### Build Order

1. `app/admin/router/router.go`
2. `app/admin/router/sys_router.go`
3. `app/admin/router/init_router.go`
4. `app/admin/apis/captcha.go`
5. `app/admin/service/sys_user_test.go`
6. `app/admin/service/sys_role_test.go`
7. `app/admin/apis/sys_user_api_test.go`
8. `common/middleware/handler/auth_test.go`

### Milestone Demo

- full admin router boots
- captcha endpoint works
- login plus protected CRUD flow works
- tests pass in front of another engineer

### Test Targets

- service tests for user and role logic
- API integration tests for user endpoints
- auth handler tests
- final smoke test for login -> token -> protected route

## Demo Script for Interviews

Use this order in interviews:

1. explain the architecture in one minute
2. show config and boot path
3. show login flow and JWT issuance
4. show user CRUD
5. show role and menu assignment
6. show permission allowed and denied behavior
7. show the test suite

## Optional Files to Add Later

Add these only after the core admin story is solid.

### Optional Config Variants

- `config/settings.sqlite.yml`
- `config/settings.demo.yml`
- `config/settings.full.yml`

### Optional Admin Features

- `app/admin/apis/go_admin.go`
- `app/admin/models/initdb.go`
- `app/admin/models/sys_post.go`
- `app/admin/service/dto/sys_post.go`
- `app/admin/service/sys_post.go`
- `app/admin/apis/sys_post.go`
- `app/admin/router/sys_post.go`
- `app/admin/models/sys_config.go`
- `app/admin/service/dto/sys_config.go`
- `app/admin/service/sys_config.go`
- `app/admin/apis/sys_config.go`
- `app/admin/router/sys_config.go`
- `app/admin/models/sys_dict_type.go`
- `app/admin/service/dto/sys_dict_type.go`
- `app/admin/service/sys_dict_type.go`
- `app/admin/apis/sys_dict_type.go`
- `app/admin/models/sys_dict_data.go`
- `app/admin/service/dto/sys_dict_data.go`
- `app/admin/service/sys_dict_data.go`
- `app/admin/apis/sys_dict_data.go`
- `app/admin/router/sys_dict.go`
- `app/admin/models/sys_opera_log.go`
- `app/admin/service/dto/sys_opera_log.go`
- `app/admin/service/sys_opera_log.go`
- `app/admin/apis/sys_opera_log.go`
- `app/admin/router/sys_opera_log.go`
- `app/admin/models/sys_login_log.go`
- `app/admin/service/dto/sys_login_log.go`
- `app/admin/service/sys_login_log.go`
- `app/admin/apis/sys_login_log.go`
- `app/admin/router/sys_login_log.go`

### Optional File Storage

- `common/file_store/interface.go`
- `common/file_store/initialize.go`
- `common/file_store/oss.go`
- `common/file_store/obs.go`
- `common/file_store/kodo.go`

### Optional Tools Module

- `cmd/api/other.go`
- `app/other/models/tools/db_tables.go`
- `app/other/models/tools/db_columns.go`
- `app/other/models/tools/sys_tables.go`
- `app/other/models/tools/sys_columns.go`
- `app/other/service/dto/sys_tables.go`
- `app/other/apis/file.go`
- `app/other/apis/sys_server_monitor.go`
- `app/other/apis/tools/db_tables.go`
- `app/other/apis/tools/db_columns.go`
- `app/other/apis/tools/sys_tables.go`
- `app/other/apis/tools/gen.go`
- `app/other/router/file.go`
- `app/other/router/gen_router.go`
- `app/other/router/init_router.go`
- `app/other/router/monitor.go`
- `app/other/router/router.go`
- `app/other/router/sys_server_monitor.go`

### Optional Jobs Module

- `cmd/api/jobs.go`
- `app/jobs/type.go`
- `app/jobs/jobbase.go`
- `app/jobs/examples.go`
- `app/jobs/models/sys_job.go`
- `app/jobs/service/dto/sys_job.go`
- `app/jobs/service/sys_job.go`
- `app/jobs/apis/sys_job.go`
- `app/jobs/router/router.go`
- `app/jobs/router/int_router.go`
- `app/jobs/router/sys_job.go`

### Optional Migration and Scaffolding

- `cmd/migrate/server.go`
- `cmd/migrate/migration/init.go`
- `cmd/migrate/migration/version-local/doc.go`
- `cmd/migrate/migration/version/1599190683659_tables.go`
- `cmd/migrate/migration/version/1653638869132_migrate.go`
- `cmd/app/server.go`
- `template/cmd_api.template`
- `template/router.template`
- `template/migrate.template`
- `template/api_migrate.template`

### Optional Dev and Repo Polish

- `Dockerfile`
- `Makefile`
- `docker-compose.yml`
- `README.md`
- `README.Zh-cn.md`
