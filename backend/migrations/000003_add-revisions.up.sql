ALTER TABLE message ADD COLUMN revisions JSONB NOT NULL;
ALTER TABLE message DROP COLUMN content;
ALTER TABLE message DROP COLUMN title;

CREATE OR REPLACE FUNCTION get_effective_username(author_id INT, revisions JSONB)
RETURNS TEXT AS $$
DECLARE
    author_username_override TEXT;
BEGIN
    author_username_override := revisions -> -1 ->> 'author_username_override';

    RETURN COALESCE(author_username_override, (SELECT username FROM author WHERE author.id = author_id));
END;
$$ LANGUAGE plpgsql;

INSERT INTO author (username) VALUES
	('John'),
	('Jane'),
	('Paul');

INSERT INTO message (parent_id, author_id, revisions) VALUES
	(NULL, 1, '[{"title": "First!", "content": "Hello World!"}]'::jsonb),
	(NULL, 2, '[{"title": "Second!", "content": "Hello John!"}]'::jsonb),
	(1, 2, '[{"title": "", "content": "<b>Accidentally made</b> a completely new post :)"}]'::jsonb),
	(3, 1, '[{"title": "", "content": "I get it. The UX of this board is kind of terrible."}]'::jsonb),
	(1, 3, '[{"title": "", "content": "Hi there, John!"}]'::jsonb);
