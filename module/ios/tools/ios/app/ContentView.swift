import SwiftUI

struct ContentView: View {
    let u: URLRequest

    var body: some View {
        WebView(url: self.u) //.edgesIgnoringSafeArea(.all)
    }

    init(url: URLRequest) {
        self.u = url
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        let u = URL.init(string: "http://localhost:{{{ .Port }}}")!
        ContentView(url: URLRequest(url: u))
    }
}
