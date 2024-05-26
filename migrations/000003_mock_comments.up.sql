-- comments to posts
INSERT INTO comments (post_id, parent_id, author_id, content, created_at)
VALUES (1, NULL, 1, 'Wow, this is some breaking news indeed!', NOW()),
       (1, NULL, 2, 'I cannot believe this is happening!', NOW()),
       (2, NULL, 3, 'Welcome to Reddit! Looking forward to your posts.', NOW()),
       (2, NULL, 1, 'Nice to see new faces around here.', NOW()),
       (3, NULL, 2, 'Again? This is getting interesting.', NOW()),
       (3, NULL, 1, 'I wonder what will happen next.', NOW());

-- comments to comments
INSERT INTO comments (post_id, parent_id, author_id, content, created_at)
VALUES (1, 1, 2, 'I agree, this is indeed breaking news!', NOW()),
       (1, 2, 3, 'Me too', NOW()),
       (2, 3, 1, 'Thank you for the warm welcome!', NOW()),
       (2, 4, 2, 'It is always nice to see new people!', NOW()),
       (3, 5, 3, 'It is indeed getting interesting.', NOW()),
       (3, 6, 1, 'I am also curious about that', NOW());
