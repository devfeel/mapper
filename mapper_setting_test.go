package mapper

import "testing"

func TestNewSetting(t *testing.T) {
	checkEnabledTypeChecking := true
	checkEnabledAutoTypeConvert := false

	setting := NewSetting(CTypeChecking(checkEnabledTypeChecking), CAutoTypeConvert(checkEnabledAutoTypeConvert))
	if setting.EnabledTypeChecking != checkEnabledTypeChecking || setting.EnabledAutoTypeConvert != checkEnabledAutoTypeConvert {
		t.Error("NewSetting error: [", checkEnabledTypeChecking, ",", setting.EnabledTypeChecking, "],[", checkEnabledAutoTypeConvert, ",", setting.EnabledAutoTypeConvert, "]")
	} else {
		t.Log("NewSetting success")
	}
}
