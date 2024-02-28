// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ping": {
            "get": {
                "description": "Pings the database to check its connectivity",
                "summary": "Ping database",
                "responses": {
                    "200": {
                        "description": "Database ping successful",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/update/{metricType}/{metricName}/{metricValue}": {
            "post": {
                "description": "Inserts or updates the value of a metric specified by its type, name, and value.",
                "produces": [
                    "application/json"
                ],
                "summary": "Insert or update metric value",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of the metric ('gauge' or 'counter')",
                        "name": "metricType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the metric",
                        "name": "metricName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Value of the metric",
                        "name": "metricValue",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Metric value inserted or updated successfully",
                        "schema": {
                            "$ref": "#/definitions/storage.Metric"
                        }
                    },
                    "400": {
                        "description": "Bad request. Invalid metric parameters or JSON payload",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/updates": {
            "post": {
                "description": "Bulk inserts or updates metric values.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Bulk insert or update metrics",
                "parameters": [
                    {
                        "description": "Array of metrics to insert or update",
                        "name": "metrics",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/storage.Metric"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Metrics inserted or updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request. Invalid JSON payload or metric parameters",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/value/{metricType}/{metricName}": {
            "get": {
                "description": "Retrieves the value of a metric specified by its type and name.",
                "summary": "Retrieve metric value by type and name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of the metric ('gauge' or 'counter')",
                        "name": "metricType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the metric",
                        "name": "metricName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Metric value retrieved successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request. Either metric type is unsupported or value is missing",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Metric value not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Retrieves the value of a metric specified by its type and name.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Retrieve metric value by type and name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of the metric ('gauge' or 'counter')",
                        "name": "metricType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the metric",
                        "name": "metricName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Metric value retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/storage.Metric"
                        }
                    },
                    "400": {
                        "description": "Bad request. Invalid JSON payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Metric value not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "storage.Metric": {
            "type": "object",
            "properties": {
                "delta": {
                    "description": "Delta value (applicable for counter type)",
                    "type": "integer"
                },
                "id": {
                    "description": "Metric ID",
                    "type": "string"
                },
                "type": {
                    "description": "Metric type: gauge or counter",
                    "type": "string"
                },
                "value": {
                    "description": "Value (applicable for gauge type)",
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "petstore.swagger.io",
	BasePath:         "/v2",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
