# grace-period

Give the members of a group a deadline to finish a setup step, and take the
deadline away once they have.

A member who has not finished is given an account expiry date. Past it,
`pam_unix` refuses the login.

```
Your account has expired; please contact your system administrator.
```

## Usage

Run it as root, from a timer.

```
grace-period -group ops -marker-dir /etc/google-authenticator -days 7
```

`-marker-dir` names a directory holding one file per user who finished, named
after them. Whatever the setup step is, it has to leave that file behind.

| marker file | expiry date | what happens |
| --- | --- | --- |
| absent | not set | the deadline is set to today plus `-days` |
| absent | set | nothing; the deadline is already running |
| present | set | the deadline is removed |
| present | not set | nothing |

Members are read from `/etc/group` and from the primary group field in
`/etc/passwd`.

`-dry-run` reports what would change without changing anything. Worth running
once before the timer is enabled, since a wrong `-marker-dir` hands a deadline
to everybody.

## From a timer

Unit files are in `systemd/`. Fill in the group and the directory in
`ExecStart`, then

```
cp systemd/grace-period.* /etc/systemd/system/
systemctl enable --now grace-period.timer
```

## Undoing a deadline

```
chage -E -1 alice
```

The next run sets it again unless the marker file is there, so the user has to
finish the setup before then.

## License

This project is licensed under the [MIT License](LICENSE).
