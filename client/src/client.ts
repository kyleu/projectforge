// Content managed by Project Forge, see [projectforge.md] for details.
import "./client.css";
import {JSX} from "./jsx"; // eslint-disable-line @typescript-eslint/no-unused-vars
import {audit} from "./audit";
import {menuInit} from "./menu";
import {modeInit} from "./mode";
import {flashInit} from "./flash";
import {linkInit} from "./link";
import {timeInit} from "./time";
import {autocompleteInit} from "./autocomplete";
import {modalInit} from "./modal";
import {tagsInit} from "./tags";
import {editorInit} from "./editor";
import {formInit} from "./form";
import {themeInit} from "./theme";
import {socketInit} from "./socket";
import {appInit} from "./app";

declare global {
  interface Window { // eslint-disable-line @typescript-eslint/consistent-type-definitions
    "projectforge": {
      relativeTime: (time: string, el?: HTMLElement) => string;
      autocomplete: (el: HTMLInputElement, url: string, field: string, title: (x: unknown) => string, val: (x: unknown) => string) => void;
      setSiblingToNull: (el: HTMLElement) => void;
      initForm: (frm: HTMLFormElement) => void;
      flash: (key: string, level: "success" | "error", msg: string) => void;
      tags: (el: HTMLElement) => void;
      Socket: unknown;
    };
    audit: (s: string, ...args: any) => void; // eslint-disable-line @typescript-eslint/no-explicit-any
    JSX: (tag: string, attrs: unknown[]) => HTMLElement;
  }
}

export function init(): void {
  const [s, i] = formInit();
  window.projectforge = {
    relativeTime: timeInit(),
    autocomplete: autocompleteInit(),
    setSiblingToNull: s,
    initForm: i,
    flash: flashInit(),
    tags: tagsInit(),
    Socket: socketInit()
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
