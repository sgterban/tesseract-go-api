# Tesseract-OCR API

# Install

## Ubunutu

```
apt-get install tesseract-ocr
apt-get install libtesseract-dev
apt-get install libleptonica-dev
apt-get install libpng-dev
apt-get install libjpeg-dev
apt-get install libmagickwand-dev
```

## Common

```
go get github.com/otiai10/gosseract
go get gopkg.in/gographics/imagick.v2/imagick
go get github.com/sgterban/tesseract-go-api
```

# Examples

```
go run hello/hello.go "PATH_TO_IMAGE_ON_COMPUTER"
```

```
go run optimize/optimize.go "PATH_TO_IMAGE_ON_COMPUTER"
```

```
go run server/server.go
localhost:8080/
```