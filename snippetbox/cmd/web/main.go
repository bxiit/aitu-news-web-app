package main

import (
	"alexedwards.net/snippetbox/pkg/models/postgresql"
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// Add a templateCache field to the application struct.
type application struct {
	errorLog      *log.Logger       // Понятное дело логгеры
	infoLog       *log.Logger       // Понятное дело логгеры
	session       *sessions.Session // TODO понять для чего
	news          *postgresql.NewsModel
	templateCache map[string]*template.Template // Оптимизация засчет избегания перекомпилирования
	users         *postgresql.UserModel
	comment       *postgresql.CommentModel
}

func main() {
	data := "user=bxit password=aa dbname=postgres sslmode=disable host=localhost port=5433"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Секретный ключ для шифрования и аутентификации сессионных cookie.
	// Значение по умолчанию "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge".
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// получение объекта бд
	db, err := openDB(data)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Кэш шаблон для html страниц
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Экземпляр application-а с зависимостями
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		news:          &postgresql.NewsModel{DB: db},
		templateCache: templateCache,
		users:         &postgresql.UserModel{DB: db},
		comment:       &postgresql.CommentModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want
	// the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Set the server's TLSConfig field to use the tlsConfig variable we just created
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)

	// Use the ListenAndServeTLS() method to start the HTTPS server. We
	// pass in the paths to the TLS certificate and corresponding private key as
	// the two parameters.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
