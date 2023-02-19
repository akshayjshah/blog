# Getting Started with Go

If you're starting a new project or a job that requires knowing Go, here's an
opinionated learning plan that will get you ready in a hurry. If you're
exploring Go out of idle curiosity, this probably isn't the plan for
you---there are plenty of excellent videos and conference talks with a better
balance of entertainment and technical content.

## Learn the Language

You don't need a local Go installation to get started - you can learn the
language and start writing code in your browser.

1. Watch an [hour-long talk by Rob
   Pike](http://www.infoq.com/presentations/Go-Google), one of the core
   authors, for a sense of Go's goals and aesthetic. 
2. Do the entire [Tour of Go](http://tour.golang.org) and read [How to Write Go
   Code](https://golang.org/doc/code.html).
3. If you're writing Go professionally, join the [Gophers
   Slack](https://invite.slack.golangbridge.org/) and the
   [golang-nuts](https://groups.google.com/forum/#!forum/golang-nuts) mailing
   list. If that sounds like too much, at least join
   [golang-announce](https://groups.google.com/forum/#!forum/golang-announce).
4. Install Go using your typical package manager. We'll revisit this later.
5. Install and get used to [`gofmt`](https://blog.golang.org/go-fmt-your-code),
   Go's code formatter. Use it on all your code.

After those three steps, spend some time writing a bit of toy code.

## Learn the Idioms

As a community, Go programmers tend to prefer strong conventions over freedom
of expression. Before you start writing Go professionally, get acquainted with
the community norms.

You don't need to read these documents in one sitting, but you should *read*
them...don't just skim the headings and move on.

1. As soon as you've read and written a little code, read [Effective
   Go](https://golang.org/doc/effective_go.html). This is canonical advice, and
   *everything* in it is widely accepted in the Go community.
2. For authoritative but slightly less canonical advice, browse the Go team's
   [common code review
   feedback](https://github.com/golang/go/wiki/CodeReviewComments) and their
   [guidance on package names](https://blog.golang.org/package-names).
3. Learn the common [idioms for working with
   slices](https://github.com/golang/go/wiki/SliceTricks).
4. Read some core material from the Go blog. Start with the [overall approach
   to concurrency](https://blog.golang.org/share-memory-by-communicating) and
   [error handling](https://blog.golang.org/error-handling-and-go) (including
   [part two](https://blog.golang.org/errors-are-values)), then move on to Go's
   take on [try-catch-finally
   blocks](https://blog.golang.org/defer-panic-and-recover).
5. Unit testing is built into the language, so get familiar with the [standard
   library's testing package](https://golang.org/pkg/testing/),
   [sub-tests](https://blog.golang.org/subtests), the [race
   detector](https://blog.golang.org/race-detector), and [how to measure code
   coverage](https://blog.golang.org/cover).
6. Documentation generation is also built into the language---familiarize
   yourself with the [comment
   conventions](https://blog.golang.org/godoc-documenting-go-code) and
   [testable examples](https://blog.golang.org/examples). Bookmark
   ~~[godoc.org](http://godoc.org)~~ [pkg.go.dev](https://pkg.go.dev), which
   serves documentation for the standard library and all open-source Go
   packages. 

## Learn the Details

Once you've picked up the basics and some common idioms, dive a little deeper.

1. At this point, it makes sense to stop relying on your package manager's Go
   distribution. Instead, use upstream directly---it'll be much easier to test
   beta releases and debug on older compilers. I use [Travis CI's
   gimme](https://github.com/travis-ci/gimme) script. It's a plain shell
   script, and it's the same script that manages Go versions in Travis builds.
2. Learn the details of Go [constants](https://blog.golang.org/constants),
   [strings](https://blog.golang.org/strings),
   [slices](https://blog.golang.org/go-slices-usage-and-internals) ([part
   two](https://blog.golang.org/slices)), and
   [maps](https://blog.golang.org/go-maps-in-action).
3. Grok some [more advanced concurrency
   patterns](https://blog.golang.org/advanced-go-concurrency-patterns),
   including [pipelines and cancellation](https://blog.golang.org/pipelines).
4. Gophers solve lots of problems by generating code, so read a little about
   [the `go generate` command](https://blog.golang.org/generate).
5. Try [profiling](https://blog.golang.org/profiling-go-programs) a running
   program.
6. Read about the design of the [garbage
   collector](https://blog.golang.org/go15gc).
7. Browse a bit of the [language specification](https://golang.org/ref/spec).
   It's very readable, so don't be afraid to consult it if you have specific
   questions later on.

If you're done with all that and still want more material, work your way
through [The Go Programming
Language](https://www.amazon.com/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440)
and take Bill Kennedy's [Ultimate
Go](https://www.safaribooksonline.com/library/view/ultimate-go-programming/9780134757476/)
class. Both resources are long but worthwhile.
