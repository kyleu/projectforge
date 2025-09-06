import "./client.css";
import { appInit } from "./app";
import { audit } from "./audit";
import { autocompleteInit } from "./autocomplete";{{{ if .HasModule "richedit" }}}
import { editorInit } from "./editor";{{{ end }}}
import { flashInit } from "./flash";
import { formInit } from "./form";{{{ if .HasModule "jsx" }}}
import { JSX } from "./jsx";{{{ end }}}
import { linkInit } from "./link";
import { menuInit } from "./menu";
import { modalInit } from "./modal";
import { modeInit } from "./mode";{{{ if .HasModule "websocket" }}}
import { SocketMessage, socketInit } from "./socket";{{{ if .HasModule "process" }}}
import { socketLog } from "./socketlog";{{{ end }}}{{{ end }}}
import { tagsInit } from "./tags";
import { themeInit } from "./theme";
import { timeInit } from "./time";

declare global {
  // eslint-disable-next-line @typescript-eslint/consistent-type-definitions
  interface Window {
    {{{ .CleanKey }}}: {
      wireTime: (el: HTMLElement) => void;
      relativeTime: (el: HTMLElement) => string;
      autocomplete: (
        el: HTMLInputElement,
        url: string,
        field: string,
        title: (x: unknown) => string,
        val: (x: unknown) => string
      ) => void;
      setSiblingToNull: (el: HTMLElement) => void;
      initForm: (frm: HTMLFormElement) => void;
      flash: (key: string, level: "success" | "error", msg: string) => void;
      tags: (el: HTMLElement) => void;{{{ if .HasModule "websocket" }}}
      Socket: unknown;{{{ if .HasModule "process" }}}
      socketLog: (
        debug: boolean,
        parentEl: HTMLElement,
        terminal: boolean,
        url: string,
        extraHandlers: Array<(m: SocketMessage) => void>
      ) => void;{{{ end }}}{{{ end }}}
    };
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    audit: (s: string, ...args: any) => void;{{{ if .HasModule "jsx" }}}
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
