# CREATE

## create_shop
```json
{
    "name": "Shop Name",
    "mail": "shop@domain.jp",
    "phone": "XXXX-XX-XXXX",
    "address": "Tokyo"
}
```

## create_group
```json
{
    "shop_id": 1,
    "name": "Group Name",
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
    "name": "Product Name",
    "max_people": 1,
    "qty": 10,
    "remark": "remark"
}
```

## create_price
```json
{
    "product_id": 1,
    "name": "Price Name",
    "value": 1000,
    "tax": 10,
    "remark": "remark"
}
```

## create_seat
```json
{
    "product_id": 1,
    "row": "A",
    "column": "1",
    "is_enable": true,
    "remark": "remark"
}
```

## create_stock
```json
{
    "price_id": 1,
    "name": "Stock Name",
    "qty": 10,
    "start_at": "2024-01-20T16:00:00+09:00",
    "end_at": "2024-01-20T18:00:00+09:00",
    "is_enable": true
}
```

## create_customer
```json
{
    "name": "Customer Name",
    "mail": "customer@domain.jp",
    "phone": "XXXX-XX-XXXX",
    "password": "PASSWORD",
    "address": "Tokyo",
    "payment_info": "Pay"
}
```