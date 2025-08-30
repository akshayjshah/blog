# Go Time: Protobuf & gRPC

<audio data-theme="night" data-src="https://changelog.com/gotime/256/embed"
src="https://op3.dev/e/https://cdn.changelog.com/uploads/gotime/256/go-time-256.mp3"
preload="none" class="changelog-episode" controls></audio>
<script async src="//cdn.changelog.com/embed.js"></script>

I was interviewed on [episode 256 of the Go Time
podcast](https://changelog.com/gotime/256)! We talked about Protocol Buffers,
gRPC, and some common misconceptions about both. I tried to emphasize that:

* Protobuf and gRPC are useful because they reduce developer toil. Any
  performance improvements are nice, but they're secondary.
* Protobuf isn't just for binary data --- it works just as well with JSON.
* gRPC is a simple, HTTP-based protocol. Google's gRPC implementations optimize
  for Google's problems, but there's plenty room for [gRPC implementations that
  work better](https://github.com/bufbuild/connect-go) for most Gophers.

Overall, I found it surprisingly hard to stay on track during the discussion.
I have *so* much more respect for previous Go Time guests!
