// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/compress_image": {
            "post": {
                "description": "Compress image with specified compression level (1-5)",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Compress image",
                "operationId": "compress_image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "compression level (1-5)",
                        "name": "compression_level",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/convert_png_to_jpeg": {
            "post": {
                "description": "Convert image from PNG format to JPEG format",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Convert PNG to JPEG",
                "operationId": "convert_png_to_jpeg",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/resize_image": {
            "post": {
                "description": "Resize image to the specified width and height",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Resize image",
                "operationId": "resize_image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "width",
                        "name": "width",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "height",
                        "name": "height",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Image Converter API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
