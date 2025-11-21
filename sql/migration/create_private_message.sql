-- +goose Up
-- +goose StatementBegin
CREATE TABLE private_messages (
                                                id SERIAL PRIMARY KEY,
                                                sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_name VARCHAR(255) NOT NULL,
    receiver_name VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                                                                                                    );

CREATE INDEX idx_private_messages_sender_receiver ON private_messages(sender_id, receiver_id);
CREATE INDEX idx_private_messages_receiver_read ON private_messages(receiver_id, is_read);
CREATE INDEX idx_private_messages_created_at ON private_messages(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE private_messages;
-- +goose StatementEnd