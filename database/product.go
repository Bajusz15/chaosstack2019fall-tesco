package database

type Product struct {
	Name   string
	Rating int
}

func GetProducts(product string, rating int) ([]Product, error) {
	product = "%" + product + "%"
	rows, err := postgres.Query("SELECT name, rating FROM products WHERE name LIKE $1 AND rating=$2", product, rating)
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
	_, err := postgres.Exec("INSERT INTO products (name, rating) VALUES ($1, $2) ", product.Name, product.Rating)
	if err != nil {
		_, err := postgres.Exec("UPDATE products SET rating=$1 WHERE name=$2", product.Rating, product.Name)
		return err
	}
	return nil
}
