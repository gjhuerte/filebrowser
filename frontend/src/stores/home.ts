import { defineStore } from "pinia";
import type { ScanSummary } from "@/utils/fileStats";

/*
 * Last dashboard result for this session. Lets Home show the previous
 * numbers instantly when revisited while fresh ones load in the background.
 */
export interface HomeCache {
  owner: string;
  usage: { used: number; total: number } | null;
  scan: ScanSummary | null;
  shares: number | null;
  users: number | null;
}

export const useHomeStore = defineStore("home", {
  state: (): { cache: HomeCache | null } => ({ cache: null }),
  actions: {
    cacheFor(owner: string): HomeCache {
      if (!this.cache || this.cache.owner !== owner) {
        this.cache = {
          owner,
          usage: null,
          scan: null,
          shares: null,
          users: null,
        };
      }
      return this.cache;
    },
  },
});
