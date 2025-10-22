CREATE TABLE IF NOT EXISTS roadmap_columns (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    position INT NOT NULL DEFAULT 0,
    is_visible_to_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, slug)
);

CREATE INDEX idx_roadmap_columns_tenant_position ON roadmap_columns(tenant_id, position);
