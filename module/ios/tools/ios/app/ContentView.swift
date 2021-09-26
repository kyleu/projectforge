import SwiftUI

struct ContentView: View {
    var body: some View {
        Text("{{{ .Name }}}!").padding()
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
