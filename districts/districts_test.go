package districts

import (
	"testing"
	"districts"
)

func TestToRedis(t *testing.T) {
	/*areaId  := "1"
	level   := 1
	members := []string{"2", "20", "38"}

	if err := districts.ToRedis(areaId, level, members); err != nil {
		t.Fatalf("%v", err)
	}*/
}

func TestSource(t *testing.T) {
	_, _ = districts.Source()
}
