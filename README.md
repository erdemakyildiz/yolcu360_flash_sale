# Flash Sale Management

## Overview

Flash Sale Management is an e-commerce system designed to handle flash sales for high-demand products. It supports real-time discount campaigns where users can take advantage of limited-time offers. The system ensures data consistency and prevents stock depletion under high load.

## Technologies

- **Go**: Programming language used for backend development.
- **Postgres**: Relational database for storing sale data.
- **Redis**: Caching system to manage and synchronize stock levels.
- **Docker**: Containerization for deployment.

## Features

- Create, read, update, and delete flash sales.
- Manage stock levels and ensure consistency during high demand.
- Track sales and manage active/inactive status of flash sales.
- Support for percentage-based discounts.
- Documented REST API with Swagger.

## Business Requirements

- Each flash sale is linked to a single product with a specified stock limit.
- Flash sales are only applicable as percentage discounts.
- Sales are active only during specified start and end times.
- All sales transactions are recorded in the system.
- Sales cannot proceed if the product stock is zero.
- Ensure data consistency and prevent stock overselling during concurrent requests.

## Technical Requirements

- Dockerize the service and provide a `docker-compose` setup.
- Document all endpoints using Swagger.
- Share the code repository on GitHub or Bitbucket.

## API Endpoints

### 1. Create Flash Sale

Create a new flash sale.

```bash
curl --location 'http://127.0.0.1:3000/flash-sales' \
--header 'Content-Type: application/json' \
--data '{
  "discount": 10,
  "endTime": "2024-09-26T11:04",
  "product_id": 1,
  "startTime": "2024-09-16T11:04",
  "saleStock": 5
}'
```

**Response:**

```json
{
  "product_id": 1,
  "saleStock": 5,
  "discount": 10,
  "startTime": "2024-09-16T11:04:00Z",
  "endTime": "2024-09-26T11:04:00Z",
  "active": false
}
```

### 2. Update Flash Sale

Update an existing flash sale.

```bash
curl --location --request PUT 'http://127.0.0.1:3000/flash-sales' \
--header 'Content-Type: application/json' \
--data '{
  "id" : 1,
  "discount": 20,
  "endTime": "2024-09-26T11:04",
  "product_id": 1,
  "startTime": "2024-09-16T11:04",
  "saleStock": 5,
  "active": true 
}'
```

**Response:**

```json
{
  "product_id": 1,
  "saleStock": 5,
  "discount": 20,
  "startTime": "2024-09-16T11:04:00Z",
  "endTime": "2024-09-26T11:04:00Z",
  "active": true
}
```

### 3. Get Flash Sale by ID

Retrieve details of a flash sale by its ID.

```bash
curl --location 'http://127.0.0.1:3000/flash-sales/1' \
--header 'accept: application/json'
```

**Response:**

```json
{
  "product_id": 1,
  "saleStock": 5,
  "discount": 20,
  "startTime": "2024-09-16T11:04:00Z",
  "endTime": "2024-09-26T11:04:00Z",
  "active": true
}
```

### 4. Get All Flash Sales

Retrieve a list of all active flash sales.

```bash
curl --location 'http://127.0.0.1:3000/flash-sales' \
--header 'accept: application/json'
```

**Response:**

```json
[
  {
    "product_id": 1,
    "saleStock": 5,
    "discount": 20,
    "startTime": "2024-09-16T11:04:00Z",
    "endTime": "2024-09-26T11:04:00Z",
    "active": true
  }
]
```


### 5. Delete Flash Sale

Delete an existing flash sale by ID.

```bash
curl --location --request DELETE 'http://127.0.0.1:3000/flash-sales/1' \
--header 'accept: application/json'
```

**Response:**

```json
200 OK
```


### 6. Sale Product

Purchase a product from an active flash sale.

```bash
curl --location --request POST 'http://127.0.0.1:3000/flash-sales/2/buy?wait=1' \
--header 'accept: application/json'
```

**Response:**

```json
{
  "ID": 1,
  "ProductID": 1,
  "RemainingSaleStock": 4,
  "RemainingProductStock": 9,
  "Price": 40,
  "CreatedAt": "2024-09-18T05:23:05.714762+03:00"
}
```

## Setup and Running

1. Clone the repository from GitHub or Bitbucket.
2. Ensure Docker and Docker Compose are installed.
3. Run `docker-compose up` to start the service and its dependencies.
4. Access the service at `http://127.0.0.1:3000`.
5. Access the Swagger documentation at `http://127.0.0.1:3000/swagger/index.html`.

## Documentation

Swagger documentation is available for all endpoints to facilitate integration and testing.
