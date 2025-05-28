package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"motion/pkgs/config"
	"motion/pkgs/models"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		c.String(http.StatusForbidden, "Sem assinatura")
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao ler corpo da requisição")
		return
	}

	// Verifica assinatura
	if !verifySignature(signature, body, []byte(config.All.Secret)) {
		c.String(http.StatusForbidden, "Assinatura inválida")
		return
	}

	var payload struct {
		Ref        string `json:"ref"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.String(http.StatusBadRequest, "Payload inválido")
		return
	}

	repo := payload.Repository.FullName
	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	if repoConf, ok := config.All.Repos[repo]; ok && repoConf.Branch == branch {
		go deploy(repoConf)
	}

	c.Status(http.StatusOK)
}

func verifySignature(signature string, body, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func deploy(repo models.RepoConfig) {
	fmt.Println("Executando deploy para", repo.Path)

	port1 := strconv.Itoa(repo.Ports[0])
	port2 := strconv.Itoa(repo.Ports[1])

	cmd := exec.Command("bash", "./up-docker.sh", port1, port2)
	cmd.Dir = "./"

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erro:", err)
	}
	fmt.Println(string(output))
}