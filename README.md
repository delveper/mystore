# mystore

API based on REST and Clean Architecture principles.

## Task

[Test-Task](TT.md)

## SQL

[Script](./sql/init_schema.sql) will be initialized upon deploying the Postgres image.

## Database diagram

```mermaid
erDiagram
    merchants {
        Int id PK
        String merchant_name
        Int phone
        String password
        String email
        timestamp created_at
    }

    buyers {
        Int id PK
        String full_name
        Int phone
        String email
        String password
        timestamp created_at
    }

    products {
        Int id PK
        Int merchant_id FK
        String name
        Text description
        BigInt price
        products_status status
        timestamp created_at
    }

    orders {
        Int id PK
        Int buyer_id FK
        order_status status
        timestamp created_at
        timestamp deleted_at
    }

    order_items {
        Int order_id FK
        Int product_id FK
        Int quantity
    }

    merchants ||--o{ products : has
    buyers ||--o{ orders : make
    products ||--o{ order_items : has
```

## Endpoints

Based on the REST design conventions it seems reasonable to follow a common pattern of using nested resources.
This approach makes API more organized and easier to navigate for clients and also reflect the relationship between
resources.

- `/api/merchants/:id/products`
- `/api/buyers/:id/orders`
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
| `/buyers`        | `POST`   | Creates a new buyer                     |
| `/buyers`        | `GET`    | Fetches a list of all buyers            |
| `/buyers/:id`    | `GET`    | Returns the buyer by id                 |
| `/buyers/:id`    | `PUT`    | Updates the buyer by id                 |
| `/buyers/:id`    | `DELETE` | Deletes the buyer by id                 |
| `/orders`        | `POST`   | Creates a new order                     |
| `/orders`        | `GET`    | Fetches a list of all orders            |
| `/orders/:id`    | `GET`    | Returns the order by id                 |
| `/orders/:id`    | `PUT`    | Updates the order by id                 |
| `/orders/:id`    | `DELETE` | Deletes the order by id                 |


