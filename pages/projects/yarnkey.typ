#set document(
  title: [Yarnkey/Kyokukey],
	author: ("Ken Shibata"),
	description: [Yarnkey (formerly Kyokukey) is a wearable input device that can be used in various postures with one hand using conductive thread. It can be used in a narrow space because the sensor has an area of a few cm#super[2].],
	keywords: ("HCI", "Experts in Information Science"),
  date: datetime(year: 2025, month: 1, day: 22)
)

#html.elem("img", attrs: (src: "/projects/yarnkey/hero.png", alt: "Inputting \"hello\" using Yarnkey."))

Contributors: Ken Shibata, #link("https://narumi.me/profile")[Koya Narumi]

With the spread of mobile devices, the demand for wearable character input devices is increasing.
Various methods have been proposed so far, but can be improved such as their visual invasiveness, and bulkiness.
Yarnkey attempts to solve the problems as a mobile character input device.
We propose a wearable keyboard that can be used with one hand, is inconspicuous, is not bulky, and can be used quietly in various situations such as in a crowded train in a natural posture.
The new wearable device provides 6-point touch sensing with conductive threads sewn near the thighs on trousers.
The microcomputer recognizes characters drawn in a 2D pattern on the sensor.
It can be used in a narrow space because the sensor has an area of a few cm#super[2].

- #link("https://yarnkey.nyiyui.ca")[User Evaluation App] (#link("https://github.com/nyiyui/kyokukey")[source]#footnote[The server code is a simple PHP script.])
- #link("https://github.com/nyiyui/yarnkey")[Device (microcontroller) source code and support files]
- #link("https://docs.google.com/presentation/d/1tFMlMbAf0-CJZBCnUkoLLcBch4otMCAlcAMVbbjPzi4/edit#slide=id.p")[Slides for DICOMO 2022]
- #link("/projects/yarnkey/yarnkey_dicomo-2022_ja.pdf")[Paper for DICOMO 2022 (in Japanese)]
// - #link("/projects/yarnkey/yarnkey_dicomo-2022_en.pdf")[Paper for DICOMO 2022]
- #link("https://projectboard.world/ysc/project/kyokukey-wearable-input-device-using-conductive-threads")[TSF 2022/CWSF 2022 Exhibition]

#html.elem("iframe", attrs: (src: "https://www.youtube-nocookie.com/embed/dMT6_UEli-g", title: "YouTube video player", width: "560", height: "315", allowfullscreen: "true"))
