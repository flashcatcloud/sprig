package sprig

import "testing"

func Test_normalizeMessage(t *testing.T) {
	tpl := `{{normalizeMessage "ConnID{id='6f2f2bac-ae63-4d32-83e7-d07ce14ac537', clusterType=SERVICE_DISCOVER_CLUSTER"}}`
	expected := `ConnID{id='{HASH}', clusterType=SERVICE_DISCOVER_CLUSTER`
	if err := runt(tpl, expected); err != nil {
		t.Error(err)
	}
}
