-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN image_url VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN image_url;
-- +goose StatementEnd
