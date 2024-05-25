CREATE INDEX idx_post_id ON comments(post_id);
CREATE INDEX idx_parent_id ON comments(parent_id);
CREATE INDEX idx_created_at ON comments(created_at);