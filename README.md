# mystore

API based on REST and Clean Architecture principles.

## Usage

```shell
make docker # Build and deploy docker image with Postgres  
```

```shell
make run # Run server
```

```shell
make test # Run test
```

```shell
make bench # Run benchmark for task* (separate PR)
```



## SQL

[Script](./sql/init_schema.sql) will be initialized upon deploying the Postgres image.

## Database diagram

```mermaid
erDiagram
    merchants {
        int id PK
        string name
        int phone
        string password
        string email
        timestamp created_at
        timestamp deleted_at
    }

    customers {
        int id PK
        string full_name
        int phone
        string email
        string password
        timestamp created_at
        timestamp deleted_at
    }

    products {
        int id PK
        int merchant_id FK
        string name
        text description
        bigint price
        products_status status
        timestamp created_at
        timestamp deleted_at
    }

    orders {
        int id PK
        int customer_id FK
        order_status status
        timestamp created_at
        timestamp deleted_at
    }

    order_items {
        int order_id FK
        int product_id FK
        int quantity
        int price
    }

    merchants ||--o{ products : has
    customers ||--o{ orders : places 
    orders ||--|{ order_items : contains 
    products ||--|{ order_items :contains 
```

## Endpoints

Based on the REST design conventions it seems reasonable to follow a common pattern of using nested resources.
This approach makes API more organized and easier to navigate for clients and also reflect the relationship between
resources.

- `/api/merchants/:id/products`
- `/api/customers/:id/orders`
- `/api/orders/:id/items`
- etc.

Flattened structure will look like this.

| Endpoint         | Method   | Description                             |
|------------------|----------|-----------------------------------------|
| `/auth/signup`   | `POST`   | Create a new user account               |
| `/auth/login`    | `POST`   | Auth user with email/phone and password |
| `/auth/login`    | `POST`   | Log out the current user                |
| `/merchants`     | `POST`   | Creates a new merchant                  |
| `/merchants`     | `GET`    | Fetch list of all merchants             |
| `/merchants/:id` | `PUT`    | Updates the merchant by id              |
| `/merchants/:id` | `GET`    | Returns the merchant by id              |
| `/merchants/:id` | `DELETE` | Deletes the merchant by id              |
| `/products`      | `POST`   | Creates a new product                   |
| `/products`      | `GET`    | Fetches a list of all products          |
| `/products/:id`  | `GET`    | Returns the product by id               |
| `/products/:id`  | `PUT`    | Updates the product by id               |
| `/products/:id`  | `DELETE` | Deletes the product by id               |
| `/customers`     | `POST`   | Creates a new customer                  |
| `/customers`     | `GET`    | Fetches a list of all customers         |
| `/customers/:id` | `GET`    | Returns the customer by id              |
| `/customers/:id` | `PUT`    | Updates the customer by id              |
| `/customers/:id` | `DELETE` | Deletes the customer by id              |
| `/orders`        | `POST`   | Creates a new order                     |
| `/orders`        | `GET`    | Fetches a list of all orders            |
| `/orders/:id`    | `GET`    | Returns the order by id                 |
| `/orders/:id`    | `PUT`    | Updates the order by id                 |
| `/orders/:id`    | `DELETE` | Deletes the order by id                 |

## Graph
not finished...
```mermaid
graph TD;
  subgraph "Interactors Layer"
    PI((ProductInteractor)) -->|uses| PR(ProductRepo)
    PI -->|uses| L(Logger)
    subgraph "Entities Layer"
      E(Product) --- PI
    end
  end

  subgraph "Repository Layer"
    PR -->|implements| PIR(ProductRepo)
    subgraph "Database Layer"
      D(Database) --- PIR
    end
  end

  subgraph "Transport Layer"
    subgraph "API Handlers"
      H1[HTTP Handler 1] -->|uses| PI
      H2[HTTP Handler 2] -->|uses| PI
    end

    subgraph "Router"
      R[HTTP Router] -->|routes| H1
      R -->|routes| H2
    end
  end
```
