# go-core

## Requirements

- [libvips@8.3+](https://libvips.github.io/libvips/install.html)
  - Installation: https://github.com/h2non/bimg#libvips
- [mozjpeg](https://github.com/mozilla/mozjpeg)
  - Installation:
    ```php
    $ cd /tmp && \
    wget https://github.com/mozilla/mozjpeg/archive/v3.3.1.tar.gz && \
    tar -zxvf v3.3.1.tar.gz && \
    cd mozjpeg-3.3.1 && \
    autoreconf -fiv && \
    mkdir build && cd build && \
    sh ../configure && \
    make install && \
    ln -s /opt/mozjpeg/bin/* /usr/local/bin/
    ```
- [pngquant](https://github.com/kornelski/pngquant)
  - Installation: `apt install pngquant`
- [gifsicle](https://github.com/kohler/gifsicle)
  - Installation: `apt install gifsicle`

## Installation

Download Project

```php
$ GO111MODULE=on go get -u github.com:eiicon-company/go-core
$ GO111MODULE=on go get -u github.com:eiicon-company/go-core/v1
$ GO111MODULE=on go get -v github.com/eiicon-company/go-core@master
$ GO111MODULE=on go get -u github.com/eiicon-company/go-core@develop
$ GO111MODULE=on go get -v github.com/eiicon-company/go-core@feature_branch
$ GO111MODULE=on go get -u github.com/eiicon-company/go-core@feature/branch
$ GO111MODULE=on go get -v github.com/eiicon-company/go-core@bbb2610aea46e47f09f71d6bbc5b75dc3f585b08
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


