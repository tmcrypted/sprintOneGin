-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20) NOT NULL DEFAULT 'worker',  -- owner | worker | moderator
    fio           VARCHAR(255) NOT NULL,
    photo_url     VARCHAR(512),
    lat           DECIMAL(10, 8),
    lng           DECIMAL(11, 8),
    rating_avg    DECIMAL(3, 2) DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ПВЗ владельца. На модерации до approved; только approved видны и доступны для листингов
CREATE TABLE pvz (
    id              BIGSERIAL PRIMARY KEY,
    owner_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending | approved | rejected
    city            VARCHAR(100) NOT NULL,
    address         VARCHAR(255) NOT NULL,
    company_name    VARCHAR(255) NOT NULL,
    description     TEXT,
    contact_phone   VARCHAR(50) NOT NULL,
    contact_telegram VARCHAR(100),
    lat             DECIMAL(10, 8),
    lng             DECIMAL(11, 8),
    moderated_at    TIMESTAMPTZ,
    moderated_by    BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pvz_owner ON pvz(owner_id);
CREATE INDEX idx_pvz_status ON pvz(status);

-- Вакансия на смену: владелец выбирает свой ПВЗ (из уже добавленных и approved)
CREATE TABLE listings (
    id            BIGSERIAL PRIMARY KEY,
    owner_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pvz_id        BIGINT NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    cells_count   INT NOT NULL,
    pay_per_shift INT NOT NULL,
    shift_date    DATE NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'active',  -- active | closed
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_listings_owner ON listings(owner_id);
CREATE INDEX idx_listings_pvz ON listings(pvz_id);
CREATE INDEX idx_listings_shift_date ON listings(shift_date);
CREATE INDEX idx_listings_status ON listings(status);

-- Заявка на вакансию. status=accepted => сделка (deal)
CREATE TABLE applications (
    id           BIGSERIAL PRIMARY KEY,
    listing_id   BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    applicant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending | accepted | rejected
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(listing_id, applicant_id)
);

CREATE INDEX idx_applications_listing ON applications(listing_id);
CREATE INDEX idx_applications_applicant ON applications(applicant_id);
CREATE INDEX idx_applications_status ON applications(status);

CREATE TABLE messages (
    id         BIGSERIAL PRIMARY KEY,
    deal_id    BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    sender_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_deal ON messages(deal_id);

-- Отзыв по отработанной сделке: о владельце/ПВЗ, привязка к PВЗ для страницы ПВЗ
CREATE TABLE reviews (
    id             BIGSERIAL PRIMARY KEY,
    deal_id        BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    pvz_id         BIGINT NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    author_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating         SMALLINT NOT NULL
    body           TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(deal_id, author_id)
);

CREATE INDEX idx_reviews_target ON reviews(target_user_id);
CREATE INDEX idx_reviews_pvz ON reviews(pvz_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS listings;
DROP TABLE IF EXISTS pvz;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
