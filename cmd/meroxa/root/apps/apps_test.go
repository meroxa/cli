package apps

import "testing"

func processError(t *testing.T, given error, wanted error) {
	if given != nil {
		if wanted == nil {
			t.Fatalf("unexpected error \"%s\"", given)
		} else if wanted.Error() != given.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", wanted, given)
		}
	}
}
