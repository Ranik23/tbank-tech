-- +goose Up
CREATE TABLE IF NOT EXISTS link_users (
    link_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (link_id, user_id),
    CONSTRAINT fk_link FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS link_users;
