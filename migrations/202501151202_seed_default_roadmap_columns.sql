-- Insert default roadmap columns for all existing tenants
INSERT INTO roadmap_columns (tenant_id, name, slug, position, is_visible_to_public)
SELECT 
    id as tenant_id,
    'Staging' as name,
    'staging' as slug,
    0 as position,
    false as is_visible_to_public
FROM tenants
WHERE status = 1; -- Only active tenants

INSERT INTO roadmap_columns (tenant_id, name, slug, position, is_visible_to_public)
SELECT 
    id as tenant_id,
    'Now' as name,
    'now' as slug,
    1 as position,
    true as is_visible_to_public
FROM tenants
WHERE status = 1;

INSERT INTO roadmap_columns (tenant_id, name, slug, position, is_visible_to_public)
SELECT 
    id as tenant_id,
    'Next' as name,
    'next' as slug,
    2 as position,
    true as is_visible_to_public
FROM tenants
WHERE status = 1;

INSERT INTO roadmap_columns (tenant_id, name, slug, position, is_visible_to_public)
SELECT 
    id as tenant_id,
    'Later' as name,
    'later' as slug,
    3 as position,
    true as is_visible_to_public
FROM tenants
WHERE status = 1;
