#set document(
  title: [All-Wayland Setup on a Touchscreen-only Tablet],
	author: ("Ken Shibata"),
	keywords: ("Wayland"),
  date: datetime(year: 2025, month: 1, day: 30)
)

= The Premise

I managed to dig out a #link("https://www.lenovo.com/us/en/p/tablets/windows-tablets/thinkpad-series/thinkpad-x1-tablet-(3rd-gen)/22tp2cp0113")[ThinkPad X1 Tablet Gen 3].
Since the keyboard was broken, I decided to take this as a challenge to setup a usable Wayland environment on a touchscreen-only device.

This entry is composed of a semi-chronological description of each layer, starting from the BIOS to the window manager (desktop environment).
Although each section builds on previous ones, each can-be read separately.

= Secure Boot and Disk Encryption

Since a tablet is the child of a smartphone and a laptop, it is certainly intended for mobile usage.
Therefore, I needed to get disk encryption (to protect data from thieves) and Secure Boot (to protect the kernel from backdoors) working first.
After installing NixOS with disk encryption and my standard Sway-based layer-shell, I started on Secure Boot.
A while ago, I read that NixOS Secure Boot wasn't really ready nor usable. However, to my surprise the fantastic contributors nanaged to develop an easy solution.
I just followed the directions on the Lanzaboote README, and everything worked the first try (quite a surprise for Linux users like me)!

= Display Manager and Greeter

It was surprisingly smooth so far.
Now is where the real yak shaving starts.
The first problem was the display manager.
Most of the common options (lightdm, tuigreet, etc) didn't support on-screen keyboards.
Although there were a few that did, they did not work on the ThinkPad Tablet for one reason or another.
SDDM's keyboard just didn't show up, even when the button was pressed.
GDM's keyboard actually worked, but either Mutter or GTK (the problem was common on GNOME) had a problem with the touchscreen.
After a minute or so, it would not stop responding to touch input. However, mouse input worked perfectly.

While looking for other options, I found #link("https://sr.ht/~kennylevinsen/greetd/")[greetd].
Since it can launch basically any program as a greeter, I figured it could also run #link("https://git.sr.ht/~kennylevinsen/gtkgreet")[gtkgreet] with an on-screen keyboard, #link("https://github.com/jjsullivan5196/wvkbd")[wvkbd].
As it turns out, this worked great. Since I used Sway in the greeter too, I could also run #link("https://sr.ht/~kennylevinsen/wlsunset/")[wlsunset] to have "night light" on the login screen (especially useful when rebooting your tablet in the middle of the night…).

*An aside:* Although configuration tends not to be straightforward for low-layer things like display managers, greetd made it very simple - just run sway with a few startup commands:

```nix
{
  services.greetd = {
    enable = true;
    settings.default_session =
      let
        swayConfig = pkgs.writeText "greetd-sway-config" ''
          exec "${pkgs.greetd.gtkgreet}/bin/gtkgreet -l; swaymsg exit"
          exec "${pkgs.wvkbd}/bin/wvkbd-mobintl -L 256"
          exec "${pkgs.wlsunset}/bin/wlsunset -L -79.38 -T 6500 -g 1.000000 -l 43.65 -t 2000"
        '';
      in
      {
        command = "${pkgs.sway}/bin/sway --config ${swayConfig}";
        user = "greeter";
      };
  };
  environment.etc."greetd/environments" = {
    enable = true;
    text = ''
      sway
    '';
  };
}
```

= Window Manager

Now, we've gotten everything until the window manager working well on touch.
Although I've already used Sway for the greeter, now it has to handle moving windows around, and not just have a single layout.
As it turns out, Sway's dragging doesn't really work well on touch (the touch area is not extended beyond the visual border), and gestures weren't really a thing.

I figured I don't use it in tablet mode when managing multiple windows (as that usually means I'm developing something with a keyboard).
Maybe in a future blog post...

== Touchscreen Gestures using lisgd

#link("https://sr.ht/~mil/lisgd/")[lisgd] is a nice little tool that reads touchscreen input events (from e.g. `/dev/input/event9`) and runs commands based on them.
Excited, I added swipe gestures to move workspaces:
```sh
lisgd -d '/dev/input/event9' -v -m 800 \
  -g '1,LR,L,*,*,swaymsg workspace prev' \
  -g '1,RL,R,*,*,swaymsg workspace next'
```
And it worked...until the next reboot.
As it turns out, the path was reset every boot, and the touchscreen device would be `/dev/input/event9` one boot, and `/dev/input/event13` the next.

I ended up running `libinput list-devices` during startup to find the touchscreen by name:
```sh
device=$(libinput list-devices | grep -A 1 '${cfg.touchscreenDeviceName}' | grep /dev/input | awk '{print $2}')
lisgd -d "$device" -v -m 800 \
  -g '1,LR,L,*,*,swaymsg workspace prev' \
  -g '1,RL,R,*,*,swaymsg workspace next'
```
It's not...the most elegant solution, but it works.

I also added a swipe from the top to close a window.
The final command is:
```sh
device=$(libinput list-devices | grep -A 1 '${cfg.touchscreenDeviceName}' | grep /dev/input | awk '{print $2}')
lisgd -d "$device" -v -m 800 \
  -g '1,LR,L,*,*,swaymsg workspace prev' \
  -g '1,RL,R,*,*,swaymsg workspace next' \
  -g '1,UD,U,*,*,swaymsg kill'
```

Package this as a user service, and you're good to go!
