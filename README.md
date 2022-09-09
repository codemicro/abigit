# go-fiber-sql

A project template for a web application based on the Fiber web framework. This incorporates:

* Error handling with `github.com/rs/zerlog` and `github.com/pkg/errors`
* Database migrations (requires Go v1.18 or newer)
* Tools to load YAML-based configurations

## To use this template:

* Add a `LICENSE`
* Update the package name in `go.mod` and in the imports within the application
* Rename the `application` directory to your application's name
* Define your database DSN and driver
* Write some code and make an app, I guess?