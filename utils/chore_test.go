package utils

import (
	"testing"

	"github.com/google/uuid"
)

// TestGenerateUUID tests the GenerateUUID function
func TestGenerateUUID(t *testing.T) {
	uuidStr := GenerateUUID()
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("Generated UUID is invalid: %s", uuidStr)
	}
}

// TestFormatHexString tests the FormatHexString function
func TestFormatHexString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123abc", "0x123abc"},
		{"0x123abc", "0x123abc"},
	}

	for _, test := range tests {
		actual := FormatHexString(test.input)
		if actual != test.expected {
			t.Errorf("FormatHexString(%s) = %s; expected %s", test.input, actual, test.expected)
		}
	}
}

// TestCleanHexString tests the CleanHexString function
func TestCleanHexString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0x123abc", "123abc"},
		{"123abc", "123abc"},
	}

	for _, test := range tests {
		actual := CleanHexString(test.input)
		if actual != test.expected {
			t.Errorf("CleanHexString(%s) = %s; expected %s", test.input, actual, test.expected)
		}
	}
}

//	func TestCamelToSnake(t *testing.T) {
//		tests := []struct {
//			input    string
//			expected string
//		}{
//			{"orderConfirmation", "order_confirmation"},
//			{"cancelAllConfirmation", "cancel_all_confirmation"},
//			{"tradeNotification", "trade_notification"},
//			{"clickDepthInputAmount", "click_depth_input_amount"},
//			{"priceTime24h", "price_time_24h"},
//		}
//
//		for _, test := range tests {
//			actual := CamelToSnake(test.input)
//			if actual != test.expected {
//				t.Errorf("CamelToSnake(%s) = %s; expected %s", test.input, actual, test.expected)
//			}
//		}
//	}
//
//	func TestStructToSnakeCaseMap(t *testing.T) {
//		type Demo struct {
//			OrderConfirmation     bool   `json:"orderConfirmation"`
//			CancelAllConfirmation bool   `json:"cancelAllConfirmation"`
//			TradeNotification     bool   `json:"tradeNotification"`
//			ClickDepthInputAmount bool   `json:"clickDepthInputAmount"`
//			PriceTime24H          string `json:"priceTime24H"`
//		}
//
//		input := Demo{
//			OrderConfirmation:     true,
//			CancelAllConfirmation: false,
//			TradeNotification:     true,
//			ClickDepthInputAmount: false,
//			PriceTime24H:          "24:00",
//		}
//
//		expected := map[string]interface{}{
//			"order_confirmation":       true,
//			"cancel_all_confirmation":  false,
//			"trade_notification":       true,
//			"click_depth_input_amount": false,
//			"price_time_24h":           "24:00",
//		}
//
//		actual, err := StructToSnakeCaseMap(input)
//		if err != nil {
//			t.Fatalf("StructToSnakeCaseMap returned an error: %v", err)
//		}
//
//		if !reflect.DeepEqual(actual, expected) {
//			t.Errorf("StructToSnakeCaseMap returned %v; expected %v", actual, expected)
//		}
//	}
