INSERT INTO
    application_configurations (
        config_key,
        config_value,
        value_type,
        description,
        is_active
    )
VALUES
    (
        'app_version',
        '1.0.0',
        'TEXT',
        'Current version of the application.',
        TRUE
    ),
    (
        'app_name',
        'palebluedot4',
        'TEXT',
        'The display name for the application.',
        TRUE
    );