CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE songs
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "group" TEXT,
    song TEXT,
    release_date VARCHAR(50),
    "text" TEXT,
    link TEXT
);

