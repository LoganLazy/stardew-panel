package handler

import "testing"

func TestValidateServerSetting(t *testing.T) {
	valid := []struct {
		key   string
		value interface{}
	}{
		{"maxPlayers", float64(4)},
		{"autoStartNewDay", true},
		{"newDayWaitSeconds", float64(30)},
	}
	for _, tc := range valid {
		if err := validateServerSetting(tc.key, tc.value); err != nil {
			t.Errorf("valid setting rejected: %s: %v", tc.key, err)
		}
	}

	invalid := []struct {
		key   string
		value interface{}
	}{
		{"maxPlayers", float64(17)},
		{"maxPlayers", 4},
		{"autoStartNewDay", "true"},
		{"newDayWaitSeconds", float64(-1)},
		{"unknown", true},
	}
	for _, tc := range invalid {
		if err := validateServerSetting(tc.key, tc.value); err == nil {
			t.Errorf("invalid setting accepted: %s=%v", tc.key, tc.value)
		}
	}
}
