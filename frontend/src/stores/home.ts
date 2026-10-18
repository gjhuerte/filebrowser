import { defineStore } from "pinia";
import { files as filesApi, summary as summaryApi } from "@/api";
import type { SummaryStatus } from "@/api/summary";
import { fetchJSON, StatusError } from "@/api/utils";
import { useAuthStore } from "@/stores/auth";
import type { ScanSummary } from "@/utils/fileStats";

/*
 * Home dashboard data.
 *
 * The last result is kept in the browser, so opening Home or reloading the
 * page shows it straight away without rescanning. A new scan runs only when
 * the user presses Refresh, on the auto-refresh schedule (if it is on), or
 * the very first time Home has nothing to show.
 *
 * The file scan runs on the server as a background job: we start it, then
 * poll until it is done. Nothing waits on one long request, so a reverse
 * proxy's timeout cannot cut it off, and only the small summary is sent.
 *
 * Sections load independently and in parallel. Scheduled scans run in the
 * background while the app is open (even on other pages) and pause while the
 * browser tab is hidden.
 */

export type SectionStatus = "loading" | "ready" | "error";
export interface Section<T> {
  status: SectionStatus;
  data: T | null;
}
export interface Usage {
  used: number;
  total: number;
}

/** Choices offered in the UI, in minutes. 0 = off. */
export const REFRESH_CHOICES = [0, 5, 15, 30, 60] as const;
const STORAGE_KEY = "home-refresh-interval";
// Also cleared on logout in utils/auth.ts
const CACHE_KEY = "home-dashboard-cache";
const TICK_MS = 30_000;
const POLL_MS = 1500;
/** Polls that may fail in a row (e.g. a proxy hiccup) before giving up. */
const MAX_POLL_FAILURES = 5;

/** Resolve after ms, or reject as soon as signal aborts. */
function wait(ms: number, signal: AbortSignal) {
  return new Promise<void>((resolve, reject) => {
    const id = window.setTimeout(resolve, ms);
    signal.addEventListener(
      "abort",
      () => {
        window.clearTimeout(id);
        reject(new StatusError("000 No connection", 0, true));
      },
      { once: true }
    );
  });
}

function readInterval(): number {
  try {
    const v = Number(localStorage.getItem(STORAGE_KEY));
    return (REFRESH_CHOICES as readonly number[]).includes(v) ? v : 0;
  } catch {
    return 0;
  }
}

const emptySection = <T>(): Section<T> => ({ status: "loading", data: null });

/** What is saved between page loads: the data of each section. */
interface SavedDashboard {
  owner: string;
  lastUpdated: number;
  usage: Usage | null;
  scan: ScanSummary | null;
  shares: number | null;
  users: number | null;
}

function readSaved(owner: string): SavedDashboard | null {
  try {
    const saved = JSON.parse(localStorage.getItem(CACHE_KEY) ?? "null");
    return saved && saved.owner === owner ? (saved as SavedDashboard) : null;
  } catch {
    return null;
  }
}

const restored = <T>(data: T | null): Section<T> =>
  data === null ? emptySection<T>() : { status: "ready", data };

interface HomeState {
  owner: string;
  usage: Section<Usage>;
  scan: Section<ScanSummary>;
  shares: Section<number>;
  users: Section<number>;
  refreshing: boolean;
  /** What the running server scan has found so far, or null. */
  progress: SummaryStatus["progress"] | null;
  lastUpdated: number | null;
  /** Auto-refresh interval in minutes (0 = off). */
  interval: number;
  /** Ticks every 30s so "Updated 3 min ago" stays current. */
  now: number;
}

// Browser objects stay outside reactive state.
let controller: AbortController | null = null;
let timer: number | null = null;

