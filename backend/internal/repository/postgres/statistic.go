package postgres

import (
	"database/sql"
	"fmt"
	"prediction_service/models"
	"time"
)

type PostgresStatisticRepository struct {
	db *sql.DB
}

func NewPostgresStatisticRepository(db *sql.DB) *PostgresStatisticRepository {
	return &PostgresStatisticRepository{db: db}
}

func (r *PostgresStatisticRepository) GetStatistic() (*models.StatisticResponse, error) {
	query := `
	SELECT 
		p.name,
		r.date,
		SUM(CASE WHEN r.rating = 'положительно' THEN 1 ELSE 0 END) AS positive,
		SUM(CASE WHEN r.rating = 'негативно' THEN 1 ELSE 0 END) AS negative,
		SUM(CASE WHEN r.rating = 'нейтрально' THEN 1 ELSE 0 END) AS neutral
	FROM Products p
	INNER JOIN ProductReviews r
		ON r.product_id = p.id
	GROUP BY p.name, r.date
	ORDER BY p.name;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productMap := make(map[string][]models.TimePoint)

	for rows.Next() {
		var name string
		var positive, negative, neutral int
		var date time.Time
		var dateSql sql.NullTime

		if err := rows.Scan(&name, &dateSql, &positive, &negative, &neutral); err != nil {
			return nil, err
		}

		if !dateSql.Valid {
			fmt.Println("REP: ", name)
		} else {
			date = dateSql.Time
		}

		productMap[name] = append(productMap[name], models.TimePoint{
			Date:     date,
			Positive: positive,
			Negative: negative,
			Neutral:  neutral,
		})
	}

	var response models.StatisticResponse
	for name, timeline := range productMap {
		response.Products = append(response.Products, models.Product{
			Name:     name,
			TimeLine: timeline,
		})
	}

	return &response, nil
}
