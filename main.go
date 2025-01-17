package main

import (
	"database/sql"

	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var store = sessions.NewCookieStore([]byte("Fuck-This-Shit"))
type Article struct {
    ID          int
    Topic       string
    Description string
    CreatedAt   string
    Image       string
	
}


type ArticleUser struct {
    ID          int
    Topic       string
    Description string
    CreatedAt   string
    Image       string
	Username    string
	Tags        []string
}
var user_id=0
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



func helloHandleFunc (w http.ResponseWriter, r *http.Request) {
    fmt.Fprint (w, "Hello, World!")
}
func registerHandleFunc (w http.ResponseWriter, r *http.Request) {
    fmt.Fprint (w, "register")
}
func routes(){
	tpl, _ = template.ParseFiles("templates/news-aggregator.html","templates/login.html","templates/register.html","templates/social-newsfeed-v1.html","templates/news-aggregatortag.html","templates/Article.html")
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

    http.HandleFunc("/writearticle", writearticleHandler)
	http.HandleFunc("/article", articleHandler)
    http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
    fmt.Println("Server is running on :3000...")
    http.ListenAndServe(":3000", nil)
}


var tpl *template.Template

func main() {
    database()
	routes()
	

    
}



func writearticleHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the query parameters
    queryValues := r.URL.Query()

    // Create variables based on the query parameters
    title := queryValues.Get("title")
    description := queryValues.Get("desc")
    tags := queryValues["tags"]
    image := queryValues.Get("image")

    // Print the extracted values (for debugging)
    fmt.Printf("Title: %s\n", title)
    fmt.Printf("Description: %s\n", description)
    fmt.Printf("Tags: %v\n", tags)
    fmt.Printf("Image: %s\n", image)


	session, _ := store.Get(r, "User") // Replace "session-name" with your session name

    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		print("Unauthorized")
		print(session.Values["authenticated"])

		
        tpl.ExecuteTemplate(w, "login.html" ,nil)
        return
    }else{
        if title != "" && description != ""  &&  image != ""{
                        db, err := sql.Open("sqlite3", "./database2.db")
                        if err != nil {
                            http.Error(w, "Database connection error", http.StatusInternalServerError)
                            return
                        }
                        defer db.Close()
                    
                        // Your code to save the article data to the database goes here
                        
                
                        query := "SELECT id FROM Users WHERE email = ?"

                        // Execute the query with the provided email
                        err = db.QueryRow(query, session.Values["email"].(string)).Scan(&user_id)
                        // Insert the article information into the Articles table
                        _, err = db.Exec("INSERT INTO Articles (user_id, topic, description, image, create_at) VALUES (?, ?, ?, ?, datetime('now'))", user_id, title, description, image)
                        if err != nil {
                            http.Error(w, "Error saving the article", http.StatusInternalServerError)
                            return
                        }

                        // Retrieve the last inserted row ID
                        var articleID int64
                        err = db.QueryRow("SELECT last_insert_rowid()").Scan(&articleID)
                        if err != nil {
                            http.Error(w, "Error retrieving article ID", http.StatusInternalServerError)
                            return
                        }

                        // Insert the tag relationships into the TagArticle table
                        for _, tagID := range tags {
                            _, err = db.Exec("INSERT INTO TagArticle (tag_id, article_id) VALUES (?, ?)", tagID, articleID)
                            if err != nil {
                                http.Error(w, "Error saving tag relationships", http.StatusInternalServerError)
                                return
                            }
                        }
       
                        // tpl.ExecuteTemplate(w,"" ,nil)
                        http.Redirect(w, r, "/success", http.StatusSeeOther)
                        
       
    }
    
    tpl.ExecuteTemplate(w, "Article.html" ,nil)

}



    
  
}





