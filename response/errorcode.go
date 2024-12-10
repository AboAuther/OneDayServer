package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/golang-cz/nilslice"
)

var (
	// General Errors
	InternalServerError = ErrorCode{
		Code:           "10000",
		Message:        "Internal server error",
		HttpStatusCode: 500,
	}
	BadGateway = ErrorCode{
		Code:           "10001",
		Message:        "Bad Gateway",
		HttpStatusCode: 502,
	}
	UnableToConnectToServer = ErrorCode{
		Code:           "10002",
		Message:        "Unable to connect to server",
		HttpStatusCode: 503,
	}
	RequestTimestampFarFromServerTime = ErrorCode{
		Code:           "10003",
		Message:        "Request timestamp is far from the server time",
		HttpStatusCode: 504,
	}
	Timeout = ErrorCode{
		Code:           "10004",
		Message:        "Timeout",
		HttpStatusCode: 504,
	}
	UnauthorizedApiKey = ErrorCode{
		Code:           "10005",
		Message:        "Unauthorized API key",
		HttpStatusCode: 401,
	}

	UnauthorizedUserPassword = ErrorCode{
		Code:           "10006",
		Message:        "Unauthorized user password",
		HttpStatusCode: 401,
	}

	UnauthorizedJWTAccessToken = ErrorCode{
		Code:           "10007",
		Message:        "Invalid or expired access token",
		HttpStatusCode: 401,
	}

	ApiNotFound = ErrorCode{
		Code:           "10008",
		Message:        "API not found",
		HttpStatusCode: 404,
	}
	RequestMethodNotAllowed = ErrorCode{
		Code:           "10009",
		Message:        "Request method is not allowed",
		HttpStatusCode: 405,
	}
	IpAddressRateLimitExceeded = ErrorCode{
		Code:           "10010",
		Message:        "IP address rate limit exceeded",
		HttpStatusCode: 429,
	}
	UserRateLimitExceeded = ErrorCode{
		Code:           "10011",
		Message:        "User rate limit exceeded",
		HttpStatusCode: 429,
	}
	InternalMaintaining = ErrorCode{
		Code:           "10012",
		Message:        "Internal maintaining",
		HttpStatusCode: 500,
	}

	UnauthorizedJWTRefreshToken = ErrorCode{
		Code:           "100014",
		Message:        "Invalid or expired refresh token",
		HttpStatusCode: 401,
	}

	InvalidJWTTokenClaims = ErrorCode{
		Code:           "10015",
		Message:        "Invalid jwt token claims",
		HttpStatusCode: 401,
	}

	InvalidJWTTokenFormat = ErrorCode{
		Code:           "10016",
		Message:        "Invalid jwt token format, must start with 'Bearer '",
		HttpStatusCode: 401,
	}

	UnauthorizedJWTAccessTokenExpired = ErrorCode{
		Code:           "10016",
		Message:        "Invalid jwt token, access token expired",
		HttpStatusCode: 401,
	}

	SMSDeliveryFailed = ErrorCode{
		Code:           "10017",
		Message:        "SMS delivery failed",
		HttpStatusCode: 500,
	}
	InvalidOrExpiredVerificationCode = ErrorCode{
		Code:           "10018",
		Message:        "Invalid or expired verification code",
		HttpStatusCode: 401,
	}

	// Header and Parameter Errors
	MissingRequiredHeader = ErrorCode{
		Code:           "11001",
		Message:        "Missing required header '%s'",
		HttpStatusCode: 400,
	}
	MissingRequiredParameter = ErrorCode{
		Code:           "11002",
		Message:        "Missing required parameter '%s'",
		HttpStatusCode: 400,
	}
	ParameterHeaderNameEmpty = ErrorCode{
		Code:           "11003",
		Message:        "Header '%s' was empty",
		HttpStatusCode: 400,
	}
	ParameterParamNameEmpty = ErrorCode{
		Code:           "11004",
		Message:        "Parameter '%s' was empty",
		HttpStatusCode: 400,
	}
	InvalidParamNameType = ErrorCode{
		Code:           "11005",
		Message:        "Invalid '%s'. Value must be of type (%s)",
		HttpStatusCode: 400,
	}
	InvalidParamNameEnumValues = ErrorCode{
		Code:           "11006",
		Message:        "Invalid '%s'. Value not in allowed enum values [%s]",
		HttpStatusCode: 400,
	}
	InvalidParamNameMaxValue = ErrorCode{
		Code:           "11007",
		Message:        "Invalid '%s'. Value exceeds maximum allowed (%s)",
		HttpStatusCode: 400,
	}
	InvalidParamNameMinValue = ErrorCode{
		Code:           "11008",
		Message:        "Invalid '%s'. Value does not reach minimum required (%s)",
		HttpStatusCode: 400,
	}
	InvalidParamNameExpectedPrefix = ErrorCode{
		Code:           "11009",
		Message:        "Invalid '%s'. Value must start with '%s'",
		HttpStatusCode: 400,
	}
	InvalidParamInsufficientAccuracy = ErrorCode{
		Code:           "11010",
		Message:        "Invalid '%s'. '%s' must be an integral multiple of (%s)",
		HttpStatusCode: 400,
	}
	UnknownParameterParamName = ErrorCode{
		Code:           "11011",
		Message:        "Unknown parameter '%s'",
		HttpStatusCode: 400,
	}
	InvalidRequestBody = ErrorCode{
		Code:           "11012",
		Message:        "Invalid request body",
		HttpStatusCode: 400,
	}

	UserNotFound = ErrorCode{
		Code:           "11013",
		Message:        "User not found",
		HttpStatusCode: 400,
	}
)

