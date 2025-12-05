package commonTool

import (
	"testing"
	"time"
)

func TestStructToSignString(t *testing.T) {
	var req struct {
		TraceId      int64  `json:"traceId"`
		MerchantCode string `json:"merchantCode"`
		Chain        string `json:"chain"`  //主链名称 TRON、ETH、SOLANA
		Remark       string `json:"remark"` //备注
		Sign         string `json:"sign"`
	} = struct {
		TraceId      int64  `json:"traceId"`
		MerchantCode string `json:"merchantCode"`
		Chain        string `json:"chain"`
		Remark       string `json:"remark"`
		Sign         string `json:"sign"`
	}{TraceId: time.Now().Unix(), MerchantCode: "10000", Chain: "TRON", Remark: "123000", Sign: "jjjjwiiie"}

	signString := StructToSignString(req, "sign")
	t.Log(signString)
}
