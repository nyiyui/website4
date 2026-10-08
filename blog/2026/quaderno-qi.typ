#set document(
  title: [I added Wireless Charging to my A4 E-Reader],
  author: ("Ken Shibata"),
  description: [How I (jankily) added wireless charging to my Quaderno e-reader off-the-shelf wireless charging receiver. Also discuss charging performance and challenges.],
  keywords: ("Tinkering with Hardware"),
  date: datetime(year: 2026, month: 10, day: 7),
)

== Some Motivation

I have a #link("https://www.fmworld.net/digital-paper/gen2/product.html")[Fujitsu Quaderno A4 (Gen.2)] e-reader that I've been using for a few years now.
It's light and thin enough that it feels like ten sheets of cardstock and has a paper-like textured screen.
There's also an #link("https://github.com/janten/dpt-rp1-py")[open-source library] which can sync a local directory with it (Thanks to everyone who worked on it!).
I take notes#footnote[I think it's a great for digital notes at meetings without seeming inattentive. The third person can see that I'm not reading an email while pretending to take notes.], reading papers and books, and to quickly sketch designs with it, and it mostly works great (ignoring the closed-source firmware).

However, it's not without downsides.
I usually use up the battery in a few days. That is, not a constant number of days but sometimes 3, often 4, and rarely 7.
Having to unlock it and see how much charge it's at or having to charge it every night is cumbersome not counting the following reason.
More importantly, the device's build quality is not great. It's bendy and a thin plastic piece around the USB-C port broke off.

== Disassembly

#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-disas.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-disas.webp"), type: "image/webp"))
  #html.img(src: "/2026/quaderno-disas.jpg", alt: "Taking off the back (plastic) cover off the Quaderno; two credit cards on the top corners of the image (the tablet is flipped upside down). The a piece of plastic around the USB-C port is broken.", width: 1024, height: 1180)
]

Opening the Quaderno up seems quite scary (given how frail it feels), but it's actually quite easy to losslessly disassemble.
Slide in a couple credit cards (three to be safe) in the bottom corners (and middle).
Then, making sure that the entire bottom few centimetres of the plastic are out of the way, smear the display down and off the plastic cover.
After a few centimetres more, the two halves should come apart.

The electronics are all on the display side.

#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-internals.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-internals.webp"), type: "image/webp"))
  #html.img(src: "/2026/quaderno-internals.jpg", alt: "The internals! There is a mainboard, I/O board (good for repair!), reemingly two batteries, and a bunch of cables. The display side is diecast and the back cover ir basically all white plastic.", width: 1024, height: 663)
]

#html.elem("aside")[
  I found it interesting that it seems to have two batteries (the one annotated in cyan and the big and obvious one). Not sure what's it's used for…
  
  Also, the I/O board is separate from the main one, so it seems like the flimsy USB-C connector can be replaced. For me, I can take it off and fiddle with it when I'm adding the wireless charger for real.

  As another side note, the two modules besides the I/O board are the Wi-Fi and Bluetooth antennas, and the NFC antenna is to the left of that.
]

== Adding the Receiver Coil

I got #link("https://www.amazon.com/dp/B08RDMRZL5")[a set of Qi receiver coils] with a PCB that outputs DC 5V.
I sacrificed a USB-A to USB-C cable, and took the power lines and soldered it to the coils'.

// quaderno-coil-battery.avif
#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-coil-battery.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-coil-battery.webp"), type: "image/webp"))
  #html.img(src: "/2026/quaderno-coil-battery.jpg", alt: "The diecast part of the tablet, with the coil on top of the battery. It's all held together with masking tape…masking tape fixes everything!", width: 1024, height: 1497)
]


Initially, I placed the coils on the battery. However, it didn't charge well (i.e., even with the coil as close as it could be). Also…turns out the battery was the thickest part of the internals, and the case didn't close…good thing I didn't force the case closed!

#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-coil-final.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-coil-final.webp"), type: "image/webp"))
  #html.img(src: "/2026/quaderno-coil-final.jpg", alt: "The diecast part of the tablet, with the coil nestled between the two batteries (above the larger battery, to the right of the smaller one, and below the I/O board). Also it's masking taped on there but we don't talk about that.", width: 1024, height: 1366)
]

