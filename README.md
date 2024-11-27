# HTTP Server in Go

This repository contains a simple http server written in Go. The server provides a RESTful API for retrieving and adding quotes to local json file.

## Features

* Retrieve a random quote
* Retrieve a quote by author
* Add a new quote

## API Endpoints

* `GET /quotes`: Retrieves a random quote
* `GET /quotes?author=<author>`: Retrieves a quote by author
* `POST /quotes`: Adds a new quote

## Quote Data

The server stores quotes in a JSON file named `quotes.json`. The file contains an array of quote objects with the following structure:

```json
{
  "author": "Author Name",
  "quote": "Quote text"
}
```

## Running the Server

To run the server, navigate to the repository directory and execute the following command:

```
go run main.go
```

The server will start listening on ports 2222 and 4444. You can use a tool like `curl` to test the API endpoints.

## Adding New Quotes

To add a new quote, send a POST request to the `/quotes` endpoint with a JSON payload containing the author and quote text:

```json
{
  "author": "New Author",
  "quote": "New quote text"
}
```

Note: This repository is based on the provided code snippets and may not be a complete or fully functional application.