// Content managed by Project Forge, see [projectforge.md] for details.
import "./client.css"
import {menuInit} from "./menu";
import {modeInit} from "./mode";
import {flashInit} from "./flash";
import {linkInit} from "./link";
import {modalInit} from "./modal";
import {editorInit} from "./editor";
import {themeInit} from "./theme";
import {appInit} from "./app";

export function init(): void {
  (window as any).projectforge = {};
  menuInit();
  modeInit();
  flashInit();
  linkInit();
  modalInit();
  editorInit();
  themeInit();
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
