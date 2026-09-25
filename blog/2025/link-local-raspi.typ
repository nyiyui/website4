#set document(
  title: [Link-Local Connection to Raspberry Pi],
	author: ("Ken Shibata"),
	keywords: ("NetworkManager"),
  date: datetime(year: 2025, month: 2, day: 24),
)

_This is a living document; content will be updated as necessary (after the initial post date)._

I needed to have a Ethernet connection from my laptop to a Raspberry Pi without any extra hardware (e.g. a router).
Turns out this isn't as smooth as I expected (especially with NetworkManager), so the following is a list of steps to get it working.

A few notes:
- Steps below have been tested on a Raspberry Pi 3B with a ThinkPad X1 Carbon Gen 10 (running NixOS with NetworkManager).
- The Ethernet cable type (patch/crossover) doesn't seem to matter (I suppose this is because the USB-to-Ethernet adapter and Pi both support #link("https://en.wikipedia.org/wiki/Medium-dependent_interface")[Auto-MDI-X]).

== On The Raspberry Pi

Make a new connection profile that uses Link-Local configuration for IPv4.
You can either do this using the `nmtui` terminal interface, or through `nmcli`:

```sh
nmcli connection add type ethernet con-name "Link-Local Ethernet" ifname eth0 ipv4.method link-local ipv6.method ignore autoconnect yes
```

Note: If you're running the Pi headless, you should specify autoconnect yes, so NetworkManager retries without any interaction.

#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2025/nmtui-ll.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2025/nmtui-ll.webp"), type: "image/webp"))
  #html.img(src: "/2025/nmtui-ll.png", alt: "a new Ethernet connection profile with default settings, except IPv4 set to Link-Local and IPv6 set to Ignore", width: 1000, height: 563)
]


Make sure there are no other (competing) Ethernet connections (or they have a lower priority) (nmcli connection | grep -E 'TYPE|ethernet' or by using the TUI).
If there are, you can either delete them with nmcli connection delete <name> or change their priority with nmcli connection modify <name> connection.autoconnect-priority <priority>.

== On A Linux Laptop

Make a new connection profile that uses Link-Local configuration for IPv4.
Essentially, we do the same but our make sure you have the correct Ethernet interface: ifconfig should list all interfaces the device has.
If you're using a USB-C-to-Ethernet adapter, try plugging it into different USB ports.
In my case, the USB port that the adapter works in is not idempotent (!).

== On a Mac

Make sure to unplug and other Ethernet cables.
Then, plug in the Ethernet cable to the Raspberry Pi and your Mac.
macOS should prompt you to connect to the new network.

If it doesn't work, make sure there aren't any other network connections that are interfering:
go to Network Settings → Other Services and delete ones that you aren't using at this moment.
If you're prompted to "re-add on connection," you can click either option.

The same troubleshooting tip for Linux applies: try plugging it in again, and make sure to try other Ethernet ports!

