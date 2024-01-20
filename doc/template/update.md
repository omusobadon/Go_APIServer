# UPDATE
- 各テーブルのIDは必須
- それ以外は任意で、jsonに項目があるもののみ書き換えられる

## update_shop
```json
{
    "id": 1,
    "name": "Shop Name",
    "mail": "shop@domain.jp",
    "phone": "XXXX-XX-XXXX",
    "address": "Tokyo"
}
```

## update_group
```json
{
    "id": 1,
    "shop_id": 1,
    "name": "Group Name",
    "start_before": 24,
    "invalid_duration": 0,
    "unit_time": 60,
    "max_time": 72,
    "interval": 3
}
```

## update_product
```json
{
    "id": 1,
    "group_id": 1,
    "name": "Product Name",
    "max_people": 1,
    "qty": 10,
    "remark": "remark",
    "img_data": "Encoded data",
}
```

## update_price
```json
{
    "id": 1,
    "product_id": 1,
    "name": "Price Name",
    "value": 1000,
    "tax": 10,
    "remark": "remark"
}
```

## update_seat
```json
{
    "id": 1,
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
    "id": 1,
    "price_id": 1,
    "name": "Stock Name",
    "qty": 10,
    "start_at": "2024-01-20T16:00:00+09:00",
    "end_at": "2024-01-20T18:00:00+09:00",
    "is_enable": true
}
```

## update_customer
```json
{
    "id": 1,
    "name": "test",
    "mail": "test",
    "phone": "test",
    "password": "test",
    "address": "test",
    "payment_info": "test"
}
```

## update_order
```json
{
    "id": 1,
    "customer_id": 1,
    "start_at": "2024-01-20T16:00:00+09:00",
    "end_at": "2024-01-20T18:00:00+09:00",
    "is_accepted": false,
    "is_pending": false,
    "remark": "remark"
}
```

## update_order_detail
```json
{
    "id": 1,
    "order_id": 1,
    "stock_id": 1,
    "seat_id": 1,
    "number_people": 1,
    "qty": 1,
}
```