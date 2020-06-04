package db

import (
	"encoding/base64"
	"fmt"
)

// DSN ...
type DSN struct {
	Name   string
	Scheme string
	User   *Userinfo // username and password information
	Host   *Host     // host or host:port
	Path   string    // path (relative paths may omit leading slash)
}

// String reassembles the URL into a valid URL string.
// The general form of the result is one of:
//
//	scheme:opaque?query#fragment
//	scheme://userinfo@host/path?query#fragment
//
// If u.Opaque is non-empty, String uses the first form;
// otherwise it uses the second form.
// Any non-ASCII characters in host are escaped.
// To obtain the path, String uses u.EscapedPath().
//
// In the second form, the following rules apply:
//	- if u.Scheme is empty, scheme: is omitted.
//	- if u.User is nil, userinfo@ is omitted.
//	- if u.Host is empty, host/ is omitted.
//	- if u.Scheme and u.Host are empty and u.User is nil,
//	   the entire scheme://userinfo@host/ is omitted.
//	- if u.Host is non-empty and u.Path begins with a /,
//	   the form host/path does not add its own /.
//	- if u.RawQuery is empty, ?query is omitted.
//	- if u.Fragment is empty, #fragment is omitted.
func (u *DSN) String() string {
	if u.Scheme == "" {
		u.Scheme = "tcp"
	}
	// for now ....
	// return fmt.Sprintf("root:pass@tcp(127.0.0.1:3306)/%s", "db01")
	return fmt.Sprintf("%s@%s(%s)/%s", u.User.String(), u.Scheme, u.Host.String(), u.Path)
}

// Key ...
func (u *DSN) Key() string {
	keyData := fmt.Sprintf("%s@%s/%s", u.User.Username, u.Host.String(), u.Path)
	return base64.StdEncoding.EncodeToString([]byte(keyData))
}

// The Userinfo type is an immutable encapsulation of username and
// password details for a URL. An existing Userinfo value is guaranteed
// to have a username set (potentially empty, as allowed by RFC 2396),
// and optionally a password.
type Userinfo struct {
	Username    string
	Password    string
	passwordSet bool
}

// String returns the encoded userinfo information in the standard form
// of "username[:password]".
func (u *Userinfo) String() string {
	return fmt.Sprintf("%s:%s", u.Username, u.Password)
}

// Host ...
type Host struct {
	Host string
	Port string
}

// String returns the encoded userinfo information in the standard form
// of "username[:password]".
func (h *Host) String() string {
	return fmt.Sprintf("%s:%s", h.Host, h.Port)
}

// NewDSN ...
func NewDSN(username, password, host, port, path string) *DSN {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	return &DSN{
		User: &Userinfo{
			Username: username,
			Password: password,
		},
		Host: &Host{
			Host: host,
			Port: port,
		},
		Path: path,
	}
}
