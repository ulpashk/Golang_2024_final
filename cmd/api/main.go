package main

import (
	"context" 
	"database/sql" 
	"flag"
	"fmt"
	"log"
	"net/http"
	// "strings"
	"os"
	"time"
	"goproject/internal/data"
	// "goproject/internal/jsonlog"


	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)


const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
type config struct {
	port int
	env string
	db struct {
		dsn string
		// maxOpenConns int
		// maxIdleConns int
		// maxIdleTime string
	}
	// limiter struct {
	// 	enabled bool
	// 	rps float64
	// 	burst int
	// }
	// smtp struct {
	// 	host string
	// 	port int
	// 	username string
	// 	password string
	// 	sender string
	// }
	// cors struct {
	// 	trustedOrigins []string
	// }
	
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *log.Logger
	models data.Models
}
	

func main() {
	// Declare an instance of the config struct.
	var cfg config
	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	
	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	// flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	// flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	// flag.StringVar(&cfg.db.dsn, "db-dsn",  os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "user=postgres password='Ulp@sh05' dbname=kkpop sslmode=disable", "PostgreSQL DSN")

	// flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	// flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	// flag.StringVar(&cfg.smtp.username, "smtp-username", "0f1d85c09e6d8e", "SMTP username")
	// flag.StringVar(&cfg.smtp.password, "smtp-password", "e89654b1c53c45", "SMTP password")
	// flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")


	// flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
	// 	cfg.cors.trustedOrigins = strings.Fields(val)
	// 	return nil
	// })

	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)


	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Printf("database connection pool established")

	app := &application {
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}


	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr: 			fmt.Sprintf(":%d", cfg.port),
		Handler: 		app.routes(),
		IdleTimeout:	time.Minute,
		ReadTimeout: 	10 * time.Second,
		WriteTimeout: 	30 * time.Second,
	}
	
	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	// Because the err variable is now already declared in the code above, we need
	// to use the = operator here, instead of the := operator.
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
	