I found a better spot, above the battery, and it fortunately didn't conflict with the plastic fins (or the plastic cover, probably for structural integrity).
I realized while writing that putting the coil against the diecast is not optimal. While testing the setup, I found it gets quite hot and it stops charging (!) after half an hour or so—presumably due do some thermal cutoff.
Instead, I put the coil against the plastic case so that it's close as possible to the charging coil, hopefully reducing thermal loss.

#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-coil-final2.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-coil-final2.webp"), type: "image/webp"))
  #html.img(src: "/2026/quaderno-coil-final2.jpg", alt: "The coil attached to the plastic instead", width: 1024, height: 1364)
]

#html.elem("aside")[
  Although the best result is to have all of the electronics fit inside the case, I wasn't totally on board with the placement and didn't have a solder rework station handy. The plan is to test this out for a few months first.
  Once I'm ready, I plan to remove the USB-C ports entirely and seal up the port hole ih the plastic cover.
]

Now, time to close the cover (for real this time)!

#html.elem("picture")[
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-charging.avif"), type: "image/avif"))
  #html.elem("source", attrs: (srcset: ("/2026/quaderno-charging.webp"), type: "image/webp"))
  #html.img(src: "/2026/quaderno-charging.jpg", alt: "Close up with the charging inicator on (red LED) and the Qi charger on the tablet, flipped down.", width: 1024, height: 1024)
]

// TODO: 実際に使ってみてどう？
// TODO: split into another post on why big paper e-readers are better

== Charging Performance

I left it on the charger overnight and jotted down the battery percentage whenever I happened to glance at it. Charging is fairly linear, at roughly 13 percentage points an hour, and it topped out at 100% a bit under four hours after starting at 51%.

#import "@preview/cetz:0.4.0"
#import "@preview/cetz-plot:0.1.2": plot

#let readings = (
  (0, 51), (8, 53), (24, 56), (44, 60), (59, 64), (75, 67),
  (106, 73), (109, 74), (137, 80), (141, 81), (205, 94), (236, 100),
)

#let make-chart(axis-color, line-color) = html.frame({
  text(font: "Atkinson Hyperlegible", fill: axis-color, {
    cetz.canvas({
      import cetz.draw: *
      set-style(axes: (
        stroke: axis-color,
        tick: (stroke: axis-color),
        grid: (stroke: axis-color.lighten(70%)),
      ))
      plot.plot(
        size: (13, 7),
        x-label: [Elapsed time (minutes)],
        y-label: [Battery (%)],
        x-tick-step: 30,
        y-tick-step: 10,
        x-min: 0, x-max: 236,
        y-min: 0, y-max: 100,
        mark-style: (stroke: line-color, fill: line-color),
        {
          plot.add(readings, mark: "o", mark-size: 0.12, style: (stroke: line-color, fill: line-color))
        }
      )
    })
  })
})

#let light-axis = rgb("#111111")
#let dark-axis = rgb("#eaeaea")
#let light-line = rgb("#c2255c")
#let dark-line = rgb("#ff9ecb")

#let chart-alt = "Line chart: battery percentage climbs linearly from 51% to 100% over about 4 hours of wireless, so approx. 13% per hour."

#html.elem("figure")[
  #html.elem("div", attrs: (class: "qi-chart-light", role: "img", "aria-label": chart-alt))[#make-chart(light-axis, light-line)]
  #html.elem("div", attrs: (class: "qi-chart-dark", role: "img", "aria-label": chart-alt, "aria-hidden": "true"))[#make-chart(dark-axis, dark-line)]
  #let chart-css = ".qi-chart-light svg, .qi-chart-dark svg { width: 100% !important; max-width: 100% !important; height: auto !important; } .qi-chart-dark { display: none; } @media (prefers-color-scheme: dark) { .qi-chart-light { display: none; } .qi-chart-dark { display: block; } }"
  #html.elem("style")[#chart-css]
  #html.elem("figcaption")[Battery percentage over time while charging wirelessly. Started charging from 51%.]
]

I'll spend a few months with this setup, and write another post on how it's like using it. Be sure to subscribe to the Atom feed so you don't miss it!
#html.span(role: "img", aria-label: "waving bye bye")[(\*´∀｀)ﾉ~~♪]

