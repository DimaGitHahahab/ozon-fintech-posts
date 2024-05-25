-- comments to posts
INSERT INTO comments (id, post_id, parent_id, author_id, content, created_at)
VALUES (1, 1, NULL, 1, 'Wow, this is some breaking news indeed!', NOW()),
       (2, 1, NULL, 2, 'I cannot believe this is happening!', NOW()),
       (3, 2, NULL, 3, 'Welcome to Reddit! Looking forward to your posts.', NOW()),
       (4, 2, NULL, 1, 'Nice to see new faces around here.', NOW()),
       (5, 3, NULL, 2, 'Again? This is getting interesting.', NOW()),
       (6, 3, NULL, 1, 'I wonder what will happen next.', NOW());

-- comments to comments
INSERT INTO comments (id, post_id, parent_id, author_id, content, created_at)
VALUES (7, 1, 1, 2, 'I agree, this is indeed breaking news!', NOW()),
       (8, 1, 2, 3, 'Me too', NOW()),
       (9, 2, 3, 1, 'Thank you for the warm welcome!', NOW()),
       (10, 2, 4, 2, 'It is always nice to see new people!', NOW()),
       (11, 3, 5, 3, 'It is indeed getting interesting.', NOW()),
       (12, 3, 6, 1, 'I am also curious about that', NOW());
