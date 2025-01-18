CREATE TABLE urls (
    id TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    click_count INTEGER NOT NULL DEFAULT 0,
    created_date DATETIME NOT NULL
);