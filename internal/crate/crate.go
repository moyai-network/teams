package crate

import "github.com/moyai-network/moose/crate"

func init() {
	crate.Register(koth{})
	crate.Register(partner{})
	crate.Register(pharaoh{})
	crate.Register(menes{})
	crate.Register(ramses{})
}
