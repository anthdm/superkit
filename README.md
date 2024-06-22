# SUPERKIT

Build high-performance apps swiftly with minimal team resources in Go programming language. **SUPERKIT** is a full-stack web framework that provides a set of tools and libraries to help you build web applications quickly and efficiently. **SUPERKIT** is built on top of the Go programming language and is designed to be simple and easy to use.

> The project (for now) is in the **experimental** phase.

## Table of content

- [SUPERKIT](#superkit)
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
- [Creating views with Templ](#creating-views-with-templ)
- [Validations](#validations)
- [Testing](#testing)
  - [Testing handlers](#testing-handlers)
- [Create a production release](#create-a-production-release)

## Installation

To create a new **SUPERKIT** project, you can run the following command:

```bash
# Create your SUPERKIT project in a single command:
go run github.com/anthdm/superkit@master [yourprojectname]

# You can now navigate to your project:
cd [myprojectname]

# Run npm install to install both tailwindcss and esbuild locally.
npm install

# If you run into dependency issues you can run:
go clean -modcache && go get -u ./...  

# If you have the authentication plugin enabled you need to migrate your database.
make db-up 
```

## Getting started

## Application structure

The **SUPERKIT** project structure is designed to be simple and easy to understand. The project structure is as follows:

```bash
├── bootstrap
│   ├── app
│     ├──  assets
│     ├──  conf
│     ├──  db
│       ├── migrations 
│     ├──  events
│     ├──  handlers
│     ├──  types
│     ├──  views
│       ├── components
│       ├── errors
│       ├── landing
│       ├── layouts
│   ├── cmd
│     ├── app
│     ├── scripts
│       ├── seed
│   ├── plugins
│     ├── auth
│   ├── public
│     ├── assets
│   ├── env.local
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   ├── package-lock.json
│   ├── package.json
│   ├── tailwind.config.js
├── db
├── event 
├── kit
│   ├── middleware
├── validate
├── view
├── go.mod
├── install.go
├── README.md
```

### assets

Assets are stored in the `assets` directory. This directory contains all your CSS and JavaScript files. The `assets` directory is structured as follows:

```bash
assets
├── css
│   ├── app.css
├── js
│   ├── app.js
```

### conf

Configuration. First, config.yml is read, then environment variables overwrite the yaml config if they match. The config structure is in the config.go. The env-required: true tag obliges you to specify a value (either in yaml, or in environment variables).

Reading the config from yaml contradicts the ideology of 12 factors, but in practice, it is more convenient than reading the entire config from ENV. It is assumed that default values are in yaml, and security-sensitive variables are defined in ENV.

### db

The `db` directory contains all your database related files. The `db` directory is structured as follows:

```bash
db
├── migrations
│   ├── 20210919123456_create_users_table.sql
├── seeds
│   ├── seed.go
```

### events

The `events` directory contains all your event related files. These files are used to define custom events and event handlers for the project. The `events` directory is structured as follows:

```bash
events
├── event.go
```

### handlers

The `handlers` directory contains the main handlers or controllers for the project. These handlers handle incoming requests, perform necessary actions, and return appropriate responses. They encapsulate the business logic and interact with other components of the project, such as services and data repositories.

It is important to note that the project structure described here may not include all the directories and files present in the actual project. The provided overview focuses on the key directories relevant to understanding the structure and organization of the project.

### types

The `types` directory contains all your type related files. For example, you can define your models, structs, and interfaces in this directory. The `types` directory is structured as follows:

```bash
types
├── user.go
├── auth.go
```

### views

The `views` directory contains all your view related files. These files are used to render HTML templates for the project. The `views` directory is structured as follows:

```bash
views
├── home.go
├── about.go
```

## Development server

You can run the development server with the following command:

```bash
make dev 
```

## Hot reloading the browser

Hot reloading is configured by default when running your application in development.

> NOTE: on windows or on in my case (WSL2) you might need to run `make watch-assets` in another terminal to watch for CSS and JS file changes.

## Migrations

### Create a new migration

```bash
make db-mig-create add_users_table
```

The command will create a new migration SQL file located at `app/db/migrations/add_users_table.sql`

### Migrate the database

```bash
make db-up
```

### Reset the database

```bash
make db-reset
```

## Seeds

```bash
make db-seed
```

This command will run the seed file located at `cmd/scripts/seed/main.go`

## Creating views with Templ

superkit uses Templ as its templating engine. Templ allows you to create type safe view components that renders fragments of HTML. In-depth information about Templ can be found here:
[Templ documentation](https://templ.guide)

## Validations

todo

## Testing

### Testing handlers

## Create a production release

superkit will compile your whole application including its assets into a single binary. To build your application for production you can run the following command:

```bash
make build
```

This will create a binary file located at  `/bin/app_prod`.

Make sure you also set the correct application environment variable in your `.env` file.

```bash
SUPERKIT_ENV = production
```
