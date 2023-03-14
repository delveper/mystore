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

```mermaid
graph LR;
    subgraph Entities
        E1[User]
        E2[Product]
    end
    subgraph Interactors
        IU1[UserInteractor]
        IP1[ProductInteractor]
    end
    subgraph Repository
        R1[PostgreSQL DB]
        R2[Database Connection]
    end
    subgraph InterfaceAdapters
        subgraph Controllers
            C1[UserController]
            C2[ProductController]
        end
        subgraph Presenters
            P1[UserPresenter]
            P2[ProductPresenter]
        end
    end
    subgraph Infrastructure
        IN1[Logger]
        IN2[Exception]
    end
    subgraph Transport
        T1[RESTful HTTP]
    end

    E1 --> IU1
    E2 --> IP1
    IU1 --> R1
    IP1 --> R1
    R2 --> R1
    R1 --> IU1
    R1 --> IP1
    IU1 --> C1
    IP1 --> C2
    C1 --> P1
    C2 --> P2
    P1 --> T1
    P2 --> T1
    IN1 --> T1
    IN2 --> T1
```
