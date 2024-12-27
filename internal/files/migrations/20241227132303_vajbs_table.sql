-- +goose Up
-- +goose StatementBegin
CREATE TYPE region AS ENUM ('pardubicky', 'liberecky', 'kralovehradecky', 'moravskoslezsky', 'olomoucky', 'praha', 'stredocesky', 'ustecky', 'vysocina', 'zlinsky', 'jihocesky', 'jihomoravsky', 'karlovarsky', 'plzensky');

CREATE TABLE vajbs (
  id SERIAL PRIMARY KEY,
  creator_id INT REFERENCES users(id),
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  date TIMESTAMP NOT NULL,
  region region NOT NULL,
  address VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vajbs;
DROP TYPE region;
-- +goose StatementEnd
