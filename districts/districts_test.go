package districts

import (
	"testing"
	"districts"
)

func TestSource(t *testing.T) {
	dists, err := districts.Source()
	if err != nil {
		t.Fatal(err)
	}
	if err := districts.ToRedis(dists); err != nil {
		t.Fatal(err)
	}
}
