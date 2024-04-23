DO $$
DECLARE
    v_item_id VARCHAR(255);
    v_count INTEGER;
BEGIN
    SELECT item_id INTO v_item_id FROM accounts WHERE id = '%v';

    SELECT COUNT(*) INTO v_count FROM accounts WHERE item_id = v_item_id;

    IF v_count = 1 THEN
        DELETE FROM items WHERE id = v_item_id;
    ELSE
        DELETE FROM accounts WHERE id = '%v';
    END IF;
END $$;
