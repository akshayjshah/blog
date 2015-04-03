+++
title = "Hakyll and Bootstrap"
date = "2012-08-04"
+++

<p class="update"><em>Update:</em> It's probably no surprise, but diving
straight into Haskell wasn't the smoothest transition; without a better
understanding of the language as a whole, I spent more time fiddling with the
code than writing. I experimented with a few different site generation tools,
including a custom Python script, before settling on Steve Francia's <a
href="http://gohugo.io/" title="Hugo: A Fast & Modern Static Website
Generator">Hugo</a>.
</p>

Though I've run a few [Wordpress][] sites in the past, I've always found the
software irritating to use. In particular, I *loathe* the composer--the
rich text editor never works quite as I expect it to, and the <abbr
class="initialism">HTML</abbr> editor mangles my markup without warning. The
installation is dead simple, the theme support is fantastic, and I love [Matt
Mullenweg's support of free software][gpl], but Wordpress just isn't for me.
At the same time, I don't miss hand-coding entire sites; I've done that exactly
twice, and maintaining scads of duplicated markup is a nightmare I don't care
to repeat.  Since I'm starting this site with a clean slate, I have the luxury
of choosing whatever tools I like; after looking at everything from [Drupal][]
to [Tumblr][], I decided to keep things simple and build with [Hakyll][] and
[Bootstrap][].

## Hakyll: Static Site Generation

Hakyll is a static site generator, which means that it's only a few short steps
away from writing markup by hand. Those few steps make all the difference,
though--I can still control exactly what the final <abbr>HTML</abbr> will look
like, but Hakyll lets me write articles in [Markdown][] and use templates for
site-wide elements like navigation and footers.  Static sites are all the rage
among geeks these days, so [a lot's][nanoc] already [been written][stevelosh]
about their technical advantages.  In short, static sites are:

* fast by default, even on inexpensive hardware;
* immune to many common security exploits;
* easily saved and versioned; and
* editable offline.

Most importantly for me, Hakyll is comfortable. I can write templates and posts
in [vim][], keep versions and branches in [git][], and generally work with
whatever text manipulation tools make me happy. Hakyll even comes with a
built-in webserver, so it's easy to see a live preview of any changes I'm
making. Since writing on this site is supposed to be a fun side project,
comfortable tools are priority zero.

Building a site without a database isn't all roses and puppies, though. There's
no easy way for me to include an automated "Popular Posts" widget in the
sidebar, for example, and any future search widgets will need to rely on an
external search engine. Most importantly, it's impossible for me to store and
manage reader comments. Services like [Disqus][] and [IntenseDebate][] offer
easy Javascript-based workarounds, but I'll need to do a little more due
diligence before I'm comfortable trusting them with critical data.

Hakyll also has its own set of challenges, mostly because it's written and
configured in [Haskell][]. To put it mildly, I'm a Haskell neophyte--I've been
interested in the language for months, but haven't done anything more than a
few [Project Euler][] questions. Since I don't have a strong background in
category theory, monads and arrows are *blowing my mind*.  There's something
really amazing and elegant going on, but I'm only catching glimpses of it
between compiler errors. Nevertheless, the Hakyll [documentation][hakyll-docs]
is excellent, the [mailing list][hakyll-list] is active, and the
[author][jasper] is exceptionally helpful, so my first foray into practical
functional programming has been more enlightening than infuriating.

## Bootstrap: Clean Design, No Fuss

I have trouble matching my clothing, let alone the dozens of small elements
that make up a website, so creating an attractive design for my new site was a
daunting task. Luckily, nobody needs to see my first efforts--I decided to use
Twitter's [Bootstrap][] framework instead. It's clean, attractive, and
mobile-friendly, and it's teaching me some of <abbr>HTML</abbr>5's new tricks.
[{less}][less], the <abbr>CSS</abbr> meta-language Boostrap uses, is also
wonderful: it's close enough to vanilla <abbr>CSS</abbr> that it's easy to
learn, but it makes my stylesheets much more modular and consistent.

While I haven't tweaked Bootstrap's default styling much, I *had* to do
something about the fonts. I like Helvetica, especially on visually intense
marketing sites--but Bootstrap's tiny default font size combined with
Helvetica's clinical modernism made blocks of text downright hostile. After a few
hours poking through [Google Web Fonts][webfonts] and testing different styles,
I settled on [Omnibus Type's][omnibus] [Rosario][rosario]. To my eye, it
manages to be a little more playful and human than Helvetica without
distracting from the words themselves.

My efforts to choose a different color scheme, though, have been a complete
failure. The defaults are nice enough, but they lack soul, and my efforts to
change them usually end in a neon-tinted nightmare. There's hope on the
horizon, though--I just read Ian Taylor's ["Never Use Black,"][black]
and I may try mixing some blue or red into the default grays.

I'm not much of a programmer or a designer, so I'm always in the market for
suggestions! If coding's your thing, take a look at the [source code on
GitHub](http://github.com/akshayjshah/akshayjshah.github.io); otherwise, send
me a [tweet](http://twitter.com/akshayshah) and let me know what you think.

[Bootstrap]: http://twitter.github.com/bootstrap/ "Twitter Bootstrap"
[Disqus]: http://disqus.com "Disqus"
[Drupal]: http://drupal.org/ "Drupal"
[Hakyll]: http://jaspervdj.be/hakyll/ "Hakyll"
[Haskell]: http://www.haskell.org/haskellwiki/Haskell "HaskellWiki"
[IntenseDebate]: http://intensedebate.com "IntenseDebate"
[Markdown]: http://daringfireball.net/projects/markdown/ "Markdown"
[Project Euler]: http://projecteuler.net/ "Project Euler"
[Tumblr]: https://www.tumblr.com/ "Tumblr"
[WordPress]: http://www.wordpress.com "WordPress"
[black]: http://ianstormtaylor.com/design-tip-never-use-black/ "Design Tip: Never User Black"
[git]: http://git-scm.com/ "Git"
[gpl]: http://ma.tt/tag/gpl/ "Matt Mullenweg on the GPL"
[hakyll-docs]: http://jaspervdj.be/hakyll/tutorials.html "Hakyll Tutorials"
[hakyll-list]: http://groups.google.com/group/hakyll "Hakyll Google Group"
[jasper]: http://jaspervdj.be "Jasper Van der Jeugt"
[less]: http://lesscss.org "{less}"
[nanoc]: http://nanoc.stoneship.org/docs/1-introduction/ "Nanoc Documentation"
[omnibus]: http://www.omnibus-type.com/ "Omnibus Type"
[rosario]: http://www.google.com/webfonts/specimen/Rosario "Rosario"
[stevelosh]: http://stevelosh.com/blog/2010/01/moving-from-django-to-hyde/ "Steve Losh: Moving from Django to Hyde"
[vim]: http://stevelosh.com/blog/2010/09/coming-home-to-vim/ "Steve Losh: Coming Home to Vim"
[webfonts]: http://www.google.com/webfonts "Google Web Fonts"