func loginHandler(w http.ResponseWriter, r *http.Request) {
    // func (t *Template) Execute Template (wr io.Writer, name string, data interface{}) error
	queryValues := r.URL.Query()

    // Get the values of 'email' and 'password'
    email := queryValues.Get("email")
    password := queryValues.Get("password")

    // Now, you have the values of 'email' and 'password' from the URL.
    // You can perform any necessary processing with this data.

    // Example: Print the received data
    fmt.Printf("Email: %s\n", email)
    fmt.Printf("Password: %s\n", password)

	if email != "" && password != "" {
		db, err := sql.Open("sqlite3", "./database2.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()
	
		
		var storedPassword string
		
		err = db.QueryRow("SELECT password FROM Users WHERE email = ?", email).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "User not found or database error", http.StatusUnauthorized)
			return
		}
	
		// Check if the stored password matches the provided password
		if storedPassword == password {
			// Passwords match; the user is authenticated
			// You can respond with a success message or redirect to a protected page.
			print("Authentication successful")
			// If authentication is successful, save user's session
			session, _ := store.Get(r, "User") // Replace "session-name" with your session name
			session.Values["email"] = email
			session.Values["authenticated"] = true
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, "Session save error", http.StatusInternalServerError)
				return
			}
		
			// Redirect the user to an authenticated page or provide access

			http.Redirect(w, r, "/success", http.StatusSeeOther)
		} else {
			// Passwords do not match; authentication failed
			print("Authentication failed")
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
		}
	}
	
    // Perform  authentication logic here


    // Redirect or respond with a success message
   
    tpl.ExecuteTemplate(w, "login.html" ,nil)
	
	
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()

    // Get the values of 'email' and 'password'
    email := queryValues.Get("email")
    password := queryValues.Get("password")
    repass := queryValues.Get("Repeatpassword")

    // Now, you have the values of 'email' and 'password' from the URL.
    // You can perform any necessary processing with this data.

    // Example: Print the received data
    fmt.Printf("Email: %s\n", email)
    fmt.Printf("Password: %s\n", password)
    fmt.Printf("Password: %s\n", repass)

	if email != "" && password != "" {
        // Open a database connection
        db, err := sql.Open("sqlite3", "./database2.db")
        if err != nil {
            http.Error(w, "Database connection error", http.StatusInternalServerError)
            return
        }
        defer db.Close()
    
        // Ensure that the password and repeated password match
        if password != repass {
            http.Error(w, "Passwords do not match", http.StatusBadRequest)
            return
        }
    
        // Insert the user's email and password into the Users table without creating new variables
        _, err = db.Exec("INSERT INTO Users (email, password) VALUES (?, ?)", email, password)
        if err != nil {
            http.Error(w, "Error creating the user", http.StatusInternalServerError)
            return
        }
    
        // User created successfully
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    


	
    // func (t *Template) Execute Template (wr io.Writer, name string, data interface{}) error
	session, _ := store.Get(r, "User") // Replace "session-name" with your session name

    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		print("Unauthorized")
		print(session.Values["authenticated"])

		tpl.ExecuteTemplate(w, "register.html" ,nil)
        return
    }


	print("session email",session.Values["email"].(string))
	http.Redirect(w, r, "/success", http.StatusSeeOther)
  
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    id := queryValues.Get("id")
    fmt.Println("id", id) // Print the ID for debugging

    db, err := sql.Open("sqlite3", "./database2.db")
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // First query to fetch article and user information
    query := `
    SELECT Articles.id, Articles.topic, Articles.description, Articles.create_at, Articles.image, Users.name
    FROM Articles
    INNER JOIN Users ON Articles.user_id = Users.id
    WHERE Articles.id = ?
    `
    fmt.Println("Query:", query) // Print the SQL query for debugging

    row := db.QueryRow(query, id)

    var article ArticleUser
    err = row.Scan(&article.ID, &article.Topic, &article.Description, &article.CreatedAt, &article.Image, &article.Username)

    if err != nil {
        http.Error(w, "Article not found", http.StatusNotFound)
        return
    }

    // Second query to fetch tags associated with the article
    tagsQuery := `
    SELECT Tags.name
    FROM Tags
    INNER JOIN TagArticle ON Tags.id = TagArticle.tag_id
    WHERE TagArticle.article_id = ?
    `
    fmt.Println("Tags Query:", tagsQuery) // Print the SQL query for debugging

    rows, err := db.Query(tagsQuery, id)
    if err != nil {
        http.Error(w, "Error fetching tags", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var tags []string
    for rows.Next() {
        var tag string
        err := rows.Scan(&tag)
        if err != nil {
            http.Error(w, "Error scanning tags", http.StatusInternalServerError)
            return
        }
        tags = append(tags, tag)
    }

    fmt.Println("Tags:", tags) // Print the tags for debugging

    // Combine article and tags data
    article.Tags = tags

    err = tpl.ExecuteTemplate(w, "social-newsfeed-v1.html", article)
    if err != nil {
        http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
    }
}





func logoutHandler(w http.ResponseWriter, r *http.Request) {
    // Get the session and delete it
    session, _ := store.Get(r, "User") // Replace "User" with your session name
    session.Options.MaxAge = -1 // Set the session cookie to expire immediately
    err := session.Save(r, w)
    if err != nil {
        http.Error(w, "Session delete error", http.StatusInternalServerError)
        return
    }

    // Redirect the user to a logout success page or any other desired page
    http.Redirect(w, r, "/logout-success", http.StatusSeeOther)
}



func indexHandler(w http.ResponseWriter, r *http.Request) {
   
	queryValues := r.URL.Query()
    // Get the values of 'email' and 'password'
    tagID := queryValues.Get("tag_id")
    db, err := sql.Open("sqlite3", "./database2.db")
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var query string
    var args []interface{}
    if tagID != "" {
        // If tag ID is provided, filter articles by tag ID and get the latest data
        query = `
            SELECT a.id, a.topic, a.description, a.create_at, a.image
            FROM Articles a
            INNER JOIN TagArticle ta ON a.id = ta.article_id
            WHERE ta.tag_id = ?
            ORDER BY a.create_at DESC`
        args = []interface{}{tagID}
    } else {
        // If no tag ID is provided, fetch all articles
        query = "SELECT id, topic, description, create_at, image FROM Articles ORDER BY create_at DESC"
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var articles []Article
    for rows.Next() {
        var a Article
        err := rows.Scan(&a.ID, &a.Topic, &a.Description, &a.CreatedAt, &a.Image)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        articles = append(articles, a)
    }
	if tagID != ""{

		err = tpl.ExecuteTemplate(w, "news-aggregatortag.html", articles)

	}else{
        err = tpl.ExecuteTemplate(w, "news-aggregator.html", articles)

    }
    // Execute the HTML template with articles as data
    
}




