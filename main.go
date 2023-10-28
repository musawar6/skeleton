package main

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

// Embedd all assets and templates into the final binary.
//
//go:embed assets
var assetFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

// Template is the template renderer echo needs
type Template struct {
	templates *template.Template
}

// Render is the real render which has to implement the renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Migrate prepares the database and creates all tables and data which are required for the app to run
func Migrate(db *sqlx.DB) {
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS "messages" (
		"id"	INTEGER,
		"message"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`)
}

// Message is the message data model
type Message struct {
	ID      int    `db:"id"`
	Message string `db:"message"`
}

// CreateMessage inserts a new message into the messages table
func CreateMessage(db *sqlx.DB, message string) error {
	_, err := db.Exec("INSERT INTO messages(message) VALUES(?)", message)
	return err
}

// LoadMessages loads all messages from the database
func LoadMessages(db *sqlx.DB) ([]Message, error) {
	var messages []Message
	err := db.Select(&messages, "SELECT id, message FROM messages")
	return messages, err
}

// indexHandler renders the index page
func indexHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	}
}

// messagesHandler renders the messages
func messagesHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		messages, err := LoadMessages(db)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.Render(http.StatusOK, "messages.html", messages)
	}
}

// createMessageHandler handles the create message request and renders the messages as result
func createMessageHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		message := c.FormValue("message")
		err := CreateMessage(db, message)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		messages, err := LoadMessages(db)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.Render(http.StatusOK, "messages.html", messages)
	}
}

func main() {
	// init database
	db := sqlx.MustConnect("sqlite3", "skeleton.db")
	Migrate(db)

	// init echo
	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseFS(templateFiles, "templates/*.html")),
	}

	// register routes
	e.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(assetFiles))))
	e.GET("/", indexHandler(db))
	e.GET("/messages", messagesHandler(db))
	e.POST("/messages", createMessageHandler(db))

	// run webserver
	e.Start("127.0.0.1:3000")
}
