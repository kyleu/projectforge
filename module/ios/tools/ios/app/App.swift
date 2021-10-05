import SwiftUI
import {{{ .KeyProper }}}Server

@main
struct Project: App {
    init() {
        print("starting {{{ .Name }}}...")
        let path = NSSearchPathForDirectoriesInDomains(.libraryDirectory, .userDomainMask, true)
        let port = {{{ .KeyProper }}}Server.CmdLib(path[0])
        print("{{{ .Name }}} started on port [\(port)]")
        let url = URL.init(string: "http://localhost:\(port)/")!
        self.cv = ContentView(url: URLRequest(url: url))
    }

    var cv: ContentView

    var body: some Scene {
        WindowGroup {
            cv
        }
    }
}
