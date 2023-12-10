// テーブル編集用のメソッドまとめ
package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"Go_APIServer/db"
)

// Order Insert
func (o *Order) Insert(c *db.PrismaClient) error {
	ctx := context.Background()

	// Insert
	_, err := c.Order.CreateOne(
		db.Order.Customer.Set(o.Customer),
		db.Order.Product.Set(o.Product),
		db.Order.Start.Set(o.Start),
		db.Order.End.Set(o.End),
		db.Order.Num.Set(o.Num),
		db.Order.Time.Set(o.Time),
		db.Order.State.Set(o.State),
		db.Order.Note.Set(o.Note),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Stock Update
func (s *Stock) Update(c *db.PrismaClient) error {
	ctx := context.Background()

	// Stock構造体の値の有無で場合分け
	var p []db.StockSetParam
	if s.ID == nil {
		return errors.New("エラー : 在庫IDがありません")
	}
	if s.Start != nil {
		p = append(p, db.Stock.Start.Set(*s.Start))
	}
	if s.End != nil {
		p = append(p, db.Stock.End.Set(*s.End))
	}
	if s.Interval != nil {
		p = append(p, db.Stock.Interval.Set(*s.Interval))
	}
	if s.Num != nil {
		p = append(p, db.Stock.Num.Set(*s.Num))
	}
	if s.State != nil {
		p = append(p, db.Stock.State.Set(*s.State))
	}

	// Update
	_, err := c.Stock.FindUnique(
		db.Stock.ID.Equals(*s.ID),
	).Update(p...).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Stock Insert
func (s *Stock) Insert(c *db.PrismaClient) error {
	ctx := context.Background()

	// 値がnilの場合は初期化
	// if s.Category == nil {
	// 	category := ""
	// 	s.Category = &category
	// }
	// if s.Name == nil {
	// 	name := ""
	// 	s.Name = &name
	// }
	// if s.Interval == nil {
	// 	interval := ""
	// 	s.Interval = &interval
	// }
	// if s.Value == nil {
	// 	value := 0
	// 	s.Value = &value
	// }
	// if s.Num == nil {
	// 	num := 0
	// 	s.Num = &num
	// }
	// if s.Note == nil {
	// 	note := ""
	// 	s.Note = &note
	// }

	// Insert
	_, err := c.Stock.CreateOne(
		db.Stock.Product.Set(*s.Product),
		db.Stock.Start.Set(*s.Start),
		db.Stock.End.Set(*s.End),
		db.Stock.Interval.Set(*s.Interval),
		db.Stock.Num.Set(*s.Num),
		db.Stock.State.Set(*s.State),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Stock Delete
func (s *Stock) Delete(c *db.PrismaClient) error {
	ctx := context.Background()

	if s.ID == nil {
		return errors.New("IDがありません")
	}

	// Delete
	_, err := c.Stock.FindUnique(
		db.Stock.ID.Equals(*s.ID),
	).Delete().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// EditInfo Insert
func (e *EditInfo) Insert(c *db.PrismaClient) error {
	ctx := context.Background()

	// DB登録のためMap型をjsonに変換
	info_json, err := json.Marshal(e.Info)
	if err != nil {
		return err
	}

	// Insert
	_, err = c.EditInfo.CreateOne(
		db.EditInfo.Table.Set(e.Table),
		db.EditInfo.Type.Set(e.Type),
		db.EditInfo.Info.Set(info_json),
		db.EditInfo.Time.Set(e.Time),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
