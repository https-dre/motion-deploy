package cmd

import (
	"motion/routes"
	"motion/pkgs/config"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o servidor HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		gin.SetMode(gin.ReleaseMode)

		r := gin.Default()
		r.POST("/webhook", routes.HandleWebhook)

		log.Println("Servidor ouvindo na porta", config.General.CurrentPort)
		if err := r.Run(":" + config.General.CurrentPort); err != nil {
			log.Fatal("Erro ao iniciar servidor:", err)
		}
	},
}
