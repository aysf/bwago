# bwago
building modern web application with Go

# How to use

clone this repo, go to project directory and run

```
go run cmd/web/*.go
```

# Run test

opsi 1

```
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
```

opsi 2

```
go test -cover
```

opsi 2

```
go test -v
```

# 3rd party library

- [Chi Route](https://github.com/go-chi/chi)
- [Justinas nosurf](https://github.com/justinas/nosurf)
- Session package [Alexedwards scs](https://github.com/alexedwards/scs)
- Popup Message [sweetalert2](https://github.com/sweetalert2/sweetalert2) 
- Notification [Notie](https://jaredreich.com/notie/)
- [vanilajs-datepicker](https://mymth.github.io/vanillajs-datepicker/#/)

# Motivation

this is a personal project and it is used for learning and educational purpose
the content of this repo is adapted from https://www.udemy.com/course/building-modern-web-applications-with-go/
