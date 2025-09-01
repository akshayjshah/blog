---
title: "A decade of CAP"
description: "Slides and resources from a talk on Eric Brewer's CAP theorem."
created: 2016-03-01
---

> Consistency, availability, and tolerance to network partitions: you can have
> at most two of these properties for any shared-data system.
> <cite>[Eric Brewer](http://www.cs.berkeley.edu/~brewer/cs262b-2004/PODC-keynote.pdf)</cite>

Distributed systems folks think that everyone should read
[FLP](https://groups.csail.mit.edu/tds/papers/Lynch/jacm85.pdf), but engineers
in the trenches know the truth: the CAP theorem is king. (Whether or not the
many alternate formulations of the safety-liveness tradeoff are somehow
*better* than CAP is beside the point. They lost this round of marketing, so
they're not widely-known enough to be a useful *lingua franca*.) Despite
[all](http://cs-www.cs.yale.edu/homes/dna/papers/abadi-pacelc.pdf)
[the](http://www.bailis.org/papers/hat-hotos2013.pdf)
[competition](http://arxiv.org/pdf/1509.05393v2.pdf), Brewer's conjecture is
the most common framework for working programmers to think through the
trade-offs inherent in distributed systems.

But even though we name-check CAP all the time, most of us haven't fully
grokked [Gilbert and Lynch's
proof](https://www.comp.nus.edu.sg/~gilbert/pubs/BrewersConjecture-SigAct.pdf)
of Brewer’s conjecture. If you're not sure exactly what *consistency* and
*availability* mean, you're not alone; I recently gave a talk called [A Decade
of CAP](https://speakerdeck.com/akshayjshah/a-decade-of-cap) to a group of
engineers at work. (The slides for my talk include a great illustration of
linearizability that I stole from [Martin Kleppman's
blog](https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html).
Thanks, Martin!) Even in a group of experienced engineers, not everyone was
familiar with the details of the paper. (If you'd prefer, you can also download
a [PDF version of the slides](/static/decade-of-cap/decade-of-cap.pdf).)

It's trendy to [criticize CAP](http://arxiv.org/pdf/1509.05393v2.pdf) these
days, but I find most of the criticism unconvincing (though interesting). CAP's
weakness is simple: its notions of consistency and availability are so
stringent, and the system it proposes is so adversarial, that the final result
doesn't shed much light on our everyday work. Put differently, most useful
systems are neither CP nor AP. That doesn't make CAP any less true, it just
means that reality is complicated. Until an alternative emerges that's both
simpler *and* more useful, I'll continue to anchor my architecture discussions
in CAP, and so should you.

## Readable references

If you're interested in distributed systems but don't know where to start, you
may find some of these resources useful. Most of them take great pains to avoid
the dense jargon that makes formal distributed systems writing so hard to
approach.

- The [slides](http://www.cs.berkeley.edu/~brewer/cs262b-2004/PODC-keynote.pdf)
  from Eric Brewer's original keynote. Without a video recording, they're a
  little hard to interpret.
- Gilbert and Lynch's
  [proof](https://www.comp.nus.edu.sg/~gilbert/pubs/BrewersConjecture-SigAct.pdf).
  It's more approachable than you'd think.
- Gilbert and Lynch's
  [retrospective](https://groups.csail.mit.edu/tds/papers/Gilbert/Brewer2.pdf)
  on CAP. This is a great bridge from CAP to the wider world of formal
  distributed systems.
- [Tyler Treat](http://bravenewgeek.com/) and [Martin
  Kleppman's](https://martin.kleppmann.com/)
  blogs.
