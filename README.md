# SUPERKIT

Build high-performance apps swiftly with minimal team resources in Go.

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
- [Creating Views with Templ](#creating-views-with-templ)
- [Validations](#validations)
- [Testing](#testing)
    - [Testing Handlers](#testing-handlers)
- [Create a production release](#create-a-production-release)

# Installation

Create your SUPERKIT project in a single command:

```shell
go run github.com/anthdm/superkit@master [yourprojectname]
```

Navigate to your project directory:

```shell
cd [myprojectname]
```

Install both Tailwind CSS and esbuild locally:

```shell
npm install
```

If you encounter dependency issues, you can run:

```shell
go clean -modcache && go get -u ./...
```

If you have the authentication plugin enabled, you must migrate your database:

```shell
make db-up
```

# Getting Started

## Application Structure

- **assets**: Contains CSS and JavaScript files.
- **conf**: Holds application-specific configuration files.
- **db**: Manages database migrations using Goose.
- **events**: Contains event system handlers, similar to HTTP handlers but for events.
- **handlers**: Manages HTTP request handlers.
- **types**: Defines application-specific types.
- **views**: Includes layouts and other view components.

## Development server

To start the development server, use the following command:

```shell
make dev 
```

## Hot reloading the browser

Hot reloading is enabled by default when running your application in development mode.

> **Note:** On Windows, or when using WSL2, you may need to run `make watch-assets` in
> a separate terminal to monitor changes to CSS and JS files.

# Migrations

## Create a new migration

```shell
make db-mig-create add_users_table
```

The command will generate a new SQL migration file located at `app/db/migrations/add_users_table.sql`

## Migrate the database

```shell
make db-up
```

## Reset the database

```shell
make db-reset
```

## Seeds

```shell
make db-seed
```

This command will execute the seed file located at `cmd/scripts/seed/main.go`

# Creating Views with Templ

Superkit utilizes Templ as its templating engine, enabling the creation of type-safe view components that render
fragments of HTML. Templ enhances your development process by ensuring that your views are both maintainable and robust,
reducing the likelihood of runtime errors.

### Getting Started with Templ

To create a new view component with Templ, follow these steps:

1. **Define the Component**:
   Create a Go struct that represents your view component. This struct can include any data you need to render in your
   view.

    ```go
    package views
    
    type UserProfile struct {
       Name  string
       Email string
    }
    ```

2. **Create the Template File**:
   Write an HTML template file that utilizes the Templ syntax. Place this file in your views directory.

    ```html
    <!-- views/user_profile.templ -->
    <div>
        <h1>{{.Name}}</h1>
        <p>{{.Email}}</p>
    </div>
    ```

3. **Render the View**:
   Use Templ’s rendering functions to render your view component within your application logic.

    ```go
    package main
    
    import (
        "github.com/templ/templ"
        "github.com/yourusername/yourproject/views"
        "net/http"
    )
    
    func userProfileHandler(w http.ResponseWriter, r *http.Request) {
        profile := views.UserProfile{
            Name:  "John Doe",
            Email: "john.doe@example.com",
        }
        templ.Render(w, "user_profile.templ", profile)
    }
    
    func main() {
        http.HandleFunc("/profile", userProfileHandler)
        http.ListenAndServe(":8080", nil)
    }
    ```

### Benefits of Using Templ

- **Type Safety**: Templ ensures that your view components are type-safe, reducing runtime errors.
- **Maintainability**: By separating your view logic into components, your codebase becomes easier to maintain and
  extend.
- **Performance**: Templ is optimized for performance, ensuring that your views render quickly and efficiently.

For more detailed information about Templ, please refer to the official [Templ documentation](https://templ.com/docs).

# Validations

The validation system in this project is designed to ensure that data meets specific criteria before being processed. It
consists of several components that work together to provide a comprehensive and flexible validation framework.

### Components

1. **Errors**
    - The `Errors` type is a map that holds potential validation errors.
    - Methods include:
        - `Any() bool`: Checks if there are any errors.
        - `Add(field, msg string)`: Adds an error message for a specific field.
        - `Get(field string) []string`: Retrieves all error messages for a specific field.
        - `Has(field string) bool`: Checks if a specific field has any errors.

2. **Schema**
    - The `Schema` type represents a validation schema, which is a map of field names to rule sets.
    - Methods include:
        - `Merge(schema, other Schema) Schema`: Merges two schemas into one.
        - `Rules(rules ...RuleSet) []RuleSet`: Creates a list of rule sets.

3. **Validation Functions**
    - `Validate(data any, fields Schema) (Errors, bool)`: Validates data based on the provided schema.
    - `Request(r *http.Request, data any, schema Schema) (Errors, bool)`: Parses and validates data from an HTTP request
      based on the provided schema.

### Rule Sets

Rules define the conditions that fields must meet. Several predefined rules include:

`Required` &raquo; Ensures a field is not empty.  
`Email` &raquo; Validates that a field contains a valid email address.  
`URL` &raquo; Validates that a field contains a valid URL.  
`ContainsUpper` &raquo; Checks if a string contains at least one uppercase character.  
`ContainsDigit` &raquo; Checks if a string contains at least one numeric character.  
`ContainsSpecial` &raquo; Checks if a string contains at least one special character.  
`Time` &raquo; Validates that a field contains a valid time.  
`TimeAfter(time.Time)` &raquo; Ensures a time is after a specified value.  
`TimeBefore(time.Time)` &raquo; Ensures a time is before a specified value.  
`In([]T)` &raquo; Checks if a value is within a specified set of values.  
`EQ(T)` &raquo; Checks if a field is equal to a specified value.  
`LTE(T)` &raquo; Ensures a numeric field is less than or equal to a specified value.  
`GTE(T)` &raquo; Ensures a numeric field is greater than or equal to a specified value.  
`LT(T)` &raquo; Ensures a numeric field is less than a specified value.  
`GT(T)` &raquo; Ensures a numeric field is greater than a specified value.  
`Max(int)` &raquo; Ensures a string's length is at most a specified number of characters.  
`Min(int)` &raquo; Ensures a string's length is at least a specified number of characters.

### Example Schema

Here's an example of a validation schema:

```go
var testSchema Schema = Schema{
    "createdAt": Rules(Time),
    "startedAt": Rules(TimeBefore(time.Now())),
    "deletedAt": Rules(TimeAfter(createdAt)),
    "email":     Rules(Email),
    "url":       Rules(URL),
    "password":  Rules(
        ContainsSpecial,
        ContainsUpper,
        ContainsDigit,
        Min(7),
        Max(50),
    ),
    "age":      Rules(GTE(18)),
    "bet":      Rules(GT(0), LTE(10)),
    "username": Rules(Required),
}
```

# Testing

## Testing Handlers

todo

# Create a production release

Superkit will compile your entire application, including its assets, into a single binary.
To build your application for production, use the following command:

```shell
make build
```

This command will create a binary file located at  `/bin/app_prod`.

Make sure that you set the appropriate application environment variables in your `.env` file.

```shell
SUPERKIT_ENV = production
```

