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
          <div class="sidebar__heading-row">
            <p class="sidebar__heading">{{ $t("sidebar.locations") }}</p>
            <button
              v-if="locations.length > 0"
              type="button"
              class="sidebar__heading-btn"
              :class="{ 'sidebar__heading-btn--on': editingLocations }"
              :aria-pressed="editingLocations"
              :aria-label="
                editingLocations
                  ? $t('sidebar.doneEditing')
                  : $t('sidebar.editLocations')
              "
              :title="
                editingLocations
                  ? $t('sidebar.doneEditing')
                  : $t('sidebar.editLocations')
              "
              @click="toggleEditing"
            >
              <i class="material-icons" aria-hidden="true">{{
                editingLocations ? "check" : "tune"
              }}</i>
            </button>
          </div>
          <p v-if="editingLocations" class="sidebar__reorder-hint">
            {{ $t("sidebar.editHint") }}
          </p>
          <router-link
            to="/files/"
            class="action"
            :class="{ 'action--current': currentLocation === '/files/' }"
            @click="closeHovers"
          >
            <i class="material-icons">dns</i>
            <span>{{ $t("sidebar.myFiles") }}</span>
          </router-link>
          <TransitionGroup
            ref="locList"
            tag="div"
            name="sidebar-reorder"
            class="sidebar__locations"
            :class="{ 'sidebar__locations--reordering': editingLocations }"
          >
            <div
              v-for="(loc, index) in visibleLocations"
              :key="loc.url"
              class="sidebar__item"
              :class="{
                'sidebar__item--dragging': dragUrl === loc.url,
                'sidebar__item--hidden': isHidden(loc),
              }"
            >
              <router-link
                :to="loc.url"
                class="action"
                :class="{ 'action--current': currentLocation === loc.url }"
                :title="loc.name"
                :tabindex="editingLocations ? -1 : undefined"
                @click.capture="onLocClick"
                @click="closeHovers"
                @pointerdown="onLocPointerDown($event, loc)"
                @pointermove="onLocPointerMove"
                @pointerup="cancelLongPress"
                @pointercancel="cancelLongPress"
                @pointerleave="cancelLongPress"
                @contextmenu="onLocContextMenu($event, loc)"
              >
                <i class="material-icons">folder</i>
                <span>{{ loc.name }}</span>
              </router-link>
              <button
                v-if="editingLocations"
                type="button"
                class="sidebar__visibility"
                :aria-pressed="!isHidden(loc)"
                :aria-label="
                  isHidden(loc)
                    ? $t('sidebar.showLocation', { name: loc.name })
                    : $t('sidebar.hideLocation', { name: loc.name })
                "
                :title="
                  isHidden(loc)
                    ? $t('sidebar.showInSidebar')
                    : $t('sidebar.hideFromSidebar')
                "
                @click="toggleHidden(loc)"
              >
                <i class="material-icons" aria-hidden="true">{{
                  isHidden(loc) ? "visibility_off" : "visibility"
                }}</i>
              </button>
              <button
                v-if="editingLocations"
                type="button"
                class="sidebar__handle"
                :data-url="loc.url"
                :aria-label="
                  $t('sidebar.dragToReorder', {
                    name: loc.name,
                    position: index + 1,
                    total: visibleLocations.length,
                  })
                "
                :title="$t('sidebar.dragToReorderShort')"
                @pointerdown="onHandlePointerDown($event, loc)"
                @pointermove="onHandlePointerMove"
                @pointerup="onHandlePointerUp"
                @pointercancel="onHandlePointerUp"
                @keydown="onHandleKeydown($event, loc)"
              >
                <i class="material-icons" aria-hidden="true">drag_indicator</i>
              </button>
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
                <button
                  type="button"
                  class="sidebar__ctx-item"
                  @mousedown.prevent="locCtxHide(loc)"
                >
                  <i class="material-icons" aria-hidden="true"
                    >visibility_off</i
                  >
                  <span>{{ $t("sidebar.hideFromSidebar") }}</span>
                </button>
              </div>
            </div>
          </TransitionGroup>
          <button
            v-if="!editingLocations && hiddenCount > 0"
            type="button"
            class="sidebar__hidden-link"
            @click="startEditing"
          >
            <i class="material-icons" aria-hidden="true">visibility_off</i>
            <span>{{
              $t("sidebar.hiddenCount", { count: hiddenCount }, hiddenCount)
            }}</span>
          </button>
          <span class="sidebar__sr-only" aria-live="polite">{{
            locationAnnouncement
          }}</span>
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
            to="/shares"
            class="action"
            :class="{ 'action--current': onRoute('/shares') }"
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
          <a v-else rel="noopener noreferrer" target="_blank" href="#"
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

