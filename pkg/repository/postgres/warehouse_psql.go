package postgres

import (
	"database/sql"
	"fmt"
	"lamoda_test_task/pkg/models"
)

const warehouseTableName = "warehouse"

type WarehosuePSQL struct {
	conn *sql.DB
}

func NewWarehousePSQL(conn *sql.DB) *WarehosuePSQL {
	return &WarehosuePSQL{conn: conn}
}

func (r *WarehosuePSQL) CreateNewWarehouse(name string, available bool) (int, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var warehouseId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (name, available) values ($1, $2) RETURNING id", warehouseTableName)

	row := tx.QueryRow(createItemQuery, name, available)
	err = row.Scan(&warehouseId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return warehouseId, tx.Commit()
}

func (r *WarehosuePSQL) GetWarehouse(id int) (models.Warehouse, error) {
	var wh models.Warehouse

	query := fmt.Sprintf(`SELECT id, name, available FROM %s where id = $1`, warehouseTableName)
	row := r.conn.QueryRow(query, id)
	err := row.Scan(&wh.ID, &wh.Name, &wh.IsAvailable)

	if err != nil {
		return models.Warehouse{}, err
	}

	return wh, nil
}

func (r *WarehosuePSQL) GetWarehouseList() ([]models.Warehouse, error) {
	var whs []models.Warehouse
	query := fmt.Sprintf(`SELECT id, name, available FROM %s`, warehouseTableName)

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		wh := models.Warehouse{}
		if err = rows.Scan(&wh.ID, &wh.Name, &wh.IsAvailable); err != nil {
			return nil, err
		}

		whs = append(whs, wh)
	}

	return whs, nil
}
