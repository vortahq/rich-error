# Rich Error

Rich Error is an error library that:

- allows nested errors
- makes it easy to store extra information along errors
- provides an easy to use API
- follows go1.13 errors conventions

## Installation

As we're publishing this module in our private repo, you'll need to do following steps before running `go get`:

1. Add repo address to GOPRIVATE env, to prevent go from using the proxy server:

    ```bash
    go env -w GOPRIVATE=github.com/vortahq/rich-error
    ```

    This is a shortcut for `GONOPROXY` and `GONOSUMDB`.

2. Add your Gitlab credentials to allow git to access to your gitlab account using http API

    After creating a new `Personal Access Token` with `write_repository` permission from "Profile" -> "Access Tokens",
    create a `~/.netrc` file with the following content

    ```
    machine gitlab.com
    login oauth2
    password <your access token>

    ```

    `~/.netrc` is the file that `curl` uses for authentication purposes. As `git` uses `curl` internally this will allow
    git to access to your repository by git-over-http

## Public API

For more information about public API checkout [contract.go](./contract.go)

## Usage

Using rich error is easy, you can create a new RichError using `richerror.New("error message")`. It automatically adds
runtime information to your error (like line number, file name, etc.). If you wish to add extra information to your error
you can use following methods. You can chain them together and except `NilIfNoError` their ordering is not
important.

### WithFields & WithField

These methods add a metadata (or a number of metadata) to your error. You can access them using `RichError.Metadata()`
method. These metadata are a good place to store information like request ID, user ID, etc. and allow you to debug errors
more efficiently.

### WithKind & WithLevel

These methods allow you to assign a kind and level to your errors. Kind and Level are predefined enums and you have to
choose from provided options. This allows you to hint to caller functions that the error is recoverable or not, or what
kind of issue caused the error.

### WithOperation

Operation is a hint that you can store in error to make debugging and grouping of errors easier.

### WithType

WithType lets you assign a type to your error. The type stores information that you're going to show to the user.
Because each service has it's own error types, You have to define your type struct yourself, it just have to implements
`Type` interface. You'll probably end up with something like:

```go
type MyErrorType string

func (t *MyErrorType) string {
    return string(t)
}

const (
    ErrX = MyErrorType("x")
)
```

### WithError

WithError allows you to wrap another error inside your error. It follows go 1.13 conventions and supports Unwrap, Is,
and As methods.

WithError tries to fill type, level, type, operation, etc. if they haven't been filled explicitly.

### NilIfNoError

NilIfNoError returns `nil` if the underling error is not present. It helps you avoid `if err != nil` check as much as possible.

### JsonMode

JsonMode is a flag controlling the format of generated output, the default format is string. You can change the output format by setting the `JsonMode` to `true`.

## Helpers

This package provides a set of helper function and structs to help users to utilize the full power of the RichError.
Currently, there are following constructs:
- **gRPC interception** which uses RichError's Kind to determine gRPC's status code.
- **Logger** which tries to log RichErrors in the most complete way (based on the logger given to it).
- **Sentry** which reports errors to sentry using `sentry-go` and uses RichErrors metadata to enrich the reported errors.
