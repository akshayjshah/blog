+++
title = "License Your Code!"
date = "2013-08-04"
+++

If you're a [GitHub][] user, you probably
want to let other people use your code in their own projects. Just throwing
your code online isn't enough, though---you *have* to release it under an
open-source license.  If you're not sure why, read on: in this post, I'll
explain the absolute minimum that every developer should know about copyright
and software licensing.

Even though I'm not a lawyer, I'll start with a few
disclaimers: this post isn't legal advice, I'm not going to address patents at
all, and some of what I'm about to explain is specific to the United States.

## Closed by Default

In the US, every shred of code you've ever written is automatically protected
by copyright. That means, broadly speaking, that it's your property.  If you
leave your bike unlocked outside a coffee shop, I'm not allowed to walk by and
steal it; similarly, if you post your code online, I'm not automatically
allowed to copy and use it. Don't worry, though---seventy years after you die,
the copyright expires and your code enters the *public domain*, at which point
I can finally use it however I like.

Under American law, it's very difficult to reliably circumvent the century-long
wait before your copyright expires. To let other developers legally copy,
modify, and reuse your code while it's still relevant, you have to offer it to
them under a *license*.

## Permissive Licensing

Licenses are just agreements between you, the author and owner of the code, and
all the other developers who'd like to use your code. Even if you've never
exchanged emails or chatted on the phone, you can sue developers who use your
code but don't follow the rules in your license.

Some open-source software licenses are quite straightforward: they let other
developers do whatever they'd like with your code as long as they don't hold
you liable for the consequences. They can package your library into a
closed-source application and sell subscriptions, they can modify the code and
distribute it under a different license, and they can even sell copies of your
code without changing a single line. If your code has a bug that causes a
nuclear meltdown, though, they're on their own. Because they let other
developers do whatever they'd like, these licenses are often called
*permissive*; the [MIT license][] is a popular example.

## Copyleft

But what if you don't want your code integrated into closed-source systems?
Richard Stallman and the [Free Software Foundation][fsf] are right there with
you. In a clever legal hack, they wrote the [GPL][], a license which lets other
developers run, read, change, and share your source code, but requires that any
copies (modified or not) be licensed in a way that preserves the same freedoms.
Licensing your work under the GPL allows other people to integrate your code
into a paid product, but it forces them to license their product under the GPL,
too (or something that's essentially identical). In short, your code and all
its descendants will be free forever, but many companies won't be willing to
use it. Because it uses a copyright hack to attack the very notion of
software ownership, the GPL and its variants are often called *copyleft*
licenses.

And with that, you know more than 99% of programmers about software licensing.
Stand proud, grasshopper, and remember to include an explicit license when
sharing your code online.

## More Questions?

* Need a refresher? Looking for some actual legalese to include in your
  project? Check out GitHub's [license picker][].
* Want to know more about intellectual property and open source? Read Richard
  Stallman's [Why Software Should Not Have Owners][no-owners]. It's a seminal
  piece, and it's well worth your time.
* Wondering why most folks these days call it "open-source" and not "free"
  software? Read Eric Raymond's [1998 call to action][esr-opensource]. While
  you're there, you may want to check out [The Cathedral and the Bazaar][catb].
  (Ignore anything related to guns, sex, or politics.) For Stallman's response,
  read [Why Open Source Misses the Point of Free Software][rms-opensource].
* Curious how all this copyright stuff relates to software patents? It's
  complicated. Patents aren't the same thing as copyright, though some
  open-source licenses attempt to address them. Think of patent protection as
  another axis that's orthogonal to the permissive-copyleft axis.

[GitHub]: http://github.com
[Bitbucket]: http://bitbucket.org
[Gitorious]: http://www.gitorious.com
[MIT license]: https://en.wikipedia.org/wiki/MIT_License
[fsf]: http://fsf.org
[GPL]: http://en.wikipedia.org/wiki/GNU_General_Public_License
[no-owners]: http://www.gnu.org/philosophy/why-free.html
[four-freedoms]: http://www.gnu.org/philosophy/free-sw.html
[license picker]: http://choosealicense.com/
[esr-opensource]: http://catb.org/~esr/open-source.html
[catb]: http://catb.org/esr/writings/cathedral-bazaar/cathedral-bazaar/
[rms-opensource]: http://www.gnu.org/philosophy/open-source-misses-the-point.html
