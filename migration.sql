CREATE TABLE label_templates (
    id UUID PRIMARY KEY,
    manufacturer_organization_name VARCHAR(255) NOT NULL
--     manufacturer_address VARCHAR(255),
--     manufacturer_phone VARCHAR(255),
--     manufacturer_email VARCHAR(255),
--     manufacturer_site VARCHAR(255),
--     manufacturer_logo VARCHAR(255)
);

CREATE TABLE label_templates_events (
    aggregate_id UUID,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX label_templates_aggregate_id_idx ON label_templates_events (aggregate_id);