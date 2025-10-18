import "./client.css";
import { appInit } from "./app";
import { audit } from "./audit";
import { autocompleteInit } from "./autocomplete";
import { editorInit } from "./editor";
import { flashInit } from "./flash";
import { formInit } from "./form";
import { JSX } from "./jsx";
import { linkInit } from "./link";
import { menuInit } from "./menu";
import { modalInit } from "./modal";
import { modeInit } from "./mode";
import { SocketMessage, socketInit } from "./socket";
import { socketLog } from "./socketlog";
import { tagsInit } from "./tags";
import { themeInit } from "./theme";
import { timeInit } from "./time";

declare global {
  interface Window {
    projectforge: {
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
      tags: (el: HTMLElement) => void;
      Socket: unknown;
      socketLog: (
        debug: boolean,
        parentEl: HTMLElement,
        terminal: boolean,
        url: string,
        extraHandlers: ((m: SocketMessage) => void)[]
      ) => void;
    };
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    audit: (s: string, ...args: any) => void;
    JSX: (tag: string, attrs: Record<string, unknown>, ...args: Node[]) => HTMLElement;
  }
}

export function init(): void {
  const [s, i] = formInit();
  const [wireTime, relativeTime] = timeInit();
  window.projectforge = {
    wireTime: wireTime,
    relativeTime: relativeTime,
    autocomplete: autocompleteInit(),
    setSiblingToNull: s,
    initForm: i,
    flash: flashInit(),
    tags: tagsInit(),
    Socket: socketInit(),
    socketLog: socketLog
  };
  menuInit();
  modeInit();
  linkInit();
  modalInit();
  themeInit();
  editorInit();
  window.audit = audit;
  window.JSX = JSX;
  appInit();
}

document.addEventListener("DOMContentLoaded", init);
