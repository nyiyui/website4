// Builds a <picture> element with AVIF and WebP <source>s plus a fallback
// <img>, for images that (per render/cmd/imgprep) have .avif/.webp/.jpg (or
// .png) siblings generated from the same base path.
#let picture(src, alt, width: none, height: none, ext: "jpg", fetchpriority: none) = {
  let img-attrs = (src: src + "." + ext, alt: alt)
  if width != none { img-attrs.width = str(width) }
  if height != none { img-attrs.height = str(height) }
  if fetchpriority != none { img-attrs.fetchpriority = fetchpriority }

  html.elem("picture")[
    #html.elem("source", attrs: (srcset: src + ".avif", type: "image/avif"))
    #html.elem("source", attrs: (srcset: src + ".webp", type: "image/webp"))
    #html.elem("img", attrs: img-attrs)
  ]
}
