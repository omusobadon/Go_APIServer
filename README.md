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

## ファイル一覧
- db/           prisma-client-goが作成したフォルダ。DB操作用のパッケージ等
- handlers/     api_serverから呼び出されるハンドラ群
- api_server    APIServerの本体
- get_time      NTP時刻取得処理
- schema        prismaの設定ファイル。DBのURLやテーブルの定義など


## GET
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

## POST /post
- POSTされた注文情報を取得して注文処理

- json形式
```json
{
    "customer": 1,
    "product": 1,
    "start": "2023-11-10T10:10:00+09:00",
    "end": "2023-11-10T18:10:00+09:00",
    "num": 1
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
