// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/persons": {
            "get": {
                "description": "Получение списка людей при помощи фильтров и пагинацией. Примечание name и surname работает через поиск подстроки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "persons"
                ],
                "summary": "List persons",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Фильтр по имени",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по фамилии",
                        "name": "surname",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Фильтр по возрасту",
                        "name": "age",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по полу",
                        "name": "gender",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по национальности",
                        "name": "nationality",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Размер страницы",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ListPersonsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.BadRequest"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.NotFound"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.InternalServerError"
                        }
                    }
                }
            },
            "post": {
                "description": "Создаёт нового человека, обогащает информацию с помощью внешнего API и сохраняет в PostgreSQL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "persons"
                ],
                "summary": "Add person",
                "parameters": [
                    {
                        "description": "Параметры для создания",
                        "name": "person",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreatePersonRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Сущность человека",
                        "schema": {
                            "$ref": "#/definitions/dto.PersonResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/dto.BadRequest"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/dto.InternalServerError"
                        }
                    },
                    "503": {
                        "description": "Сервис недоступен",
                        "schema": {
                            "$ref": "#/definitions/dto.ServiceUnavailable"
                        }
                    }
                }
            }
        },
        "/persons/{id}": {
            "get": {
                "description": "Возвращает информацию о человеке по его уникальному идентификатору",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "persons"
                ],
                "summary": "Get person by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор человека",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Найдённая сущность Person",
                        "schema": {
                            "$ref": "#/definitions/dto.PersonResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат ID",
                        "schema": {
                            "$ref": "#/definitions/dto.BadRequest"
                        }
                    },
                    "404": {
                        "description": "Person не найден",
                        "schema": {
                            "$ref": "#/definitions/dto.NotFound"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/dto.InternalServerError"
                        }
                    }
                }
            },
            "put": {
                "description": "Обновляет информацию о человеке (любые поля: name, surname, patronymic, age, gender, nationality)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "persons"
                ],
                "summary": "Update person",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор человека",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Поля для обновления",
                        "name": "person",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdatePersonRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Обновлённая сущность человека",
                        "schema": {
                            "$ref": "#/definitions/dto.PersonResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос или параметры",
                        "schema": {
                            "$ref": "#/definitions/dto.BadRequest"
                        }
                    },
                    "404": {
                        "description": "Человек не найден",
                        "schema": {
                            "$ref": "#/definitions/dto.NotFound"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/dto.InternalServerError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет человека по идентификатору",
                "tags": [
                    "persons"
                ],
                "summary": "Delete person",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.BadRequest"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.NotFound"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.InternalServerError"
                        }
                    }
                }
            },
            "patch": {
                "description": "Обновляет информацию о человеке (любые поля: name, surname, patronymic, age, gender, nationality)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "persons"
                ],
                "summary": "Update person",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор человека",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Поля для обновления",
                        "name": "person",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdatePersonRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Обновлённая сущность человека",
                        "schema": {
                            "$ref": "#/definitions/dto.PersonResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос или параметры",
                        "schema": {
                            "$ref": "#/definitions/dto.BadRequest"
                        }
                    },
                    "404": {
                        "description": "Человек не найден",
                        "schema": {
                            "$ref": "#/definitions/dto.NotFound"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/dto.InternalServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.BadRequest": {
            "description": "Bad request",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "invalid request data"
                }
            }
        },
        "dto.CreatePersonRequest": {
            "type": "object",
            "required": [
                "name",
                "surname"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                },
                "patronymic": {
                    "type": "string",
                    "maxLength": 50
                },
                "surname": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                }
            }
        },
        "dto.InternalServerError": {
            "description": "Internal server error",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "internal error"
                }
            }
        },
        "dto.ListPersonsResponse": {
            "description": "Ответ списка Person",
            "type": "object",
            "properties": {
                "pagination": {
                    "$ref": "#/definitions/dto.PaginationInfo"
                },
                "persons": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.PersonResponse"
                    }
                }
            }
        },
        "dto.NotFound": {
            "description": "Resource not found",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "person not found"
                }
            }
        },
        "dto.PaginationInfo": {
            "description": "Метаданные пагинации",
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "dto.PersonResponse": {
            "description": "основная информация о пользователе в json",
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "dto.ServiceUnavailable": {
            "description": "Service unavailable",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "service unavailable"
                }
            }
        },
        "dto.UpdatePersonRequest": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer",
                    "maximum": 110,
                    "minimum": 0
                },
                "gender": {
                    "type": "string",
                    "enum": [
                        "male",
                        "female"
                    ]
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                },
                "nationality": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string",
                    "maxLength": 50
                },
                "surname": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Persons API",
	Description:      "Сервис выполненный по тз https://drive.google.com/file/d/1zUU44O1ye5-3yYRdhLMEhpG9JXz-TWeQ/view",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