// How long a press on a location must last to start editingLocations
const LONG_PRESS_MS = 500;

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
      // Long-press and drag bookkeeping; not rendered, so not reactive
      longPressTimer: null,
      longPressFired: false,
      pressStart: null,
      dragTop: 0,
      dragRowHeight: 0,
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

      // Editing locations: order and visibility
      editingLocations: false,
      dragUrl: null,
      locationAnnouncement: "",
      // Names of locations the user hid from the sidebar
      hiddenLocations: [],
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
    // Hidden locations stay listed while editing, so they can be shown again
    visibleLocations() {
      if (this.editingLocations) return this.locations;
      return this.locations.filter((l) => !this.isHidden(l));
    },
    hiddenCount() {
      return this.locations.filter((l) => this.isHidden(l)).length;
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
      const sorted = (items || [])
        .filter((i) => i.isDir && !(hideDot && i.name.startsWith(".")))
        .sort((a, b) => a.name.localeCompare(b.name))
        .map((i) => ({ name: i.name, url: i.url }));
      // Folders the user has arranged come first, in their order; new ones
      // follow alphabetically. Saved names that no longer exist are ignored.
      const order = this.loadLocationOrder();
      const rank = new Map(order.map((name, i) => [name, i]));
      this.locations = sorted.sort(
        (a, b) =>
          (rank.get(a.name) ?? Infinity) - (rank.get(b.name) ?? Infinity)
      );
    },
    // --- Location order (stored per user in this browser) ---
    locationOrderKey() {
      return `sidebar-location-order:${this.user?.username ?? ""}`;
    },
    loadLocationOrder() {
      try {
        const saved = JSON.parse(
          localStorage.getItem(this.locationOrderKey()) ?? "[]"
        );
        return Array.isArray(saved) ? saved : [];
      } catch {
        return [];
      }
    },
    saveLocationOrder() {
      try {
        localStorage.setItem(
          this.locationOrderKey(),
          JSON.stringify(this.locations.map((l) => l.name))
        );
      } catch {
        // Private mode or storage full: the order lasts for this visit only
      }
    },
    hiddenLocationsKey() {
      return `sidebar-hidden-locations:${this.user?.username ?? ""}`;
    },
    loadHiddenLocations() {
      try {
        const saved = JSON.parse(
          localStorage.getItem(this.hiddenLocationsKey()) ?? "[]"
        );
        this.hiddenLocations = Array.isArray(saved) ? saved : [];
      } catch {
        this.hiddenLocations = [];
      }
    },
    saveHiddenLocations() {
      try {
        localStorage.setItem(
          this.hiddenLocationsKey(),
          JSON.stringify(this.hiddenLocations)
        );
      } catch {
        // Private mode or storage full: the change lasts for this visit only
      }
    },
    isHidden(loc) {
      return this.hiddenLocations.includes(loc.name);
    },
    toggleHidden(loc) {
      const hidden = this.isHidden(loc);
      this.hiddenLocations = hidden
        ? this.hiddenLocations.filter((n) => n !== loc.name)
        : [...this.hiddenLocations, loc.name];
      this.saveHiddenLocations();
      this.locationAnnouncement = this.$t(
        hidden ? "sidebar.locationShown" : "sidebar.locationHidden",
        { name: loc.name }
      );
    },
    locCtxHide(loc) {
      this.closeCtxMenu();
      if (!this.isHidden(loc)) this.toggleHidden(loc);
    },
    startEditing() {
      if (this.editingLocations) return;
      this.closeCtxMenu();
      this.editingLocations = true;
      document.addEventListener("keydown", this.onEditEscape);
    },
    stopEditing() {
      if (!this.editingLocations) return;
      this.editingLocations = false;
      this.dragUrl = null;
      this.saveLocationOrder();
      document.removeEventListener("keydown", this.onEditEscape);
    },
    toggleEditing() {
      if (this.editingLocations) this.stopEditing();
      else this.startEditing();
    },
    onEditEscape(event) {
      if (event.key === "Escape") this.stopEditing();
    },
    moveLocation(url, toIndex) {
      const from = this.locations.findIndex((l) => l.url === url);
      const to = Math.max(0, Math.min(toIndex, this.locations.length - 1));
      if (from === -1 || from === to) return false;
      const [item] = this.locations.splice(from, 1);
      this.locations.splice(to, 0, item);
      this.locationAnnouncement = this.$t("sidebar.movedTo", {
        name: item.name,
        position: to + 1,
        total: this.locations.length,
      });
      return true;
    },
    // Long press on a location turns edit mode on
    onLocPointerDown(event, loc) {
      if (this.editingLocations || event.button !== 0) return;
      this.cancelLongPress();
      this.longPressFired = false;
      this.pressStart = { x: event.clientX, y: event.clientY };
      this.longPressTimer = window.setTimeout(() => {
        this.longPressTimer = null;
        this.longPressFired = true;
        this.startEditing();
        navigator.vibrate?.(15);
        this.$nextTick(() => this.focusHandle(loc.url));
      }, LONG_PRESS_MS);
    },
    onLocPointerMove(event) {
      if (!this.longPressTimer || !this.pressStart) return;
      const dx = event.clientX - this.pressStart.x;
      const dy = event.clientY - this.pressStart.y;
      // Moving the finger means scrolling, not pressing
      if (Math.hypot(dx, dy) > 8) this.cancelLongPress();
    },
    cancelLongPress() {
      if (this.longPressTimer) window.clearTimeout(this.longPressTimer);
      this.longPressTimer = null;
    },
    // Links don't navigate while editingLocations, or on the click that ends a
    // long press.
    onLocClick(event) {
      if (this.editingLocations || this.longPressFired) {
        event.preventDefault();
        event.stopImmediatePropagation();
        this.longPressFired = false;
      }
    },
    onLocContextMenu(event, loc) {
      // A long press on touch screens also fires contextmenu
      if (this.editingLocations || this.longPressFired) {
        event.preventDefault();
        return;
      }
      event.preventDefault();
      this.openLocCtxMenu(event, loc);
    },
    // Dragging by the handle. Rows are equal height, so the target slot is
    // worked out from the pointer position within the list.
    onHandlePointerDown(event, loc) {
      if (event.button !== 0) return;
      event.preventDefault();
      event.currentTarget.setPointerCapture(event.pointerId);
      const list = this.$refs.locList?.$el;
      const row = event.currentTarget.closest(".sidebar__item");
      if (!list || !row) return;
      this.dragUrl = loc.url;
      this.dragTop = list.getBoundingClientRect().top;
      this.dragRowHeight = row.offsetHeight + 2; // + the 2px gap
    },
    onHandlePointerMove(event) {
      if (!this.dragUrl) return;
      const index = Math.floor(
        (event.clientY - this.dragTop) / this.dragRowHeight
      );
      this.moveLocation(this.dragUrl, index);
    },
    onHandlePointerUp() {
      if (!this.dragUrl) return;
      this.dragUrl = null;
      this.saveLocationOrder();
    },
    // Keyboard: arrow keys move the focused location
    onHandleKeydown(event, loc) {
      const index = this.locations.findIndex((l) => l.url === loc.url);
      let target = null;
      if (event.key === "ArrowUp") target = index - 1;
      else if (event.key === "ArrowDown") target = index + 1;
      else if (event.key === "Home") target = 0;
      else if (event.key === "End") target = this.locations.length - 1;
      else if (event.key === "Enter") {
        event.preventDefault();
        this.stopEditing();
        return;
      }
      if (target === null) return;
      event.preventDefault();
      if (this.moveLocation(loc.url, target)) {
        this.saveLocationOrder();
        this.$nextTick(() => this.focusHandle(loc.url));
      }
    },
    focusHandle(url) {
      const handles =
        this.$refs.locList?.$el.querySelectorAll(".sidebar__handle");
      for (const h of handles ?? []) {
        if (h.dataset.url === url) h.focus();
      }
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
    // Closing the mobile sidebar ends editingLocations
    active(open) {
      if (!open && window.innerWidth <= 736) this.stopEditing();
    },
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
        this.loadHiddenLocations();
        this.fetchLocations();
      },
      immediate: true,
    },
  },
  unmounted() {
    this.abortOngoingFetchUsage();
    this.locationsAbortController.abort();
    this.cancelLongPress();
    document.removeEventListener("keydown", this.onEditEscape);
  },
};
</script>
