#set document(
  title: [Komorebi: HTTP middleware to dynamically resize images],
	author: ("Ken Shibata"),
	keywords: ("Web"),
  date: datetime(year: 2025, month: 1, day: 22)
)

#html.elem("img", attrs: (src: "https://maclyonsden.com/media/featured_image/b8a57d0a4ca4439987c735affd7f3859.png?fmt=webp&w=500", alt: "Example Image dynamically resized and encoded to WebP using Komorebi"))

On the #link("/projects/metropolis")[Metropolis] website, static media was directly served from the uploaded version, which hwere usually unnecessarily larger and were not in a optimal format (e.g. PNG).

In order to speed up loading times, I developed this to transparently1 reduce bandwidth used. A quick search did not turn up a simple reverse proxy that I could put transparently (i.e. disables itself if not explicitly used) and has the least amount of features I need.

#link("https://github.com/wlmac/komorebi")[Source code]
