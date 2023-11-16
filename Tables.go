// 各テーブル用の構造体まとめ
package main

import "time"

// 注文情報テーブル
type Order struct {
	Order_id   int       `json:"order_id"`
	Cust_id    int       `json:"cust_id"`
	Product_id int       `json:"product_id"`
	Num        int       `json:"order_num"`
	Time       time.Time `json:"order_time"`
}

// 在庫テーブル
type Stock struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Num  int    `json:"num"`
}
