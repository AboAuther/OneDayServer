package sms

import (
	"context"
	"errors"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"

	internalRedis "one-day-server/internal/db/redis"
	"one-day-server/utils"
)

var (
	ErrCodeNotFound = errors.New("verification code not found or expired")
	ErrCodeMismatch = errors.New("verification code does not match")
)

func sendAliyunSMS(phone, templateCode string) error {
	return nil
	//client, err := dysmsapi.NewClientWithAccessKey(
	//	"cn-hangzhou",
	//	configs.MustGetEnv("ALIYUN_ACCESS_KEY_ID"),
	//	configs.MustGetEnv("ALIYUN_ACCESS_KEY_SECRET"),
	//)
	//if err != nil {
	//	return err
	//}
	//
	//request := dysmsapi.CreateSendSmsRequest()
	//request.Scheme = "https"
	//request.PhoneNumbers = phone
	//request.SignName = "YourSignName"           // 阿里云短信签名
	//request.TemplateCode = templateCode         // 模板 Code
	//request.TemplateParam = `{"code":"123456"}` // 参数（实际项目中动态生成）
	//
	//_, err = client.SendSms(request)
	//return err
}

// verifyAliyunSMS 验证用户提交的验证码
func verifyAliyunSMS(ctx context.Context, phone string, code string) error {
	redisKey := "sms:code:" + phone

	// 从 Redis 获取验证码
	storedCode, err := internalRedis.GetClient().GetResult(ctx, redisKey)
	if errors.Is(err, internalRedis.NilError) {
		return ErrCodeNotFound
	} else if err != nil {
		return err // 其他 Redis 错误
	}

	// 对比验证码
	if storedCode != code {
		return ErrCodeMismatch
	}

	if err := internalRedis.GetClient().Delete(ctx, redisKey); err != nil {
		return err
	}

	return nil
}

// SetVerificationCode 存储验证码到 Redis
func setVerificationCode(ctx context.Context, phone string, code string, ttl time.Duration) error {
	redisKey := "sms:code:" + phone
	return internalRedis.GetClient().WriteResultWithTTL(ctx, redisKey, code, ttl)
}

func SendSMSAndStoreCode(phone string) error {
	code := utils.GenerateRandomCode(6)
	err := sendAliyunSMS(phone, code)
	if err != nil {
		return err
	}
	logger.Infof("phone: %s, code: %s", phone, code)
	return setVerificationCode(context.Background(), phone, code, 5*time.Minute)
}

func VerifyUserCode(phone string, code string) error {
	if err := verifyAliyunSMS(context.Background(), phone, code); err != nil {
		if errors.Is(err, ErrCodeNotFound) {
			return fmt.Errorf("code expired or not found")
		} else if errors.Is(err, ErrCodeMismatch) {
			return fmt.Errorf("code mismatch")
		}
		return fmt.Errorf("unexpected error: %w", err)
	}

	// 验证通过，返回成功
	return nil
}
