# gRPC Doesn't Need Trailers

Among gRPC's many imperfections, none is worse than its dogged insistence on
HTTP trailers. When detractors highlight this wart, apologists inevitably argue
that trailers aren't a problem --- in fact, they're essential! Most recently, a
former gRPC maintainer published <em>[Why Does gRPC Insist on
Trailers](https://carlmastrangelo.com/blog/why-does-grpc-insist-on-trailers)</em>.
In it, he argues that trailers protect clients from abruptly dropped TCP
connections and compensate for oddities in the binary encoding of Protocol
Buffers. Bluntly, both of these arguments are nonsense.

**gRPC doesn't need HTTP trailers.** In this post, I'll explain:

1. what trailers are and how gRPC uses them,
2. why they're unnecessary,
3. how they degrade developer experience, and
4. how gRPC could migrate to a better approach.

## What are trailers?

Trailers have been lurking in the HTTP specification since at least HTTP/1.1,
released in 1997. Simply put, they're headers that come _after_ the request or
response body. The recent [HTTP Semantics
RFC](https://www.rfc-editor.org/rfc/rfc9110.html#section-6.5) suggests that
they "can be useful for supplying message integrity checks, digital signatures,
delivery metrics, or post-processing status information." If you've been
slinging HTTP for decades and have never heard of trailers, you're not alone
--- most HTTP/1.1 implementations don't support them, so software using HTTP
tends to avoid them.

Why, then, would gRPC rely on trailers? Because gRPC supports streaming
responses, in which the server writes multiple records to the response body.
Imagine that a well-intentioned gRPC server is preparing to stream a large
collection of records to a client. The server connects to the database,
executes a query, and begins inspecting the results. Everything is copacetic,
so the server sends a `200 OK` status code and a handful of headers. One by
one, the server begins reading records from the database and writing them to
the response body. Then the database crashes. How should the server tell the
client that something has gone terribly wrong? The client has already received
a `200 OK` HTTP status code, so it's too late to send a `500 Internal Server
Error`. Because the server has already started sending the response body, it's
also too late to send additional headers. In the parlance of HTTP RFCs, the
database crash is clearly "post-processing status information." The server's
only options are to send the error as the last portion of the response body or
to send it in trailers.

gRPC chooses to use trailers. Responses include a gRPC-specific error code in
the `grpc-status` trailer and a description of the error in the `grpc-message`
trailer. (For now, avert your eyes from the undocumented trash fire of
`grpc-status-details-bin`.)

## Are trailers necessary?

Of course not! But I didn't write all this just to make unsupported assertions,
so I'll directly address the two arguments above.

1. **Clients don't need trailers to detect dropped TCP connections.** In this
   argument, gRPC apologists claim that trailers help clients detect incomplete
   responses. "What if," they say, "the server --- or some proxy --- crashes
   partway through a streaming response and drops the TCP connection? By
   insisting on an explicit status code in the HTTP trailers, clients can
   detect a prematurely-terminated response body." This is plausible-sounding,
   especially when accompanied by an HTTP/1.1 example, but it's nonsense. gRPC
   requires at least HTTP/2, and both HTTP/2 and HTTP/3 handle this explicitly:
   [every HTTP/2 frame includes a byte of bitwise
   flags](https://www.rfc-editor.org/rfc/rfc9113#section-4.1), and the frame
   types used for headers, trailers, _and_ body data all include an explicit
   `END_STREAM` flag used to cleanly terminate the response. gRPC's designers
   were clearly aware of this mechanism, because the
   [protocol](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md)
   states that "for responses end-of-stream is indicated by the
   presence of the `END_STREAM` flag on the last received `HEADERS` frame that
   carries Trailers." (HTTP/2 uses `HEADERS` frames for both headers and
   trailers.) Applied more liberally, the same mechanism can unambiguously
   catch server and proxy crashes without trailers --- if the client sees the
   TCP connection drop before it receives an HTTP/2 frame with `END_STREAM`
   set, it knows that the response is incomplete.
2. **Nothing about Protocol Buffers requires trailers.** In this variant of the
   first argument, gRPC apologists argue that detecting dropped TCP connections
   is _especially_ important when using Protocol Buffers.
   "If gRPC only supported JSON," they say, "clients would detect many
   incomplete responses by noticing an unpaired `{`. But Protocol Buffer
   messages don't have explicit delimiters, so we _really_ need to rely on
   trailer

   But not only does
   HTTP/2 provide an unambiguous way to detect dropped connections, the gRPC
   protocol doesn't rely on encoding-specific delimiters to find message
   boundaries. Instead, it [prefixes each message in a stream with its
   length](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#:~:text=The%20repeated%20sequence%20of%20Length%2DPrefixed%2DMessage%20items%20is%20delivered%20in%20DATA%20frames).
   Clients can easily detect response bodies that end before delivering the
   promised quantity of data.
