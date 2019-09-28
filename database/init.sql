GRANT ALL PRIVILEGES ON DATABASE postgres TO postgres;

DROP TABLE  products;

CREATE TABLE products (
                       name VARCHAR(500) PRIMARY KEY UNIQUE ,
                       rating INT NOT NULL
);