# GOTHKIT
Create interactive applications with Golang, HTMX, and Templ


## Table of content
- [GOTHKIT](#gothkit)
	- [Table of content](#table-of-content)
- [Installation](#installation)
- [Getting started](#getting-started)
	- [Application structure](#application-structure)
		- [assets](#assets)
		- [conf](#conf)
		- [db](#db)
		- [events](#events)
		- [handlers](#handlers)
		- [types](#types)
		- [views](#views)
	- [Development server](#development-server)
	- [Hot reloading the browser](#hot-reloading-the-browser)
- [Migrations](#migrations)
	- [Create a new migration](#create-a-new-migration)
	- [Migrate the database](#migrate-the-database)
	- [Reset the database](#reset-the-database)
	- [Seeds](#seeds)
- [Validations](#validations)
- [Testing](#testing)
	- [Testing handlers](#testing-handlers)
- [Create a production release](#create-a-production-release)

# Installation
```
go install github.com/anthdm/gothkit@master
```

After installation you can create a new project by running: 
```
gothkit [myprojectname]
```

You can now navigate to your project:
```
cd [myprojectname]
```

# Getting started
## Application structure
### assets
### conf
### db
### events
### handlers
### types
### views

## Development server
Run the development server with the following command:
```
make dev 
```

## Hot reloading the browser
Hot reloading is configured by default when running your application in development.

> NOTE: on windows you might need to run `make assets` to watch for CSS and JS changes in another terminal.

# Migrations
## Create a new migration
```
make db-mig add_user_table
```

Will create a new migration SQL file inside `app/db/migrations`

## Migrate the database 
```
make db-up
```

## Reset the database 
```
make db-reset
```

## Seeds
```
make db-seed
```
This command will run the seeds file located at `cmd/scripts/seed/main.go`

# Validations

# Testing
## Testing handlers

# Create a production release
Gothkit will compile your whole application including its assets into a single binary. To build your application for production you can run the following command:
```
make build
```
This will create a binary file located at  `/bin/app_release`.

Make sure you also set the correct application environment variable in your `.env` file.
```
APP_ENV	= production
```


