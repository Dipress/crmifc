CREATE TABLE IF NOT EXISTS categories (
	id	SERIAL PRIMARY KEY,
	name	VARCHAR (50) UNIQUE NOT NULL,

	/* timestamp */
	created_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
