# metrix

This project is an exercise to learn HTMX+Go. The goal was to build a modern web app without having to write
a single line of Javascript.

The project uses:

- [PicoCSS](https://picocss.com/docs) for basic styling
- [HTMX](https://htmx.org/) as the presentation layer
- [Go](https://go.dev/) as the backend technology
- [templ](https://github.com/a-h/templ) as the HTML template engine
- [sqlite](https://sqlite.org/docs.html) for the persistence layer
- [Playwright](https://github.com/playwright-community/playwright-go) for integration testing (Go version of the framework)

To run locally:

```sh
go run main.go db migrate
go run main.go serve
```

To run the containerised version:

```sh
make run-prod
```
