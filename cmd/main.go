package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/citadel-corp/segokuning-social-app/internal/common/db"
	"github.com/citadel-corp/segokuning-social-app/internal/common/middleware"
	"github.com/citadel-corp/segokuning-social-app/internal/image"
	"github.com/citadel-corp/segokuning-social-app/internal/posts"
	"github.com/citadel-corp/segokuning-social-app/internal/user"
	"github.com/gorilla/mux"
)

func main() {
	slogHandler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(slogHandler))

	// Connect to database
	env := os.Getenv("ENV")
	sslMode := "disable"
	if env == "production" {
		sslMode = "verify-full sslrootcert=ap-southeast-1-bundle.pem"
	}
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), sslMode)
	db, err := db.Connect(dbURL)
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot connect to database: %v", err))
		os.Exit(1)
	}

	// Create migrations
	// err = db.UpMigration()
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Up migration failed: %v", err))
	// 	os.Exit(1)
	// }

	// initialize user domain
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// initialize image domain
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("S3_ID"), os.Getenv("S3_SECRET_KEY"), ""),
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot create AWS session: %v", err))
		os.Exit(1)
	}
	imageService := image.NewService(sess)
	imageHandler := image.NewHandler(imageService)

	// initialize posts domain
	postsRepository := posts.NewRepository(db)
	postsService := posts.NewService(postsRepository)
	postsHandler := posts.NewHandler(postsService)

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		io.WriteString(w, "Service ready")
	})

	v1 := r.PathPrefix("/v1").Subrouter()

	// user routes
	ur := v1.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/register", middleware.PanicRecoverer(userHandler.CreateUser)).Methods(http.MethodPost)
	ur.HandleFunc("/login", middleware.PanicRecoverer(userHandler.Login)).Methods(http.MethodPost)
	ur.HandleFunc("/link/email", middleware.PanicRecoverer(middleware.Authorized(userHandler.LinkEmail))).Methods(http.MethodPost)
	ur.HandleFunc("/link/phone", middleware.PanicRecoverer(middleware.Authorized(userHandler.LinkPhoneNumber))).Methods(http.MethodPost)
	ur.HandleFunc("", middleware.PanicRecoverer(middleware.Authorized(userHandler.Update))).Methods(http.MethodPatch)

	// image routes
	ir := v1.PathPrefix("/image").Subrouter()
	ir.HandleFunc("", middleware.PanicRecoverer(middleware.Authorized(imageHandler.UploadToS3))).Methods(http.MethodPost)

	// posts routes
	pr := v1.PathPrefix("/post").Subrouter()
	pr.HandleFunc("", middleware.PanicRecoverer(middleware.Authorized(postsHandler.CreatePost))).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:     ":8000",
		Handler:  r,
		ErrorLog: slog.NewLogLogger(slogHandler, slog.LevelError),
	}

	go func() {
		slog.Info(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("HTTP server error: %v", err))
		}
		slog.Info("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	slog.Info(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error: %v", err)
	}
	slog.Info("Shutdown complete.")
}
