# CREATE

## create_shop
```json
{
    "name": "test",
    "mail": "test",
    "phone": "test",
    "password": "test",
    "address": "test"
}
```

## create_group
```json
{
    "shop_id": 1,
    "name": "test",
    "start_before": 24,
    "invalid_duration": 0,
    "unit_time": 60,
    "max_time": 72,
    "interval": 3
}
```

## create_product
```json
{
    "group_id": 1,
    "name": "test",
    "max_people": 1,
    "qty": 5,
    "remark": "test"
}
```

## create_price
```json
{
    "product_id": 1,
    "name": "test",
    "value": 1000,
    "tax": 10,
    "remark": "test"
}
```

## create_seat
```json
{
    "product_id": 1,
    "row": "test",
    "column": "test",
    "is_enable": true,
    "remark": "test"
}
```

## create_stock
```json
{
    "price_id": 1,
    "name": "test",
    "qty": 4,
    "start_at": true,
    "start_at": true,
    "is_enable": true
}
```

## create_customer
```json
{
    "name": "test",
    "mail": "test",
    "phone": "test",
    "password": "test",
    "address": "test",
    "payment_info": "test"
}
```