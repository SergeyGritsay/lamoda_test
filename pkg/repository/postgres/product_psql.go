package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"lamoda_test_task/pkg/models"
	"log"
)

const productTableName = "product"

type ProductPSQL struct {
	conn *sql.DB
}

func NewProductPSQL(conn *sql.DB) *ProductPSQL {
	return &ProductPSQL{conn: conn}
}

func (r *ProductPSQL) CreateNewProduct(name string, size int, value int) (int, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var productId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (name, size, value) values ($1, $2, $3) RETURNING code", productTableName)

	row := tx.QueryRow(createItemQuery, name, size, value)
	err = row.Scan(&productId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return productId, tx.Commit()
}

func (r *ProductPSQL) GetProduct(code int) (models.Product, error) {
	var pr models.Product

	query := fmt.Sprintf(`SELECT code, name, size, value FROM %s where id = $1`, productTableName)
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
		if err = rows.Scan(&pr.Code, &pr.Name, &pr.Size, &pr.Value); err != nil {
			return nil, err
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *ProductPSQL) GetProductsCountByWarehouseId(stockId string, code int) (int64, error) {
	var count int64
	q := `select sum(value) from product g
			inner join warehouse s on s.id = g.stock_id
		  where g.stock_id = $1 and s.available`

	q = q + ` and code::text = $2`

	if err := r.conn.QueryRow(q, stockId, code).Scan(&count); err != nil {
		log.Fatalf("request execution error: %s query: %s", err, q)
		return -1, nil
	}

	if err := r.conn.QueryRow(q, stockId).Scan(&count); err != nil {
		log.Fatalf("request execution error: %s query: %s", err, q)
		return -1, nil
	}

	return count, nil
}

func (r *ProductPSQL) ReservationProduct(code int, stockId string, value int64) error {
	if stockId == "" || value == 0 {
		return fmt.Errorf("result: code = %d, stock_id = %s, value = %d. must not be equal to code == '' or stock id == '' or value == 0", code, stockId, value)
	}

	var q string

	c, err := r.GetProductsCountByWarehouseId(stockId, code)
	if err != nil {
		log.Fatalf("reservation good error. couldn't get value of goods: %s", err)
		return err
	}

	if c < value {
		log.Fatal("is not possible to reserve a good because it is not in stock")
		return errors.New("cannot reserve 0 goods")
	}

	// Открываем транзакцию. Обновляем значения в goods и res_cen
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

	// Вычетаем количество резервируемого товара из таблицы goods
	q = `update product set value = (value - $1) where code::text = $2 and stock_id = $3`
	_, err = t.Exec(q, value, code, stockId)

	chErr <- err

	// Создаем новую строку в таблице res_cen
	q = `insert into res_cen (good_code, stock_id, value) values ($1, $2, $3)`
	_, err = t.Exec(q, code, stockId, value)

	chErr <- err
	// Фиксируем транзакцию, если все окей
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

	var q string = `select good_code, value, stock_id::text from res_cen rc where rc.id::text = $1`

	var stock_id string
	var good_code, res_vl int64

	if err = t.QueryRow(q, resId).Scan(&good_code, &res_vl, &stock_id); err != nil {
		return err
	}

	chErr := make(chan error)

	go func(errs chan error) {
		for err := range errs {
			log.Fatal(err)
		}
	}(chErr)

	q = `delete from res_cen where id::text = $1`
	_, err = t.Exec(q, resId)

	chErr <- err

	q = `update product set value = (select value from product where code = $2) + $1 where code = $2`
	_, err = t.Exec(q, res_vl, good_code)

	chErr <- err
	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductPSQL) AddProduct(code int, stockId string, value int64, dynamic bool) error {
	if stockId == "" || value == 0 {
		return errors.New("code, id stock or value cannot be empty")
	}

	var q string

	q = `select available from warehouse where id::text = $1`

	t, err := r.conn.Begin()
	if err != nil {
		return err
	}
	log.Println("Adding product")
	var s models.Warehouse

	if err = t.QueryRow(q, stockId).Scan(&s.IsAvailable); err != nil {
		log.Println(err)
		return err
	}

	if s.IsAvailable {
		q = `update product set value = (select value from product where code::text = $2::text) + $1 where code::text = $2::text and stock_id = $3`

		_, err = t.Exec(q, value, code, stockId)
	} else {
		switch dynamic {
		case true:
			q = `select s.id from warehouse s
							inner join product g on g.code::text = $1 and g.stock_id = $2
						where s.id != $2 and s.available limit 1`

			if err = t.QueryRow(q, code, stockId).Scan(&s.ID); err != nil {
				return err
			}

			q = `update product set value = (select value from product where code::text = $2::text) + $1 where stock_id = $3` // TODO: Пофиксить запрос

			_, err = t.Exec(q, value, code, s.ID)
			return errors.New("failed to add goods to the stock")
		}
	}

	if err := t.Commit(); err != nil {
		return err
	}
	return nil
}
