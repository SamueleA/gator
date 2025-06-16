-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL,
  user_id UUID NOT NULL,
  CONSTRAINT fk_feed_follows_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
  CONSTRAINT fk_feed_follows_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;