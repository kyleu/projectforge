import SwiftUI
import ProjectforgeServer

@main
struct Project: App {
    init() {
        print("starting Project Forge...")
        let path = NSSearchPathForDirectoriesInDomains(.libraryDirectory, .userDomainMask, true)
        let port = ProjectforgeServer.CmdLib(path[0])
        print("Project Forge started on port [\(port)]")
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
