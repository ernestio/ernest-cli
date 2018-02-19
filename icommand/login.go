package icommand

import "github.com/abiosoft/ishell"

func loginICmd(shell *ishell.Shell) func(c *ishell.Context) {
	return func(c *ishell.Context) {
		// disable the '>>>' for cleaner same line input.
		c.ShowPrompt(false)
		defer c.ShowPrompt(true) // yes, revert after login.

		cfg, err := login(c)
		if err != nil {
			c.Println("Invalid credentials")
			return
		}

		userChain = append(userChain, cfg)
		updatePrompt(shell)
	}
}
