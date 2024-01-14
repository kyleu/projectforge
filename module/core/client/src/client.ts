import "./client.css";{{{ if .HasModule "jsx" }}}
import {JSX} from "./jsx"; // eslint-disable-line @typescript-eslint/no-unused-vars{{{ end }}}
import {audit} from "./audit";
import {menuInit} from "./menu";
import {modeInit} from "./mode";
import {flashInit} from "./flash";
import {linkInit} from "./link";
import {timeInit} from "./time";
import {autocompleteInit} from "./autocomplete";
import {modalInit} from "./modal";
import {tagsInit} from "./tags";{{{ if .HasModule "richedit" }}}
import {editorInit} from "./editor";{{{ end }}}
import {formInit} from "./form";
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
    };
    audit: (s: string, ...args: any) => void; // eslint-disable-line @typescript-eslint/no-explicit-any{{{ if .HasModule "jsx" }}}
    JSX: (tag: string, attrs: unknown[]) => HTMLElement;{{{ end }}}
  }
}

export function init(): void {
  const [s, i] = formInit();
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
  themeInit();{{{ if .HasModule "richedit" }}}
  editorInit();{{{ end }}}
  window.audit = audit;{{{ if .HasModule "jsx" }}}
  window.JSX = JSX;{{{ end }}}
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
