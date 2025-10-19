package main

import (
	"trivy-scan-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// .tf .tfvars 스캔용
	r.POST("/scan", handlers.ScanHandler)
	// Terraform plan JSON 스캔용
	r.POST("/scan/plan", handlers.PlanScanHandler)

	r.Run(":8080")
}
