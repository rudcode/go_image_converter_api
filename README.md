# GO IMAGE CONVERTER API

## Instructions
```
git clone https://github.com/rudcode/go_image_converter_api
cd go_image_converter_api
go run ./cmd/api
```

The server should start at http://localhost:8000

Swagger available at http://localhost:8000/docs/index.html


## Developers
Test available with following command:
```
go test -v ./...
```

To regenerate swagger docs after making changes run this command:
```
swag init -g .\cmd\api\main.go
```

If command above doesn't work install swagger first:
```
go install github.com/swaggo/swag/cmd/swag@latest
```