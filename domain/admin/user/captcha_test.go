package user

import (
	"errors"
	"testing"

	"github.com/mojocn/base64Captcha"
)

func TestGenerateCaptchaRejectsUnsafeDimensionsBeforeStoreAccess(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{name: "极小尺寸", width: 1, height: 1},
		{name: "宽度低于下限", width: CaptchaMinWidth - 1, height: CaptchaMinHeight},
		{name: "宽度超过上限", width: CaptchaMaxWidth + 1, height: CaptchaMinHeight},
		{name: "高度低于下限", width: CaptchaMinWidth, height: CaptchaMinHeight - 1},
		{name: "高度超过上限", width: CaptchaMinWidth, height: CaptchaMaxHeight + 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GenerateCaptcha(test.width, test.height)
			if !errors.Is(err, ErrInvalidCaptchaSize) {
				t.Fatalf("GenerateCaptcha() error = %v, want ErrInvalidCaptchaSize", err)
			}
		})
	}
}

func TestCaptchaBoundaryDimensionsCanBeRendered(t *testing.T) {
	tests := []struct {
		width  int
		height int
	}{
		{width: CaptchaMinWidth, height: CaptchaMinHeight},
		{width: CaptchaMinWidth, height: CaptchaMaxHeight},
		{width: CaptchaMaxWidth, height: CaptchaMinHeight},
		{width: CaptchaMaxWidth, height: CaptchaMaxHeight},
	}

	for _, test := range tests {
		if !validCaptchaSize(test.width, test.height) {
			t.Fatalf("边界尺寸被错误拒绝: %dx%d", test.width, test.height)
		}
		driver := base64Captcha.NewDriverDigit(test.height, test.width, 6, 0.2, 50)
		_, content, _ := driver.GenerateIdQuestionAnswer()
		if _, err := driver.DrawCaptcha(content); err != nil {
			t.Fatalf("渲染边界尺寸 %dx%d 失败: %v", test.width, test.height, err)
		}
	}
}
