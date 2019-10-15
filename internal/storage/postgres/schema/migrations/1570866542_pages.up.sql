CREATE TABLE IF NOT EXISTS pages (
	id	SERIAL PRIMARY KEY,
	user_id	INT NOT NULL,
	title	VARCHAR (255) NOT NULL,
	body	TEXT NOT NULL,
	category_id INT NOT NULL,

	/* timestamp */
	created_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

