CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       author TEXT NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT now(),
                       comments_enabled BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
                          parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
                          text TEXT NOT NULL CHECK (char_length(text) <= 2000),
                          author TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT now()
);