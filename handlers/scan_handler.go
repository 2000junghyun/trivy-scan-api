package handlers

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ScanRequest struct {
	Target string `json:"target"`
}

func ScanHandler(c *gin.Context) {
	var req ScanRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	baseDir := "scan-results"
	dateDir := time.Now().Format("2006-01-02")
	saveDir := filepath.Join(baseDir, dateDir)
	os.MkdirAll(saveDir, 0755)

	targetInfo, err := os.Stat(req.Target)
	if err != nil {
		fmt.Println("Scan failed (invalid target)")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target path"})
		return
	}

	// ──────────────── [1] 단일 파일 스캔 ────────────────
	if !targetInfo.IsDir() {
		fileName := filepath.Base(req.Target)
		baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		resultFile := filepath.Join(saveDir, fmt.Sprintf("%s-scan-result.json", baseName))

		cmd := exec.Command("trivy", "config", "--format", "json", "--quiet", req.Target)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Scan failed:", fileName)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := os.WriteFile(resultFile, output, 0644); err != nil {
			fmt.Println("Save failed:", fileName)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Scan succeeded & saved:", resultFile)
		c.JSON(http.StatusOK, gin.H{"status": "success", "file": resultFile})
		return
	}

	// ──────────────── [2] 디렉토리 내 파일별 스캔 (.tf + .tfvars) ────────────────
	var results []string
	err = filepath.WalkDir(req.Target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// .tf 또는 .tfvars 파일만 스캔
		if !d.IsDir() {
			ext := filepath.Ext(d.Name())
			if ext == ".tf" || ext == ".tfvars" {
				fileName := strings.TrimSuffix(d.Name(), ext)
				resultFile := filepath.Join(saveDir, fmt.Sprintf("%s-scan-result.json", fileName))

				cmd := exec.Command("trivy", "config", "--format", "json", "--quiet", path)
				output, err := cmd.Output()
				if err != nil {
					fmt.Println("Scan failed:", d.Name())
					return err
				}

				if err := os.WriteFile(resultFile, output, 0644); err != nil {
					fmt.Println("Save failed:", d.Name())
					return err
				}

				fmt.Println("Scan succeeded & saved:", resultFile)
				results = append(results, resultFile)
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "partial_failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "saved_files": results})
}
