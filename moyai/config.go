package moyai

import "github.com/df-mc/dragonfly/server"

type Config struct {
	server.UserConfig
	Nether, End struct {
		Folder string
	}

	// Moyai contains fields specific to moyai
	Moyai struct {
		// Tebex is the Tebex API key.
		Tebex string
		// Whitelisted is true if the server is whitelisted.
		Whitelisted bool
		// Season is the current season of the server.
		Season int
		// Start is the date the season started.
		Start string
		// End is the date the season ends.
		End string
	}

	Proxy struct {
		Enabled bool
		Name    string
		Secret  string
	}

	Oomph struct {
		Enabled bool
		// Other shit
	}
}

// DefaultConfig returns a default config for the server.
func DefaultConfig() Config {
	c := Config{}
	c.UserConfig = server.DefaultConfig()
	c.Moyai.Whitelisted = true
	c.Moyai.Season = 1
	c.Moyai.Start = "Edit this!"
	c.Moyai.End = "Edit this!"
	c.Proxy.Enabled = true
	c.Oomph.Enabled = true
	return c
}
