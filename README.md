# home

This is a server and client for interacting with the Trådfri lamps from IKEA. This is just a wrapper around the code from [eriklupander/tradfri-go](https://github.com/eriklupander/tradfri-go). You might be better of by checking out that project first.

The motivation behind this project is that I want an easy way to turn my lights on and off grouped by room and name.

To run the server and/or client you need a config file named `gohome.yaml` in the directory `$XDG_CONFIG_HOME` or `~/.config`.

Example gohome.yaml
```yaml
server:
    host: "192.168.1.35"
    port: 50051

ikea:
    gateway: "192.168.1.28:5684"
    clientID: "uranium"
    psk: "your-psk-goes-here"
```

The PSK can you get by using [eriklupander/tradfri-go](https://github.com/eriklupander/tradfri-go#user-content-psk-exchange).

The client is used like this:

```text
Usage:

	home-cli <room> [<lamp>] <action> [value]

	<room> = room name or all (if all, skip <lamp>
	<lamp> = lamp name or all
	<action> =
		On
		Off
		Toggle
		Dim = [1-254]
```

So you can turn off all bedroom lamps with `home-cli bedroom all off`
