// テーブル編集用のメソッドまとめ
package main

import (
	"context"
	"encoding/json"

	"Go_APIServer/db"
)

// Stock Update
func (s Stock) Update(c *db.PrismaClient) error {
	ctx := context.Background()

	// Update
	_, err := c.Stock.FindUnique(
		db.Stock.ProductID.Equals(s.ID),
	).Update(
		db.Stock.ProductName.Set(s.Name),
		db.Stock.StockNum.Set(s.Num),
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
		db.Stock.ProductName.Set(s.Name),
		db.Stock.StockNum.Set(s.Num),
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
		db.Stock.ProductID.Equals(s.ID),
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
		db.EditInfo.Info.Set(info_json),
		db.EditInfo.Time.Set(e.Time),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
