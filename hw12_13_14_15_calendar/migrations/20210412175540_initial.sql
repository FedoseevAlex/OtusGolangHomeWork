-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id 			uuid PRIMARY KEY NOT NULL,
	title		text NOT NULL,
	start_time  timestamp NOT NULL,
	end_time	timestamp NOT NULL,
	description	text,
	owner_id 	uuid NOT NULL,
	notify_at	timestamp
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd