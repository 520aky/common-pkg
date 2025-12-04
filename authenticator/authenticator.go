package authenticator

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"
)

func GenerateTOTPSecret() (string, error) {
	// 20 字节随机（160-bit）
	secretBytes := make([]byte, 20)

	// 使用 crypto/rand 保证安全级别
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", err
	}

	// Base32 编码（Google Authenticator 要求大写、无=填充）
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secretBytes)

	return secret, nil
}

func BuildOtpAuthUrl(secret, issuer, account string) string {
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s&digits=6&period=30",
		issuer, account, secret, issuer,
	)
}

// VerifyTOTP 验证单个 TOTP 密钥是否匹配玩家输入的 code
// 参数：
//
//	secretBase32: base32 编码的 secret（通常从数据库读取）
//	code: 玩家输入的 6 位数字验证码
//	period: 时间步长，标准为 30 秒
//	digits: 验证码位数，常见是 6
//	window: 时间偏移容差（比如 1=前后各 1 个 time-step）
//
// 返回：true=验证通过
func VerifyTOTP(secretBase32, code string) (bool, error) {
	var (
		period = 30
		digits = 6
		window = 10
	)
	//if period <= 0 {
	//	period = 30
	//}
	//if digits <= 0 {
	//	digits = 6
	//}
	//if window < 0 {
	//	window = 0
	//}

	// 规范化 secret（去掉空格，大写）
	secretBase32 = strings.ToUpper(strings.ReplaceAll(secretBase32, " ", ""))

	// 1. Base32 解码
	secret, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secretBase32)
	if err != nil {
		return false, fmt.Errorf("decode secret error: %w", err)
	}
	if len(secret) == 0 {
		return false, errors.New("empty secret")
	}

	// 2. 当前时间步（counter）
	now := time.Now().Unix()
	timestep := now / int64(period)

	// 3. 在 [-window, +window] 范围内尝试
	for i := -window; i <= window; i++ {
		counter := timestep + int64(i)
		if counter < 0 {
			continue
		}
		expected := generateTOTPCode(secret, counter, digits)
		if expected == code {
			return true, nil
		}
	}

	return false, nil
}

// generateTOTPCode 根据 secret + counter 生成 digits 位 TOTP
func generateTOTPCode(secret []byte, counter int64, digits int) string {
	// 8 字节大端序 counter
	var b [8]byte
	for i := 7; i >= 0; i-- {
		b[i] = byte(counter & 0xff)
		counter >>= 8
	}

	// HMAC-SHA1
	h := hmac.New(sha1.New, secret)
	h.Write(b[:])
	hash := h.Sum(nil)

	// 动态截断（RFC 4226）
	offset := hash[len(hash)-1] & 0x0f
	binCode := (int(hash[offset])&0x7f)<<24 |
		(int(hash[offset+1])&0xff)<<16 |
		(int(hash[offset+2])&0xff)<<8 |
		(int(hash[offset+3]) & 0xff)

	mod := 1
	for i := 0; i < digits; i++ {
		mod *= 10
	}
	token := binCode % mod

	// 补零格式化为固定位数
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, token)
}
