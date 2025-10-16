-- +goose Up
CREATE TABLE dbmappings (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS dbmappings;