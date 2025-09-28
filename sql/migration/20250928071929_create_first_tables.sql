-- +goose Up
-- +goose StatementBegin
CREATE TABLE authors (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(50) NOT NULL,
                         biography TEXT,
                         photo_url VARCHAR(255),
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reading_statuses (
                                  id SERIAL PRIMARY KEY,
                                  name VARCHAR(50) NOT NULL UNIQUE -- например: "Читаю", "В планах", "Прочитано", "Брошено"
);

CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       cover_image_url VARCHAR(255),
                       file_url VARCHAR(255) NOT NULL,
                       author_id INTEGER NOT NULL REFERENCES authors(id) ON DELETE CASCADE, -- Прямая ссылка на автора
                       uploaded_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
                         id SERIAL PRIMARY KEY,
                         book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
                         user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
                         title VARCHAR(255),
                         content TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         UNIQUE(book_id, user_id)
);

CREATE TABLE bookmarks (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
                           status_id INTEGER NOT NULL REFERENCES reading_statuses(id) DEFAULT 2, -- Статус "В планах" по умолчанию
                           favorite BOOLEAN DEFAULT FALSE,
                           personal_rating SMALLINT CHECK (personal_rating >= 1 AND personal_rating <= 10),
                           personal_notes TEXT,
                           current_page INTEGER DEFAULT 0, -- Текущая страница/глава
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           UNIQUE(user_id, book_id)
);

CREATE TABLE reading_sessions (
                                  id SERIAL PRIMARY KEY,
                                  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                  book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
                                  start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                  end_time TIMESTAMP WITH TIME ZONE,
                                  pages_read INTEGER DEFAULT 0,
                                  duration INTERVAL GENERATED ALWAYS AS (end_time - start_time) STORED
);

CREATE TABLE saved_authors (
                               user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               author_id INTEGER NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               PRIMARY KEY (user_id, author_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authors, reading_statuses, books, reviews, bookmarks, reading_sessions, saved_authors;
-- +goose StatementEnd
