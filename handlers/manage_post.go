// テーブル編集情報のPOST
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DB編集リクエストを変換する構造体
// Table テーブル名
// Type 1:Update, 2:Insert, 3:Delete
// Info 更新内容
type MPostRequestBody struct {
	Table string         `json:"table"`
	Type  int            `json:"type"`
	Info  map[string]any `json:"info"`
}

// レスポンスに変換する構造体
type MPostResponseBody struct {
	Message  string   `json:"message"`
	EditInfo EditInfo `json:"edit_info"`
}

var mpost_cnt int // managePostのカウント用

func ManagePost(w http.ResponseWriter, r *http.Request) {
	var req MPostRequestBody
	var edit EditInfo
	var status int
	var message string
	mpost_cnt++

	fmt.Printf("### Manage Post No.%d ###\n", mpost_cnt)

	// リクエスト処理後のレスポンス処理
	defer func() {
		// レスポンスボディの作成
		res := MPostResponseBody{
			Message:  message,
			EditInfo: edit,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		// 処理結果メッセージの表示（サーバ側）
		if status == 0 || message == "" {
			fmt.Println("ステータスコードまたはメッセージがありません")
		} else {
			fmt.Printf("[%d] %s\n", status, message)
		}

		fmt.Printf("### Manage Post No.%d End ###\n", mpost_cnt)

	}()

	// 更新情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	fmt.Println("編集情報 :", req)

	// リクエストをEditInfoテーブルにコピー
	edit.Table = req.Table
	edit.Type = req.Type
	edit.Info = req.Info

	// 編集時刻を取得
	edit.Time = GetTime()

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("クライアント接続エラー : ", err)
		return
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー : ", err))
		}
	}()

	ctx := context.Background()

	// mapを各テーブル用の構造体に変換するため、一度jsonに変換
	edit_json, err := json.Marshal(edit.Info)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("infoエンコードエラー :", err)
		return
	}

	// 各テーブルごとに処理を分岐
	if edit.Table == "stock" {
		var stock *Stock

		// 変換したjsonをStockに変換
		if err := json.Unmarshal(edit_json, &stock); err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("infoデコードエラー :", err)
			return
		}

		// 編集タイプごとに処理を分岐
		// Type 1:Update, 2:Insert, 3:Delete
		if edit.Type == 1 {
			if err := stock.Update(client); err != nil {
				status = http.StatusBadRequest
				message = fmt.Sprint("Stock Updateエラー :", err)
				return
			}

		} else if edit.Type == 2 {
			if err := stock.Insert(client); err != nil {
				status = http.StatusBadRequest
				message = fmt.Sprint("Stock Insertエラー :", err)
				return
			}

		} else if edit.Type == 3 {
			if err := stock.Delete(client); err != nil {
				status = http.StatusBadRequest
				message = fmt.Sprint("Stock Deleteエラー :", err)
				return
			}

		} else {
			status = http.StatusBadRequest
			message = "エラー : Type is not found"
			return
		}
	}

	// 処理が正常終了したらManageテーブルに登録
	if err := edit.Insert(client); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("EditInfo Insertエラー :", err)
		return
	}

	// 時刻をもとにテーブルを検索してIDを取得
	edit_info, err := client.EditInfo.FindFirst(
		db.EditInfo.Time.Equals(edit.Time),
	).Exec(ctx)
	if err != nil {
		status = http.StatusInternalServerError
		message = fmt.Sprint("編集ID取得エラー : ", err)
		return
	}

	edit.ID = edit_info.ID
	fmt.Printf("編集終了 : %+v\n", edit)

	status = http.StatusOK
	message = "正常終了"
}
