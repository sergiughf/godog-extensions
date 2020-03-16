# Godog Extensions

This library contains the following set of Godog extensions:

- Postgres Cleanup:
  - Truncates all the tables in the public schemaname before each scenarios
  - Must have an already running postgres db
  
- Gomega Fail Handler:
  - A matcher/assertion library to use with godog
  
- Wiremock:
  - An extension to be able to mock requests and responses in godog tests
  - Must have an already running WireMock server

## Installation

```BASH
go get -v github.com/sergiughf/godog-extensions
```

## Usage

```go
package main

import (
    "github.com/cucumber/godog"
    "github.com/sergiughf/godog-extensions"
)

func FeatureContext(s *godog.Suite) {
	extension.NewGomegaFailHandler(s)
	extension.NewPostgresCleanup(s, postgresDSN)
	extension.NewWireMock(s, wireMockServerURL)
}
```

Inside the step definitions to set up a wiremock response/request:

```go

wmClient := extension.WireMockClient()

wmClient.Request = extension.WireMockRequest{
    Method: http.MethodGet,
    URL:    "/v1/recipes?country=" + country,
}

wmClient.Response = extension.WireMockResponse{
    Status:  http.StatusOK,
    Headers: map[string]string{"Content-Type": "application/json"},
    Body:    string(body),
}

wmClient.SendMocks()
```
