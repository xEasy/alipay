package alipay

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Merchant struct {
	AppId  string
	Env    string
	logger *log.Logger

	privateKey   *rsa.PrivateKey
	aliPublicKey *rsa.PublicKey
}

func NewMerchant(appid string, prikeyPath, aliPublicKeyPath string, env string) (*Merchant, error) {
	priKey, err := LoadPrivateKey(prikeyPath)
	if err != nil {
		return nil, err
	}

	aliPublicKey, err := LoadPublicKey(aliPublicKeyPath)
	if err != nil {
		return nil, err
	}

	return &Merchant{
		logger:       log.New(os.Stdout, "["+appid+"]", log.LstdFlags),
		AppId:        appid,
		Env:          env,
		privateKey:   priKey,
		aliPublicKey: aliPublicKey,
	}, nil
}

func (m *Merchant) MicroPayOrder(reqParams map[string]string) (aliResult *MicroPayOrderResponse, data []byte, err error) {
	totalFee, err := strconv.ParseInt(reqParams["total_amount"], 10, 32)
	if err != nil {
		err = errors.New("total_fee 类型错误")
		return
	}
	bizContent := map[string]interface{}{
		"store_id":        reqParams["store_id"],
		"out_trade_no":    reqParams["out_trade_no"],                  // 商户订单号
		"total_amount":    fmt.Sprintf("%.2f", float64(totalFee)/100), // 总金额
		"subject":         reqParams["subject"],                       // 商品名称
		"body":            reqParams["body"],
		"scene":           reqParams["scene"],
		"auth_code":       reqParams["auth_code"],
		"timeout_express": reqParams["timeout_express"],
	}
	if m.Env != "sandbox" {
		bizContent["sub_merchant"] = map[string]string{"merchant_id": reqParams["sub_merchant_id"]}
	}
	data, err = m.BizRequest(m.gatewayUrl(), "alipay.trade.pay", reqParams["notify_url"], bizContent)
	if err != nil {
		err = errors.New("支付返回数据格式错误")
		return
	}

	err = json.Unmarshal(data, &aliResult)
	if err != nil {
		err = errors.New("支付宝返回数据格式错误")
		return
	}

	m.Debug("place order resp:", aliResult)
	err = m.VerifyResponse(data, aliResult.Sign, "alipay_trade_pay_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
	}
	return
}

// 预下单API, totalFee单位为分
// https://app.alipay.com/market/document.htm?name=saomazhifu#page-14
func (m *Merchant) PlaceOrder(reqParams map[string]string) (aliResult *PlaceOrderResponse, data []byte, err error) {
	totalFee, err := strconv.ParseInt(reqParams["total_amount"], 10, 32)
	if err != nil {
		err = errors.New("total_fee 类型错误")
		return
	}
	bizContent := map[string]interface{}{
		"store_id":        reqParams["store_id"],
		"out_trade_no":    reqParams["out_trade_no"],                  // 商户订单号
		"total_amount":    fmt.Sprintf("%.2f", float64(totalFee)/100), // 总金额
		"subject":         reqParams["subject"],                       // 商品名称
		"body":            reqParams["body"],
		"timeout_express": reqParams["timeout_express"],
	}
	if m.Env != "sandbox" {
		bizContent["sub_merchant"] = map[string]string{"merchant_id": reqParams["sub_merchant_id"]}
	}
	data, err = m.BizRequest(m.gatewayUrl(), "alipay.trade.precreate", reqParams["notify_url"], bizContent)
	if err != nil {
		err = errors.New("支付返回数据格式错误")
		return
	}

	err = json.Unmarshal(data, &aliResult)
	if err != nil {
		err = errors.New("支付宝返回数据格式错误")
		return
	}

	m.Debug("place order resp:", aliResult)
	err = m.VerifyResponse(data, aliResult.Sign, "alipay_trade_precreate_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
	}
	return
}

/*
预下单API, totalFee单位为分
https://doc.open.alipay.com/docs/api.htm?spm=a219a.7395905.0.0.Zq03xn&docType=4&apiId=1046
*/
func (m *Merchant) JsapiOrder(reqParams map[string]string) (aliResult *JsapiOrderResponse, data []byte, err error) {
	totalFee, err := strconv.ParseInt(reqParams["total_amount"], 10, 32)
	if err != nil {
		err = errors.New("total_fee 类型错误")
		return
	}
	bizContent := map[string]interface{}{
		"store_id":        reqParams["store_id"],
		"out_trade_no":    reqParams["out_trade_no"],                  // 商户订单号
		"total_amount":    fmt.Sprintf("%.2f", float64(totalFee)/100), // 总金额
		"subject":         reqParams["subject"],                       // 商品名称
		"buyer_id":        reqParams["buyer_id"],
		"buyer_logon_id":  reqParams["buyer_logon_id"],
		"body":            reqParams["body"],
		"timeout_express": reqParams["timeout_express"],
	}
	if m.Env != "sandbox" {
		bizContent["sub_merchant"] = map[string]string{"merchant_id": reqParams["sub_merchant_id"]}
	}
	data, err = m.BizRequest(m.gatewayUrl(), "alipay.trade.create", reqParams["notify_url"], bizContent)
	if err != nil {
		err = errors.New("支付返回数据格式错误")
		return
	}

	err = json.Unmarshal(data, &aliResult)
	if err != nil {
		err = errors.New("支付宝返回数据格式错误")
		return
	}

	m.Debug("place order resp:", aliResult)
	err = m.VerifyResponse(data, aliResult.Sign, "alipay_trade_create_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
	}
	return
}

