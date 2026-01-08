package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"krapper/internal/global"
	"krapper/internal/httpsrv"
	"krapper/internal/k8s"
	"krapper/internal/misc"
	"krapper/internal/wrapstore"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

var serveParams struct {
	logConfig   misc.LogConfig
	httpConfig  httpsrv.Config
	wrapsFolder string
}

func init() {
	serveCmd.PersistentFlags().StringVarP(&serveParams.logConfig.Mode, "logMode", "", "text", "Log mode ('text' or 'json')")
	serveCmd.PersistentFlags().StringVarP(&serveParams.logConfig.Level, "logLevel", "l", "INFO", "Log level(DEBUG, INFO, WARN, ERROR)")

	serveCmd.PersistentFlags().BoolVarP(&serveParams.httpConfig.Tls, "tls", "t", false, "enable TLS")
	serveCmd.PersistentFlags().IntVar(&serveParams.httpConfig.DumpExchanges, "dumpExchanges", 0, "Dump http server req/resp (0, 1, 2 or 3)")
	serveCmd.PersistentFlags().StringVarP(&serveParams.httpConfig.BindAddr, "bindAddr", "a", "0.0.0.0", "Bind Address")
	serveCmd.PersistentFlags().IntVarP(&serveParams.httpConfig.BindPort, "bindPort", "p", 7777, "Bind port")
	serveCmd.PersistentFlags().StringVarP(&serveParams.httpConfig.CertDir, "certDir", "", "", "Certificate Directory")
	serveCmd.PersistentFlags().StringVar(&serveParams.httpConfig.CertName, "certName", "tls.crt", "Certificate Directory")
	serveCmd.PersistentFlags().StringVar(&serveParams.httpConfig.KeyName, "keyName", "tls.key", "Certificate Directory")
	serveCmd.PersistentFlags().StringVar(&serveParams.wrapsFolder, "wrapsFolder", "", "Path to wraps directory")
	_ = serveCmd.MarkPersistentFlagRequired("wrapsFolder")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "launch http server",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := misc.NewLogger(&serveParams.logConfig)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to load logging configuration: %v\n", err)
			os.Exit(2)
		}
		logger.Info("Starting krapper server", slog.String("logLevel", serveParams.logConfig.Level), slog.String("version", global.Version), slog.String("build", global.BuildTs))

		// Create and start HTTP server
		store, err := wrapstore.New(serveParams.wrapsFolder, logger)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to load wraps from '%s': %v\n", serveParams.wrapsFolder, err)
			os.Exit(2)
		}

		// Inject logger into context
		ctx := logr.NewContextWithSlogLogger(context.Background(), logger)

		// Initialize K8s client
		k8sClient, err := k8s.NewClient(logger)
		if err != nil {
			logger.Warn("Failed to initialize K8s client. K8s features will be disabled.", "error", err)
		}

		mux := http.NewServeMux()

		mux.HandleFunc("GET /api/v1/wraps", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(store.GetCatalog()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		mux.HandleFunc("GET /api/v1/wraps/{name}", func(w http.ResponseWriter, r *http.Request) {
			name := r.PathValue("name")
			wrap := store.GetWrap(name)
			if wrap == nil {
				http.Error(w, "Wrap not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(wrap); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		mux.HandleFunc("GET /api/v1/resources/{wrapName}", func(w http.ResponseWriter, r *http.Request) {
			name := r.PathValue("wrapName")
			wrap := store.GetWrap(name)
			if wrap == nil {
				http.Error(w, "Wrap not found", http.StatusNotFound)
				return
			}

			if k8sClient == nil {
				http.Error(w, "K8s client not initialized", http.StatusServiceUnavailable)
				return
			}

			// Determine namespace
			ns := wrap.Source.Namespace
			if wrap.Source.ClusterScoped {
				ns = "" // Ignore namespace for cluster scoped resources
			}

			list, err := k8sClient.ListResources(
				r.Context(),
				wrap.Source.ApiVersion,
				wrap.Source.Kind,
				ns,
				wrap.Source.Selector,
			)
			if err != nil {
				logger.Error("Failed to list resources", "error", err, "wrap", wrap.Name)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Clean up resources
			for i := range list.Items {
				list.Items[i].SetManagedFields(nil)
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(list.Items); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		httpServer := httpsrv.New("krapper", &serveParams.httpConfig, mux)

		if err := httpServer.Start(ctx); err != nil {
			logger.Error("Error starting HTTP server", "error", err)
			os.Exit(1)
		}
	},
}
