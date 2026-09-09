package main

import (
	"os"
	"path/filepath"
	"testing"
)

// A member listed nowhere is a member who never gets a deadline, and /etc/group
// names only the ones who joined on top of their own group
func TestMembersCoverBothWaysOfJoining(t *testing.T) {
	dir := t.TempDir()
	group := write(t, dir, "group", "ops:x:2000:alice,bob\nother:x:2001:carol\n")
	passwd := write(t, dir, "passwd", ""+
		"alice:x:3001:3001::/home/alice:/bin/bash\n"+
		"bob:x:3002:3002::/home/bob:/bin/bash\n"+
		"dave:x:3003:2000::/home/dave:/bin/bash\n"+ // ops is dave's own group
		"carol:x:3004:2001::/home/carol:/bin/bash\n")

	got, err := members(passwd, group, "ops")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alice", "bob", "dave"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

// An empty group would otherwise take every account whose gid field failed to
// parse along with it
func TestUnknownGroupHasNoMembers(t *testing.T) {
	dir := t.TempDir()
	group := write(t, dir, "group", "ops:x:2000:alice\n")
	passwd := write(t, dir, "passwd", "root:x:0:0::/root:/bin/bash\n")

	got, err := members(passwd, group, "nosuchgroup")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got %v", got)
	}
}

func write(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}
