# Rental Property API

A read-only REST API for browsing rental property listings, built with **Go** and the **Beego** framework. Property data is loaded once from a JSON file into memory at startup — there is no database and no external API calls. The API exposes two endpoints: a filterable property list and a single-property lookup by ID.

---

## Features

- List properties with **9 combinable AND filters** (price range, star rating, review score, review count, published status, property type, feed, bedroom count)
- **OR-based amenities filter** — match properties with any one of several requested amenities
- Optional **result limiting** via a `limit` query parameter
- Single property lookup by ID
- Strict **query parameter validation** with `400 Bad Request` on invalid input
- Consistent JSON error shape across all endpoints
- In-memory data store — loaded once at startup, no repeated file reads
- Interactive **Swagger UI** for exploring and testing the API
- Table-driven unit tests covering transformation, filtering, and lookup logic

---

## Technology Stack

| Component | Technology |
|---|---|
| Language | Go 1.22+ |
| Framework | [Beego v2](https://github.com/beego/beego) (API mode) |
| Data source | Local JSON file, loaded into memory |
| JSON handling | `encoding/json` (standard library) |
| API documentation | Swagger 2.0 (via `bee generate docs`) |
| Testing | Go's built-in `testing` package (table-driven tests) |

---

## Project Structure

```
Beego-Api-Project/
├── conf
│   └── app.conf                
├── controllers
│   ├── property.go             
│   └── base.go                
├── data
│   └── rental_properties.json  
├── models
│   ├── source.go                
│   ├── response.go              
│   └── filter.go                
├── routers
│   ├── router.go                
│   └── commentsRouter.go        
├── services
│   ├── property.go               
│   └── property_test.go          
├── swagger
│   └── v1  # Generated Swagger UI + swagger.json/swagger.yml and other files
├── main.go                       
├── go.mod 
|──go.sum               
└── README.md
```

---

## Architecture

The project follows a strict layered flow:

```
Router → Controller → Service → In-Memory Store
```

- **Router** (`routers/`) — maps HTTP method + path to controller methods.
- **Controller** (`controllers/`) — parses and validates HTTP input (query params, path params), and writes HTTP responses. Contains no business logic.
- **Service** (`services/`) — owns all transformation, filtering, and lookup logic. Operates on plain Go values, independent of HTTP.
- **In-memory store** — a package-level slice populated once at startup by `LoadData()`.

This separation keeps controllers thin and testable service functions free of HTTP concerns — the service layer can be unit tested directly without spinning up a server.

---

## Prerequisites

- [Go](https://go.dev/dl/) 1.22 or higher
- [Bee CLI](https://github.com/beego/bee) (for running, doc generation, and route generation):
  ```bash
  go install github.com/beego/bee/v2@latest
  ```

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/SakibShehan/Beego-Api-Project
   cd Beego-Api-Project
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Confirm the data file is present at:
   ```
   data/rental_properties.json
   ```

---

## Running the Application

Start the server with Bee's live-run tool:

```bash
bee run
```

Or build and run manually:

```bash
go build -o rental-property-api .
./rental-property-api
```

The server starts on the port configured in `conf/app.conf` (default `8080`). Property data is loaded into memory once at startup — check the logs for confirmation:

```
loaded 100 properties from data/rental_properties.json
```

---

## Swagger API Documentation

Interactive API documentation is available via Swagger UI once the server is running:

```
http://localhost:8080/swagger/
```

From there you can browse both endpoints, inspect all supported query parameters, and use **Try it out** to fire real requests against the running server.

To regenerate the documentation after changing controller annotations:

```bash
bee generate docs
```

---

## API Base URL

```
http://localhost:8080/v1
```

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/v1/properties` | List properties, with optional filters and `limit` |
| `GET` | `/v1/properties/:id` | Get a single property by its ID |

### List properties

```bash
curl "http://localhost:8080/v1/properties"
```

### Get a property by ID

```bash
curl "http://localhost:8080/v1/properties/BC-1000001"
```

---

## Property Filtering

All filters apply only to `GET /v1/properties`. Filtering happens in this order: **filter → transform → apply limit → build response**.

### AND filters (all supplied filters must match)

| Query param | Type | Rule |
|---|---|---|
| `min_price` | float | `usd_price >= min_price` |
| `max_price` | float | `usd_price <= max_price` |
| `min_star_rating` | int | `star_rating >= min_star_rating` |
| `min_review_score` | float | `review_score_general >= min_review_score` |
| `min_reviews` | int | `number_of_review >= min_reviews` |
| `published` | bool | `published` equals `true` or `false` |
| `property_type` | string | Exact match (case-sensitive): `Hotel`, `House`, `Apartment`, `Villa`, `Resort`, `Hostel` |
| `feed` | int | Exact match: `11`, `12`, `22`, or `24` |
| `min_bedroom` | int | `bedroom_count >= min_bedroom` |

```bash
curl "http://localhost:8080/v1/properties?min_price=50&max_price=150&property_type=Hotel"
```

### OR filter — amenities

```
amenities=Internet,Parking
```

Matches a property if it has **at least one** of the listed amenities. Matching is case-sensitive. Combined with other filters, the amenities OR-check must still pass alongside every AND filter.

```bash
curl "http://localhost:8080/v1/properties?feed=11&amenities=Internet,Parking"
```

### Limiting results

```bash
curl "http://localhost:8080/v1/properties?limit=5"
```

### Invalid parameters

Any invalid parameter value returns `400 Bad Request` with a descriptive error:

```bash
curl "http://localhost:8080/v1/properties?min_price=abc"
```

---

## Testing

Unit tests live in `services/property_test.go` and follow a table-driven pattern. They cover:

- Property transformation (source → API response mapping)
- AND filter combinations
- Amenities OR filter
- Combined AND + OR filtering
- Empty result handling (never returns `null`)
- Lookup by ID (found and not found)

Run all tests:

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test ./services/... -cover
```

### Generate a Coverage Profile

```bash
go test ./services/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### View Coverage in the Browser

```bash
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in any browser to see line-by-line coverage.

---

## Static Code Analysis

Run Go's built-in vet tool to catch suspicious constructs before committing:

```bash
go vet ./...
```

---

## API Response Format

### List response

```json
{
  "Result": {
    "Count": 2,
    "Items": [
      {
        "ID": "BC-1000001",
        "Feed": 11,
        "GeoInfo": { "...": "..." },
        "Property": { "...": "..." },
        "Published": true
      }
    ]
  }
}
```

`Items` is always an array — never `null` — even when no properties match.

### Single property response

```json
{
  "ID": "BC-1000001",
  "Feed": 11,
  "GeoInfo": { "...": "..." },
  "Property": { "...": "..." },
  "Published": true
}
```

---

## Error Handling

All errors share a single consistent shape:

```json
{
  "Error": "Human-readable message"
}
```

| Status | When |
|---|---|
| `400 Bad Request` | An invalid or malformed query parameter was supplied |
| `404 Not Found` | The requested property ID does not exist |

---

## Development Commands

| Command | Purpose |
|---|---|
| `bee run` | Run the app with live reload |
| `go build ./...` | Compile the project |
| `go test ./... -v` | Run all tests verbosely |
| `go vet ./...` | Run static analysis |
| `go mod tidy` | Clean up module dependencies |
| `bee generate routers` | Regenerate route registrations from `@router` annotations |
| `bee generate docs` | Regenerate Swagger documentation |

---

## Author

**Sakib Hossen Shehan**
 W3 Engineers Intern Assignment — Rental Property API