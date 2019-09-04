# go-core


>
> A go@1.12.9
>

This system force required to: You **have to must** have `GO111MODULE` environment variable that activates `go mod` module.

Therefore, it set into.

```php
$ export GO111MODULE=on
```

or use direnv.

```php
$ brew install direnv
```

## Installation

Download Project

```php
$ go get -u github.com:eiicon-company/go-core
$ go get -u github.com:eiicon-company/go-core/v1
$ go get -u github.com:eiicon-company/go-core@rom5kdxv4kfq0uhfq1hfq4
```

## Usage

Make sure what tasks are existing.

```php
$ make help
```

Update Golang dependencies. This command would be help you when something went wrong happended.

```php
$ make gomodule
```

## DATA I/O

- storage: s3, filesystam, gcs, etc..
- rdb: redis


