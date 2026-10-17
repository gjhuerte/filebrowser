<template>
  <div v-show="active" @click="closeHovers" class="overlay"></div>
  <nav :class="{ active }" class="sidebar">
    <template v-if="isLoggedIn">
      <div class="sidebar__main">
        <div class="sidebar__section">
          <router-link
            to="/home"
            class="action"
            :class="{ 'action--current': onHome }"
            :aria-current="onHome ? 'page' : undefined"
            @click="closeHovers"
          >
            <i class="material-icons">space_dashboard</i>
            <span>{{ $t("sidebar.home") }}</span>
          </router-link>
        </div>

        <!-- Favorites -->
        <div class="sidebar__group">
          <p class="sidebar__heading">{{ $t("sidebar.favorites") }}</p>
          <div v-for="fav in favorites" :key="fav.path" class="sidebar__item">
            <router-link
              v-if="renamingPath !== fav.path"
              :to="fav.url"
              class="action"
              :class="{ 'action--current': isCurrent(fav.url, true) }"
              :title="fav.path"
              @click="closeHovers"
              @dblclick.prevent="startRename(fav)"
              @contextmenu.prevent="openCtxMenu($event, fav)"
            >
              <i class="material-icons">{{
                fav.isDir ? "folder" : "description"
              }}</i>
              <span>{{ fav.customName || fav.name }}</span>
            </router-link>
            <div v-else class="sidebar__rename-row">
              <i class="material-icons">{{
                fav.isDir ? "folder" : "description"
              }}</i>
              <input
                ref="renameInput"
                class="sidebar__rename-input"
                type="text"
                v-model="renameValue"
                @keydown.enter="commitRename(fav.path)"
                @keydown.esc="cancelRename"
                @blur="commitRename(fav.path)"
              />
            </div>
            <div
              v-if="ctxMenuPath === fav.path"
              class="sidebar__ctx"
              :style="{ top: ctxMenuY + 'px', left: ctxMenuX + 'px' }"
            >
              <button
                type="button"
                class="sidebar__ctx-item"
                @mousedown.prevent="ctxRename(fav)"
              >
                <i class="material-icons" aria-hidden="true">edit</i>
                <span>{{ $t("sidebar.editName") }}</span>
              </button>
              <button
                type="button"
                class="sidebar__ctx-item sidebar__ctx-item--danger"
                @mousedown.prevent="ctxRemove(fav.path)"
              >
                <i class="material-icons" aria-hidden="true">delete</i>
                <span>{{ $t("sidebar.removeFavorite") }}</span>
              </button>
            </div>
          </div>
          <p v-if="favorites.length === 0" class="sidebar__empty">
            <i class="material-icons" aria-hidden="true">star_border</i>
            <span>{{ $t("sidebar.noFavorites") }}</span>
          </p>
        </div>

        <!-- Locations: the root and its top-level folders -->
        <div class="sidebar__group">
          <p class="sidebar__heading">{{ $t("sidebar.locations") }}</p>
          <router-link
            to="/files/"
            class="action"
            :class="{ 'action--current': currentLocation === '/files/' }"
            @click="closeHovers"
          >
            <i class="material-icons">dns</i>
            <span>{{ $t("sidebar.myFiles") }}</span>
          </router-link>
          <div
            v-for="loc in visibleLocations"
            :key="loc.url"
            class="sidebar__item"
          >
            <router-link
              :to="loc.url"
              class="action"
              :class="{ 'action--current': currentLocation === loc.url }"
              :title="loc.name"
              @click="closeHovers"
              @contextmenu.prevent="openLocCtxMenu($event, loc)"
            >
              <i class="material-icons">folder</i>
              <span>{{ loc.name }}</span>
            </router-link>
            <div
              v-if="locCtxMenuUrl === loc.url"
              class="sidebar__ctx"
              :style="{ top: ctxMenuY + 'px', left: ctxMenuX + 'px' }"
            >
              <button
                v-if="canShare"
                type="button"
                class="sidebar__ctx-item"
                @mousedown.prevent="locCtxShare(loc)"
              >
                <i class="material-icons" aria-hidden="true">link</i>
                <span>{{ $t("buttons.share") }}</span>
              </button>
              <button
                type="button"
                class="sidebar__ctx-item"
                @mousedown.prevent="locCtxFavorite(loc)"
              >
                <i class="material-icons" aria-hidden="true">{{
                  isFavorited(loc) ? "star" : "star_border"
                }}</i>
                <span>{{
                  isFavorited(loc)
                    ? $t("sidebar.removeFavorite")
                    : $t("sidebar.addFavorite")
                }}</span>
              </button>
              <button
                type="button"
                class="sidebar__ctx-item"
                @mousedown.prevent="locCtxInfo(loc)"
              >
                <i class="material-icons" aria-hidden="true">info</i>
                <span>{{ $t("prompts.fileInfo") }}</span>
              </button>
            </div>
          </div>

        </div>
      </div>

      <div class="sidebar__footer">
        <!-- Administration -->
        <div v-if="user.perm.admin || canShare" class="sidebar__group">
          <p v-if="user.perm.admin" class="sidebar__heading">
            {{ $t("sidebar.administration") }}
          </p>
          <router-link
            v-if="user.perm.admin"
            to="/admin/settings/general"
            class="action"
            :class="{ 'action--current': onRoute('/admin/settings') }"
            @click="closeHovers"
          >
            <i class="material-icons">tune</i>
            <span>{{ $t("sidebar.systemSettings") }}</span>
          </router-link>
          <router-link
            v-if="user.perm.admin"
            to="/admin/users"
            class="action"
            :class="{ 'action--current': onRoute('/admin/users') }"
            @click="closeHovers"
          >
            <i class="material-icons">group</i>
            <span>{{ $t("sidebar.users") }}</span>
          </router-link>
          <router-link
            v-if="canShare"
            to="/settings/shares"
            class="action"
            :class="{ 'action--current': onRoute('/settings/shares') }"
            @click="closeHovers"
          >
            <i class="material-icons">link</i>
            <span>{{ $t("sidebar.sharedLinks") }}</span>
          </router-link>
        </div>

        <!-- Tools -->
        <div v-if="canUseShell" class="sidebar__group">
          <p class="sidebar__heading">{{ $t("sidebar.tools") }}</p>
          <button
            type="button"
            class="action"
            :class="{ 'action--current': showShell }"
            :aria-pressed="showShell"
            @click="openShell"
          >
            <i class="material-icons">terminal</i>
            <span>{{ $t("sidebar.terminal") }}</span>
          </button>
        </div>

        <!-- Storage -->
        <div v-if="!disableUsedPercentage" class="sidebar__usage">
          <p class="sidebar__usage-title">
            <i class="material-icons" aria-hidden="true">storage</i>
            {{ $t("sidebar.storage") }}
          </p>
          <progress-bar :val="usage.usedPercentage" size="small"></progress-bar>
          <p class="sidebar__usage-text">
            {{
              $t("sidebar.diskUsed", { used: usage.used, total: usage.total })
            }}
          </p>
        </div>

        <!-- Account -->
        <div class="sidebar__account">
          <button
            class="sidebar__profile"
            :class="{ 'sidebar__profile--current': onProfile }"
            @click="toAccountSettings"
            :title="$t('settings.profileSettings')"
          >
            <span class="avatar" aria-hidden="true">{{ initial }}</span>
            <span class="sidebar__username">{{ user.username }}</span>
          </button>
          <button
            v-if="canLogout"
            @click="confirmLogout"
            class="action sidebar__logout"
            id="logout"
            :aria-label="$t('sidebar.logout')"
            :title="$t('sidebar.logout')"
          >
            <i class="material-icons">logout</i>
          </button>
        </div>

        <p class="sidebar__credits">
          <span v-if="disableExternal">File Browser</span>
          <a
            v-else
            rel="noopener noreferrer"
            target="_blank"
            href="#"
            >File Browser</a
          >
          <span>{{ version }}</span>
          <span aria-hidden="true">·</span>
          <a href="#" @click.prevent="help">{{ $t("sidebar.help") }}</a>
        </p>
      </div>
    </template>

    <div v-else class="sidebar__section">
      <router-link
        v-if="!hideLoginButton"
        class="action"
        to="/login"
        :aria-label="$t('sidebar.login')"
        :title="$t('sidebar.login')"
      >
        <i class="material-icons">login</i>
        <span>{{ $t("sidebar.login") }}</span>
      </router-link>

      <router-link
        v-if="signup"
        class="action"
        to="/login"
        :aria-label="$t('sidebar.signup')"
        :title="$t('sidebar.signup')"
      >
        <i class="material-icons">person_add</i>
        <span>{{ $t("sidebar.signup") }}</span>
      </router-link>
    </div>
  </nav>
