-- +goose Up
-- +goose StatementBegin
ALTER TABLE vajbs ADD COLUMN header_image VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vajbs DROP COLUMN header_image;
-- +goose StatementEnd
