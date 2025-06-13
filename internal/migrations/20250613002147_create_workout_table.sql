-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "workout" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    drills JSON,
    status VARCHAR(255) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "workout";
-- +goose StatementEnd