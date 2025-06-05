package caddy_plugin

import (
	"context"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

var isProd *bool

func IsProd() bool {
	if isProd != nil {
		return *isProd
	}
	value := !slices.Contains([]string{"dev", "development"}, os.Getenv("APP_ENV"))
	isProd = &value
	return *isProd
}

func getClientUrl() string {
	if IsProd() {
		return "caddy-valkey-1.caddy_default:6379"
	} else {
		return "127.0.0.1:6380"
	}
}

var client *valkey.Client

func getClient(logger *zap.Logger) valkey.Client {

	if client != nil {
		return *client
	}

	clientUrl := getClientUrl()

	client_, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{
			clientUrl,
		},
	})
	if err != nil {
		logger.Error("Failed to create Valkey client: " + clientUrl)
		panic(err)
	}
	client = &client_
	// defer client_.Close()
	return client_
}

func ServeHttp(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler, m Clipboard) error {

	client := getClient(m.logger)

	if r.Method == http.MethodGet && r.URL.Path == "/get" {
		getKeyValue(w, r, client, m.logger)
	} else if r.Method == http.MethodPost && r.URL.Path == "/set" {
		setKeyValue(w, r, client, m.logger)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("omit"))
	}

	return nil
}

var ctx context.Context = context.Background()

func getKeyValue(w http.ResponseWriter, r *http.Request, client valkey.Client, logger *zap.Logger) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	value, err := client.Do(ctx, client.B().Get().Key(key).Build()).AsBytes()
	if value == nil || err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(value)
}

func setKeyValue(w http.ResponseWriter, r *http.Request, client valkey.Client, logger *zap.Logger) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request body is required", http.StatusBadRequest)
		return
	}

	err = client.Do(ctx, client.B().Set().Key(key).Value(string(value)).ExSeconds(3600).Build()).Error()
	if err != nil {
		http.Error(w, "Failed to save key-value into valkey", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(value)
}
