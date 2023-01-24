// Content managed by Project Forge, see [projectforge.md] for details.
import {els} from "./dom";

export function timeInit() {
  (window as any).projectforge.relativeTime = relativeTime;
  els(".reltime").forEach(el => {
    relativeTime(el.dataset["time"] || "", el);
  })
}

export function utc(date: Date) {
  var utc = Date.UTC(
    date.getUTCFullYear(), date.getUTCMonth(),
    date.getUTCDate(), date.getUTCHours(),
    date.getUTCMinutes(), date.getUTCSeconds()
  );
  return new Date(utc).toISOString().substring(0, 19).replace("T", " ");
}

export function relativeTime(time: string, el?: HTMLElement): string {
  const str = (time || "").replace(/-/g, "/").replace(/[TZ]/g, " ") + " UTC";
  const date = new Date(str);
  const diff = (((new Date()).getTime() - date.getTime()) / 1000);
  const day_diff = Math.floor(diff / 86400);
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();

  if (isNaN(day_diff) || day_diff < 0 || day_diff >= 31) {
    return year.toString() + '-' + ((month < 10) ? '0' + month.toString() : month.toString()) + '-' + ((day < 10) ? '0' + day.toString() : day.toString());
  }

  let ret = "";
  let timeoutSeconds = 0;

  if (day_diff == 0) {
    if (diff < 5) {
      timeoutSeconds = 1;
      ret = "just now";
    } else if (diff < 60) {
      timeoutSeconds = 1;
      ret = Math.floor(diff) + " seconds ago";
    } else if (diff < 120) {
      timeoutSeconds = 10;
      ret = "1 minute ago";
    } else if (diff < 3600) {
      timeoutSeconds = 30;
      ret = Math.floor(diff / 60) + " minutes ago";
    } else if (diff < 7200) {
      timeoutSeconds = 60;
      ret = "1 hour ago";
    } else {
      timeoutSeconds = 60;
      ret = Math.floor(diff / 3600) + " hours ago";
    }
  } else if (day_diff == 1) {
    timeoutSeconds = 600;
    ret = "yesterday";
  } else if (day_diff < 7) {
    timeoutSeconds = 600;
    ret = day_diff + " days ago";
  } else {
    timeoutSeconds = 6000;
    ret = Math.ceil(day_diff / 7) + " weeks ago"
  }
  if (el) {
    el.innerText = ret;
    setTimeout(() => relativeTime(time, el), timeoutSeconds * 1000);
  }
  return ret;
}
