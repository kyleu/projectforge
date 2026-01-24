import { els } from "./dom";

export function utc(date: Date) {
  const u = Date.UTC(
    date.getUTCFullYear(),
    date.getUTCMonth(),
    date.getUTCDate(),
    date.getUTCHours(),
    date.getUTCMinutes(),
    date.getUTCSeconds()
  );
  return new Date(u).toISOString().substring(0, 19).replace("T", " ");
}

export function wireTime(el: HTMLElement) {
  const tsStr = el.dataset.timestamp;
  if (tsStr === "") {
    return;
  }
  if (!tsStr) {
    console.log("invalid timestamp [" + (tsStr ?? "") + "]", el);
    return;
  }
  const ts = new Date(tsStr);

  const tsEl = document.createElement("span");
  tsEl.title = ts.toISOString();
  if (el.classList.contains("millis")) {
    tsEl.innerText = ts.toISOString();
  } else {
    tsEl.innerText = ts.toLocaleString();
  }
  el.replaceChildren(tsEl);
}

function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  return (
    year.toString() +
    "-" +
    (month < 10 ? "0" + month.toString() : month.toString()) +
    "-" +
    (day < 10 ? "0" + day.toString() : day.toString())
  );
}

export function relativeTime(el: HTMLElement): string {
  const timestamp = el.dataset.timestamp;
  if (!timestamp) {
    el.innerText = "-";
    return "-";
  }
  const str = (timestamp ?? "").replace(/-/gu, "/").replace(/[TZ]/gu, " ") + " UTC";
  const date = new Date(str);
  if (isNaN(date.getTime())) {
    el.innerText = "-";
    return "-";
  }
  const diff = (new Date().getTime() - date.getTime()) / 1000;
  const dayDiff = Math.floor(diff / 86400);

  if (isNaN(dayDiff) || dayDiff < 0 || dayDiff >= 31) {
    const abs = formatDate(date);
    el.innerText = abs;
    return abs;
  }

  let ret: string;
  let timeoutSeconds: number;
  if (dayDiff === 0) {
    if (diff < 5) {
      timeoutSeconds = 1;
      ret = "just now";
    } else if (diff < 60) {
      timeoutSeconds = 1;
      ret = Math.floor(diff).toString() + " seconds ago";
    } else if (diff < 120) {
      timeoutSeconds = 10;
      ret = "1 minute ago";
    } else if (diff < 3600) {
      timeoutSeconds = 30;
      ret = Math.floor(diff / 60).toString() + " minutes ago";
    } else if (diff < 7200) {
      timeoutSeconds = 60;
      ret = "1 hour ago";
    } else {
      timeoutSeconds = 60;
      ret = Math.floor(diff / 3600).toString() + " hours ago";
    }
  } else if (dayDiff === 1) {
    timeoutSeconds = 600;
    ret = "yesterday";
  } else if (dayDiff < 7) {
    timeoutSeconds = 600;
    ret = Math.floor(dayDiff).toString() + " days ago";
  } else {
    timeoutSeconds = 6000;
    ret = Math.ceil(dayDiff / 7).toString() + " weeks ago";
  }
  el.innerText = ret;
  if (el.isConnected) {
    setTimeout(() => {
      if (el.isConnected) {
        relativeTime(el);
      }
    }, timeoutSeconds * 1000);
  }
  return ret;
}

export function timeInit(): [(el: HTMLElement) => void, (el: HTMLElement) => string] {
  els(".timestamp").forEach((el) => {
    wireTime(el);
  });
  els(".reltime").forEach((el) => {
    relativeTime(el);
  });
  return [wireTime, relativeTime];
}
