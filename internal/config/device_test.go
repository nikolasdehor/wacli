package config

import (
	"testing"

	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
)

func TestApplyDeviceConfigFromEnv(t *testing.T) {
	// Reset store state before each test if possible, but store global state is persistent.
	// Since tests run sequentially (unless t.Parallel()), we rely on setting env vars.
	// Note: We cannot easily reset store globals completely without manual resets.
	// We will just verify specific changes.

	t.Run("Default", func(t *testing.T) {
		t.Setenv("WACLI_DEVICE_LABEL", "")
		t.Setenv("WACLI_DEVICE_PLATFORM", "")
		
		// Reset to default manually for test baseline
		store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()
		store.SetOSInfo("Go", [3]uint32{0, 1, 0}) // default in whatsmeow

		ApplyDeviceConfigFromEnv()

		if store.DeviceProps.GetPlatformType() != waCompanionReg.DeviceProps_CHROME {
			t.Errorf("expected default platform CHROME, got %v", store.DeviceProps.GetPlatformType())
		}
	})

	t.Run("OverridePlatform", func(t *testing.T) {
		t.Setenv("WACLI_DEVICE_PLATFORM", "SAFARI")
		ApplyDeviceConfigFromEnv()

		if store.DeviceProps.GetPlatformType() != waCompanionReg.DeviceProps_SAFARI {
			t.Errorf("expected platform SAFARI, got %v", store.DeviceProps.GetPlatformType())
		}
	})

	t.Run("OverrideLabel", func(t *testing.T) {
		label := "TestDevice"
		t.Setenv("WACLI_DEVICE_LABEL", label)
		ApplyDeviceConfigFromEnv()

		if store.DeviceProps.GetOs() != label {
			t.Errorf("expected OS label %q, got %q", label, store.DeviceProps.GetOs())
		}
		
		// Verify protobuf fields
		if got := store.BaseClientPayload.UserAgent.GetDevice(); got != label {
			t.Errorf("expected UserAgent.Device %q, got %q", label, got)
		}
		if got := store.BaseClientPayload.UserAgent.GetManufacturer(); got != label {
			t.Errorf("expected UserAgent.Manufacturer %q, got %q", label, got)
		}
	})

	t.Run("InvalidPlatformDefaultsToChrome", func(t *testing.T) {
		t.Setenv("WACLI_DEVICE_PLATFORM", "INVALID_PLATFORM_XYZ")
		// Manually set to something else first to ensure change happens (or stays default)
		// Wait, if it's invalid, it should reset to Chrome? Or just ignore?
		// Logic: if platformRaw != "", it calls parsePlatformType.
		// parsePlatformType returns CHROME for invalid input.
		// So it should reset to CHROME.
		
		// Set to IE first
		store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_IE.Enum()
		
		ApplyDeviceConfigFromEnv()

		if store.DeviceProps.GetPlatformType() != waCompanionReg.DeviceProps_CHROME {
			t.Errorf("expected fallback to CHROME for invalid platform, got %v", store.DeviceProps.GetPlatformType())
		}
	})
}

func TestParsePlatformType(t *testing.T) {
	tests := []struct {
		input    string
		expected waCompanionReg.DeviceProps_PlatformType
	}{
		{"", waCompanionReg.DeviceProps_CHROME},
		{"chrome", waCompanionReg.DeviceProps_CHROME},
		{"CHROME", waCompanionReg.DeviceProps_CHROME},
		{"safari", waCompanionReg.DeviceProps_SAFARI},
		{"FIREFOX", waCompanionReg.DeviceProps_FIREFOX},
		{"ie", waCompanionReg.DeviceProps_IE},
		{"unknown", waCompanionReg.DeviceProps_UNKNOWN},
		{"GARBAGE_VALUE", waCompanionReg.DeviceProps_CHROME},
		{"  safari  ", waCompanionReg.DeviceProps_SAFARI},
	}

	for _, tc := range tests {
		got := parsePlatformType(tc.input)
		if got != tc.expected {
			t.Errorf("parsePlatformType(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}
