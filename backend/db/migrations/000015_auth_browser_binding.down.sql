ALTER TABLE auth_transactions
    DROP CONSTRAINT IF EXISTS auth_transactions_browser_binding_hash_length;

ALTER TABLE auth_transactions
    DROP COLUMN IF EXISTS browser_binding_hash;
