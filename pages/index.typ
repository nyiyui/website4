#set document(
  title: [Ken Shibata's Corner of The Internet],
  date: datetime(year: 2026, month: 9, day: 23)
)

I am an undergraduate studying CS at #link("https://gatech.edu")[Georgia Tech] (#link("https://www.cc.gatech.edu/academics/threads/systems-architecture")[Systems & Architecture] and #link("https://www.cc.gatech.edu/academics/threads/theory")[Theory]) with a minor in Math.
Welcome to my website!
#html.span(role: "img", aria-label: "happily waving hello")[ヾ(^▽^)ﾉ]
My interests lie in low-level systems, programming languages (PL) and compilers, and VR (#link(<pl-and-vr>)[why?]).

At Georgia Tech, I
- spread interesting PL and compiler ideas by hosting weekly meetings of the #link("https://dtyped.gtorg.gatech.edu")[dependently-typed club] (#link("mailto:ken.shibata@gatech.edu")[email me if you want to present!]),
- support researchers using LLM tooling at the #link("https://pace.gatech.edu/")[supercomputing center],
- help run the #link("https://www.gtxr.club/")[Georgia Tech Extended Reality club] as a Directory of Technology, and
- run a #link("mailto:modeltrains@groups.gatech.edu")[model trains club].

My hobbies include model trains and riding railways.

// A complete description is on my #link("./cv")[online CV].
// TODO

== News <news>
- I've rewrote my website (again)! Pages use Typst now! I'm in the process of copying pages over from my old website, please excuse this empty site…
- I'm attending #link("https://ieeevr.org/2026/")[IEEE VR 2026] and presenting at the #link("https://sites.google.com/view/immer-sense/home")[IMMER-SENSE workshop].

== Recent Posts
// The below element is replaced w/ the actual recent posts fragment at build time.

#html.elem("div", attrs: (id: "recent-posts-here"))

== Contact and Identities <contact>
Email: #link("mailto:ken.shibata@gatech.edu") or same username but this domain instead

ORCID: #link("https://orcid.org/0009-0000-4393-9964")[0009-0000-4393-9964]

== Why PL and VR, two seemingly unrelated fields? <pl-and-vr>
I grew up surrounded by computers and spend a good chunk of my time
understanding them through experimentation and reading as I found joy in
clarifying mechanical#footnote[meaning automatic (in the sense of "automatic programming"), not necessarily related to physical machinery] processes.
Given it's not possible for one to grasp every detail of a system,
I want to make computing systems _understandable through abstraction_.
This motivation led to
- PL/compilers for clarifying programming, which is arguably the main interface to stretch a computer's capabilities and
- VR for an "ultimate" skeumorphic computing interface.

