package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"time"
)

const shadowPath = "/etc/shadow"

// deadlines reads the expiry field of every local account.
//
// chage -l would have to be parsed out of prose that changes with the locale,
// and it holds the same number this field does
func deadlines(shadowFile string) (map[string]bool, error) {
	f, err := os.Open(shadowFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	set := make(map[string]bool)
	s := bufio.NewScanner(f)
	for s.Scan() {
		field := strings.Split(s.Text(), ":")
		if len(field) < 8 {
			continue
		}
		set[field[0]] = field[7] != ""
	}
	return set, s.Err()
}

// The deadline is an expiry date, not a locked password. usermod -L only marks
// the password hash, which leaves public key authentication working, and public
// keys are how accounts like these are reached.
//
// chage does the writing rather than this file being edited in place. Editing it
// would race with any other tool holding the lock, and losing that race costs
// every account in the file
func setDeadline(name string, on time.Time) error {
	return exec.Command("chage", "-E", on.Format("2006-01-02"), "--", name).Run()
}

func clearDeadline(name string) error {
	return exec.Command("chage", "-E", "-1", "--", name).Run()
}
