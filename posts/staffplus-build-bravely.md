---
title: "Build bravely: delivering risky projects"
description: "Engineering teams shy away from ambitious projects for good reason: they’re risky, even with the best strategy, execution plan, and team. But just like baggy jeans, the best ideas of the late 90s are back, better than ever, and ready to save us. RFC 2119, adopted in 1997, gives us unambiguous language for specifying behavior. Quickcheck, released in 1999, lets us turn specifications into executable test suites. By combining them, we can shift ambiguity left, deliver working systems quickly, and ease collaboration with LLMs."
created: 2025-10-15
---

This fall, I spoke at [LeadDev][leaddev] about property-based testing --- but I
camouflaged all the Haskell in a story about well-intentioned wizards and
murderous steampunk robots. There's a [nice video of the talk][vid] and a [PDF
of my slides][slides] available on the LeadDev website.

Unfortunately, the video and PDF are behind a registration wall. The LeadDev
folks put on an excellent conference, so I don't begrudge them the opportunity
to grow their community. But if you'd rather stay anonymous, here's a version
of my slides annotated with a slightly cleaned-up transcript of the talk.

<img src="/static/staffplus-build-bravely/slide01.png" alt="" max-height="400px">

Welcome to Build Bravely, in which we'll see how two bangers from the 90s can
help us deliver risky projects on time, on spec, and *not* on fire.

<img src="/static/staffplus-build-bravely/slide02.png" alt="" max-height="400px">

I'm Akshay, and I've been an infra and platform engineer for about 15 years.
I've worked at other people's silly startups, my own silly startup, Uber, and
Microsoft. These days, I'm the Field CTO at [Antithesis][antithesis].

This is Gemini's portrayal of me as a wizard. The shell-shocked look on my face
is spot-on, because I've spent most of my career fixing horrifyingly broken
software. I'm now middle-aged, a bit grumpy, and a lot skeptical. I am not
usually brave. I am *not* into the big-bang rewrite.

<img src="/static/staffplus-build-bravely/slide03.png" alt="" max-height="400px">

I prefer gradual improvement. I like local maxima. But sometimes, requirements
change, the business grows, the unexpected arrives, and we are grudgingly
forced to build something completely new. Something risky. Something a bit
brave.

<img src="/static/staffplus-build-bravely/slide04.png" alt="" max-height="400px">

Every brave project begins with a plan. Often, a humble user story: "as a
wizard, I want autonomous defenders so I can stay in my lab and not talk to
anyone." And of course, our brave project comes with milestones, culminating in
my triumphant staff promotion.

<img src="/static/staffplus-build-bravely/slide05.png" alt="" max-height="400px">

Maybe we're a bit more formal, and our projects start with plans like this. But
whether they're long or short, our plans tend to focus on three things: why,
how, and when. These are good questions to answer. But even when the plan is
really good, and really aligned with our engineering strategy, all of you
reading it feel this creeping sense of doom.

Yeah, you can quibble with the details here. Will our investment really be
sub-linear? Is a gossip protocol the right choice? How did Svelte even become
involved? But quibbling doesn't change the most likely outcome...

<img src="/static/staffplus-build-bravely/slide06.png" alt="" max-height="400px">

...which is complete and utter disaster. Our project will almost certainly be
late. It will solve only part of the problem. After all the delays and
thrashing, half the team will quit or transfer. And then, just as the robot
automatons are murdering all of us, we'll realize that we have drifted *ever so
slightly* off-spec.

Yikes. How did we go so wrong?

<img src="/static/staffplus-build-bravely/slide07.png" alt="" max-height="400px">

Well, we treated our ambitious project like it was going to be a casual weekend
hike --- when in fact it's an expedition into a blasted hellscape littered with
the bodies of our peers.

