    INSERT INTO apps (id, name, secret)
    VALUES (1, 'test-client', 'test-secret')
    ON CONFLICT DO NOTHING;