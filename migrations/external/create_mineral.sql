DROP TABLE IF EXISTS mineral;

CREATE TABLE mineral (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    amount BIGINT NOT NULL
);

INSERT INTO mineral VALUES (1, 'Copper', '100');
INSERT INTO mineral VALUES (2, 'Coal', '133');
