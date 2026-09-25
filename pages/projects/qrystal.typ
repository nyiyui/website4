#set document(
  title: [Qrystal],
	author: ("Ken Shibata"),
	description: [Centralized WireGuard configuration utility],
	keywords: ("Networking"),
  date: datetime(year: 2025, month: 1, day: 22)
)

Qrystal sets up WireGuard tunnels on several devices from a central config (shown below) and an authn token on each device.
I use this to setup WireGuard tunnels between my servers (which are often on different cloud providers).

```nix
{
  qrystal.services.cs = {
    enable = true;
    config.central.networks.example = {
      listenPort = 12345;
      ips = [ "10.123.0.0/16" ];
    };
    config.tokens = [
      {
        name = "server0";
        hash = "hash-of-token";
        networks.example = "server0"; # name of server0 in "example" network
        canPull = true; # can server0 read the "example" network config?
      }
    ];
  };
}
```

#link("https://nyiyui.ca/qrystal")[Website]
