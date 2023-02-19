this was a debugging session for the ages, but it turns out that having both
a .sheet modifier in a view, and an @Environment reference to the presentationMode,
causes that strange nav bar title overwriting that we were seeing (only if the child
view is long enough to scroll down and compress the title into the top, and then you
hit back from that scrolled position...finding that took a while). However, it turns
out that we don't need presentationMode to dismiss ourselves if we plumb through an
isActive bool binding, which is what this does.

Simple repro code for this error is below for future reference.

```
import SwiftUI

struct ContentView: View {
    var body: some View {
        NavigationView {
            List {
                ForEach(1..<5) { i in
                    NavigationLink(destination: ChildView(i: i)) {
                        Text("test \(i)")
                    }
                }
            }
            .navigationBarTitle("list")
        }
    }
}

struct ChildView: View {
    // for some reason just holding this reference reproduces the issue
    @Environment(\.presentationMode) var presentationMode

    @State var i: Int
    @State var showSheet = false

    var body: some View {
        ScrollView {
            // it's important that this is long so that you can scroll down for the navbar title to be compressed
            ForEach(1..<100) { j in
                Text("\(self.i) ipsum lorem \(j)")
            }
        }
        // a sheet, even if never displayed, is also required for this repro
        .sheet(isPresented: $showSheet) {
            Text("Sheet)")
        }
        .navigationBarTitle("child \(i)")
    }
}
```
