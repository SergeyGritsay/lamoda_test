package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"lamoda_test_task/pkg/models"
	"log"
)

const productTableName = "product"
const reserveTableName = "reserver"

type ProductPSQL struct {
	conn *sql.DB
}

func NewProductPSQL(conn *sql.DB) *ProductPSQL {
	return &ProductPSQL{conn: conn}
}

func (r *ProductPSQL) CreateNewProduct(name string, size float64, value int, stock_id int) (int, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var productId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (name, size, value, stock_id) values ($1, $2, $3, $4) RETURNING code", productTableName)
	fmt.Println(createItemQuery, tx)
	row := tx.QueryRow(createItemQuery, name, size, value, stock_id)
	err = row.Scan(&productId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return productId, tx.Commit()
}

func (r *ProductPSQL) GetProduct(code int) (models.Product, error) {
	var pr models.Product

	query := fmt.Sprintf(`SELECT code, name, size, value FROM %s where code = $1`, productTableName)

	row := r.conn.QueryRow(query, code)
	err := row.Scan(&pr.Code, &pr.Name, &pr.Size, &pr.Value)

	if err != nil {
		return models.Product{}, err
	}

	return pr, nil
}

func (r *ProductPSQL) GetProductList() ([]models.Product, error) {
	var prs []models.Product
	query := fmt.Sprintf(`SELECT code, name, size, value, stock_id FROM %s`, productTableName)

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		pr := models.Product{}
		if err = rows.Scan(&pr.Code, &pr.Name, &pr.Size, &pr.Value, &pr.StockId); err != nil {
			return nil, err
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *ProductPSQL) GetProductsCountByWarehouseId(stockId int, code int) (int64, error) {
	var count int64
	query := fmt.Sprintf(`select sum(value) from %s g
			inner join warehouse s on s.id = g.stock_id
		  where g.stock_id = $1 and s.available`, productTableName)

	query = query + ` and code = $2`

	if err := r.conn.QueryRow(query, stockId, code).Scan(&count); err != nil {
		log.Fatalf("request execution error: %s query: %s", err, query)
		return -1, nil
	}

	if err := r.conn.QueryRow(query, stockId).Scan(&count); err != nil {
		log.Fatalf("request execution error: %s query: %s", err, query)
		return -1, nil
	}

	return count, nil
}

func (r *ProductPSQL) ReservationProduct(code int, stockId int, value int64) error {
	var query string

	count, err := r.GetProductsCountByWarehouseId(stockId, code)
	if err != nil {
		log.Fatalf("reservation good error. couldn't get value of goods: %s", err)
		return err
	}

	if count < value {
		log.Fatal("is not possible to reserve a good because it is not in stock")
		return errors.New("cannot reserve 0 goods")
	}

	t, err := r.conn.Begin()
	if err != nil {
		return err
	}

	chErr := make(chan error)

	go func(errs chan error) {
		for err := range errs {
			if err != nil {
				log.Fatal(err)
			}
		}

	}(chErr)

	query = fmt.Sprintf(`update %s set value = (value - $1) where code = $2 and stock_id = $3`, productTableName)
	_, err = t.Exec(query, value, code, stockId)

	chErr <- err

	query = fmt.Sprintf(`insert into %s (good_code, stock_id, value) values ($1, $2, $3)`, reserveTableName)
	_, err = t.Exec(query, code, stockId, value)

	chErr <- err
	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductPSQL) CancelProductReservation(resId string) error {
	if resId == "" {
		return fmt.Errorf("resId is null")
	}

	t, err := r.conn.Begin()
	if err != nil {
		return err
	}

	var query string = fmt.Sprintf(`select good_code, value, stock_id from %s rc where rc.id = $1`, reserveTableName)

	var stock_id string
	var good_code, res_vl int64

	if err = t.QueryRow(query, resId).Scan(&good_code, &res_vl, &stock_id); err != nil {
		return err
	}

	chErr := make(chan error)

	go func(errs chan error) {
		for err := range errs {
			log.Fatal(err)
		}
	}(chErr)

	query = fmt.Sprintf(`delete from %s where id = $1`, reserveTableName)
	_, err = t.Exec(query, resId)

	chErr <- err

	query = fmt.Sprintf(`update %s set value = (select value from %s where code = $2) + $1 where code = $2`, productTableName, productTableName)
	_, err = t.Exec(query, res_vl, good_code)

	chErr <- err
	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}
