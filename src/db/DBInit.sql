-- Prepares the PostgreSQL database tables, functions, etc
CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(36) NOT NULL DEFAULT '',
    created_on TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_on TIMESTAMP WITH TIME ZONE NOT NULL,
    email TEXT NOT NULL DEFAULT '' CHECK(length(email) <= 320),
    name TEXT NOT NULL DEFAULT '' CHECK(length(name) <= 64),
    bio TEXT NOT NULL DEFAULT '' CHECK(length(bio) <= 384),
    profile_photo_id TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL,
    sign_in_attempts INT NOT NULL DEFAULT 0 CHECK(sign_in_attempts >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    account_reset_token TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS poems(
    id VARCHAR(36) NOT NULL DEFAULT '',
    created_on TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_on TIMESTAMP WITH TIME ZONE NOT NULL,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    title TEXT NOT NULL DEFAULT '' CHECK(length(title) <= 256),
    text TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS comments(
    id VARCHAR(36) NOT NULL DEFAULT '',
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    poem_id VARCHAR(36) NOT NULL REFERENCES poems(id),
    comment_id VARCHAR(36) NOT NULL DEFAULT '',
    text TEXT NOT NULL CHECK(length(text) <= 384),
    created_on TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users_followings(
    id VARCHAR(36) NOT NULL DEFAULT '',
    follower_id VARCHAR(36) NOT NULL REFERENCES users(id),
    following_id VARCHAR(36) NOT NULL REFERENCES users(id),
    created_on TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (follower_id, following_id)
);

CREATE TABLE IF NOT EXISTS poems_likes(
    id VARCHAR(36) NOT NULL DEFAULT '',
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    poem_id VARCHAR(36) NOT NULL REFERENCES poems(id),
    created_on TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (user_id, poem_id)
);

CREATE INDEX IF NOT EXISTS users_txt_search_idx
    ON users
        USING GIN (to_tsvector('english', name || ' ' || bio));

CREATE INDEX IF NOT EXISTS poems_txt_search_idx
    ON poems
        USING GIN (to_tsvector('english', title || ' ' || text));
