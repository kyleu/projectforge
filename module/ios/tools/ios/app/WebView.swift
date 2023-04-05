import SwiftUI
import WebKit

class FullScreenWKWebView: WKWebView {
    override var safeAreaInsets: UIEdgeInsets {
        return UIEdgeInsets(top: 0, left: 0, bottom: 0, right: 0)
    }

}

struct WebViewWrapper : UIViewRepresentable {
    let view: FullScreenWKWebView

    init(url: URLRequest) {
        self.view = FullScreenWKWebView()
        self.view.load(url)
    }

    func makeUIView(context: Context) -> WKWebView  {
        return view
    }

    func updateUIView(_ uiView: WKWebView, context: Context) {
    }
}

struct WebView: View {
    let u: URLRequest

    var body: some View {
        WebViewWrapper(url: u)
    }

    init(url: URLRequest) {
        self.u = url
    }
}