type ErrorCode struct {
	Code           string
	Message        string
	HttpStatusCode int `json:"-"`
}

// GetMessage 获取格式化的错误消息
func (e *ErrorCode) GetMessage(params ...interface{}) string {
	template := e.Message

	// 检查是否有参数以及模板是否包含占位符
	if len(params) > 0 {
		// 使用 fmt.Sprintf 来格式化消息
		return fmt.Sprintf(template, params...)
	}

	// 如果没有参数，则返回原始消息模板
	return template
}

// ToJSON 将错误代码转换为 JSON 格式
func (e *ErrorCode) ToJSON(params ...interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"code":    e.Code,
		"message": e.GetMessage(params...),
	}
	return result
}

// ToJSONString 将错误代码转换为 JSON 字符串格式
func (e *ErrorCode) ToJSONString(params ...interface{}) string {
	// 将结果对象编码为 JSON 字符串
	jsonBytes, err := json.Marshal(e.ToJSON(params...))
	if err != nil {
		fmt.Println("Failed to encode JSON:", err)
		return ""
	}
	return string(jsonBytes)
}

func Fail(code ErrorCode, params ...interface{}) interface{} {
	return code.ToJSON(params...)
}

func SendError(c *gin.Context, code ErrorCode, params ...interface{}) {
	c.JSON(code.HttpStatusCode, Fail(code, params...))
	c.Abort()
}

func SendInternalServerError(c *gin.Context) {
	SendError(c, InternalServerError)
}

func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, nilslice.Initialize(data))
}

func SendSuccessMessage(c *gin.Context) {
	c.JSON(http.StatusOK, nilslice.Initialize(map[string]interface{}{
		"message": "success",
	}))
}

func SendSuccessWithAny(c *gin.Context, data any) {
	result := make(map[string]interface{})

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	switch dataType.Kind() {
	case reflect.Struct:
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			fieldValue := dataValue.Field(i)
			tagName := field.Tag.Get("json")
			if tagName == "" {
				tagName = field.Name
			}
			result[tagName] = fieldValue.Interface()
		}
	case reflect.Map:
		result = data.(map[string]interface{})
	default:
		result["data"] = data
	}

	result["error"] = false
	c.JSON(http.StatusOK, result)
}
