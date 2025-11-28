package snowflake

import "testing"

func TestGenerator_NextID(t *testing.T) {
	sf, err := New(1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(sf.NextID())
	t.Log(sf.NextString())
}
