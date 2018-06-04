// +build integration

package archiveorg

import (
	"testing"
)

func TestTimeMapFor(t *testing.T) {
	u := "https://jaytaylor.com/"

	timemap, err := TimeMapFor(u)
	if err != nil {
		t.Fatalf("Error getting TimeMap for url=%v: %s", u, err)
	}

	t.Logf("num mementos: %v", len(timemap.Mementos))
	t.Logf("timemap: %+v", timemap)

	if err := validateMemento(timemap.Original); err != nil {
		t.Errorf("Error validating Original: %v", err)
	}

	if err := validateMemento(timemap.Self); err != nil {
		t.Errorf("Error validating Self: %v", err)
	}

	if err := validateMemento(timemap.TimeGate); err != nil {
		t.Errorf("Error validating TimeGate: %v", err)
	}

	for i, memento := range timemap.Mementos {
		if err := validateMemento(&memento); err != nil {
			t.Errorf("[i=%v] Error validating Memento slice element: %v", i, err)
		}
	}
}
