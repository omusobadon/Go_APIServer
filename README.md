## 環境構築方法
### 1. クローンの生成
```shell
    git clone https://github.com/omusobadon/Go_APIServer.git
```
### 2. 環境変数ファイルの作成
- DBのURLを記述した環境変数ファイル(.env)を"Go_APIServer/"へコピー
- 記述例）DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@[URL]/postgres

### 3. Prisma-Client-Goのインストール
```shell
    go get github.com/steebchen/prisma-client-go
```

### 4. /Go_APIServer内で以下のコマンドを実行してDBを同期（DB操作用のパッケージが生成される）
```shell
    go run github.com/steebchen/prisma-client-go db push
```

## 在庫情報のGET
- この情報をもとにユーザが注文を行う

### 各テーブルのGET
- Shop :        /get_shop
- ProductGroup :/get_group
- Product :     /get_product
- Price :       /get_price
- Seat :        /get_seat
- Stock :       /get_stock

### リクエストパラメータ
- 各テーブルについて、1つ前のテーブルIDを使用して絞り込み
- パラメータを設定しない場合は全取得する
- 例）/get_group?id=1　→　ProductGroupテーブルの"shop_id=1"の情報を取得
- 例）/get_group　→　ProductGroupテーブルを全取得

## 予約注文のPOST
- このPOSTを受け取り注文処理が行われる

#### 車の予約注文をPOSTする際のJSONパラメータ
- customer_id   : Customer（顧客情報）テーブルID
- stock_id      : Stock（在庫情報）テーブルID
- qty           : 数量
- start_at      : 予約開始時刻
- end_at        : 予約終了時刻
- remark        : 備考

- json記述例
```json
{
    "customer_id": 1,
    "stock_id": 1,
    "qty": 5,
    "start_at": "2023-11-10T10:10:00+09:00",
    "end_at": "2023-11-10T18:10:00+09:00",
    "remark": "test"
}
```

#### 車の予約注文に対するレスポンス
- 注文成功時と失敗時の2種類のレスポンスがある
- 以下それぞれの場合のJSONパラメータとその例

- 【成功時】
- message       : メッセージ（正常終了など）
- request       : POSTされた情報そのまま
- order         : Order（注文情報）テーブルに登録された情報
- order_detail  : OrderDetail（注文詳細）テーブルに登録された情報

```json
{
    "message": "正常終了",
    "request": {
        "customer_id": 1,
        "stock_id": 1,
        "seat_id": 0,
        "qty": 1,
        "start_at": "2023-12-30T18:10:00+09:00",
        "end_at": "2023-12-30T18:10:00+09:00",
        "number_people": 0,
        "remark": "test"
    },
    "order": {
        "id": 13, // 注文番号
        "customer_id": 1,
        "start_at": "2023-12-30T09:10:00Z",
        "end_at": "2023-12-30T09:10:00Z",
        "number_people": 0,
        "is_accepted": true, // 注文が承認されたか
        "created_at": "2023-12-28T06:59:10.697Z", // 注文日時
        "remark": "test"
    },
    "order_detail": {
        "id": 12, // 注文詳細番号
        "order_id": 13, // 注文番号
        "stock_id": 1,
        "qty": 1
    }
}
```

- 【失敗時】
- 失敗時は"message"と"request"情報のみのレスポンス

```json
{
    "message": "在庫不足",
    "request": {
        "customer_id": 1,
        "stock_id": 1,
        "seat_id": null,
        "qty": 1000,
        "start_at": "2023-12-30T18:10:00+09:00",
        "end_at": "2023-12-30T18:10:00+09:00",
        "number_people": 0,
        "remark": "test"
    }
}
```

## POST /edit
- 管理用
- POSTされたテーブル編集情報を取得して各テーブルを編集
- Type  1: Update, 2: Insert, 3: Delete
- Table テーブル名
- Info  更新内容

- 例）Insert
```json
{
    "type": 2,
    "table": "stock",
    "info": {
        "category": "car",
        "name": "car1",
        "value": 8000,
        "num": 22
    }
}
```

## ファイル一覧
- db/           prisma-client-goが作成したフォルダ。DB操作用のパッケージ等
- handlers/     api_serverから呼び出されるハンドラ群
- api_server    APIServerの本体
- get_time      NTP時刻取得処理
- schema        prismaの設定ファイル。DBのURLやテーブルの定義など