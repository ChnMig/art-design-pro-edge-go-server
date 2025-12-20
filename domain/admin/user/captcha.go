package user

import (
	"github.com/mojocn/base64Captcha"

	"api-server/db/rdb/captcha"
)

type Captcha struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

func GenerateCaptcha(width, height int) (Captcha, error) {
	driver := base64Captcha.NewDriverDigit(height, width, 6, 0.2, 50)
	client := base64Captcha.NewCaptcha(driver, captcha.GetRedisStore())
	id, b64s, _, err := client.Generate()
	if err != nil {
		return Captcha{}, err
	}
	return Captcha{
		ID:    id,
		Image: b64s,
	}, nil
}

func VerifyCaptcha(captchaID, captchaValue string) bool {
	return captcha.GetRedisStore().Verify(captchaID, captchaValue, true)
}
