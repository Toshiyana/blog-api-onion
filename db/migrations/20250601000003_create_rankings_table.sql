CREATE TABLE IF NOT EXISTS rankings (
    blog_id VARCHAR(36) PRIMARY KEY,
    ranking_position INT NOT NULL,
    score INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);

CREATE INDEX idx_rankings_position ON rankings(ranking_position);
CREATE INDEX idx_rankings_score ON rankings(score);
