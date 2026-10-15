import { defineStore } from "pinia";
import { useAuthStore } from "./auth";

/*
 * Pinned files and folders, shown under "Favorites" in the sidebar.
 * Stored per user in this browser (localStorage); nothing is sent to the
 * server, so pins don't follow the user to another device.
 */

export interface Favorite {
  path: string; // e.g. "/Projects" or "/notes.txt"
  url: string; // router path, e.g. "/files/Projects/"
  name: string;
  isDir: boolean;
}

const MAX_FAVORITES = 50;

const storageKey = (username: string) => `favorites:${username}`;

const normalize = (path: string) =>
  path.length > 1 && path.endsWith("/") ? path.slice(0, -1) : path;

const isFavorite = (value: unknown): value is Favorite => {
  const f = value as Favorite;
  return (
    typeof f === "object" &&
    f !== null &&
    typeof f.path === "string" &&
    typeof f.url === "string" &&
    typeof f.name === "string" &&
    typeof f.isDir === "boolean"
  );
};

export const useFavoritesStore = defineStore("favorites", {
  state: (): { items: Favorite[]; owner: string | null } => ({
    items: [],
    owner: null,
  }),
  actions: {
    /** Load the current user's pins (no-op if already loaded for them). */
    load() {
      const username = useAuthStore().user?.username ?? "";
      if (this.owner === username) return;
      this.owner = username;
      this.items = [];
      if (!username) return;
      try {
        const parsed = JSON.parse(
          localStorage.getItem(storageKey(username)) ?? "[]"
        );
        if (Array.isArray(parsed)) {
          this.items = parsed.filter(isFavorite).slice(0, MAX_FAVORITES);
        }
      } catch {
        // unreadable or blocked storage: start with no pins
      }
    },
    save() {
      if (!this.owner) return;
      try {
        localStorage.setItem(
          storageKey(this.owner),
          JSON.stringify(this.items)
        );
      } catch {
        // storage blocked: pins last for this session only
      }
    },
    has(path: string) {
      const p = normalize(path);
      return this.items.some((f) => f.path === p);
    },
    toggle(item: Favorite) {
      this.load();
      const path = normalize(item.path);
      if (this.has(path)) {
        this.remove(path);
        return;
      }
      this.items = [{ ...item, path }, ...this.items].slice(0, MAX_FAVORITES);
      this.save();
    },
    remove(path: string) {
      const p = normalize(path);
      this.items = this.items.filter((f) => f.path !== p);
      this.save();
    },
  },
});
