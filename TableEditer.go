// テーブル編集用のメソッドまとめ
package main

import (
	"context"
	"encoding/json"

	"Go_APIServer/db"
)

// Order Insert
func (o Order) Insert(c *db.PrismaClient) error {
	ctx := context.Background()

	// Insert
	_, err := c.Order.CreateOne(
		db.Order.Customer.Set(o.Customer),
		db.Order.Product.Set(o.Product),
		db.Order.Start.Set(o.Start), // 追加
		db.Order.End.Set(o.End),     // 追加
		db.Order.Num.Set(o.Num),
		db.Order.Time.Set(o.Time),
		db.Order.State.Set(o.State), // 追加
		db.Order.Note.Set(o.Note),   // 追加
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Stock Update
func (s Stock) Update(c *db.PrismaClient) error {
	ctx := context.Background()

	// Update
	_, err := c.Stock.FindUnique(
		db.Stock.ID.Equals(s.ID),
	).Update(
		db.Stock.Name.Set(s.Name),
		db.Stock.Num.Set(s.Num),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Stock Insert
func (s Stock) Insert(c *db.PrismaClient) error {
	ctx := context.Background()

	// Insert
	_, err := c.Stock.CreateOne(
		db.Stock.Name.Set(s.Name),
		db.Stock.Num.Set(s.Num),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Stock Delete
func (s Stock) Delete(c *db.PrismaClient) error {
	ctx := context.Background()

	// Delete
	_, err := c.Stock.FindUnique(
		db.Stock.ID.Equals(s.ID),
	).Delete().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// EditInfo Insert
func (e EditInfo) Insert(c *db.PrismaClient) error {
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
		db.EditInfo.Info.Set(string(info_json)), // []byte から string への変換
		db.EditInfo.Time.Set(e.Time),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
