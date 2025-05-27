# go-core

## Requirements
- [Nix](https://nixos.org/)
- [Direnv](https://direnv.net/)

### Runtime
- [mozjpeg](https://github.com/mozilla/mozjpeg)
- [pngquant](https://github.com/kornelski/pngquant)
- [gifsicle](https://github.com/kohler/gifsicle)

## Installation

```sh
$ git clone https://github.com/eiicon-company/go-core
$ cd go-core
$ direnv allow
```

## Usage

Make sure what tasks are existing.

```php
$ make help
```

Update Golang dependencies. This command would be help you when something went wrong.

```php
$ make gomodule
```

## DATA I/O

- storage: s3, filesystam, gcs, etc..
