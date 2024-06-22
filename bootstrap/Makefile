# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifeq ($(OS),Windows_NT)
	MAIN_PATH = /tmp/bin/main.exe
	SYNC_ASSETS_COMMAND =	@go run github.com/makiuchi-d/arelo@v1.13.1 \
	--target "./public" \
	--pattern "**/*.js" \
	--pattern "**/*.css" \
	--delay "100ms" \
	--templ generate --notify-proxy
else
	MAIN_PATH = tmp/bin/main
	SYNC_ASSETS_COMMAND =	@go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "public" \
	--build.include_ext "js,css" \
	--screen.clear_on_rebuild true \
	--log.main_only true
endif

# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
templ:
	@templ generate --watch --proxy="http://localhost$(HTTP_LISTEN_ADDR)" --open-browser=false

# run air to detect any go file changes to re-build and re-run the server.
server:
	@go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build --tags dev -o ${MAIN_PATH} ./cmd/app/" --build.bin "${MAIN_PATH}" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--screen.clear_on_rebuild true \
	--log.main_only true

# run tailwindcss to generate the styles.css bundle in watch mode.
watch-assets:
	@npx tailwindcss -i app/assets/app.css -o ./public/assets/styles.css --watch   

# run esbuild to generate the index.js bundle in watch mode.
watch-esbuild:
	@npx esbuild app/assets/index.js --bundle --outdir=public/assets --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
sync_assets:
	${SYNC_ASSETS_COMMAND}

# start the application in development
dev:
	@make -j5 templ server watch-assets watch-esbuild sync_assets

# build the application for production. This will compile your app
# to a single binary with all its assets embedded.
build:
	@npx tailwindcss -i app/assets/app.css -o ./public/assets/styles.css
	@npx esbuild app/assets/index.js --bundle --outdir=public/assets
	@go build -o bin/app_prod cmd/app/main.go
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"

db-status:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) status

db-reset:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) reset

db-down:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) down

db-up:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) up

db-mig-create:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) create $(filter-out $@,$(MAKECMDGOALS)) sql

db-seed:
	@go run cmd/scripts/seed/main.go
