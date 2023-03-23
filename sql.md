# SQL tasks
## Top-10 products by popularity (number of ordered products)

```sql
SELECT p.id, p.merchant_id, p.name, p.description, p.price SUM(oi.Quantity) as ordered
FROM sales.products p
         LEFT JOIN sales.order_items oi
                   ON p.id = oi.product_id
WHERE p.deleted_at IS NULL
GROUP BY p.id
ORDER BY ordered DESC
LIMIT 10;
```

## All merchants with products `in_stock`

```sql
SELECT m.id, m.name, p.*
FROM sales.merchants m
         INNER JOIN sales.products p ON p.merchant_id= m.id
WHERE p.status= 'in_stock';
```

## All customers that spent more than 500 ¢
```sql
SELECT c.*, SUM(oi.price * oi.Quantity) as spended
FROM sales.customers c
         INNER JOIN orders o ON o.customer_id = c.id
         INNER JOIN order_items oi on o.id = oi.order_id
GROUP BY c.id
HAVING SUM(oi.price * oi.Quantity) > 500
ORDER BY spended DESC;
```


