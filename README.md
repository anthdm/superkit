# GOTHKIT
Create interactive applications with Golang, HTMX, and Templ


## Table of content
- [GOTHKIT](#gothkit)
	- [Table of content](#table-of-content)
- [Installation](#installation)
- [Getting started](#getting-started)
	- [Development server](#development-server)
	- [Hot reloading the browser](#hot-reloading-the-browser)
- [Migrations](#migrations)
	- [Create a new migration](#create-a-new-migration)
	- [Migrate the database](#migrate-the-database)
	- [Reset the database](#reset-the-database)
	- [Seeds](#seeds)
- [Testing](#testing)
	- [Testing handlers](#testing-handlers)

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
## Development server
Run the development server with the following command:
```
make dev 
```

## Hot reloading the browser
Hot reloading is configured by default when running your application in development.

> NOTE: on windows you might need to run `make assets` in another terminal for god knows why.

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


# Testing
## Testing handlers


