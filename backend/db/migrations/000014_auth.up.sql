CREATE TABLE user_identities (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    provider text NOT NULL,
    provider_subject text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_login_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT user_identities_user_fk
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT user_identities_provider_check
        CHECK (length(provider) > 0),
    CONSTRAINT user_identities_subject_check
        CHECK (length(provider_subject) > 0),
    CONSTRAINT user_identities_timestamps_check
        CHECK (updated_at >= created_at AND last_login_at >= created_at),
    CONSTRAINT user_identities_provider_subject_unique
        UNIQUE (provider, provider_subject)
);

CREATE TABLE auth_sessions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz NULL,
    user_agent text NULL,
    ip_address inet NULL,

    CONSTRAINT auth_sessions_user_fk
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT auth_sessions_token_hash_unique
        UNIQUE (token_hash),
    CONSTRAINT auth_sessions_token_hash_length
        CHECK (octet_length(token_hash) = 32),
    CONSTRAINT auth_sessions_timestamps_check
        CHECK (last_seen_at >= created_at AND expires_at > created_at AND (revoked_at IS NULL OR revoked_at >= created_at))
);

CREATE INDEX auth_sessions_user_id_idx ON auth_sessions(user_id);
CREATE INDEX auth_sessions_expires_at_idx ON auth_sessions(expires_at);

CREATE TABLE auth_transactions (
    id uuid PRIMARY KEY,
    provider text NOT NULL,
    state_hash bytea NOT NULL,
    nonce_hash bytea NOT NULL,
    code_verifier text NOT NULL,
    redirect_uri text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz NULL,

    CONSTRAINT auth_transactions_provider_check
        CHECK (length(provider) > 0),
    CONSTRAINT auth_transactions_state_hash_length
        CHECK (octet_length(state_hash) = 32),
    CONSTRAINT auth_transactions_nonce_hash_length
        CHECK (octet_length(nonce_hash) = 32),
    CONSTRAINT auth_transactions_code_verifier_check
        CHECK (length(code_verifier) > 0),
    CONSTRAINT auth_transactions_redirect_uri_check
        CHECK (length(redirect_uri) > 0),
    CONSTRAINT auth_transactions_timestamps_check
        CHECK (expires_at > created_at AND (consumed_at IS NULL OR consumed_at >= created_at)),
    CONSTRAINT auth_transactions_state_hash_unique
        UNIQUE (state_hash)
);

CREATE INDEX auth_transactions_expires_at_idx ON auth_transactions(expires_at);
