### 環境構築方法 ###
1.　クローンの生成
    git clone https://github.com/omusobadon/Go_APIServer.git

2.　環境変数ファイルの作成
    Discordの GO-API repo の環境変数ファイルを Go_APIServer/ へコピー

3.　/Go_APIServer内で以下のコマンドを実行してDBを同期（DB操作用のパッケージが生成される）
    go run github.com/steebchen/prisma-client-go db push

### ファイル一覧　###
db/             prisma-client-goが作成したフォルダ。DB操作用のパッケージ等
Go_APIServer    APIServerの本体
GetTime         時刻同期処理
schema          prismaの設定ファイル。DBのURLやテーブルの定義など
TableEditer     APIServerから実行されるテーブル編集用のメソッド群
TableMemo       作成するテーブルのメモ
Tables          各テーブル用の構造体のまとめ
test            テスト用

### NAT設定 ###
ip nat inside source static tcp 192.168.1.7 8080 interface GigabitEthernet8 8080

### POST /post ###
・POSTされた注文情報を取得して注文処理

・json形式
{
    "cust_id": 100,
    "product_id": 2,
    "order_num": 43
}

### POST /edit ###
・管理用
・POSTされたテーブル編集情報を取得して各テーブルを編集
・Type  1:Update, 2:Insert, 3:Delete
・Table テーブル名
・Info  更新内容

・json形式
{
    "type": 1,
    "table": "stock",
    "info": {
        "id": 1,
        "name": "car1",
        "num": 2
    }
}