package params

import (
	libparams "github.com/qdm12/golibs/params"
)

// GetResolveHostnames obtains 'yes' or 'no' to resolve hostnames in order to obtain
// more IP addresses, from the environment variable RESOLVE_HOSTNAMES, and defaults to no.
// If you are blocking the hostname resolution on your network, turn this feature off.
func (p *Getter) GetResolveHostnames() (resolveHostnames bool, err error) {
	return p.envParams.YesNo("RESOLVE_HOSTNAMES", libparams.Default("no"))
}
