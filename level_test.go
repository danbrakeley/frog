package frog

import "testing"

func Test_LevelStrings(t *testing.T) {
	usedLevels := make(map[string]int)
	for l := levelMin; l < levelMax; l++ {
		str := l.String()
		if len(str) == 0 {
			t.Errorf("Level %d requires a non-empty string", int(l))
		}

		usedBy, inUse := usedLevels[str]
		if inUse {
			t.Errorf("Level %d has a non-unique string value (level %d also uses \"%s\")", int(l), usedBy, str)
		}
		usedLevels[str] = int(l)
	}
}
