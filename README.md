# GOTHKIT

> NOT READY (YET)

# Installation
```
go install github.com/anthdm/gothkit@master
```

After installation you can create a new project by running: 
```
gothkit [myprojectname]
```

Navigate to your project and start the development server:
```
cd [myprojectname]

make dev
```

# Getting started
## Migrations
### Create a new migration
```
make db-mig add_user_table
```

### Migrate the database 
```
make db-up
```

### Reset the database 
```
make db-reset
```

### Seed the database 
```
make db-seed
```





