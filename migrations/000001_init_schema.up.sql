CREATE TABLE posts
(
    id                SERIAL PRIMARY KEY,
    title             TEXT NOT NULL,
    content           TEXT NOT NULL,
    author_id  INT       NOT NULL,
    created_at TIMESTAMP NOT NULL,
    comments_disabled BOOLEAN DEFAULT FALSE
);


CREATE TABLE comments
(
    id         SERIAL PRIMARY KEY,
    post_id    INT       NOT NULL,
    parent_id  INT,
    author_id  INT       NOT NULL,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments (id) ON DELETE CASCADE
);
