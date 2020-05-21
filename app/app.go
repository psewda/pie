package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/psewda/pie/session"
)

type App interface {
	Run(port uint16) error
	Dispose()
}

type httpApp struct {
	server  *http.Server
	running bool
}

func (app *httpApp) Run(port uint16) error {
	if !app.running {
		if port < 1024 {
			err := errors.New("Port number should be in private/dynamic port range.")
			return err
		}
		app.server.Addr = fmt.Sprintf(":%d", port)

		go func() {
			log.Printf("server starting @ port '%d'", port)
			if err := app.server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("server not started, error occured => [%s]", err)
			}
		}()
		app.running = true
	}
	return nil
}

func (app *httpApp) Dispose() {
	if app.running {
		app.server.Shutdown(context.Background())
		app.running = false
	}
}

func NewApp() App {
	store := session.NewSessionStore()
	r := NewRouter(store)
	return &httpApp{
		running: false,
		server: &http.Server{
			Handler: r.Handler,
		},
	}
}

func GetRandPort() uint16 {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return uint16(rnd.Intn(math.MaxUint16-1024) + 1024)
}