// 查询订单状态
// https://app.alipay.com/market/document.htm?name=saomazhifu#page-15
func (m *Merchant) QueryOrder(orderId string) (*QueryOrderResponse, []byte, error) {
	data, err := m.BizRequest(m.gatewayUrl(), "alipay.trade.query", "", map[string]interface{}{
		"out_trade_no": orderId, // 商户订单号
	})
	if err != nil {
		return nil, nil, err
	}

	var resp QueryOrderResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, nil, err
	}

	err = m.VerifyResponse(data, resp.Sign, "alipay_trade_query_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
		return nil, nil, err
	}

	return &resp, data, nil
}

/*
撤销订单
https://doc.open.alipay.com/docs/api.htm?spm=a219a.7395905.0.0.LOxDvL&docType=4&apiId=866
*/
func (m *Merchant) CancelOrder(orderId string) (*CancelOrderResponse, []byte, error) {
	data, err := m.BizRequest(m.gatewayUrl(), "alipay.trade.cancel", "", map[string]interface{}{
		"out_trade_no": orderId, // 商户订单号
	})
	if err != nil {
		return nil, nil, err
	}

	var resp CancelOrderResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, nil, err
	}

	err = m.VerifyResponse(data, resp.Sign, "alipay_trade_cancel_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
		return nil, nil, err
	}
	return &resp, data, nil
}

// 关闭订单
// https://app.alipay.com/market/document.htm?name=saomazhifu#page-16
func (m *Merchant) CloseOrder(orderId string) (*CloseOrderResponse, []byte, error) {
	data, err := m.BizRequest(m.gatewayUrl(), "alipay.trade.close", "", map[string]interface{}{
		"out_trade_no": orderId, // 商户订单号
	})
	if err != nil {
		return nil, nil, err
	}

	var resp CloseOrderResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, nil, err
	}

	err = m.VerifyResponse(data, resp.Sign, "alipay_trade_close_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
		return nil, nil, err
	}
	return &resp, data, nil
}

func (m *Merchant) RefundOrder(reqParams map[string]interface{}) (*RefundResponse, []byte, error) {
	data, err := m.BizRequest(m.gatewayUrl(), "alipay.trade.refund", "", reqParams)
	if err != nil {
		return nil, nil, err
	}

	var resp RefundResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, nil, err
	}

	err = m.VerifyResponse(data, resp.Sign, "alipay_trade_refund_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
		return nil, nil, err
	}
	return &resp, data, nil
}

func (m *Merchant) RefundOrderQuery(reqParams map[string]interface{}) (aliResult *RefundQueryResponse, data []byte, err error) {
	data, err = m.BizRequest(m.gatewayUrl(), "alipay.trade.fastpay.refund.query", "", reqParams)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &aliResult)
	if err != nil {
		return
	}

	err = m.VerifyResponse(data, aliResult.Sign, "alipay_trade_fastpay_refund_query_response")
	if err != nil {
		err = errors.New("支付宝数据来源异常")
		return
	}
	return
}

func (m *Merchant) Notify(data []byte) (TradeResult, error) {
	resp, err := m.ParseRequest(data)
	if err != nil {
		return nil, err
	}
	return TradeResult(resp), nil
}

func (m *Merchant) ParseRequest(res []byte) (Params, error) {
	resp, err := ParseParams(string(res))
	if err != nil {
		return nil, err
	}

	sig := resp["sign"]
	delete(resp, "sign")
	delete(resp, "sign_type")

	err = Verify(m.aliPublicKey, []byte(resp.Encode(false)), sig)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *Merchant) VerifyResponse(data []byte, sign, respKey string) (err error) {
	enc := strings.TrimPrefix(string(data), fmt.Sprintf(`{"%s":{`, respKey))
	idx := strings.Index(enc, `},"sign":`)
	if idx == -1 {
		err = errors.New("支付宝返回错误")
		return
	}

	enc = "{" + enc[:idx] + "}"

	//s := strings.Replace(string(enc), `/`, `\/`, -1)
	err = Verify(m.aliPublicKey, []byte(enc), sign)
	if err != nil {
		return
	}
	return
}

func (m *Merchant) Sign(data []byte) (string, error) {
	return Sign(m.privateKey, data)
}

func (m *Merchant) Verify(data []byte, sig string) error {
	return Verify(m.aliPublicKey, data, sig)
}

func (m *Merchant) Error(args ...interface{}) {
	args = append([]interface{}{"[ERR]"}, args...)
	m.logger.Println(args...)
}

func (m *Merchant) Errorf(format string, args ...interface{}) {
	m.logger.Printf("[ERR]"+format, args...)
}

func (m *Merchant) Debug(args ...interface{}) {
	args = append([]interface{}{"[DBG]"}, args...)
	m.logger.Println(args...)
}

func (m *Merchant) Debugf(format string, args ...interface{}) {
	m.logger.Printf("[DBG]"+format, args...)
}

func (m *Merchant) IsValid() bool {
	return m.AppId != "" && m.privateKey != nil && m.aliPublicKey != nil
}

func (m *Merchant) gatewayUrl() string {
	switch m.Env {
	case "sandbox":
		return sandboxGateWayUrl
	default:
		return gatewayUrl
	}
}

func (m *Merchant) BizRequest(url, method, notifyUrl string, bizData map[string]interface{}) ([]byte, error) {
	bizContent, err := json.Marshal(bizData)
	if err != nil {
		return nil, err
	}

	var req = Params{
		"app_id":      m.AppId,
		"method":      method,
		"charset":     "utf-8",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": string(bizContent),
		"sign_type":   "RSA",
	}
	if len(notifyUrl) > 0 {
		req["notify_url"] = notifyUrl
	}

	sig, err := Sign(m.privateKey, []byte(req.Encode(false)))
	if err != nil {
		return nil, err
	}

	req["sign"] = sig

	return doHttpPost(url, []byte(req.Encode(true)))
}
