import "./client.css"
import {menuInit} from "./menu";
import {flashInit} from "./flash";
import {editorInit} from "./editor";
import {linkInit} from "./link";
import {themeInit} from "./theme";
import {modeInit} from "./mode";{{{ if .HasModule "websocket" }}}
import {socketInit} from "./socket";{{{ end }}}
import {appInit} from "./app";

export function init(): void {
  (window as any).{{{ .CleanKey }}} = {};
  menuInit();
  modeInit();
  flashInit();
  linkInit();
  editorInit();
  themeInit();{{{ if .HasModule "websocket" }}}
  socketInit();{{{ end }}}
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
