package handlers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PlanScanRequest struct {
	Target string `json:"target"`
}

func PlanScanHandler(c *gin.Context) {
	var req PlanScanRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	targetInfo, err := os.Stat(req.Target)
	if err != nil || targetInfo.IsDir() {
		fmt.Println("Invalid plan file target")
		c.JSON(http.StatusBadRequest, gin.H{"error": "target must be a JSON plan file"})
		return
	}

	// JSON 파일 확장자 확인
	if filepath.Ext(req.Target) != ".json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .json Terraform plan files are supported"})
		return
	}

	baseDir := "scan-results"
	dateDir := time.Now().Format("2006-01-02")
	saveDir := filepath.Join(baseDir, dateDir)
	os.MkdirAll(saveDir, 0755)

	fileName := filepath.Base(req.Target)
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	resultFile := filepath.Join(saveDir, fmt.Sprintf("%s-plan-scan-result.json", baseName))

	cmd := exec.Command("trivy", "config", "--format", "json", "--quiet", req.Target)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Plan scan failed:", fileName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := os.WriteFile(resultFile, output, 0644); err != nil {
		fmt.Println("Save failed:", fileName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Plan scan succeeded & saved:", resultFile)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"file":   resultFile,
	})
}
