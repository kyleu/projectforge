import "./client.css";{{{ if .HasModule "jsx" }}}
import {JSX} from "./jsx";{{{ end }}}
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
import {Message, socketInit} from "./socket";{{{ end }}}{{{ if .HasModule "process" }}}
import {socketLog} from "./socketlog";{{{ end }}}
import {appInit} from "./app";

declare global {
  interface Window { // eslint-disable-line @typescript-eslint/consistent-type-definitions
    "{{{ .CleanKey }}}": {
      wireTime: (el: HTMLElement) => void;
      relativeTime: (el: HTMLElement) => string;
      autocomplete: (el: HTMLInputElement, url: string, field: string, title: (x: unknown) => string, val: (x: unknown) => string) => void;
      setSiblingToNull: (el: HTMLElement) => void;
      initForm: (frm: HTMLFormElement) => void;
      flash: (key: string, level: "success" | "error", msg: string) => void;
      tags: (el: HTMLElement) => void;{{{ if .HasModule "websocket" }}}
      Socket: unknown;{{{ end }}}{{{ if .HasModule "process" }}}
      socketLog: (debug: boolean, tbody: HTMLElement, url: string, extraHandlers: Array<(m: Message) => void>) => void;{{{ end }}}
    };
    audit: (s: string, ...args: any) => void; // eslint-disable-line @typescript-eslint/no-explicit-any{{{ if .HasModule "jsx" }}}
    JSX: (tag: string, attrs: unknown[]) => HTMLElement;{{{ end }}}
  }
}

export function init(): void {
  const [s, i] = formInit();
  const [wireTime, relativeTime] = timeInit();
  window.{{{ .CleanKey }}} = {
    wireTime: wireTime,
    relativeTime: relativeTime,
    autocomplete: autocompleteInit(),
    setSiblingToNull: s,
    initForm: i,
    flash: flashInit(),
    tags: tagsInit(){{{ if .HasModule "websocket" }}},
    Socket: socketInit(){{{ end }}}{{{ if .HasModule "process" }}},
    socketLog: socketLog{{{ end }}}
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
