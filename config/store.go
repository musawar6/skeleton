package config

import (
	"database/sql"
	"log"
)

func database() error {
	database, err := sql.Open("sqlite3", "./database2.db")
	if err != nil {
		return err
	}
	defer database.Close()

	// Create the 'Users' table
	usersTableSQL := `
	CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		password TEXT,
		email TEXT,
		description TEXT,
		create_date DATETIME
	);
	`

	_, err = database.Exec(usersTableSQL)
	if err != nil {
		return err
	}

	tagsTableSQL := `
	CREATE TABLE IF NOT EXISTS Tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);
	`

	_, err = database.Exec(tagsTableSQL)
	if err != nil {
		return err
	}
	// Create the 'Articles' table
	articlesTableSQL := `
	CREATE TABLE IF NOT EXISTS Articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		topic TEXT,
		description TEXT,
		image BLOB,
		create_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES Users (id)
	);
	`

	_, err = database.Exec(articlesTableSQL)
	if err != nil {
		return err
	}

	tagArticleTableSQL := `
	CREATE TABLE IF NOT EXISTS TagArticle (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    tag_id INTEGER,
	    article_id INTEGER,
	    FOREIGN KEY (tag_id) REFERENCES Tags (id),
	    FOREIGN KEY (article_id) REFERENCES Articles (id)
	);
	`

	_, err = database.Exec(tagArticleTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	likesTableSQL := `
		CREATE TABLE IF NOT EXISTS Likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			article_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES Users (id),
			FOREIGN KEY (article_id) REFERENCES Articles (id)
		);
	`

	_, err = database.Exec(likesTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	print("Database created")
	return nil
}