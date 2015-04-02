package bot

import "github.com/cenkalti/rpc2"

func (c *IRCRPC) Connect(args *VoidArgs, reply *bool) error {
	if err := c.Client.Connect(); err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

func (c *IRCRPC) Disconnect(args *VoidArgs, reply *string) error {
	c.Client.Quit("bye bye")
	return nil
}

func (c *IRCRPC) Announce(src *rpc2.Client, msg *PrivMsgArgs, reply *bool) error {
	c.Client.Privmsg(msg.Target, msg.Text)
	*reply = true
	return nil
}

func (c *IRCRPC) Join(src *rpc2.Client, channel string, reply *bool) error {
	c.Client.Join(channel)
	*reply = true
	return nil
}

// plugins register to receive privmsg
func (c *IRCRPC) Register(src *rpc2.Client, _ *struct{}, reply *bool) error {
	c.Listeners[src] = true
	*reply = true
	return nil
}
