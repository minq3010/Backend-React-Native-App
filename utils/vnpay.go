package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

func CreateVnpayURL(orderID string, amount int) (string, error) {
	params := url.Values{}
	params.Set("vnp_Version", "2.1.0")
	params.Set("vnp_Command", "pay")
	params.Set("vnp_TmnCode", os.Getenv("VNP_TMNCODE"))
	params.Set("vnp_Amount", fmt.Sprintf("%d", amount*100)) // nhân 100 vì VNPAY yêu cầu
	params.Set("vnp_CurrCode", "VND")
	params.Set("vnp_TxnRef", orderID)
	params.Set("vnp_OrderInfo", "Thanh toán đơn hàng "+orderID)
	params.Set("vnp_OrderType", "other")
	params.Set("vnp_Locale", "vn")
	params.Set("vnp_ReturnUrl", os.Getenv("VNP_RETURNURL"))
	params.Set("vnp_IpAddr", "127.0.0.1")
	params.Set("vnp_CreateDate", time.Now().Format("20060102150405"))

	// Sắp xếp keys theo thứ tự alphabet
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Tạo raw string theo thứ tự keys
	var rawBuilder strings.Builder
	for i, k := range keys {
		rawBuilder.WriteString(fmt.Sprintf("%s=%s", k, params.Get(k)))
		if i < len(keys)-1 {
			rawBuilder.WriteString("&")
		}
	}
	raw := rawBuilder.String()

	// Tạo chữ ký HMAC SHA512
	h := hmac.New(sha512.New, []byte(os.Getenv("VNP_HASHSECRET")))
	h.Write([]byte(raw))
	signature := hex.EncodeToString(h.Sum(nil))

	// Trả về URL thanh toán
	return fmt.Sprintf("%s?%s&vnp_SecureHash=%s", os.Getenv("VNP_URL"), raw, signature), nil
}

func GenerateOrderID(userID uint, eventID uint) string {
	return fmt.Sprintf("ORDER_%d_%d_%d", userID, eventID, time.Now().Unix())
}

func VerifyVnpaySignature(params map[string]string, receivedSignature string) bool {
	// 1. Tạo bản copy để không thay đổi map gốc
	paramsCopy := make(map[string]string)
	for k, v := range params {
		// Chỉ thêm các tham số bắt đầu với "vnp_" và không phải SecureHash
		if strings.HasPrefix(k, "vnp_") && k != "vnp_SecureHash" {
			paramsCopy[k] = v
		}
	}

	// 2. Sắp xếp keys theo thứ tự alphabet
	var keys []string
	for k := range paramsCopy {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. Ghép lại thành chuỗi raw
	var rawParts []string
	for _, k := range keys {
		// Chỉ thêm nếu value không rỗng
		if paramsCopy[k] != "" {
			rawParts = append(rawParts, fmt.Sprintf("%s=%s", k, paramsCopy[k]))
		}
	}
	raw := strings.Join(rawParts, "&")

	// 4. Tạo chữ ký HMAC SHA512 từ raw string
	h := hmac.New(sha512.New, []byte(os.Getenv("VNP_HASHSECRET")))
	h.Write([]byte(raw))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// 5. Debug log
	fmt.Println("🔐 Filtered params:", paramsCopy)
	fmt.Println("🔐 Raw string      :", raw)
	fmt.Println("🔐 Expected signature:", expectedSignature)
	fmt.Println("🔐 Received signature:", receivedSignature)

	// 6. So sánh không phân biệt hoa thường
	return strings.EqualFold(expectedSignature, receivedSignature)
}
