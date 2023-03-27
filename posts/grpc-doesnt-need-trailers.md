# gRPC Doesn't Need Trailers

Among gRPC's many imperfections, none is worse than its use of HTTP trailers.
Perversely, gRPC fans seem eager to defend this mistake. As an example, a
former gRPC maintainer recently published <em>[Why Does gRPC Insist on
Trailers][nonsense]</em>. He argues that trailers protect clients from abruptly
dropped TCP connections and compensate for oddities in the binary encoding of
Protocol Buffers. Bluntly, both of these arguments are nonsense.

**gRPC doesn't need HTTP trailers.** In this post, I'll explain:

1. What trailers are and how gRPC uses them,
2. Why they're unnecessary,
3. How they impede gRPC adoption, 
4. How Google could fix gRPC in a minor release, and
5. How we can do even better without Google.

## What are trailers?

Trailers have been lurking in the HTTP specification since at least HTTP/1.1,
released in 1997 (when they were called "footers"). Simply put, they're headers
that come _after_ the request or response body. The recent [HTTP Semantics
RFC](https://www.rfc-editor.org/rfc/rfc9110.html#section-6.5) suggests that
they "can be useful for supplying message integrity checks, digital signatures,
delivery metrics, or post-processing status information." If you've been
slinging HTTP for decades and haven't heard of trailers, you're not alone ---
most HTTP/1.1 implementations don't support them, so they're rarely used.

Why, then, would gRPC rely on trailers? Because gRPC supports streaming
responses, in which the server writes multiple records to the response body.
Imagine that a gRPC server is preparing to stream a large collection to a
client. The server connects to the database, executes a query, and begins
inspecting the results. Everything is going well, so the server sends a `200
OK` status code and some headers. One by one, the server begins reading records
from the database and writing them to the response body. Then the database
crashes. How should the server tell the client that something has gone wrong?
The client has already received a `200 OK` HTTP status code, so it's too late
to send a `500 Internal Server Error`. Because the server has already started
sending the response body, it's also too late to send more headers. The
server's only options are to send the error as the last portion of the response
body or to send it in trailers.

gRPC chooses to use trailers. Responses include a gRPC-specific status code in
the `grpc-status` trailer and a description of the error (if any) in the
`grpc-message` trailer. Even successful responses must explicitly set
`grpc-status` to 0.

## Are trailers necessary?

Of course not! But let's at least address the two arguments above.

1. **Clients don't need trailers to detect dropped TCP connections.** In this
   argument, gRPC apologists claim that trailers help clients detect incomplete
   responses. "What if," they say, "the server --- or some proxy --- crashes
   partway through a streaming response and drops the TCP connection? By
   insisting on an explicit status code in trailers, clients can detect a
   prematurely-terminated response body." This is plausible-sounding,
   especially when accompanied by an HTTP/1.1 example, but it's nonsense. gRPC
   requires at least HTTP/2, and both HTTP/2 and HTTP/3 handle this explicitly:
   [every HTTP/2 frame includes a byte of bitwise
   flags](https://www.rfc-editor.org/rfc/rfc9113#section-4.1), and the frame
   types used for headers, trailers, _and_ body data all include an explicit
   `END_STREAM` flag used to cleanly terminate the response. If the client sees a
   TCP connection drop before it receives an HTTP/2 frame with `END_STREAM`
   set, it knows that the response is incomplete --- no trailers needed.
2. **Nothing about Protocol Buffers requires trailers.** In this variant of the
   first argument, gRPC apologists argue that detecting dropped TCP connections
   is _especially_ important when using Protocol Buffers. "If gRPC only
   supported JSON," they say, "clients would detect many incomplete responses
   by noticing unbalanced curly braces. But Protocol Buffer messages don't have
   explicit delimiters, so we _really_ need to rely on trailers to detect
   dropped connections." But not only does HTTP/2 provide an unambiguous way to
   detect dropped connections, the gRPC protocol doesn't rely on
   encoding-specific delimiters to find message boundaries. Instead, it
   [prefixes each message in a stream with its
   length](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#:~:text=The%20repeated%20sequence%20of%20Length%2DPrefixed%2DMessage%20items%20is%20delivered%20in%20DATA%20frames).
   Clients can easily detect response bodies that end before delivering the
   promised quantity of data. Again, trailers don't add any safety.

## Why are trailers bad?

If trailers were unnecessary but harmless, I wouldn't be ranting to internet
strangers. But trailers aren't harmless: they make it difficult to add gRPC
APIs to existing applications. Is your Python application built with Django,
Flask, or FastAPI? Too bad --- WSGI and ASGI don't support trailers, so your
application can't handle gRPC-flavored HTTP. Trying to call your gRPC server
from an iPhone? Sorry, `URLSession` doesn't support trailers either. Rather
than adding a few new routes to your existing server and client, you're stuck
building a entirely new application for RPC.

To support trailers, your new application uses a gRPC-specific HTTP stack. But
apart from supporting trailers, your new stack is less capable than your old
one: usually, gRPC's HTTP implementation can _only_ serve RPCs over HTTP/2. If
you also want to serve an HTML page, receive a file upload, support HTTP/1.1 or
HTTP/3, or just handle an HTTP `GET`, you're out of luck. In practice, adopting
gRPC _requires_ a multi-service backend architecture. 

These pains are most acute on the web. Like many other clients, `fetch` doesn't
support trailers. Unlike mobile or backend applications, though, web
applications _can't_ bundle an alternate, gRPC-specific HTTP client. Instead,
they're forced to proxy requests through Envoy, which translates them on the
fly from a trailer-free protocol to standard gRPC. Envoy is a perfectly fine
proxy, but it's a lot to configure and manage in production if you're only
using it to work around gRPC's quirks. And of course, no web developer enjoys
running a C++ proxy during local development.

In short, relying on trailers abandons one of HTTP's key strengths: the ready
availability of interoperable servers and clients.

## Could Google fix gRPC?

When Google designed gRPC, trailer support had just been added to the `fetch`
specification. If the Chrome, Firefox, Safari, and Edge teams had followed
through and implemented the proposed APIs, other HTTP implementations might
have followed their lead. Instead, browser makers withdrew their support for
the new APIs, and they were formally removed from the specification in late
2019.

It's now 2023. Trailers aren't coming to browsers --- or to most other HTTP
implementations --- for years, if ever. Even Cloudflare, a multi-billion dollar
internet infrastructure company, [doesn't have end-to-end support for
trailers](https://blog.cloudflare.com/road-to-grpc/). The gRPC team should
confront this reality and add support for a second, trailer-free protocol to
their servers and clients.

[gRPC-Web](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md) is the
pragmatic choice for a second protocol. It's nearly identical to standard gRPC,
except that it encodes status metadata at the end of the response body rather
than in trailers. It uses a different Content-Type, so servers could
automatically handle the new protocol alongside the old. Clients could opt into
the new protocol with a configuration toggle. Implementations wouldn't need any
other user-visible API changes, so these improvements could ship in a
backward-compatible minor release. And because gRPC-Web is already under the
gRPC umbrella, we wouldn't need to convince Google to adopt any outside ideas.
(gRPC-Web also drops gRPC's strict HTTP/2 requirement, which is nice but
unnecessary to mitigate the trailers fiasco.)

If today's gRPC implementations embraced the gRPC-Web protocol, new
implementations could _only_ support gRPC-Web. All of a sudden, `grpc-rails`
and similar framework integrations would be feasible. Browsers could call gRPC
backends directly. iOS applications could drop their multi-megabyte dependency
on `SwiftNIO`. Without trailers, gRPC could meet developers where they are.

Microsoft seems to agree with this assessment: they've built support for the
gRPC-Web protocol into `grpc-dotnet`. If you'd like Google to do the same,
[upvote issue 29818 in the main gRPC
repository](https://github.com/grpc/grpc/issues/29818).

## Could we do better without Google?

gRPC-Web might be the pragmatic choice, but it still leaves a lot to be
desired. What if we were bolder? To _really_ improve upon gRPC, we'd use
different protocols for streaming and request-response RPCs. The streaming
protocol would be similar to gRPC-Web, but we'd bring the request-respose
protocol closer to familiar, resource-oriented HTTP:

* We'd support HTTP/1.1 and HTTP/2.
* We'd use meaningful HTTP status codes.
* We'd dispense with trailers and end-of-body metadata and just rely on
  headers.
* We wouldn't need to length-prefix messages, so the body could be plain JSON
  or binary Protocol Buffer. That lets us use recognizable Content-Types,
  like `application/json`.
* We'd use the standard `Accept-Encoding` header, so web applications benefit
  from compressed responses.
* We'd support `GET` requests for cacheable RPCs. With some care, we could
  avoid having these `GET` requests trigger CORS preflight from browsers.
* For servers using Protocol Buffer schemas, we'd encourage implementations to
  support both binary and JSON payloads by default (using the [canonical JSON
  mapping](https://protobuf.dev/programming-guides/proto3/#json)).

None of these changes affect the protocol's efficiency, but they eliminate most
of gRPC's fussiness. Creating a `User` becomes a cURL one-liner:

```
curl --json '{"name": "Akshay"}' https://api.acme.com/user.v1/Create
```

This protocol _just works_ because it's boring. It works with human-readable
JSON and optimized binary encodings. It works with cURL and `requests`. It
works with `fetch` and browsers' built-in debuggers. It works with `URLSession`
and Charles Proxy. It works with Rails, Django, FastAPI, Laravel, and Express.
It works with CDNs and browser caches. It works with standard penetration
testing toolkits.

I can't imagine Google embracing a protocol that's so different from today's
gRPC, especially if it requires HTTP/1.1 support, but you can try it _today_:
use [Connect](https://connect.build). Connect servers and clients support the
full gRPC protocol, gRPC-Web, _and_ the [simpler protocol][connect-protocol] we
just outlined. Implementations are available in Go, TypeScript, Swift, and
Kotlin.

[nonsense]: https://carlmastrangelo.com/blog/why-does-grpc-insist-on-trailers
[connect-protocol]: https://connect.build/docs/protocol/
