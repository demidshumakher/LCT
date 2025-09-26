CREATE TYPE rating_type AS ENUM ('положительно', 'негативно', 'нейтрально');

CREATE TABLE Products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE ProductReviews (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES Products(id),
    date DATE NOT NULL,
    rating rating_type NOT NULL
);

