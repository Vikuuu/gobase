-- Up Migration

CREATE TABLE users (
	id INTEGER,
	name TEXT,
	created_at DATETIME,
	updated_at DATETIME,
	is_member BOOLEAN
);

-- Down Migration

DROP TABLE users;