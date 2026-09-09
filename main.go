// grace-period gives every member of a group a deadline to finish a setup step,
// and takes the deadline away once they have.
//
// Both ends of this already exist on a Linux host: an account carries an expiry
// date, and pam_unix refuses a login past it. What is missing is anything that
// ties the deadline to whether the setup was actually done. Configuration
// management cannot express it either, because "creation plus seven days"
// becomes a fresh date on every run and never arrives
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	group := flag.String("group", "", "group whose members are given a deadline")
	markerDir := flag.String("marker-dir", "", "directory holding one file per user who finished")
	days := flag.Int("days", 0, "days a member is given")
	dryRun := flag.Bool("dry-run", false, "report what would change, change nothing")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	// A default here would hand a deadline to the wrong group, or read the
	// wrong directory and hand one to everybody
	if *group == "" || *markerDir == "" || *days <= 0 || flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	if err := reconcile(*group, *markerDir, *days, *dryRun); err != nil {
		log.Fatal(err)
	}
}
