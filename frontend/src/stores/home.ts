import { defineStore } from "pinia";
import { files as filesApi } from "@/api";
import { fetchJSON, fetchURL, StatusError } from "@/api/utils";
import { useAuthStore } from "@/stores/auth";
import { summarize, type ScanSummary } from "@/utils/fileStats";

/*
 * Home dashboard data, kept for the whole session.
 *
 * Sections load independently and in parallel. When auto-refresh is on,
 * the scan also re-runs in the background while the app is open (even on
 * other pages), so Home is already up to date when you come back to it.
 * Background runs pause while the browser tab is hidden.
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
const TICK_MS = 30_000;

function readInterval(): number {
  try {
    const v = Number(localStorage.getItem(STORAGE_KEY));
    return (REFRESH_CHOICES as readonly number[]).includes(v) ? v : 0;
  } catch {
    return 0;
  }
}

const emptySection = <T>(): Section<T> => ({ status: "loading", data: null });

interface HomeState {
  owner: string;
  usage: Section<Usage>;
  scan: Section<ScanSummary>;
  shares: Section<number>;
  users: Section<number>;
  refreshing: boolean;
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
    /** Forget everything if a different user is signed in. */
    ensureOwner(owner: string) {
      if (this.owner === owner) return;
      controller?.abort();
      this.owner = owner;
      this.usage = emptySection();
      this.scan = emptySection();
      this.shares = emptySection();
      this.users = emptySection();
      this.refreshing = false;
      this.lastUpdated = null;
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
        run(this.scan, async (s) => {
          const res = await fetchURL("/api/resources/recursive/", {
            signal: s,
          });
          return summarize(await res.json());
        }),
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
    },

    /** Start the background scheduler. Safe to call more than once. */
    startAutoRefresh() {
      if (timer !== null) return;
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
    },

    stopAutoRefresh() {
      if (timer !== null) window.clearInterval(timer);
      timer = null;
      controller?.abort();
      controller = null;
      this.refreshing = false;
    },
  },
});
