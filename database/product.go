package database

type Product struct {
	Name   string
	Rating int
}

func GetProducts(product string, rating int) ([]Product, error) {
	rows, err := postgres.Query("SELECT name, rating FROM products WHERE LOWER(name) LIKE LOWER($1) AND rating=$2", product, rating)
	if err != nil {
		return nil, err
	}

	products := make([]Product, 0)
	for rows.Next() {
		pr := Product{}
		err := rows.Scan(&pr.Name, &pr.Rating) //order matters
		if err != nil {
			return nil, err
		}
		products = append(products, pr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
func SaveProduct(product Product) error {
	_, err := postgres.Exec("INSERT INTO worker (name, rating) VALUES ($1, $2)", product.Name, product.Rating)
	return err
}
