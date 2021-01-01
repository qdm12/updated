package dnscrypto

import "errors"

var (
	ErrBadStatusCode    = errors.New("bad HTTP status code")
	ErrChecksumMismatch = errors.New("checksum does not match")
)
