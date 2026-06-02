CREATE TABLE IF NOT EXISTS subscriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_name VARCHAR(128) NOT NULL,
    owner_name VARCHAR(128) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


