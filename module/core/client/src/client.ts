import "./client.css"{{{ if .HasModule "jsx" }}}
import {JSX} from "./jsx";{{{ end }}}
import {menuInit} from "./menu";
import {modeInit} from "./mode";
import {flashInit} from "./flash";
import {linkInit} from "./link";
import {modalInit} from "./modal";
import {tagsInit} from "./tags";
import {editorInit} from "./editor";
import {themeInit} from "./theme";{{{ if .HasModule "websocket" }}}
import {socketInit} from "./socket";{{{ end }}}
import {appInit} from "./app";

export function init(): void {
  (window as any).{{{ .CleanKey }}} = {};{{{ if .HasModule "jsx" }}}
  (window as any).JSX = JSX;{{{ end }}}
  menuInit();
  modeInit();
  flashInit();
  linkInit();
  modalInit();
  tagsInit();
  editorInit();
  themeInit();{{{ if .HasModule "websocket" }}}
  socketInit();{{{ end }}}
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
