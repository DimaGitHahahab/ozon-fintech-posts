INSERT INTO posts (id, title, content, author_id, created_at, comments_disabled)
VALUES (1, 'Breaking news!', 'No way! Something happened!', 2, NOW(), FALSE),
       (2, 'Hello Reddit', 'This is my first post here', 1, NOW(), FALSE),
       (3, 'News', 'Something is happening again!', 2, NOW(), FALSE);