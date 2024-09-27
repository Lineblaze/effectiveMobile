CREATE TABLE songs_detail
(
    id BIGSERIAL PRIMARY KEY,
    "group" TEXT,
    song TEXT,
    release_date VARCHAR(50),
    "text" TEXT,
    link TEXT
);

