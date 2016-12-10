package alipay

import (
	"strconv"
)

type RefundResponse struct {
	RefundResult `json:"alipay_trade_refund_response"`
	Sign         string `json:"sign"`
}

type RefundQueryResponse struct {
	RefundResult `json:"alipay_trade_fastpay_refund_query_response"`
	Sign         string `json:"sign"`
}

type RefundResult struct {
	Code       string `json:"code"`                   // "10000",
	Msg        string `json:"msg"`                    // "处理成功",
	SubCode    string `json:"sub_code,omitempty"`     //: "ACQ.TRADE_NOT_EXIST", ACQ.SYSTEM_ERROR，ACQ.INVALID_PARAMETER
	SubMsg     string `json:"sub_msg,omitempty"`      //: "交易不存在"
	TradeNO    string `json:"trade_no,omitempty"`     //  "2013112011001004330000121536",
	OutTradeNO string `json:"out_trade_no,omitempty"` //: "6823789339978248",
	// 申请退款返回参数
	BuyerID              string              `json:"buyer_user_id,omitempty"`           //: "159****5620",
	BuyerLogonID         string              `json:"buyer_logon_id,omitempty"`          //: "159****5620",
	GmtRefundPay         string              `json:"gmt_refund_pay,omitempty"`          // 退款支付时间	2014-11-27 15:45:57
	RefundDetailItemList []map[string]string `json:"refund_detail_item_list,omitempty"` // 选填	-	退款使用的资金渠道
	StoreName            string              `json:"store_name,omitempty"`              // 交易在支付时候的门店名称	望湘园联洋店
	RefundFee            float64             `json:"refund_fee,string"`                 // 退款总金额	88.88
	FundChange           string              `json:"fund_change,omitempty"`             // 本次退款是否发生了资金变化	Y
	// 查询退款返回参数
	TotalAmount  float64 `json:"total_amount,string"`      //: "88.88",
	RefundAmount float64 `json:"refund_amount,string"`     // 本次退款请求，对应的退款金额	12.33
	RefundReason string  `json:"refund_reason,omitempty"`  // 发起退款时，传入的退款原因	用户退款请求
	OutRequestNO string  `json:"out_request_no,omitempty"` // 本笔退款对应的退款请求号	20150320010101001
}

func (p *RefundResult) IsSuccess() bool {
	return p.Code == "10000"
}

type JsapiOrderResponse struct {
	PlaceOrderResult `json:"alipay_trade_create_response"`
	Sign             string `json:"sign"`
}

type PlaceOrderResponse struct {
	PlaceOrderResult `json:"alipay_trade_precreate_response"`
	Sign             string `json:"sign"`
}

type CancelOrderResponse struct {
	PlaceOrderResult `json:"alipay_trade_cancel_response"`
	Sign             string `json:"sign"`
}

type CloseOrderResponse struct {
	PlaceOrderResult `json:"alipay_trade_close_response"`
	Sign             string `json:"sign"`
}

type QueryOrderResponse struct {
	PlaceOrderQueryResult `json:"alipay_trade_query_response"`
	Sign                  string `json:"sign"`
}

type MicroPayOrderResponse struct {
	PlaceOrderQueryResult `json:"alipay_trade_pay_response"`
	Sign                  string `json:"sign"`
}

type PlaceOrderQueryResult struct {
	Code              string              `json:"code"`                          // "10000",
	Msg               string              `json:"msg"`                           // "处理成功",
	SubCode           string              `json:"sub_code,omitempty"`            //: "ACQ.TRADE_NOT_EXIST", ACQ.SYSTEM_ERROR，ACQ.INVALID_PARAMETER
	SubMsg            string              `json:"sub_msg,omitempty"`             //: "交易不存在"
	TradeNO           string              `json:"trade_no"`                      //  "2013112011001004330000121536",
	OutTradeNO        string              `json:"out_trade_no"`                  //: "6823789339978248",
	TradeStatus       string              `json:"trade_status"`                  //: "TRADE_SUCCESS",
	BuyerID           string              `json:"buyer_id,omitempty"`            //: "159****5620",
	BuyerLogonID      string              `json:"buyer_logon_id,omitempty"`      //: "159****5620",
	TotalAmount       float64             `json:"total_amount,string"`           //: "88.88",
	ReceiptAmount     float64             `json:"receipt_amount,string"`         //: "8.88",
	SendPayDate       string              `json:"send_pay_date"`                 //: "2014-11-27 15:45:57",
	StoreID           string              `json:"store_id,omitempty"`            //:"NJ_S_001",
	TerminalID        string              `json:"terminal_id,omitempty"`         //:"NJ_T_001",
	FundBillList      []map[string]string `json:"fund_bill_list,omitempty"`      //: [
	VoucherDetailList []map[string]string `json:"voucher_detail_list,omitempty"` //: [
	GmtPayment        string              `json:"gmt_payment,omitempty"`
}

func (p *PlaceOrderQueryResult) IsSuccess() bool {
	return p.TradeStatus == "TRADE_SUCCESS" || p.Code == "10000"
}

func (p *PlaceOrderQueryResult) IsCodeSuccess() bool {
	return p.Code == "10000"
}

// 预付单返回
type PlaceOrderResult map[string]string

func (p PlaceOrderResult) QrCode() string {
	return p["qr_code"]
}

func (p PlaceOrderResult) IsSuccess() bool {
	return p["code"] == "10000"
}

// 交易结果
// map[seller_id:2088021244059960 trade_status:TRADE_SUCCESS gmt_payment:2016-01-10 14:15:18 point_amount:0.00 trade_no:2016011021001004500071808018 invoice_amount:1.00 notify_type:trade_status_sync receipt_amount:1.00 buyer_logon_id:xia***@gmail.com buyer_pay_amount:1.00 subject:汽车洗车-好车店 gmt_create:2016-01-10 14:14:57 seller_email:op@lovechebang.com notify_id:c5be78fb74a7d7f2336582597678a5djuw fund_bill_list:[{"amount":"1.00","fundChannel":"ALIPAYACCOUNT"}] notify_time:2016-01-10 14:15:18 buyer_id:2088002359340503 app_id:2015081700218350 total_amount:1.00 out_trade_no:M14524064880000000000003]
type TradeResult map[string]string

func (p TradeResult) IsSuccess() bool {
	return p["trade_status"] == "TRADE_SUCCESS"
}

func (n TradeResult) IsTradeSuccess() bool {
	return IsTradeSuccess(n["trade_status"])
}

func (p TradeResult) TotalFee() int64 {
	m, _ := strconv.ParseFloat(p["total_amount"], 64)
	return int64(m * 100)
}
