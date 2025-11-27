package security

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行哈希加密
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength 验证密码强度
// 要求：至少8位，包含字母和数字
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少为8位")
	}

	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
	}

	if !hasLetter {
		return errors.New("密码必须包含字母")
	}
	if !hasNumber {
		return errors.New("密码必须包含数字")
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}
	return nil
}

// ValidateUsername 验证用户名格式
// 要求：3-50位，只能包含字母、数字、下划线
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return errors.New("用户名长度必须在3-50位之间")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	return nil
}
