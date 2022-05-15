import "./client.css"
import {menuInit} from "./menu";
import {modeInit} from "./mode";
import {flashInit} from "./flash";
import {linkInit} from "./link";
import {modalInit} from "./modal";
import {editorInit} from "./editor";
import {themeInit} from "./theme";{{{ if .HasModule "websocket" }}}
import {socketInit} from "./socket";{{{ end }}}
import {appInit} from "./app";

export function init(): void {
  (window as any).{{{ .CleanKey }}} = {};
  menuInit();
  modeInit();
  flashInit();
  linkInit();
  modalInit();
  editorInit();
  themeInit();{{{ if .HasModule "websocket" }}}
  socketInit();{{{ end }}}
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
