TRUNCATE TABLE auth_transactions;

ALTER TABLE auth_transactions
    ADD COLUMN browser_binding_hash bytea NOT NULL;

ALTER TABLE auth_transactions
    ADD CONSTRAINT auth_transactions_browser_binding_hash_length
        CHECK (octet_length(browser_binding_hash) = 32);
