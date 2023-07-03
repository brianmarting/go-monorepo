DROP TABLE IF EXISTS users;

CREATE TABLE users (
                         id BIGSERIAL PRIMARY KEY,
                         external_id VARCHAR NOT NULL,
                         name VARCHAR NOT NULL,
                         password VARCHAR NOT NULL,
                         token_version BIGSERIAL NOT NULL
);

INSERT INTO users VALUES (1, '4f017282-48b1-49e4-b98f-3f1ac538ff42', 'john', 'test', 1);
