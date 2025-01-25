-- +goose Up
-- +goose StatementBegin
ALTER TABLE joined_vajbs
DROP CONSTRAINT IF EXISTS joined_vajbs_user_id_fkey,
DROP CONSTRAINT IF EXISTS joined_vajbs_vajb_id_fkey,
ADD CONSTRAINT joined_vajbs_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
ADD CONSTRAINT joined_vajbs_vajb_id_fkey FOREIGN KEY (vajb_id) REFERENCES vajbs (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE joined_vajbs
DROP CONSTRAINT IF EXISTS joined_vajbs_user_id_fkey,
DROP CONSTRAINT IF EXISTS joined_vajbs_vajb_id_fkey,
ADD CONSTRAINT joined_vajbs_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id),
ADD CONSTRAINT joined_vajbs_vajb_id_fkey FOREIGN KEY (vajb_id) REFERENCES vajbs (id);
-- +goose StatementEnd
