package cmd

import (
	"context"
	"fmt"
	"log"
	"motion/pkgs/config"
	"motion/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		gin.SetMode(gin.ReleaseMode)

		router := gin.Default()
		router.POST("/webhook", routes.HandleWebhook)

		srv := &http.Server{
			Addr:    ":" + config.All.CurrentPort,
			Handler: router,
		}

		sigs := make(chan os.Signal, 1)

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.Println("Signal received: ", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Error: ", err)
			}

			fmt.Println("Server stopped!")
		}()

		log.Println("Server listening in port: ", config.All.CurrentPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error while starting server:", err)
		}
	},
}
