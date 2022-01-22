// Content managed by Project Forge, see [projectforge.md] for details.
import "./client.css"
import {menuInit} from "./menu";
import {flashInit} from "./flash";
import {editorInit} from "./editor";
import {sortableEdit, sortableInit} from "./sortable";
import {linkInit} from "./link";
import {themeInit} from "./theme";
import {setSiblingToNull} from "./form";
import {modeInit} from "./mode";
import {appInit} from "./app";

export function init(): void {
  (window as any).projectforge = {
    "sortableEdit": sortableEdit,
    "setSiblingToNull": setSiblingToNull
  };
  menuInit();
  modeInit();
  flashInit();
  linkInit();
  editorInit();
  sortableInit();
  themeInit();
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
