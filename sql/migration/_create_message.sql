-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
                          id VARCHAR(36) PRIMARY KEY,
                          user_id VARCHAR(36) NOT NULL,
                          username VARCHAR(255) NOT NULL,
                          content TEXT NOT NULL,
                          type VARCHAR(20) NOT NULL DEFAULT 'message' CHECK (type IN ('message', 'system', 'notification')),
                          timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          room_id VARCHAR(36) NOT NULL,

                          CONSTRAINT fk_messages_user
                              FOREIGN KEY (user_id)
                                  REFERENCES users(id)
                                  ON DELETE CASCADE,

                          CONSTRAINT fk_messages_room
                              FOREIGN KEY (room_id)
                                  REFERENCES rooms(id)
                                  ON DELETE CASCADE
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_messages_room_id ON messages(room_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp);
CREATE INDEX idx_messages_user_id ON messages(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd