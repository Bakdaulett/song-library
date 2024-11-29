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
        "/": {
            "get": {
                "summary": "Get all songs",
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "List of all songs",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Songs"]
            },
            "post": {
                "summary": "Add a new song",
                "consumes": ["application/json"],
                "parameters": [
                    {
                        "name": "song",
                        "in": "body",
                        "description": "Song details",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongRequest"
                        }
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Song added",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Songs"]
            }
        },
        "/{id}": {
            "get": {
                "summary": "Get song by ID",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Song ID"
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Song found",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Songs"]
            },
            "put": {
                "summary": "Update song by ID",
                "consumes": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Song ID"
                    },
                    {
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "description": "Song details",
                        "schema": {
                            "$ref": "#/definitions/models.SongRequest"
                        }
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Song updated",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Songs"]
            },
            "delete": {
                "summary": "Delete song by ID",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Song ID"
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Song deleted",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Songs"]
            }
        },
        "/{id}/lyrics": {
            "get": {
                "summary": "Get all lyrics of a song",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Song ID"
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "All lyrics of the song",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Lyrics"]
            }
        },
        "/{id}/lyrics/{verseId}": {
            "get": {
                "summary": "Get specific verse lyrics by verse ID",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Song ID"
                    },
                    {
                        "name": "verseId",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Verse ID"
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Lyrics of the specific verse",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Lyrics"]
            }
        },
        "/{id}/lyrics/{verseStartId}-{verseEndId}": {
            "get": {
                "summary": "Get a range of lyrics from verseStartId to verseEndId",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Song ID"
                    },
                    {
                        "name": "verseStartId",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "Start verse ID"
                    },
                    {
                        "name": "verseEndId",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "End verse ID"
                    }
                ],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Range of lyrics from start to end verse",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "default": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                },
                "tags": ["Lyrics"]
            }
        }
    },
    "definitions": {
        "models.ErrorDetail": {
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            },
            "type": "object"
        },
        "models.ErrorResponse": {
            "properties": {
                "error": {
                    "$ref": "#/definitions/models.ErrorDetail"
                }
            },
            "type": "object"
        },
        "models.SongRequest": {
            "properties": {
                "group": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string",
                    "format": "date-time"
                },
                "lyrics": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                }
            },
            "type": "object"
        },
        "models.SuccessResponse": {
            "properties": {
                "result": {}
            },
            "type": "object"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/songs",
	Schemes:          []string{"http", "https"},
	Title:            "Song Library API",
	Description:      "This is the Song Library API that allows you to manage songs and their lyrics.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
