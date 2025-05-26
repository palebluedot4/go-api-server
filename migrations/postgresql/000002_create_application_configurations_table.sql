CREATE TABLE
  application_configurations (
    config_key VARCHAR(32) PRIMARY KEY,
    config_value TEXT NOT NULL,
    value_type VARCHAR(16) NOT NULL DEFAULT 'TEXT',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

COMMENT ON TABLE application_configurations IS 'Stores application-level configurations, settings, and version information.';

COMMENT ON COLUMN application_configurations.config_key IS 'The unique key for the configuration parameter.';

COMMENT ON COLUMN application_configurations.config_value IS 'The value of the configuration parameter, stored as text.';

COMMENT ON COLUMN application_configurations.value_type IS 'The intended data type of config_value, used by the application for parsing (e.g., INTEGER, BOOLEAN, JSON).';

COMMENT ON COLUMN application_configurations.description IS 'A description of the configuration parameter and its purpose.';

COMMENT ON COLUMN application_configurations.is_active IS 'Flag to indicate if the configuration parameter is currently active and should be used.';

COMMENT ON COLUMN application_configurations.created_at IS 'Timestamp of when the configuration parameter was initially created.';

COMMENT ON COLUMN application_configurations.updated_at IS 'Timestamp of when the configuration parameter was last updated.';

CREATE TRIGGER set_timestamp_application_configurations BEFORE
UPDATE ON application_configurations FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp ();

ALTER TABLE application_configurations OWNER TO palebluedot4;