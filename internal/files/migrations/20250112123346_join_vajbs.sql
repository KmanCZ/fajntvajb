-- +goose Up
-- +goose StatementBegin
CREATE TABLE joined_vajbs (
  user_id INT REFERENCES users(id),
  vajb_id INT REFERENCES vajbs(id),
  PRIMARY KEY (user_id, vajb_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE joined_vajbs;
-- +goose StatementEnd
