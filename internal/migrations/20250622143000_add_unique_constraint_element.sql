-- +goose Up
-- +goose StatementBegin
ALTER TABLE "element" ADD CONSTRAINT unique_element_user_name UNIQUE (user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "element" DROP CONSTRAINT IF EXISTS unique_element_user_name;
-- +goose StatementEnd