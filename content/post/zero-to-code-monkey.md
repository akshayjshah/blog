+++
title = "Zero to Code Monkey"
date = "2012-09-24"
+++

Last year, I'd never have called myself a programmer. Sure, I'd dabbled a bit
here and there: I took a few computer science courses in high school and
college, cobbled together some Excel macros and one-off data processing scripts
in medical school, and hand-wrote two websites in the late 90s. I'd never
written code daily, even while taking computer science classes, and most
importantly, I'd never collaborated on a software project. At best, I was a
dilettante.

Nine months later, I'm on a different plane: I'm still no expert, but I'm
definitely a programmer (or at least a code monkey). Writing code every day is
an important part of my job, and I *love* it. If you're interested in learning
to code, especially in Python, I hope my story helps you get started faster.

## Syntax Is Easy
Oddly enough, learning the syntax of a programming language --- the rules
for assembling words, curly braces, backslashes, and all the other
typographical flotsam in the corners of your keyboard into working code ---
was relatively easy. [Python][] is particularly straightforward, and Zed Shaw's
[Learn Python the Hard Way][lpthw] taught me the basics in a single weekend.
Experimenting in an [interactive environment][repl.it] was very helpful at this
stage, since the interpreter caught all my mistakes, printed a helpful error
message, and let me try again right away. Whenever I didn't understand where
I'd gone wrong, pasting the error message verbatim into Google always turned up
a relevant answer.

Only a few days in, I was able to exercise my new skills by writing a to-do
manager like [Steve Losh's t][t] --- and it *actually worked*.  Writing
even a simple utility that functioned exactly as I wanted it to was hugely
empowering, and using my little program every day kept me excited about
learning. Emulating t had an unanticipated, but immensely helpful, side effect:
it let me compare my efforts to [the actual t source code][t-src], which is
also written in Python. Taking the time to really understand how Steve's code
worked (and why it was much, much better than mine) accelerated my learning
dramatically.

## Computer Science Is Practical
Of course, what's a to-do manager without projects, sub-projects, and maybe
even sub-sub-projects? Naively, I assumed that adding this feature to my little
utility would be a quick one-hour project. When I sat down to start
writing, though, I realized that I had no clue how to handle potentially
infinite project nesting. A little research (leavened with some hazy memories
from school) revealed that managing data with this sort of parent-child
relationship is, unsurprisingly, a well-studied problem in computer science,
and after a few nights of experimenting and reading I managed to add projects
to my to-do list manager.

Along the way, I realized that I'd always be a few steps behind until I learned
some basic computer science: most common problems have established solutions,
and my complete ignorance of those solutions was forcing me to reinvent the
wheel (poorly) at every turn. It's a little like trying to organize thousands
of loose papers without using a file cabinet; I'd probably come up with some
sort of workable system eventually, but why go through all that pain when
there's a good solution that's widely available? To push the analogy
further, file cabinets also have a large ecosystem of standard add-ons
designed to solve particular problems --- it's easy to get cabinets with
built-in indexing systems, horizontal drawers, or wheels, for example.
Similarly, it's easy to find variations of the basic computer science building
blocks that are optimized for particular applications.

Nine months isn't enough to make anyone an expert in computer science, but
there are a *lot* of excellent courses available free online. I started with
MIT's [Introduction to Computer Science and Programming][mit6.00] because it's
taught in Python, doesn't assume any prior knowledge, and moves slowly enough
for me to watch while riding the train and keeping an eye out for my stop. I'm
not sure what I'll tackle next, but I'm leaning towards continuing the MIT
undergraduate sequence.

## Craftsmanship Is Crucial

A few days after I started watching the MIT lectures, I switched jobs and began
writing code at work. This was, to say the least, a big transition: I went from
being the sole author of a hundred-line program to working with a team of real
engineers building a real web application. The standard for the correctness of
my code went up, but more importantly, so did the standards for its
*readability* and *maintainability*.

Even with reams of well-written code (and the inevitable bits of crufty legacy
code) to learn from, it took me months to really understand why anyone would
prioritize anything other than accurate output. "It works right now," I'd think to
myself, "so why would anyone spend another week just making it look pretty?"
That attitude changed the first time I had to add some functionality to a small
piece of code I'd written a few months earlier. Even as the original author, I
was completely baffled --- I'd named every variable ``x``, ``sbgs``,
``data`` (yes, really), or something equally useless, and now I couldn't make
heads or tails of any of it. I ended up scrapping the whole file, starting from
scratch, and feeling lucky that none of the many bugs in my old code had seen
the light of day.

Needless to say, I'm now a convert. I'm getting better at giving functions and
variables meaningful names. I plan my code more carefully. I'm a little more
comfortable with object-oriented programming, and I'm keeping global variables
to a minimum. More abstractly, I'm slowly learning how to hide complexity so
that other people can treat my code as a pluggable component. I want to take
pride not just in my code's functionality, but in its readability, its
cleanliness, and its aesthetics --- in short, in its craftsmanship.

Taking a few months to focus on something other than basic computer science is
a luxury I can afford only because Python provides an off-the-shelf version of
most common building blocks. Since I don't absolutely need to know, for
example, how each of the many common sorting algorithms works (though it would
probably be helpful), I can afford to put reading on hold for a little while to
become a better craftsman. I've spent most of this time reading [Code
Complete][], mostly because Jeff Atwood describes it as "[the Joy of Cooking
for software developers][codinghorror]." I'll probably still make the same
rookie mistakes, but I hope that this early investment will help me learn from
my mistakes faster.

## What Next?
I'll finish Code Complete this week, but I'm not sure where I should invest my
time next. I could:

* Review linear algebra. I'm hoping to pick up a few machine learning projects
  soon, so dusting off my math skills would be helpful.
* Start another computer science course or read an introductory book on data
  structures and algorithms. Along with learning more math, this is probably
  the best long-term use of my time.
* Learn the Python numerical toolkit in more depth. This is definitely the most
  practical option, but I think it has the least enduring value.
* Read more about software craftsmanship and project management. I have copies
  of [The Pragmatic Programmer][] and [The Mythical Man-Month][], but I'm
  inclined to focus more on hard skills for the next few months.

If you've read this far, I'd appreciate your input! Send me an
[email](mailto:akshay@akshayshah.org) and let me know where I went wrong.

[LaTeX]: http://nitens.org/taraborelli/latex "The Beauty of LaTeX"
[Python]: http://www.python.org/about/
[lpthw]: http://learnpythonthehardway.org/
[repl.it]: http://repl.it/ "repl.it"
[Ruby]: http://www.ruby-lang.org/en/documentation/quickstart/ "Ruby Quickstart"
[Perl]: http://learn.perl.org/first_steps/ "Perl First Steps"
[t]: http://stevelosh.com/projects/t/
[t-src]: https://bitbucket.org/sjl/t/src/
[mit6.00]: http://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-00-introduction-to-computer-science-and-programming-fall-2008/
[Code Complete]: http://www.amazon.com/Code-Complete-Practical-Handbook-Construction/dp/0735619670
[codinghorror]: http://www.codinghorror.com/blog/2004/02/recommended-reading-for-developers.html
[The Pragmatic Programmer]: http://pragprog.com/book/tpp/the-pragmatic-programmer
[The Mythical Man-Month]: http://www.amazon.com/The-Mythical-Man-Month-Engineering-Anniversary/dp/0201835959
