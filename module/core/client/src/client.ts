import "./client.css";{{{ if .HasModule "jsx" }}}
import {JSX} from "./jsx"; // eslint-disable-line @typescript-eslint/no-unused-vars{{{ end }}}
import {menuInit} from "./menu";
import {modeInit} from "./mode";
import {flashInit} from "./flash";
import {linkInit} from "./link";
import {timeInit} from "./time";
import {autocompleteInit} from "./autocomplete";
import {modalInit} from "./modal";
import {tagsInit} from "./tags";
import {editorInit} from "./editor";
import {themeInit} from "./theme";{{{ if .HasModule "websocket" }}}
import {socketInit} from "./socket";{{{ end }}}
import {appInit} from "./app";

declare global {
  interface Window { // eslint-disable-line @typescript-eslint/consistent-type-definitions
    "{{{ .CleanKey }}}": {
      relativeTime: (time: string, el?: HTMLElement) => string;
      autocomplete: (el: HTMLInputElement, url: string, field: string, title: (x: unknown) => string, val: (x: unknown) => string) => void;
      setSiblingToNull: (el: HTMLElement) => void;
      initForm: (frm: HTMLFormElement) => void;
      flash: (key: string, level: "success" | "error", msg: string) => void;
      tags: (el: HTMLElement) => void;{{{ if .HasModule "websocket" }}}
      Socket: unknown;{{{ end }}}
    };{{{ if .HasModule "jsx" }}}
    "JSX": (tag: string, attrs: unknown[]) => HTMLElement;{{{ end }}}
  }
}

export function init(): void {
  const [s, i] = editorInit();
  window.{{{ .CleanKey }}} = {
    relativeTime: timeInit(),
    autocomplete: autocompleteInit(),
    setSiblingToNull: s,
    initForm: i,
    flash: flashInit(),
    tags: tagsInit(){{{ if .HasModule "websocket" }}},
    Socket: socketInit(){{{ end }}}
  };
  menuInit();
  modeInit();
  linkInit();
  modalInit();
  themeInit();{{{ if .HasModule "jsx" }}}
  window.JSX = JSX;{{{ end }}}
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
