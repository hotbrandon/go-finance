# go-finance
Using Go to fetch stock or cryptocurrency data

## fetch stock data from FinMind
https://api.finmindtrade.com/api/v4/data?dataset=TaiwanStockPrice&data_id=2330&start_date=2025-10-31

```json
{
  "msg": "success",
  "status": 200,
  "data": [
    {
      "date": "2025-10-31",
      "stock_id": "2330",
      "Trading_Volume": 33838761,
      "Trading_money": 51073519270,
      "open": 1515,
      "max": 1525,
      "min": 1500,
      "close": 1500,
      "spread": -5,
      "Trading_turnover": 45696
    }
  ]
}
```

## swagger support

### Download Swag for Go by using:
go install github.com/swaggo/swag/cmd/swag@latest

### Download gin-swagger by using:
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

### Import following
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files

### add comments to endpoints and run in the project root
swag init -g main.go -o docs

## scheduler
go get github.com/robfig/cron/v3

## Docker shutdown behavior
docker compose down will ask Docker to stop each container: Docker sends SIGTERM (syscall.SIGTERM, signal 15) to PID 1 inside the container.
Docker waits for a grace period (default 10s) for the process to exit. If the container is still running after the grace period, Docker sends SIGKILL (signal 9) to forcibly terminate it.
After stopping, compose removes the stopped containers and networks (and volumes if you pass the appropriate flags).
Practical notes for your Go app

Handle SIGTERM (not SIGINT only). On SIGTERM you should cancel your root context, stop schedulers, and call http.Server.Shutdown to finish in‑flight requests.
Ensure your Go binary is PID 1 in the container (ENTRYPOINT directly runs the binary). If you wrap it in a shell that doesn’t forward signals, SIGTERM may not reach your process.
You can increase the grace period:
CLI: docker compose down --timeout 30
Compose file: add stop_grace_period: "30s" under the service
Test locally: docker stop --time=30 <container> simulates the same behavior.
Summary: expect SIGTERM first, then SIGKILL after the timeout — make sure your app reacts to SIGTERM for graceful shutdown.