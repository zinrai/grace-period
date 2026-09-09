package main

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

const (
	passwdPath = "/etc/passwd"
	groupPath  = "/etc/group"
)

// members lists everyone in the group, from both sides of the way membership is
// recorded.
//
// /etc/group holds only the users who joined the group on top of their own.
// Reading it alone would skip anyone the group is the primary group of, and
// skipping is the direction that leaves an account without a deadline
func members(passwdFile, groupFile, group string) ([]string, error) {
	gid, listed, err := groupEntry(groupFile, group)
	if err != nil {
		return nil, err
	}

	found := make(map[string]bool, len(listed))
	for _, u := range listed {
		found[u] = true
	}

	primary, err := primaryMembers(passwdFile, gid)
	if err != nil {
		return nil, err
	}
	for _, u := range primary {
		found[u] = true
	}

	out := make([]string, 0, len(found))
	for u := range found {
		out = append(out, u)
	}
	sort.Strings(out)
	return out, nil
}

func groupEntry(groupFile, group string) (gid string, listed []string, err error) {
	f, err := os.Open(groupFile)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		field := strings.Split(s.Text(), ":")
		if len(field) < 4 || field[0] != group {
			continue
		}
		if field[3] == "" {
			return field[2], nil, nil
		}
		return field[2], strings.Split(field[3], ","), nil
	}
	if err := s.Err(); err != nil {
		return "", nil, err
	}
	return "", nil, nil
}

func primaryMembers(passwdFile, gid string) ([]string, error) {
	if gid == "" {
		return nil, nil
	}

	f, err := os.Open(passwdFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		field := strings.Split(s.Text(), ":")
		if len(field) < 4 {
			continue
		}
		if field[3] == gid {
			out = append(out, field[0])
		}
	}
	return out, s.Err()
}
