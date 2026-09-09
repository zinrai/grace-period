package main

import "testing"

// Reading the wrong column decides against the wrong number. Column eight is
// the expiry; the ones around it are the password ages
func TestDeadlineIsReadFromTheExpiryColumn(t *testing.T) {
	dir := t.TempDir()
	shadow := write(t, dir, "shadow", ""+
		"alice:!:20000:0:99999:7:::\n"+ // no expiry
		"bob:!:20000:0:99999:7::20500:\n"+ // expiry set
		"carol:!:20000:0:99999:7:30::\n"+ // inactive set, not expiry
		"truncated:!:20000\n")

	got, err := deadlines(shadow)
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]bool{"alice": false, "bob": true, "carol": false} {
		if got[name] != want {
			t.Errorf("%s: got %v, want %v", name, got[name], want)
		}
	}
	if _, ok := got["truncated"]; ok {
		t.Error("read a deadline out of a short line")
	}
}
