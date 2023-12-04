package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Note: Could be made an interface to allow for mocking in tests
// if one likes that kind of thing.
type Database struct {
	Conn *sqlx.DB
}

const (
	DB_USER_ENV     = "DB_USER"
	DB_PASSWORD_ENV = "DB_PASSWORD"
	DB_HOST_ENV     = "DB_HOST"
	DB_PORT_ENV     = "DB_PORT"
	DB_NAME_ENV     = "DB_NAME"

	REDACTED_USERNAME = "<REDACTED>"
	REDACTED_TITLE    = "<DELETED>"
	REDACTED_CONTENT  = "<DELETED>"
)

func ConnectDb() Database {
	dbUser := GetEnv(DB_USER_ENV)
	dbPassword := GetEnv(DB_PASSWORD_ENV)
	dbHost := GetEnv(DB_HOST_ENV)
	dbPort := GetEnv(DB_PORT_ENV)
	dbName := GetEnv(DB_NAME_ENV)

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	log.Printf("Connecting to database with connection %s", connectionString)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	return Database{Conn: db}
}

func (db *Database) getMessages(messageId *int64) ([]Message, error) {
	messages := make([]Message, 0)

	err := db.Conn.Select(&messages,
		`SELECT
			m.id,
			m.parent_id,
			a.username,
			m.revisions -> -1 ->> 'title' as title,
			m.revisions -> -1 ->> 'content' as content,
			m.created_at
		FROM message m
		JOIN author a ON m.author_id = a.id
		WHERE m.parent_id IS NOT DISTINCT FROM $1
		ORDER BY m.id ASC
		`, messageId)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (db *Database) createMessage(message InputMessage) (Message, error) {
	var createdMessage Message
	err := db.Conn.QueryRowx(
		`
		INSERT INTO message (parent_id, author_id, revisions)
		VALUES (
			$1,
			(SELECT id FROM author WHERE username = $2),
			jsonb_build_array(jsonb_build_object('title', $3::text, 'content', $4::text))
		)
		RETURNING
			id,
			parent_id,
			get_effective_username(author_id, revisions) as username,
			revisions -> -1 ->> 'title' as title,
			revisions -> -1 ->> 'content' as content,
			created_at`,
		message.ParentId, message.Author, message.Title, message.Content).StructScan(&createdMessage)

	if err != nil {
		return Message{}, err
	}

	return createdMessage, nil
}

func (db *Database) deleteMessage(messageId int64) (Message, error) {
	return db.editMessage(
		messageId,
		REDACTED_TITLE,
		REDACTED_CONTENT,
		sql.NullString{
			String: REDACTED_USERNAME,
			Valid:  true,
		})
}

func (db *Database) editMessage(
	messageId int64,
	title string,
	content string,
	authorUsernameOverride sql.NullString,
) (Message, error) {
	var editedMessage Message

	err := db.Conn.QueryRowx(
		`UPDATE message
			SET
				revisions = revisions || jsonb_build_array(jsonb_build_object(
					'title', $1::text, 'content', $2::text, 'author_username_override', $3::text))
			WHERE id = $4
			RETURNING
				id,
				parent_id,
				get_effective_username(author_id, revisions) as username,
				revisions -> -1 ->> 'title' as title,
				revisions -> -1 ->> 'content' as content,
				created_at`,
		title,
		content,
		authorUsernameOverride,
		messageId,
	).StructScan(&editedMessage)

	if err != nil {
		return Message{}, err
	}

	return editedMessage, nil
}
