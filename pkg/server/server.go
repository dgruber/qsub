package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"math/rand"

	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server"
	genserver "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server/generated"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const DefaultHost string = "127.0.0.1"
const DefaultPort int = 8088

type Config struct {
	Host     string
	Port     int
	Backend  string
	Password string
}

func Serve(c Config) error {

	if c.Host == "" {
		c.Host = DefaultHost
	}

	if c.Port == 0 {
		c.Port = DefaultPort
	}

	// backend process
	jobStore, err := simpletracker.NewPersistentJobStore("qsubjob.db")
	if err != nil {
		return fmt.Errorf("cannot create job store: %v", err)
	}

	processTracker, err := simpletracker.NewWithJobStore(
		"qsubsession", jobStore, true)
	if err != nil {
		return fmt.Errorf("cannot create process tracker: %v", err)
	}

	// connect the OpenAPI spec with the job tracker
	// interface implementation - could be anything
	impl, err := server.NewJobTrackerImpl(processTracker)
	if err != nil {
		return fmt.Errorf("cannot create job tracker implementation: %v", err)
	}

	// using chi router and logging + basic auth middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.BasicAuth(
		"qsub",
		map[string]string{
			"qsub": c.Password,
		}))

	fmt.Printf("Using password: %s", c.Password)

	m := http.NewServeMux()
	m.Handle("/qsub/",
		genserver.HandlerFromMuxWithBaseURL(
			impl, router, "/qsub"))

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", c.Host, c.Port),
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

// GetOrCreateSecret returns a secret which is stored in the user's home
// .qsub directory. If the secret file does not exist it is created.
func GetOrCreateSecret() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	qsubDir := filepath.Join(homeDir, ".qsub")
	if _, err := os.Stat(qsubDir); os.IsNotExist(err) {
		if err := os.Mkdir(qsubDir, 0700); err != nil {
			return "", err
		}
	}

	secretFile := filepath.Join(qsubDir, "secret")
	if _, err := os.Stat(secretFile); os.IsNotExist(err) {
		secret := randomString(32)
		if err := os.WriteFile(secretFile, []byte(secret), 0600); err != nil {
			return "", err
		}
		return secret, nil
	}

	secret, err := os.ReadFile(secretFile)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func DeleteSecret() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	qsubDir := filepath.Join(homeDir, ".qsub")
	if _, err := os.Stat(qsubDir); os.IsNotExist(err) {
		return nil
	}

	secretFile := filepath.Join(qsubDir, "secret")
	if _, err := os.Stat(secretFile); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(secretFile)
}

func randomString(size int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, size)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
