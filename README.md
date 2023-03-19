# bwago
building modern web application with Go

## Prerequisite

1. Go 1.18 or above
2. [Soda](https://gobuffalo.io/documentation/database/soda/)


## How to use

1. clone this repo, go to project directory 
2. migrate db `soda migrate up`
3. run script `./run.sh`
4. open browser `locahost:9000`
5. login admin email `admin@here.com`, password `password`

# Run test

Test all packages

```
go test -coverpkg=./... ./...
```

Alternative command for testing

```
go test -v -cover  ./...
```

Test all packages with total percentage

```
go test --coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

Test all packages with total percentage and display in the browser

```
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

# 3rd party library

- [Chi Route](https://github.com/go-chi/chi)
- [Justinas nosurf](https://github.com/justinas/nosurf)
- Session package [Alexedwards scs](https://github.com/alexedwards/scs)
- Popup Message [sweetalert2](https://github.com/sweetalert2/sweetalert2) 
- Notification [Notie](https://jaredreich.com/notie/)
- [vanilajs-datepicker](https://mymth.github.io/vanillajs-datepicker/#/)
- [go-simple-mail](https://github.com/xhit/go-simple-mail)

# Motivation

this is a personal project and it is used for learning and educational purpose
the content of this repo is adapted from https://www.udemy.com/course/building-modern-web-applications-with-go/
