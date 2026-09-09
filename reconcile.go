package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// reconcile brings the deadline on each account into line with whether the
// setup was done. Nothing is remembered between runs; the two states it reads
// are the whole of what it knows
func reconcile(group, markerDir string, days int, dryRun bool) error {
	users, err := members(passwdPath, groupPath, group)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		log.Printf("no members in %s", group)
		return nil
	}

	deadline, err := deadlines(shadowPath)
	if err != nil {
		return err
	}

	on := time.Now().AddDate(0, 0, days)
	for _, u := range users {
		a := decide(hasMarker(markerDir, u), deadline[u])
		if a == leave {
			continue
		}
		log.Printf("%s: %s", u, a)
		if dryRun {
			continue
		}
		if err := apply(a, u, on); err != nil {
			return err
		}
	}
	return nil
}

// hasMarker asks the directory, not the user and not the account. Whatever the
// setup step is, finishing it has to leave a file behind named after the user
func hasMarker(markerDir, name string) bool {
	_, err := os.Stat(filepath.Join(markerDir, name))
	return err == nil
}

func apply(a action, name string, on time.Time) error {
	if a == clear {
		return clearDeadline(name)
	}
	return setDeadline(name, on)
}