Yeah, we planned the how and the when. The route and the itinerary. We read
[Will Larson's books][lethain]. We followed all the best practices. But we didn't spend
enough time describing the destination. We didn't specify our new system
carefully enough and early enough, so we spent a lot of time thrashing and ---
oops! --- built robots that turned out to be murder hobos.

So what does a good up-front spec look like?

<img src="/static/staffplus-build-bravely/slide08.png" alt="" max-height="400px">

I don't know, but I am familiar with one kind of spec --- the basic integration
test. This is a critical test that we obviously forgot to write: the robots
should not murder us.

<img src="/static/staffplus-build-bravely/slide09.png" alt="" max-height="400px">

But wait! To be really thorough, we also need a test with two robots, in a
forest, with a slightly different command. But this is just the tip of the
iceberg.

<img src="/static/staffplus-build-bravely/slide10.png" alt="" max-height="400px">

There are so many possible situations. We need a whole suite of integration
tests. And they're probably going to end up outweighing the implementation. We
need a QA team, or an SDET team --- this *cannot* be our problem.

And deep inside, we all know that this is a complete sham. Even with all of
this, we've only tested the tiniest fraction of the possible states that our
code could be in. And we have burned *tons* of engineering time on this. Whole
teams have been hired just to write and maintain test suites like this one.
When the spec changes, we have to amend all of this code.

So far from giving us wings, all of this is just weighing us down. And worst of
all, these specs are very bad at communicating intent. Rather than explicitly
stating the primary constraint --- our creations should not rise up and kill us
--- they're just hinting at it, over and over and over again. They're asking
everyone on our team to mentally compile this series of examples down into the
main point.

This kind of test is not worth the effort.

<img src="/static/staffplus-build-bravely/slide11.png" alt="" max-height="400px">

The root of the problem is that the test suite --- this spec --- is mired in
trivial details. One robot. In the lab. One wizard. One specific command. It
sounds like we're playing Clue! Ideally, we would condense all of this into one
much more general test.

<img src="/static/staffplus-build-bravely/slide12.png" alt="" max-height="400px">

That all-encompassing test might look like this. Any number of robots. Any
place. Any number of wizards. Any phrase at all. Aha! Now we need fewer tests
--- maybe just this one! --- and we are clearly and explicitly communicating
intent. Rather than a test, we might call this a property.

<img src="/static/staffplus-build-bravely/slide13.png" alt="" max-height="400px">

We've also made evolution easier. Remember, I'm suggesting that we write this
property --- this spec --- really early in the planning cycle, maybe even
before we've written any code. As we talk to our users, we're going to discover
even more requirements. And when we do, it's easy to extend our properties.

And as we extend our properties, we're forced to be very clear: are these new
requirements really *required*, or are they just recommendations, or even just
nice-to-haves? If an engineer or an LLM on my team wants to amend this
property, their intention is very clear. If we weaken "*must* not murder" to
"*should* not murder," that's an extremely large red flag that's very easy to
catch.

Of course, if you're listening carefully, all of this "must" and "should" and
"may" is starting to ring a bell.

<img src="/static/staffplus-build-bravely/slide14.png" alt="" max-height="400px">

It's all starting to sound a bit like [RFC 2119][rfc], the number one banger of 1997!
In two pages, this RFC defines the terminology for requirement levels, which is
a really simple, powerful idea: we should distinguish between absolute
requirements, recommendations, and nice-to-haves. Thinking this way forces us
to navigate tradeoffs and competing priorities very early.

So now we've got this RFC-like property, which is strong, and short, and still
manages to be nuanced, but it has the same weakness as internet RFCs. It's not
executable.

<img src="/static/staffplus-build-bravely/slide15.png" alt="" max-height="400px">

Or at least it wasn't until 1999, when [QuickCheck][quickcheck] was released.
QuickCheck takes properties like ours, and it generates input data, and then it
checks that the properties still hold.

<img src="/static/staffplus-build-bravely/slide16.png" alt="" max-height="400px">

You can think of that like our robot defense system running in multiple
parallel universes, with different numbers of robots, different locations,
different enemies, and different commands. QuickCheck is generating thousands
of these universes and verifying that our robots never rise up in rebellion.

This is called property-based testing, and it is phenomenally effective. In a
[recent study][oopsla] of 40 Python projects on Github, a couple of researchers
discovered that property-based tests are *fifty* times more effective at
catching bugs than standard tests. (At the conference, I misspoke and said that
the study looked at 400 projects. Oops!)

<img src="/static/staffplus-build-bravely/slide17.png" alt="" max-height="400px">

Unfortunately, QuickCheck was written in Haskell, so only seven people in the
world cared --- and all of them were busy writing the Haskell compiler. But
over the next 25 years, we were able to decrypt these tomes of wisdom and
translate them into all the languages that all of us actually use. So now, no
matter what language you work in, you have a really good PBT library available
to you.

Every project benefits from this style of verification. It's simple. It's easy.
I promise, it doesn't hurt.

<img src="/static/staffplus-build-bravely/slide18.png" alt="" max-height="400px">

And that's my pitch to you. Use the language of RFC 2119. Use property-based
testing. And write specs for your software early --- ideally during planning,
when you're already thinking about prioritization and tradeoffs. And then, for
the life of your project, your team has jet fuel. You can evolve your code
rapidly --- with humans, with agents, with whomever you like. Because you have
a powerful, understandable test suite that keeps your project aligned with the
original plan. For sure. Guaranteed!

<img src="/static/staffplus-build-bravely/slide19.png" alt="" max-height="400px">

And that's it! Special thanks to Scott Bradner, who wrote RFC 2119, to John
Hughes and every other contributor to QuickCheck, and all the folks pushing PBT
forwards. And with that, thank you for your attention.

[antithesis]: https://antithesis.com
[leaddev]: https://leaddev.com
[lethain]: https://lethain.com/
[oopsla]: https://cseweb.ucsd.edu/~mcoblenz/assets/pdf/OOPSLA_2025_PBT.pdf
[quickcheck]: https://hackage.haskell.org/package/QuickCheck
[rfc]: https://datatracker.ietf.org/doc/html/rfc2119
[slides]: https://leaddev.com/wp-content/uploads/2025/11/Build-bravely-Delivering-risky-projects-Akshay-Shah-StaffPlus-New-York-2025.pdf
[vid]: https://leaddev.com/software-quality/build-bravely-delivering-risky-projects
