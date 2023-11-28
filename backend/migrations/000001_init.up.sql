CREATE TABLE IF NOT EXISTS author (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS message (
		id BIGSERIAL PRIMARY KEY,
		parent_id INTEGER REFERENCES "message" (id),
		author_id INTEGER REFERENCES "author" (id),
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO author (username) VALUES
	('John'),
	('Jane'),
	('Paul');

INSERT INTO message (parent_id, author_id, title, content) VALUES
	(NULL, 1, 'First!', 'Hello World!'),
	(NULL, 2, 'Second!', 'Hello John!'),
	(1, 2, '', 'Hello John! <b>Accidentally made</b> a completely new post :)'),
	(3, 1, '', 'I get it. The UX of this board is kind of terrible.'),
	(1, 3, '', 'Hi there, John!');
