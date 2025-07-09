-- +goose Up
ALTER TABLE feeds
ADD last_fetched_at TIMESTAMP;

-- +goose down
ALTER TABLE feeds
DROP last_fetched_at;
