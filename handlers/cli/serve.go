package cli

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"buf.build/gen/go/git-age/agent/connectrpc/go/agent/v1/agentv1connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/99designs/keyring"
	"github.com/alecthomas/kong"
	"github.com/coreos/go-systemd/v22/activation"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/prskr/git-age-keyring-agent/handlers/grpc"
	"github.com/prskr/git-age-keyring-agent/infrastructure/httpx"
)

type ServeCliHandler struct {
	ServiceNameFlag `embed:""`
	Http            struct {
		ListenAddress     *url.URL      `default:"unix://${XDG_RUNTIME_DIR}/git-age-keyring-agent.sock"`
		ReadHeaderTimeout time.Duration `default:"5s" help:"Read header timeout"`
		ShutdownTimeout   time.Duration `default:"5s" help:"Shutdown timeout"`
	} `embed:"" prefix:"http."`
}

func (h ServeCliHandler) Run(ctx context.Context, kr keyring.Keyring) (err error) {
	listener, err := h.prepareListener()
	if err != nil {
		return err
	}

	reflector := grpcreflect.NewStaticReflector(grpchealth.HealthV1ServiceName, agentv1connect.IdentitiesStoreServiceName)
	checker := grpchealth.NewStaticChecker(agentv1connect.IdentitiesStoreServiceName)

	agentServer := grpc.NewAgentServer(kr)

	mux := http.NewServeMux()

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(grpchealth.NewHandler(checker))
	mux.Handle(agentv1connect.NewIdentitiesStoreServiceHandler(agentServer))

	fmt.Printf(`export GIT_AGE_AGENT_HOST="%s"`, h.Http.ListenAddress.String())

	srv := http.Server{
		Handler:           h2c.NewHandler(httpx.LoggingMiddleware(mux), new(http2.Server)),
		ReadHeaderTimeout: h.Http.ReadHeaderTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if serveErr := srv.Serve(listener); serveErr != nil {
			if !errors.Is(serveErr, http.ErrServerClosed) {
				err = errors.Join(err, serveErr)
			}
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), h.Http.ShutdownTimeout)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

func (h ServeCliHandler) AfterApply(kongCtx *kong.Context) error {
	keyRingCfg := keyring.Config{
		ServiceName: h.ServiceName,
	}

	kr, err := keyring.Open(keyRingCfg)
	if err != nil {
		return err
	}

	kongCtx.BindTo(kr, (*keyring.Keyring)(nil))

	return nil
}

func (h ServeCliHandler) prepareListener() (listener net.Listener, err error) {
	listeners, err := activation.Listeners()
	if err == nil && len(listeners) > 0 {
		listener = listeners[0]
		return listener, nil
	}

	if h.Http.ListenAddress.Scheme == "unix" {
		if _, err := os.Stat(h.Http.ListenAddress.Path); err == nil {
			if err := os.Remove(h.Http.ListenAddress.Path); err != nil {
				return nil, err
			}
		}

		return net.Listen(h.Http.ListenAddress.Scheme, h.Http.ListenAddress.Path)
	} else {
		return net.Listen(h.Http.ListenAddress.Scheme, h.Http.ListenAddress.Host)
	}
}
