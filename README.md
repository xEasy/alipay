### 支付宝蚂蚁金服当面付接口

----

#### 介绍

封装了支付宝蚂蚁金服当面付接口

https://doc.open.alipay.com/docs/doc.htm?spm=a219a.7629140.0.0.11hKzV&treeId=193&articleId=105072&docType=1

#### Getting Started

-----


##### Installing

```
go get github.com/xEasy/alipay
```

#### Usage

**Example Program:**

```go
package main

import (
    "github.com/xEasy/alipay"
)

func main() {
	AlipayClient, err := alipay.NewMerchant(appid, "path_to_alipay_private.pem", "path_to_alipay_public.pem")

	if err != nil {
		panic(err)
	}


  // 扫码支付接口
	respResult, _, err := AlipayClient.PlaceOrder(
		map[string]string{
			"out_trade_no":    "900001",
			"total_amount":    "1000", // 单位为分
			"subject":         "商品名称",
			"body":            "附加信息",
			"operator_id":     "operator_id",
			"terminal_id":     "terminal_id",
			"notify_url":      "http://notify.url"
			"timeout_express": "1d",
			"sub_merchant_id": "sub_merchant_id", // Only for Alipay Bank partner
		},
	)

  if err != nil {
    fmt.Println("request Alipay FAIL:", err.Error())
  }

  fmt.Println("request Alipay SUCCESS:", respResult)


  // 处理支付异步支付通知回调
  data = "buyer_id=20880..中间省略N个字符...seller_email=zfbtest25%40service.aliyun.com"
	tradeResult, err := AlipayClient.Notify([]byte(data))
	if err != nil {
    fmt.Println("request Alipay Notify FAIL:", err.Error())
		return
	}
  if tradeResult.IsTradeSuccess() {
     // Handler your own logic
  }
}
```
