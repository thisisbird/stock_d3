## 運行步驟

1. go run 1minToOther.go
    - 一分k壓成其他分k,供策略使用,產生的csv擋在/public/data
    - 執行後前往使用圖形介面 localhost:8080/kline
    
2. go run StrategyPG.go
    - 策略邏輯的撰寫地方並執行
    - 產生csv檔在/public/data
    
