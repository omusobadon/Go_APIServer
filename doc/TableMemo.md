### テーブル一覧
- Shop              : 店舗情報
- ProductGroup      : 商品グループ情報
- Product           : 商品情報
- Price             : 価格情報
- Seat              : 座席情報
- Stock             : 在庫情報
- SeatReservation   : 座席予約ステータス
- Customer          : 顧客情報
- Order             : 注文情報
- OrderDetail       : 注文詳細情報
- PaymentState      : 決済ステータス
- ReservationCancel : 予約キャンセル受付
- ReservationEnd    : 予約終了受付
- EditInfo          : DB編集履歴

### 各テーブルについて
　各テーブル最初のカラムであるIDは、Prismaによって管理される。
　そのIDはテーブルにおいて固有の値であり、新たな行が増えるとインクリメントされたIDが自動で生成される。

### テーブル内容の見方
- (カラム概要) : (jsonタグかつDBでのカラム名)

### Shop : 店舗情報
- id        : 店舗ID（自動生成）
- name      : 店舗名
- mail      : メールアドレス
- phone     : 電話番号
- address   : 住所

### ProductGroup : 商品グループ情報
- id: 商品グループID（自動生成）
- shop_id: 店舗ID
- name: 商品グループ名
- start_before: 

### Order : 注文情報
- 注文ID : id
- 注文詳細ID : detail_id
- 顧客ID : customer_id
- 在庫ID : stock_id
- 座席ID : seat_id
- 注文数 : num
- 予約開始日時 : start
- 予約終了日時 : end
- 予約日時 : time
- 予約状態 : state（確定、保留、キャンセルなど）
- 備考 : note

### 在庫テーブル : Stock
- 在庫ID : id
- 商品ID : product
- 開始時刻 : start（時刻指定する場合。例：映画の開始時刻）
- 終了時刻 : end
- 在庫数 : num
- 在庫状態 : state（予約を受け付けるかなどの状態）

### 商品テーブル : Product
- 商品ID : id
- グループID : group_id
- 商品名 : name
- 個数 : num（デフォルトの個数。例：映画の座席数）
- 価格 : value
- 備考 : note

### 商品グループテーブル : ProductGroup
- グループID : id
- 商品カテゴリ : category
- グループ名 : name
- インターバル : interval（次の開始時刻までの間隔）

### 座席テーブル : Seat
- 座席ID : id
- 商品ID : product_id
- 座席名 : name
- 座席状態 : state

### 料金計算テーブル : Fee
- 料金計算ID : id
- 商品ID : product_id
- 税率(%) : tax
- 割引(%) : special

### 決済処理テーブル : Payment
- 決済処理ID : id
- 注文ID : order_id
- 決済状態 : state（支払い完了、未支払い、払い戻し済みなど）

### 顧客情報テーブル : Customer
- 顧客ID : id
- 氏名 : name
- メール : mail
- 電話番号 : phone
- 住所 : address
- 決済情報 : payment（クレジットカードデータなど）

### アクティブ時間テーブル : ActiveTime
- ID : id
- 注文ID : order_id
- 予約開始日時 : start
- 予約終了日時 : end

### テーブル編集履歴テーブル : EditInfo
- 編集ID : id
- テーブル名 : table
- 編集タイプ : type
- 編集内容 : info
- 編集時刻 : time