export const useHomeStore = defineStore("home", {
  state: (): HomeState => ({
    owner: "",
    usage: emptySection(),
    scan: emptySection(),
    shares: emptySection(),
    users: emptySection(),
    refreshing: false,
    progress: null,
    lastUpdated: null,
    interval: readInterval(),
    now: Date.now(),
  }),

  getters: {
    nextRefreshAt(state): number | null {
      if (!state.interval || !state.lastUpdated) return null;
      return state.lastUpdated + state.interval * 60_000;
    },
  },

  actions: {
    /**
     * Switch to owner's dashboard: restore their saved result if there is
     * one, otherwise start empty.
     */
    ensureOwner(owner: string) {
      if (this.owner === owner) return;
      controller?.abort();
      const saved = owner ? readSaved(owner) : null;
      this.owner = owner;
      this.usage = restored(saved?.usage ?? null);
      this.scan = restored(saved?.scan ?? null);
      this.shares = restored(saved?.shares ?? null);
      this.users = restored(saved?.users ?? null);
      this.refreshing = false;
      this.progress = null;
      this.lastUpdated = saved?.lastUpdated ?? null;
    },

    /** True when nothing has been loaded or saved yet. */
    isEmpty(): boolean {
      return this.lastUpdated === null && this.scan.data === null;
    },

    save() {
      if (!this.owner || this.lastUpdated === null) return;
      const saved: SavedDashboard = {
        owner: this.owner,
        lastUpdated: this.lastUpdated,
        usage: this.usage.data,
        scan: this.scan.data,
        shares: this.shares.data,
        users: this.users.data,
      };
      try {
        localStorage.setItem(CACHE_KEY, JSON.stringify(saved));
      } catch {
        // Private mode or storage full: the result lasts for this visit only
      }
    },

    setInterval(minutes: number) {
      this.interval = minutes;
      try {
        localStorage.setItem(STORAGE_KEY, String(minutes));
      } catch {
        /* private mode — keep it for this session only */
      }
    },

    async refresh() {
      const auth = useAuthStore();
      if (!auth.user) return;
      this.ensureOwner(auth.user.username);

      controller?.abort();
      controller = new AbortController();
      const { signal } = controller;
      this.refreshing = true;

      const run = async <T>(
        section: Section<T>,
        loader: (s: AbortSignal) => Promise<T>
      ) => {
        if (section.data === null) section.status = "loading";
        try {
          const value = await loader(signal);
          if (signal.aborted) return;
          section.data = value;
          section.status = "ready";
        } catch (e) {
          if (signal.aborted || (e instanceof StatusError && e.is_canceled))
            return;
          // Keep showing stale data if we have it
          if (section.data === null) section.status = "error";
        }
      };

      const perm = auth.user.perm;
      const jobs: Promise<void>[] = [
        run(this.usage, (s) => filesApi.usage("/files/", s) as Promise<Usage>),
        run(this.scan, (s) => this.scanOnServer(s)),
      ];
      if (perm.share && perm.download) {
        jobs.push(
          run(
            this.shares,
            async (s) =>
              (await fetchJSON<unknown[]>("/api/shares", { signal: s })).length
          )
        );
      }
      if (perm.admin) {
        jobs.push(
          run(
            this.users,
            async (s) =>
              (await fetchJSON<unknown[]>("/api/users", { signal: s })).length
          )
        );
      }

      await Promise.allSettled(jobs);
      if (signal.aborted) return;
      this.refreshing = false;
      this.lastUpdated = Date.now();
      this.now = this.lastUpdated;
      this.save();
    },

    /**
     * Start (or join) the server scan and poll until it finishes. If the
     * user leaves, the server keeps going and the next visit picks it up.
     */
    async scanOnServer(signal: AbortSignal): Promise<ScanSummary> {
      try {
        let st = await summaryApi.start(signal);
        let failures = 0;
        while (st.scanning) {
          this.progress = st.progress;
          await wait(POLL_MS, signal);
          try {
            st = await summaryApi.status(signal);
            failures = 0;
          } catch (e) {
            if (signal.aborted || ++failures >= MAX_POLL_FAILURES) throw e;
          }
        }
        // A failed scan leaves the previous result in place; say so instead
        // of passing it off as fresh.
        if (st.error || !st.summary) throw new Error(st.error || "No result");
        return st.summary;
      } finally {
        if (!signal.aborted) this.progress = null;
      }
    },

    /** Start the background scheduler. Safe to call more than once. */
    startAutoRefresh() {
      if (timer !== null) return;
      // Load the saved result first, so the schedule continues from the last
      // scan rather than from this page load.
      const auth = useAuthStore();
      if (auth.user) this.ensureOwner(auth.user.username);
      const tick = () => {
        this.now = Date.now();
        const due = this.nextRefreshAt;
        if (
          due !== null &&
          this.now >= due &&
          !this.refreshing &&
          document.visibilityState === "visible"
        ) {
          this.refresh();
        }
      };
      timer = window.setInterval(tick, TICK_MS);
      document.addEventListener("visibilitychange", tick);
      // A scan that fell due while the app was closed runs now
      tick();
    },

    stopAutoRefresh() {
      if (timer !== null) window.clearInterval(timer);
      timer = null;
      controller?.abort();
      controller = null;
      this.refreshing = false;
      this.progress = null;
    },
  },
});
