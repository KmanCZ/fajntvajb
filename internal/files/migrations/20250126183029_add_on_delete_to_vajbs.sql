-- +goose Up
-- +goose StatementBegin
ALTER TABLE vajbs
DROP CONSTRAINT IF EXISTS vajbs_creator_id_fkey,
ADD CONSTRAINT vajbs_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vajbs
DROP CONSTRAINT IF EXISTS vajbs_creator_id_fkey,
ADD CONSTRAINT vajbs_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES users (id);
-- +goose StatementEnd
