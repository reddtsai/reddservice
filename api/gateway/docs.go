// Package gateway Code generated by swaggo/swag. DO NOT EDIT
package gateway

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
        "/v1/signup": {
            "post": {
                "description": "註冊用戶",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "註冊",
                "parameters": [
                    {
                        "description": "raw",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "$ref": "#/definitions/main.SignUpResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "401": {
                        "description": "unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "403": {
                        "description": "forbidden",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "409": {
                        "description": "conflict",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "string"
                },
                "result": {}
            }
        },
        "main.SignUpRequest": {
            "type": "object",
            "required": [
                "account",
                "email"
            ],
            "properties": {
                "account": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "main.SignUpResponse": {
            "type": "object"
        }
    },
    "tags": [
        {
            "name": "auth"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Gateway",
	Description:      "This is a Gateway API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
