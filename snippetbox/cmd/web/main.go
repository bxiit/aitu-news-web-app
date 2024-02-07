package main

import (
	"alexedwards.net/snippetbox/pkg/models/postgresql"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Add a templateCache field to the application struct.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	//session       *sessions.Session
	news          *postgresql.NewsModel
	templateCache map[string]*template.Template
}

func main() {
	data := "user=bxit password=aa dbname=postgres sslmode=disable host=localhost port=5433"
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32
	// will be used to encrypt and authenticate session cookies). It should be 32
	// bytes long.
	//secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(data)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	//session := sessions.New([]byte(*secret))
	//session.Lifetime = 12 * time.Hour

	// And add it to the application dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		//session:       session,
		news:          &postgresql.NewsModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
