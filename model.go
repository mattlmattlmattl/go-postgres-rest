package main

import (
	"database/sql"
)

type item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Notes string  `json:"notes"`
} // item struct

func (i *item) getItem(db *sql.DB) error {
	return db.QueryRow("SELECT name, notes FROM item WHERE id=$1", i.ID).Scan(&i.Name, &i.Notes)
} // getItem

func (i *item) updateItem(db *sql.DB) error {
	_, err := db.Exec("UPDATE item SET name=$1, notes=$2 WHERE id=$3", i.Name, i.Notes, i.ID)
	return err
} // updateItem

func (i *item) deleteItem(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM item WHERE id=$1", i.ID)
	return err
} // deleteItem

func (i *item) createItem(db *sql.DB) error {
	// postgres doesn't return the last inserted ID so this is the workaround
	err := db.QueryRow(
		"INSERT INTO item(name, notes) VALUES($1, $2) RETURNING id", i.Name, i.Notes).Scan(&i.ID)
	return err
} // createItem

func getItems(db *sql.DB, start, count int) ([]item, error) {
	rows, err := db.Query("SELECT id, name, notes FROM item LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []item{}

	for rows.Next() {
		var i item
		if err := rows.Scan(&i.ID, &i.Name, &i.Notes); err != nil {
			return nil, err
		}
		items = append(items, i)
	} // for

	return items, nil
} // getItems
