# Go Tools SDK

## 簡介
這是一個用 Go 語言開發的 SDK，提供以下功能模組，方便整合到其他專案中：
- **Logger**: 提供日誌記錄功能。
- **Prometheus**: 與 Prometheus 進行整合的客戶端。
- **RabbitMQ**: 與 RabbitMQ 進行整合的客戶端。


## 📦 安裝

```bash
go get github.com/Wuli-Giao-Giao/tools@v0.2.0
```

---

## 📌 模組介紹

### 🐰 RabbitMQ Client

快速建立連線與發送訊息。

```go
client, err := tools.NewClient("amqp://guest:guest@localhost:5672/")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

err = client.Publish("my-exchange", "my-key", []byte(`{"msg":"hello"}`))
if err != nil {
    log.Printf("publish failed: %v", err)
}
```

---

### 📈 Prometheus Client (with TLS/mTLS)

支援 TLS / mTLS，包含：
- 自訂 CA 憑證
- 客戶端憑證與私鑰
- 可選擇跳過 TLS 驗證

```go
client, err := tools.NewClient(
    "https://prom.example.com",
    "username",
    "password",
    "./client.crt",
    "./client.key",
    "./ca.crt",
    false, // insecureTLS: true 表示跳過驗證
)
if err != nil {
    log.Fatal(err)
}

// 假設你有封裝 client.DoQuery(query string)
resp, err := client.DoQuery("up")
```

---

### 📋 Logrus Logger 初始化工具

快速建立一個可自訂等級與輸出的 logrus logger。

```go
log := tools.NewLogrusLogger("debug", os.Stdout)
log.Info("Logger is ready")
```

支援等級：`trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`  
輸出位置為任意符合 `io.Writer` 的實例（如檔案或 `os.Stdout`）

---

### ⚙️ Server 啟動與管理工具

提供統一介面與工具，方便同時啟動多個 server（如 gRPC、HTTP），並支援優雅關閉。

---

#### 🏗️ 基本概念

這個模組包含三部分：
- **Server 介面**  
  定義 `Start()` 和 `Stop(ctx)` 方法。
- **Runner**  
  負責管理多個 server，統一啟動、監聽系統訊號、優雅關閉。
- **內建封裝**  
  - gRPC server → 使用 `NewGRPCServer(port int, grpcServer *grpc.Server)`
  - HTTP server → 使用 `NewHTTPServer(addr string, handler http.Handler)`

---

#### 🚀 使用範例

```go
import (
    "github.com/Wuli-Giao-Giao/tools/server"
    "google.golang.org/grpc"
    "net/http"
)

func main() {
    // 建立 gRPC server
    grpcSrv := grpc.NewServer()
    grpcServer := server.NewGRPCServer(50051, grpcSrv)

    // 建立 HTTP server
    httpHandler := http.NewServeMux()
    httpHandler.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })
    httpServer := server.NewHTTPServer(":8080", httpHandler)

    // 使用 Runner 管理多個 server
    runner := server.NewRunner(grpcServer, httpServer)
    runner.Run()
}