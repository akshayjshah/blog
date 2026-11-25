---
title: "Into the labyrinth: what distributed systems have taught us about guardrails"
description: "Engineering with agents is new, but we already know how to safely change systems that no one engineer fully understands: distributed systems engineers have struggled with the same problem for decades! Sometimes they've even succeeded. Deterministic simulation testing is one of their little-known superpowers. And as agents make all our code feel more like the unknowable distributed systems of yesteryear, perhaps DST can help us ship all kinds of code faster, with fewer bugs."
created: 2026-11-16
---

<img src="/static/qcon-labyrinth/slide01.png" alt="Title slide" max-height="400px">

Engineering with agents is new, but we already know how to safely change
systems that no single human fully understands: distributed systems engineers
have struggled with the same problem for decades! Sometimes they've even
succeeded. Deterministic simulation testing is one of their little-known
superpowers. And as agents make all our code feel more like the unknowable
distributed systems of yesteryear, perhaps DST can help us ship all kinds of
code faster, with fewer bugs.

<img src="/static/qcon-labyrinth/slide02.png" alt="Title slide" max-height="400px">

[Tanya Reilly][tanya], who's wonderful and the author of [The Staff Engineer's
Path][staffeng], invited me to talk about DST and agent guardrails at this
year's QCon San Francisco. My [slides are online](slides) and [downloadable as
a PDF][slides-pdf], and a [recording of the talk][vid] is available on the QCon
website. You do need to be logged in to watch the recording, though!

<img src="/static/qcon-labyrinth/slide03.png" alt="Title slide" max-height="400px">

I wish I'd spent a little more time rehearsing my timing, but I was happy with
the talk and especially happy with my craft-project slides. They took forever
(and some help from my kids) to make, but they did stand out from the typical
Claude-slop.

[slides-pdf]: /static/qcon-labyrinth/slides.pdf
[slides]: https://docs.google.com/presentation/d/1Phf1HfeNeaWru2xJ5eF2o4m8I6w75VJRzmZgEoK9BjE/edit?slide=id.p#slide=id.p
[staffeng]: https://bookshop.org/p/books/the-staff-engineer-s-path-a-guide-for-individual-contributors-navigating-growth-and-change-tanya-reilly/d83a627675d8b033?ean=9781098118730
[tanya]: https://www.noidea.dog/about
[vid]: https://qconsf.com/presentation/nov2026/labyrinth-what-distributed-systems-have-taught-us-about-guardrails
