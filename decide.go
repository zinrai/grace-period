package main

type action int

const (
	leave action = iota
	clear
	set
)

func (a action) String() string {
	switch a {
	case clear:
		return "clearing the deadline"
	case set:
		return "setting the deadline"
	}
	return ""
}

// decide reconciles two pieces of state that already exist on the host: the
// file that shows the setup was done, and the expiry field the account carries.
//
// Nothing is remembered between runs, and nobody is locked out here. An account
// past its expiry is refused by pam_unix at the next login
func decide(hasMarker, hasDeadline bool) action {
	switch {
	case hasMarker && hasDeadline:
		return clear
	case !hasMarker && !hasDeadline:
		return set
	}
	return leave
}
