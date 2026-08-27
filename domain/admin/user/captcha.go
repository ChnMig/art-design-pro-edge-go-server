package user

import (
	"fmt"

	"github.com/mojocn/base64Captcha"

	"api-server/db/rdb/captcha"
)

const (
	CaptchaMinWidth  = 80
	CaptchaMaxWidth  = 600
	CaptchaMinHeight = 40
	CaptchaMaxHeight = 200
)

type Captcha struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

func GenerateCaptcha(width, height int) (Captcha, error) {
	if !validCaptchaSize(width, height) {
		return Captcha{}, fmt.Errorf("%w: width=%d height=%d", ErrInvalidCaptchaSize, width, height)
	}
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

func validCaptchaSize(width, height int) bool {
	return width >= CaptchaMinWidth && width <= CaptchaMaxWidth &&
		height >= CaptchaMinHeight && height <= CaptchaMaxHeight
}

func VerifyCaptcha(captchaID, captchaValue string) bool {
	return captcha.GetRedisStore().Verify(captchaID, captchaValue, true)
}
