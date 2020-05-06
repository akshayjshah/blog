# Lazy-Loading Data with SwiftUI and Combine

Apple clearly wants us to pair SwiftUI with Combine for asynchronous network
operations, but using the two frameworks together is surprisingly awkward. The
`@ObservedObject` property wrapper makes it easy for SwiftUI views to
auto-update in response to changing data, but there's no obvious way to *start*
fetching data.

Kicking off network requests in initializers is tempting, but it's common for
views to be initialized and never presented: for example,  `NavigationLink`
eagerly constructs the destination view. If you're building a list of
navigation links, making network requests in `init` can easily load thousands
of unused resources.

Instead of making network requests in `init`, we could push them into each
view's
[`onAppear`](https://developer.apple.com/documentation/swiftui/text/3276931-onappear)
hook. This avoids fetching unused data, but it scatters imperative networking
code in hooks throughout your code. It also takes some extra bookkeeping to
avoid unnecessarily refetching resources shared between views.

Ideally, we'd have the best of both worlds: resources automatically fetched on
demand, but without constant use of the `onAppear` hook. In [Swift Talk
160](http://talk.objc.io/episodes/S01E160-lazy-data-loading), Florian Kugler
and Chris Eidhof dove into the beta releases of Combine and SwiftUI and
discovered a way to defer fetching data until a subscriber is waiting. With the
changes made before Combine's GA release, their approach can be made even
simpler.

```swift
import Combine
import Foundation

public final class Lazy<A>: ObservableObject {
    // We want to lazy-load data automatically; that is, we should
    // defer network requests until the data is required to render
    // an on-screen view. We do this by wrapping the actual publisher
    // in an event handler that fetches data on subscription.
    private let changes = ObservableObjectPublisher()
    private var subscribers = 0
    public var objectWillChange = ObservableObjectPublisher().handleEvents()

    // Intentionally oversimplified loading - use an abstraction from
    // your networking layer here.
    private let load: ()->A
    public var value: A? {
        willSet {
            DispatchQueue.main.async {
                self.changes.send()
            }
        }
    }

    public init(load: @escaping ()->A) {
        self.load = load
        self.objectWillChange = self.changes.handleEvents(
            receiveSubscription: { _ in
                let isFirst = self.subscribers == 0
                self.subscribers += 1
                if isFirst {
                    self.reload()
                }
            },
            receiveCancel: {
                self.subscribers -= 1
            }
        )
    }

    public func reload() {
        value = load()
    }
}
```

Using the `Lazy` class is short and declarative, and data is loaded only when
the observing view is presented.

```swift
import SwiftUI

struct ExpensiveView: View {
    @ObservedObject var n = Lazy<Int>() {
        print("Simulating network latency...")
        sleep(2)
        return Int.random(in: 1...100)
    }

    var body: some View {
        if n.value != nil {
            return Text("\(n.value!)")
        } else {
            return Text("Loading...")
        }
    }
}
```

This implementation doesn't include any synchronization, but it seems likely
that `value` and `subscribers` need to be protected by a `DispatchSemaphore`.
I'm new to Apple's developer ecosystem, but I've been unable to confirm this
suspicion.
