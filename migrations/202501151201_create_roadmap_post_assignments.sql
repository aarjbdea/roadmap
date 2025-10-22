CREATE TABLE IF NOT EXISTS roadmap_post_assignments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    column_id INT NOT NULL,
    tenant_id INT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by_id INT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (column_id) REFERENCES roadmap_columns(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(post_id, tenant_id)
);

CREATE INDEX idx_roadmap_assignments_column_position ON roadmap_post_assignments(column_id, position);
CREATE INDEX idx_roadmap_assignments_post ON roadmap_post_assignments(post_id);
