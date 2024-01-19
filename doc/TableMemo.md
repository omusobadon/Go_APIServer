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

### Shop : 店舗情報
- id        : 店舗ID（自動生成）
- name      : 店舗名
- mail      : メールアドレス
- phone     : 電話番号
- address   : 住所

### ProductGroup : 商品グループ情報
- id                : 商品グループID（自動生成）
- shop_id           : 店舗ID
- name              : 商品グループ名
- start_before      : 予約開始可能時間(h)（詳細はTime.pngを参照）
- invalid_duration  : 予約無効期間(h)（詳細はTime.pngを参照）
- unit_time         : 時間単位(m)
- max_time          : 最大予約時間(h)
- interval          : インターバル

### Product : 商品情報
- id        : 商品ID（自動生成）
- group_id  : 商品グループID
- name      : 商品名
- max_people: 最大予約人数
- qty       : デフォルトの数量
- remark    : 備考
- img_data  : 商品画像データ

### Price : 価格情報
- id        : 価格ID（自動生成）
- product_id: 商品ID
- name      : 価格名
- value     : 価格
- tax       : 税率
- remark    : 備考

### Seat : 座席情報
- id        : 座席ID（自動生成）
- product_id: 商品ID
- row       : 座席（行）
- column    : 座席（列）
- is_enable : 座席の有効・無効
- remark    : 備考

### Stock : 在庫情報
- id        : 在庫ID（自動生成）
- price_id  : 価格ID
- name      : 在庫名
- qty       : 数量
- start_at  : 予約期間の開始時刻
- end_at    : 予約期間の終了時刻
- is_enable : 在庫の有効・無効

### SeatReservation : 座席予約ステータス（自動生成）
- id            : 座席予約ステータスID
- stock_id      : 在庫ID
- seat_id       : 座席ID
- is_reserved   : ステータス

### Customer : 顧客情報
- id            : 顧客ID（自動生成）
- name          : 顧客名
- mail          : メールアドレス
- phone         : 電話番号
- password      : パスワード
- address       : 住所
- payment_info  : 決済情報

### Order : 注文情報
- id            : 注文ID（自動生成）
- customer_id   : 顧客ID
- start_at      : 予約期間の開始時刻
- end_at        : 予約期間の終了時刻
- is_accepted   : 承認ステータス
- is_pending    : 保留ステータス
- created_at    : 注文日時
- remark        : 備考

### PaymentState : 決済ステータス
- id            : 決済ステータスID
- order_id      : 注文ID
- is_accepted   : 決済ステータス
- message       : メッセージ
- created_at    : 作成日時

### ReservationCancel : 予約キャンセル受付
- id            : 予約キャンセル受付ID
- order_id      : 注文ID
- is_accepted   : キャンセルステータス
- created_at    : 作成日時

### ReservationEnd : 予約終了受付
- id            : 予約終了受付ID
- order_id      : 注文ID
- is_accepted   : 終了ステータス
- created_at    : 作成日時