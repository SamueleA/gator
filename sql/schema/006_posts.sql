-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_posts_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE posts;
