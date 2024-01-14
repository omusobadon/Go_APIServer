package srcbk

// // 商品・在庫テーブルが空の場合、自動生成するかどうか（テスト用）
// const auto_insert bool = false

// // SeatReservationの自動生成
// func generateSeatReservation() error {

// 	// データベース接続用クライアントの作成
// 	client := db.NewClient()
// 	if err := client.Prisma.Connect(); err != nil {
// 		return fmt.Errorf("クライアント接続エラー : %w", err)
// 	}
// 	defer func() {
// 		// クライアントの切断
// 		if err := client.Prisma.Disconnect(); err != nil {
// 			panic(fmt.Sprint("クライアント切断エラー : ", err))
// 		}
// 	}()

// 	ctx := context.Background()

// 	// Stockの取得
// 	stock, err := client.Stock.FindMany().Exec(ctx)
// 	if err != nil {
// 		return fmt.Errorf("stock取得エラー : %w", err)
// 	}

// 	for _, st := range stock {

// 		// Seatの取得
// 		seat, err := client.Seat.FindMany().Exec(ctx)
// 		if err != nil {
// 			return fmt.Errorf("seat取得エラー : %w", err)
// 		}

// 		for _, se := range seat {

// 			// SeatReservationにインサート
// 			_, err := client.SeatReservation.CreateOne(
// 				db.SeatReservation.Stock.Link(
// 					db.Stock.ID.Equals(st.ID),
// 				),
// 				db.SeatReservation.Seat.Link(
// 					db.Seat.ID.Equals(se.ID),
// 				),
// 				db.SeatReservation.IsReserved.Set(false),
// 			).Exec(ctx)
// 			if err != nil {
// 				return fmt.Errorf("SeatReservationインサートエラー : %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }

// func init() {

// 	// テーブルに情報がない場合に自動インサート（テスト用）
// 	if auto_insert {
// 		if err := AutoInsert(); err != nil {
// 			fmt.Println(err)
// 			fmt.Println("処理を続行します")
// 		}
// 	}

// 	// SeatReservationの作成
// 	if post.OPTIONS.Seat_enable {
// 		err := generateSeatReservation()
// 		if err != nil {
// 			panic(fmt.Sprintln("SeatReservation作成エラー :", err))
// 		}
// 	}

// 	// 各テーブルのチェック
// 	// 未実装

// }
