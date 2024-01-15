package scheduler

import (
	"Go_APIServer/db"
	"Go_APIServer/funcs"
	"Go_APIServer/ini"
	"context"
	"fmt"
	"time"
)

type order struct {
	orders []db.OrderModel
}

type stock struct {
	id int
}

func Scheduler() {
	var cnt int
	delay := time.Duration(ini.OPTIONS.Delay) * time.Second

	fmt.Println("[Scheduler start] delay time :", delay)

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(fmt.Sprint("クライアント接続エラー : ", err))
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー : ", err))
		}
	}()

	ctx := context.Background()

	for {
		cnt++
		var update_time time.Time

		// ユーザによる時刻指定が有効な場合はOrderのend_atを監視
		// そうでない場合はStockを監視
		if ini.OPTIONS.Time_free_enable {

			// Orderテーブルで、現在時刻よりも後の情報を取得
			orders, err := client.Order.FindMany(
				db.Order.EndAt.After(funcs.GetTime()),
			).Exec(ctx)
			if err != nil {
				fmt.Println("Orderテーブル取得エラー :", err)
				return
			}

			// 現在より後の情報がない場合
			if len(orders) == 0 {
				fmt.Printf("[Sceduler.%d] 更新予定なし\n", cnt)
				time.Sleep(delay)
				continue
			}

			// 比較用にOrderテーブルの最初の1行を基準としてセット
			update_time, _ = orders[0].EndAt()

			for _, o := range orders {

				compared_time, _ := o.EndAt()

				// 終了時刻がより早い場合はその行を新たにセット
				if compared_time.Before(update_time) {
					update_time = compared_time
				}
			}

		} else {

			// Stockテーブルで、現在時刻よりも後の情報を取得
			stock, err := client.Stock.FindMany(
				db.Stock.EndAt.After(funcs.GetTime()),
			).Exec(ctx)
			if err != nil {
				fmt.Println("Stockテーブル取得エラー :", err)
				return
			}

			// 現在より後の情報がない場合
			if len(stock) == 0 {
				fmt.Printf("[Sceduler.%d] 更新予定なし\n", cnt)
				time.Sleep(delay)
				continue
			}

			// 比較用にStockテーブルの1行をセット
			stock_one := stock[0]
			var ve time.Time
			var se time.Time

			for _, value := range stock {
				// fmt.Printf("index : %d, value : %+v\n", index, value)

				ve, _ = value.EndAt()
				se, _ = stock_one.EndAt()

				// 終了時刻がより早い場合はその行を新たにセット
				if ve.Before(se) {
					stock_one = value
				}
			}
		}

		// 現在時刻との間隔を求める
		duration := update_time.Sub(funcs.GetTime())

		// durationがdelayよりも短い場合
		// その間隔分遅延し、遅延後にタスク処理を実行
		if duration < delay {
			fmt.Printf("[Sceduler.%d] 更新予定: %v後\n", cnt, duration)
			time.Sleep(duration)

			if ini.OPTIONS.Time_free_enable {

				// update_timeに一致するOrderを取得
				var order order
				order.orders, _ = client.Order.FindMany(
					db.Order.EndAt.Equals(update_time),
				).Exec(ctx)

				if err := order.task(client); err != nil {
					fmt.Printf("[Sceduler.%d] 更新エラー : %s\n", cnt, err)
					return
				}
			}

			// if err := task(client, stock_id); err != nil {
			// 	fmt.Printf("[Sceduler.%d] 更新エラー : %s\n", cnt, err)
			// 	return
			// }

			fmt.Printf("[Sceduler.%d] 更新完了\n", cnt)

		} else {
			fixed_time := update_time.In(time.Local)
			fmt.Printf("[Sceduler.%d] 更新予定: %v (%v後)\n", cnt, fixed_time, duration)
			time.Sleep(delay)
		}
	}
}
