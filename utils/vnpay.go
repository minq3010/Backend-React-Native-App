package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sort"
	"time"
)

func CreateVnpayURL(orderID string, amount int) (string, error) {
	params := url.Values{} 
	params.Set("vnp_Version", "2.1.0")
	params.Set("vnp_Command", "pay")
	params.Set("vnp_TmnCode", os.Getenv("VNP_TMNCODE"))
	params.Set("vnp_Amount", fmt.Sprintf("%d", amount*100))
	params.Set("vnp_CurrCode", "VND")
	params.Set("vnp_TxnRef", orderID)
	params.Set("vnp_OrderInfo", "Thanh toán đơn hàng " + orderID)
	params.Set("vnp_OrderType", "other")
	params.Set("vnp_Locale", "vn")
	params.Set("vnp_ReturnUrl", os.Getenv("VNP_RETURNURL"))
	params.Set("vnp_IpAddr", "127.0.0.1")
	params.Set("vnp_CreateDate", time.Now().Format("20060102150405"))

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var raw string
	for i, k := range keys {
		raw += fmt.Sprintf("%s=%s", k, params.Get(k))
		if i < len(keys)-1 {
			raw += "&"
		}
	}

	h := hmac.New(sha512.New, []byte(os.Getenv("VNP_HASHSECRET")))
	h.Write([]byte(raw))
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s?%s&vnp_SecureHash=%s", os.Getenv("VNP_URL"), params.Encode(), signature), nil
}