</template>

<script>
import { reactive } from "vue";
import { mapActions, mapState } from "pinia";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { useFavoritesStore } from "@/stores/favorites";

import {
  version,
  signup,
  hideLoginButton,
  disableExternal,
  disableUsedPercentage,
  noAuth,
  logoutPage,
  loginPage,
  enableExec,
} from "@/utils/constants";
import { files as api } from "@/api";
import ProgressBar from "@/components/ProgressBar.vue";
import prettyBytes from "pretty-bytes";

const USAGE_DEFAULT = { used: "0 B", total: "0 B", usedPercentage: 0 };

// Compare router paths regardless of encoding and trailing slash
const normalizePath = (path) => {
  let p = path;
  try {
    p = decodeURIComponent(path);
  } catch {
    // keep as-is if it isn't valid percent-encoding
  }
  return p.endsWith("/") ? p : p + "/";
};

export default {
  name: "sidebar",
  setup() {
    const usage = reactive({ ...USAGE_DEFAULT });
    return {
      usage,
      usageAbortController: new AbortController(),
      locationsAbortController: new AbortController(),
    };
  },
  data() {
    return {
      locations: [],

      usageLoaded: false,
      renamingPath: null,
      renameValue: "",
      ctxMenuPath: null,
      ctxMenuX: 0,
      ctxMenuY: 0,
      locCtxMenuUrl: null,
    };
  },
  components: {
    ProgressBar,
  },
  inject: ["$showError"],
  computed: {
    ...mapState(useAuthStore, ["user", "isLoggedIn"]),
    ...mapState(useFileStore, ["isFiles", "reload", "req"]),
    ...mapState(useLayoutStore, ["currentPromptName", "showShell"]),
    ...mapState(useFavoritesStore, { favorites: "items" }),
    active() {
      return this.currentPromptName === "sidebar";
    },
    signup: () => signup,
    hideLoginButton: () => hideLoginButton,
    version: () => version,
    disableExternal: () => disableExternal,
    disableUsedPercentage: () => disableUsedPercentage,
    canLogout: () => !noAuth && (loginPage || logoutPage !== "/login"),
    canShare() {
      return !!(this.user?.perm.share && this.user?.perm.download);
    },
    canUseShell() {
      return enableExec && !!this.user?.perm.execute;
    },
    onHome() {
      return this.$route.path === "/home";
    },
    onProfile() {
      return this.$route.path.startsWith("/settings/profile");
    },
    visibleLocations() {
      return this.locations;
    },
    // The location the current page lives in (longest matching prefix)
    currentLocation() {
      if (!this.$route.path.startsWith("/files")) return null;
      const here = normalizePath(this.$route.path);
      const match = this.visibleLocations.find((loc) =>
        here.startsWith(normalizePath(loc.url))
      );
      return match ? match.url : "/files/";
    },
    initial() {
      return (this.user?.username ?? "?").charAt(0).toUpperCase();
    },
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers", "showHover", "toggleShell"]),
    ...mapActions(useFavoritesStore, {
      loadFavorites: "load",
      removeFavorite: "remove",
      renameFavorite: "rename",
      toggleFavorite: "toggle",
    }),
    onRoute(prefix) {
      return this.$route.path.startsWith(prefix);
    },
    isCurrent(url) {
      return normalizePath(this.$route.path) === normalizePath(url);
    },
    abortOngoingFetchUsage() {
      this.usageAbortController.abort();
    },
    async fetchUsage(routePath) {
      const path = routePath.endsWith("/") ? routePath : routePath + "/";
      let usageStats = USAGE_DEFAULT;
      if (this.disableUsedPercentage) {
        return Object.assign(this.usage, usageStats);
      }
      try {
        this.abortOngoingFetchUsage();
        this.usageAbortController = new AbortController();
        const usage = await api.usage(path, this.usageAbortController.signal);
        usageStats = {
          used: prettyBytes(usage.used, { binary: true }),
          total: prettyBytes(usage.total, { binary: true }),
          usedPercentage: Math.round((usage.used / usage.total) * 100),
        };
        this.usageLoaded = true;
      } finally {
        return Object.assign(this.usage, usageStats);
      }
    },
    setLocations(items) {
      const hideDot = this.user?.hideDotfiles;
      this.locations = (items || [])
        .filter((i) => i.isDir && !(hideDot && i.name.startsWith(".")))
        .sort((a, b) => a.name.localeCompare(b.name))
        .map((i) => ({ name: i.name, url: i.url }));
    },
    async fetchLocations() {
      if (!this.isLoggedIn) return;
      try {
        this.locationsAbortController.abort();
        this.locationsAbortController = new AbortController();
        const root = await api.fetch(
          "/files/",
          this.locationsAbortController.signal
        );
        this.setLocations(root.items);
      } catch {
        // Locations are a convenience; the list simply stays as it was
      }
    },
    openShell() {
      this.toggleShell();
      this.closeHovers();
    },
    toAccountSettings() {
      this.$router.push({ path: "/settings/profile" });
      this.closeHovers();
    },
    help() {
      this.showHover("help");
    },
    // Ask first; the dialog does the actual logout
    confirmLogout() {
      this.showHover("logout");
    },
    openCtxMenu(event, fav) {
      this.ctxMenuPath = fav.path;
      this.ctxMenuX = event.clientX;
      this.ctxMenuY = event.clientY;
      const dismiss = (e) => {
        // Close on any click/contextmenu outside the menu
        if (!e.target.closest(".sidebar__ctx")) {
          this.closeCtxMenu();
        }
        document.removeEventListener("mousedown", dismiss, true);
        document.removeEventListener("contextmenu", dismiss, true);
      };
      this.$nextTick(() => {
        document.addEventListener("mousedown", dismiss, true);
        document.addEventListener("contextmenu", dismiss, true);
      });
    },
    closeCtxMenu() {
      this.ctxMenuPath = null;
      this.locCtxMenuUrl = null;
    },
    ctxRename(fav) {
      this.closeCtxMenu();
      this.startRename(fav);
    },
    ctxRemove(path) {
      this.closeCtxMenu();
      this.removeFavorite(path);
    },
    // --- Location context menu ---
    openLocCtxMenu(event, loc) {
      this.locCtxMenuUrl = loc.url;
      this.ctxMenuX = event.clientX;
      this.ctxMenuY = event.clientY;
      const dismiss = (e) => {
        if (!e.target.closest(".sidebar__ctx")) {
          this.closeCtxMenu();
        }
        document.removeEventListener("mousedown", dismiss, true);
        document.removeEventListener("contextmenu", dismiss, true);
      };
      this.$nextTick(() => {
        document.addEventListener("mousedown", dismiss, true);
        document.addEventListener("contextmenu", dismiss, true);
      });
    },
    isFavorited(loc) {
      const path = loc.url.replace(/^\/files/, "").replace(/\/$/, "") || "/";
      return useFavoritesStore().has(path);
    },
    locCtxFavorite(loc) {
      this.closeCtxMenu();
      const path = loc.url.replace(/^\/files/, "").replace(/\/$/, "") || "/";
      this.toggleFavorite({
        path,
        url: loc.url,
        name: loc.name,
        isDir: true,
      });
    },
    locCtxShare(loc) {
      this.closeCtxMenu();
      this.$router.push({ path: loc.url });
      this.$nextTick(() => {
        this.showHover("share");
      });
    },
    locCtxInfo(loc) {
      this.closeCtxMenu();
      this.$router.push({ path: loc.url });
      this.$nextTick(() => {
        this.showHover("info");
      });
    },
    startRename(fav) {
      this.renamingPath = fav.path;
      this.renameValue = fav.customName || fav.name;
      this.$nextTick(() => {
        const input = this.$refs.renameInput;
        if (input) {
          const el = Array.isArray(input) ? input[0] : input;
          el.focus();
          el.select();
        }
      });
    },
    commitRename(path) {
      if (this.renamingPath !== path) return;
      this.renameFavorite(path, this.renameValue);
      this.renamingPath = null;
      this.renameValue = "";
    },
    cancelRename() {
      this.renamingPath = null;
      this.renameValue = "";
    },
  },
  watch: {
    $route: {
      handler(to) {
        // Usage reflects the disk of the folder you're in; elsewhere, the root
        if (to.path.startsWith("/files")) {
          this.fetchUsage(to.path);
        } else if (!this.usageLoaded) {
          this.fetchUsage("/files/");
        }
      },
      immediate: true,
    },
    // Keep locations fresh whenever the root folder is (re)loaded
    req(value) {
      if (value?.isDir && value.path === "/") this.setLocations(value.items);
    },
    user: {
      handler(value) {
        if (!value) return;
        this.loadFavorites();
        this.fetchLocations();
      },
      immediate: true,
    },
  },
  unmounted() {
    this.abortOngoingFetchUsage();
    this.locationsAbortController.abort();
  },
};
</script>
