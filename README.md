# trivy-scan-api

## 1️⃣ Go 프로젝트 초기화
먼저 작업 디렉토리를 만들고 Go 모듈을 초기화해.
```
mkdir trivy-scan-api
cd trivy-scan-api
go mod init trivy-scan-api
go mod tidy
```

## 2️⃣ Trivy 설치 (CLI 실행 확인)
Trivy를 로컬에 설치해야 CLI 호출이 가능해.
macOS 기준으로는 다음 명령어 👇
```
brew install aquasecurity/trivy/trivy
```

설치 확인:
```
trivy -v
```

## 3️⃣ Trivy CLI 테스트
테스트용 Terraform 파일 하나 만들어서 스캔해보자 👇
```
mkdir test-tf
cd test-tf
echo 'resource "aws_s3_bucket" "public" { acl = "public-read" }' > main.tf
cd ..
```

Trivy 실행 테스트:
```
trivy config --format json ./test-tf > result.json
cat result.json
```

정상적으로 JSON 결과가 나오면 OK.
 이 결과를 나중에 Go 코드로 파싱할 거야.

## 4️⃣ 간단한 API 서버 만들기
Gin을 설치하고 API 스캐폴드 작성 👇
```
go get github.com/gin-gonic/gin
```

main.go 파일 생성:
```
package main

import (
    "github.com/gin-gonic/gin"
    "trivy-scan-api/handlers"
)

func main() {
    r := gin.Default()
    r.POST("/scan", handlers.ScanHandler)
    r.Run(":8080")
}
```


## 5️⃣ Trivy CLI 호출 로직
handlers/scan_handler.go 생성:
```
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"

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

    // Trivy 실행
    cmd := exec.Command("trivy", "config", "--format", "json", req.Target)
    output, err := cmd.Output()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Trivy error: %v", err)})
        return
    }

    // 결과 파싱
    var result interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid JSON from Trivy"})
        return
    }

    c.JSON(http.StatusOK, result)
}
```


## 6️⃣ 실행 및 테스트
서버 실행:
```
go run main.go
```

다른 터미널에서 요청 테스트:
```
curl -X POST http://localhost:8080/scan \
  -H "Content-Type: application/json" \
  -d '{"target":"./test-tf"}'
```

정상적으로 Trivy 결과(JSON)가 응답되면 성공 🎉

## 7️⃣ (선택) Docker 환경 구성
나중에 Tiger에서 연동하려면 Docker로 묶는 게 좋아.
 Dockerfile 예시 👇
```
FROM golang:1.22

Install Trivy
RUN apt-get update && apt-get install -y wget \
    && wget https://github.com/aquasecurity/trivy/releases/latest/download/trivy_0.54.1_Linux-64bit.tar.gz \
    && tar zxvf trivy_0.54.1_Linux-64bit.tar.gz \
    && mv trivy /usr/local/bin/

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main .

EXPOSE 8080
CMD ["./main"]
```

✅ 여기까지 하면
이 시점에서 네 API는 이미 다음 기능을 갖고 있어:
/scan POST 요청 수신


입력 받은 Terraform 디렉토리에 대해
 trivy config --format json 실행


결과(JSON)를 그대로 반환


즉, Trivy를 API 형태로 감싸는 최소 MVP 구조 완성 🎯

## 다음 단계로는:
models/trivy_result.go 만들어서 JSON 파싱을 구조체로 변환

Severity 기준 필터링/요약 기능 추가

Tiger 연동 방식 (자동 Git clone → 스캔)으로 확장


Trivy can also scan Terraform plan files in JSON format, which can be useful for analyzing the planned infrastructure changes before applying them. Generate a Terraform plan file.

### Terminal:
```
terraform init
terraform plan -out tfplan.binary
```

Convert the binary plan file to JSON:

### Terminal:
```
terraform show -json tfplan.binary > tfplan.json
```

### Terminal:
```
lsof -i :8080
kill -9 48159
go run main.go
```

### Terminal:
```
curl -X POST http://localhost:8080/scan -H "Content-Type: application/json" -d '{"target":"./test/tf"}'
```
```
curl -X POST http://localhost:8080/scan -H "Content-Type: application/json" -d '{"target":"./test/tf/test_08.tf"}'
```
```
trivy-scan-api % curl -X POST http://localhost:8080/scan/plan \
```
-H "Content-Type: application/json" \
-d '{"target":"./test/json/tfplan_02.json"}'
