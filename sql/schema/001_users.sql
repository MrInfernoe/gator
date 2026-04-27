-- +goose Up
CREATE TABLE users (
    id          UUID        PRIMARY KEY,
    created_at  TIMESTAMP   NOT NULL,
    updated_at  TIMESTAMP   NOT NULL,
    name        TEXT        NOT NULL,
    UNIQUE(name)
);

CREATE TABLE feeds (
    id          UUID        PRIMARY KEY,
    created_at  TIMESTAMP   NOT NULL,
    updated_at  TIMESTAMP   NOT NULL,
    name        TEXT        NOT NULL,
    url         TEXT        NOT NULL,
    user_id     UUID        NOT NULL REFERENCES users
                                        ON DELETE CASCADE
                                        ON UPDATE CASCADE,
    UNIQUE(url),
    FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE feeds;
DROP TABLE users;