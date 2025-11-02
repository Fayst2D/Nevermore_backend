# 1. Описание проекта 
Nevermore — это проект по созданию веб-платформы, призванной преодолеть ключевую проблему современного цифрового чтения — фрагментацию опыта. Сегодня читатели вынуждены использовать несколько разрозненных сервисов: одно приложение для чтения электронных книг, другое — для ведения читательского дневника, третье (например, форум или социальная сеть) — для обсуждения прочитанного. Nevermore интегрирует все эти функции в единое, безопасное и интеллектуальное пространство, превращая чтение из уединенного занятия в насыщенный социальный и аналитический опыт.

Основные возможности и функции:

1. Удобный интерфейс для чтения с адаптацией под разные устройства.

2. Создание и участие в виртуальных книжных клубах.

3. Совместное комментирование, обсуждения, марафоны чтения.

4. Отслеживание времени чтения, количества прочитанных книг, целей, визуализация прогресса.

5. Регулярное обновление функционала на основе отзывов.

# 2. стек используемых технологий:

- **Фреймворки:** `Gin`
- **Слой данных:**
    - База данных: `PostgreSQL`
    - ORM: `goose`
    - Хранилище объектов: `AWS S3 (minIO)`
- **Авторизация:** `JWT Bearer`
- **Представление:**
    - Документация API: `OpenAPI (Swagger)`
    - Связь в реальном времени: `gorilla/websocket`
- **Контейнеризация:** `Docker`
- **Языки программирования:** `Go`, `SQL`
- **Инструменты и IDE:** `Goland`, `VS Code`, `postman`

# 3. Роли пользователей и  описание их действий в системе

<img width="1094" height="871" alt="drawio (1)" src="https://github.com/user-attachments/assets/4ec3a8b4-bf84-4416-aecf-c8332c260667" />

# 4. Схема БД

```mermaid
erDiagram
    users {

        integer id PK
        varchar name
        varchar email
        varchar phone_number
        varchar password
        varchar role
        varchar photo
        timestamp created_at
        timestamp deleted_at
    }

    authors {
        integer id PK
        varchar name
        text biography
        varchar photo_url
        timestamp created_at
        timestamp updated_at
    }

    books {
        integer id PK
        varchar title
        text description
        varchar cover_image_url
        varchar file_url
        integer author_id FK
        integer uploaded_by FK
        timestamp created_at
        timestamp updated_at
    }

    reading_statuses {
        integer id PK
        varchar name
    }

    bookmarks {
        integer id PK
        integer user_id FK
        integer book_id FK
        integer status_id FK
        boolean favorite
        smallint personal_rating
        text personal_notes
        integer current_page
        timestamp created_at
        timestamp updated_at
    }

    reading_sessions {
        integer id PK
        integer user_id FK
        integer book_id FK
        timestamp start_time
        timestamp end_time
        integer pages_read
        interval duration
    }

    reviews {
        integer id PK
        integer book_id FK
        integer user_id FK
        smallint rating
        varchar title
        text content
        timestamp created_at
        timestamp updated_at
    }

    saved_authors {
        integer user_id PK,FK
        integer author_id PK,FK
        timestamp created_at
    }

    goose_db_version {
        integer id PK
        bigint version_id
        boolean is_applied
        timestamp tstamp
    }

    users ||--o{ books : "uploaded_by"
    users ||--o{ bookmarks : ""
    users ||--o{ reading_sessions : ""
    users ||--o{ reviews : ""
    users ||--o{ saved_authors : ""

    authors ||--o{ books : ""
    authors ||--o{ saved_authors : ""

    books ||--o{ bookmarks : ""
    books ||--o{ reading_sessions : ""
    books ||--o{ reviews : ""

    reading_statuses ||--o{ bookmarks : ""
```


# 5. API
Документация будет доступна по адресу: http://localhost:3000/docs/

# 6. Организация сетевого взаимодействия

1. Скачать любой VPN (например Proton VPN)
2. скачать ngrok
3. запустить VPN,
4. написать в ngrok ngrok http 3000
5. скомпилировать проект
6. передать полученную ссылку на фронт
<img width="1837" height="905" alt="image" src="https://github.com/user-attachments/assets/80bb587a-042a-44c1-9bdf-23216caf60f7" />

