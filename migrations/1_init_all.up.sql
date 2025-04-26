CREATE TABLE persons (
    id           SERIAL PRIMARY KEY,
    name         TEXT NOT NULL,
    surname      TEXT NOT NULL,
    patronymic   TEXT,
    age          INT    CHECK (age BETWEEN 0 AND 110),
    gender       TEXT   CHECK (gender IN ('male', 'female')),
    nationality  TEXT, 
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_persons_name_surname ON persons (lower(name), lower(surname));
CREATE INDEX idx_persons_created_at ON persons (created_at);