CREATE TABLE label_templates (
    id UUID PRIMARY KEY,
    manufacturer_organization_name VARCHAR(255) NOT NULL,
    manufacturer_organization_address VARCHAR(255),
    manufacturer_email VARCHAR(255),
    manufacturer_site VARCHAR(255)
);

CREATE TABLE label_templates_events (
    aggregate_id UUID,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE INDEX label_templates_aggregate_id_idx ON label_templates_events (aggregate_id);

CREATE TABLE label_templates_history (
    label_template_id UUID NOT NULL,
    action VARCHAR(255) NOT NULL,
    new_manufacturer_organization_name VARCHAR(255),
    new_manufacturer_organization_address VARCHAR(255),
    new_manufacturer_phone VARCHAR(255),
    new_manufacturer_email VARCHAR(255),
    new_manufacturer_site VARCHAR(255),
    created_at TIMESTAMP NOT NULL
);
CREATE INDEX label_templates_history_label_template_id_idx ON label_templates_history (label_template_id);
CREATE UNIQUE INDEX label_templates_history_label_template_id_order_key_idx ON label_templates_history (label_template_id, created_at);