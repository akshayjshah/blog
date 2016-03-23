+++
date = "2016-03-22T16:53:55-07:00"
title = "A Decade of CAP"

+++

{{% blockquote
    footer="Eric Brewer, <cite>[Towards Robust Distributed Systems](http://www.cs.berkeley.edu/~brewer/cs262b-2004/PODC-keynote.pdf)</cite>"
    epigraph="true"
%}}
Consistency, availability, and tolerance to network partitions: you can have at
most two of these properties for any shared-data system.
{{% /blockquote %}}

Distributed systems folks think that everyone should read [FLP][], but engineers
in the trenches know the truth: the CAP theorem is king.
{{< sidenote "ubiquity" >}}
Whether or not the many alternate formulations of the safety-liveness tradeoff
are somehow <em>better</em> than CAP is beside the point. They lost this round
of marketing, so they're not widely-known enough to be a useful lingua franca.
{{< /sidenote >}}
Despite [all][PACELC] [the][HAT] [competition][delay-sensitivity], Brewer's
conjecture is---by far---the most common framework for working programmers to
think through the tradeoffs inherent in distributed systems.

But even though we name-check CAP all the time, most of us haven't fully grokked
[Gilbert and Lynch's proof][proof] of Brewer's conjecture. If you're not sure
exactly what *consistency* and *availability* mean, you're not alone; I recently
gave a talk called "[A Decade of CAP][speakerdeck]" to a group of engineers at
work.
{{< sidenote "kleppman" >}}
The slides for my talk include a great illustration of linearizability that I
stole from <a
href="https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html">Martin
Kleppman&rsquo;s blog</a>. Thanks, Martin!
{{< /sidenote >}}
Even in a group of experienced engineers, not everyone was familiar with the
details of the paper. (If you'd prefer, you can also download a [PDF version of
the slides][preso].)

It's trendy to [criticize CAP][delay-sensitivity] these days, but I find most of
the criticism unconvincing (though interesting). CAP's weakness is simple: its
notions of consistency and availability are so stringent, and the system it
proposes is so adversarial, that the final result doesn't shed much light on our
everyday work.  Put differently, most useful systems are neither CP nor AP. That
doesn't make CAP any less true, it just means that reality is complicated. Until
an alternative emerges that's both simpler *and* more useful, I'll continue to
anchor my architecture discussions in CAP, and so should you.

## Readable References

If you're interested in distributed systems but don't know where to start, you
may find some of these resources useful. Most of them take great pains to avoid
the dense jargon that makes formal distributed systems writing so hard to
approach.

* The [slides][cap] from Eric Brewer's original keynote. Without a video
  recording, they're a little hard to interpret.
* Gilbert and Lynch's [proof][]. It's more approachable than you'd think.
* Gilbert and Lynch's [retrospective][] on CAP. This is a great bridge from CAP
  to the wider world of formal distributed systems.
* [Tyler Treat][] and [Martin Kleppman's][] blogs.

[FLP]: https://groups.csail.mit.edu/tds/papers/Lynch/jacm85.pdf
[HAT]: http://www.bailis.org/papers/hat-hotos2013.pdf
[Martin Kleppman's]: https://martin.kleppmann.com/
[PACELC]: http://cs-www.cs.yale.edu/homes/dna/papers/abadi-pacelc.pdf
[Tyler Treat]: http://bravenewgeek.com/
[cap]: http://www.cs.berkeley.edu/~brewer/cs262b-2004/PODC-keynote.pdf
[delay-sensitivity]: http://arxiv.org/pdf/1509.05393v2.pdf
[preso]: /docs/decade-of-cap.pdf
[proof]: https://www.comp.nus.edu.sg/~gilbert/pubs/BrewersConjecture-SigAct.pdf
[retrospective]: https://groups.csail.mit.edu/tds/papers/Gilbert/Brewer2.pdf
[speakerdeck]: https://speakerdeck.com/akshayjshah/a-decade-of-cap
