// Content managed by Project Forge, see [projectforge.md] for details.
package dev.projectforge

import android.os.Bundle
import android.util.Log
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.appcompat.app.AppCompatActivity
import android.content.Intent
import android.net.Uri

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Log.i("projectforge", "Project Forge is starting")
        val path = getFilesDir().getAbsolutePath()
        val port = cmd.Cmd.lib(path)
        Log.i("projectforge", "Project Forge has started with path [${path}] on port [${port}]")
        setContentView(R.layout.activity_main)

        val webView: WebView = findViewById(R.id.webview)
        val client = object : WebViewClient() {
            override fun shouldOverrideUrlLoading(view: WebView, url: String): Boolean {
                return if (url.startsWith("http://localhost")) {
                    false
                } else {
                    view.context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(url)))
                    true
                }
            }
        }
        webView.setWebViewClient(client)
        val settings = webView.getSettings();

        settings.loadsImagesAutomatically = true;
        settings.javaScriptEnabled = true;
        settings.domStorageEnabled = true;

        webView.loadUrl("http://localhost:${port}/")
    }
}
