package authenticator

import "testing"

func TestGenerateTOTPSecret(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Errorf("Error generating TOTP secret: %s", err)
		return
	}

	t.Logf("TOTP secret: %s", secret)
}

func TestBuildOtpAuthUrl(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Errorf("Error generating TOTP secret: %s", err)
		return
	}

	url := BuildOtpAuthUrl(secret, "TZ", "HY0001")
	t.Logf("TOTP url: %s", url)
	//otpauth://totp/TZ:HY0001?secret=CFVIZUL3VLZHAI66CCJPJILH6655JSWV&issuer=TZ&digits=6&period=30
}

func TestVerifyTOTP(t *testing.T) {
	secret := "CFVIZUL3VLZHAI66CCJPJILH6655JSWV"
	totp, err := VerifyTOTP(secret, "162067", 30, 6, 1)
	if err != nil {
		t.Errorf("Error verifying TOTP: %s", err)
		return
	}
	t.Logf("TOTP: %v", totp)
}
