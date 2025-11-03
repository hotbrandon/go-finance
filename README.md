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