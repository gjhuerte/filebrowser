<template>
  <div v-show="active" @click="closeHovers" class="overlay"></div>
  <nav :class="{ active }" class="sidebar">
    <template v-if="isLoggedIn">
      <div class="sidebar__section">
        <button
          class="action"
          :class="{ 'action--current': onFiles }"
          :aria-current="onFiles ? 'page' : undefined"
          @click="toRoot"
          :aria-label="$t('sidebar.myFiles')"
          :title="$t('sidebar.myFiles')"
        >
          <i class="material-icons">folder</i>
          <span>{{ $t("sidebar.myFiles") }}</span>
        </button>

        <button
          v-if="user.perm.admin"
          class="action"
          :class="{ 'action--current': onSettings }"
          :aria-current="onSettings ? 'page' : undefined"
          @click="toGlobalSettings"
          :aria-label="$t('sidebar.settings')"
          :title="$t('sidebar.settings')"
        >
          <i class="material-icons">settings</i>
          <span>{{ $t("sidebar.settings") }}</span>
        </button>
      </div>

      <div v-if="user.perm.create" class="sidebar__section">
        <button
          @click="showHover('newDir')"
          class="action"
          :aria-label="$t('sidebar.newFolder')"
          :title="$t('sidebar.newFolder')"
        >
          <i class="material-icons">create_new_folder</i>
          <span>{{ $t("sidebar.newFolder") }}</span>
        </button>

        <button
          @click="showHover('newFile')"
          class="action"
          :aria-label="$t('sidebar.newFile')"
          :title="$t('sidebar.newFile')"
        >
          <i class="material-icons">note_add</i>
          <span>{{ $t("sidebar.newFile") }}</span>
        </button>
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

    <div class="sidebar__footer">
      <div v-if="isFiles && !disableUsedPercentage" class="sidebar__usage">
        <p class="sidebar__usage-text">
          {{ $t("sidebar.diskUsed", { used: usage.used, total: usage.total }) }}
        </p>
        <progress-bar :val="usage.usedPercentage" size="small"></progress-bar>
      </div>

      <div v-if="isLoggedIn" class="sidebar__account">
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
          @click="logout"
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
          href="https://github.com/filebrowser/filebrowser"
          >File Browser</a
        >
        <span>{{ version }}</span>
        <span aria-hidden="true">·</span>
        <a href="#" @click.prevent="help">{{ $t("sidebar.help") }}</a>
      </p>
    </div>
  </nav>
</template>

<script>
import { reactive } from "vue";
import { mapActions, mapState } from "pinia";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import * as auth from "@/utils/auth";
import {
  version,
  signup,
  hideLoginButton,
  disableExternal,
  disableUsedPercentage,
  noAuth,
  logoutPage,
  loginPage,
} from "@/utils/constants";
import { files as api } from "@/api";
import ProgressBar from "@/components/ProgressBar.vue";
import prettyBytes from "pretty-bytes";

const USAGE_DEFAULT = { used: "0 B", total: "0 B", usedPercentage: 0 };

export default {
  name: "sidebar",
  setup() {
    const usage = reactive(USAGE_DEFAULT);
    return { usage, usageAbortController: new AbortController() };
  },
  components: {
    ProgressBar,
  },
  inject: ["$showError"],
  computed: {
    ...mapState(useAuthStore, ["user", "isLoggedIn"]),
    ...mapState(useFileStore, ["isFiles", "reload"]),
    ...mapState(useLayoutStore, ["currentPromptName"]),
    active() {
      return this.currentPromptName === "sidebar";
    },
    signup: () => signup,
    hideLoginButton: () => hideLoginButton,
    version: () => version,
    disableExternal: () => disableExternal,
    disableUsedPercentage: () => disableUsedPercentage,
    canLogout: () => !noAuth && (loginPage || logoutPage !== "/login"),
    onFiles() {
      return this.$route.path.startsWith("/files");
    },
    onSettings() {
      return (
        this.$route.path.startsWith("/settings") &&
        !this.$route.path.startsWith("/settings/profile")
      );
    },
    onProfile() {
      return this.$route.path.startsWith("/settings/profile");
    },
    initial() {
      return (this.user?.username ?? "?").charAt(0).toUpperCase();
    },
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers", "showHover"]),
    abortOngoingFetchUsage() {
      this.usageAbortController.abort();
    },
    async fetchUsage() {
      const path = this.$route.path.endsWith("/")
        ? this.$route.path
        : this.$route.path + "/";
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
      } finally {
        return Object.assign(this.usage, usageStats);
      }
    },
    toRoot() {
      this.$router.push({ path: "/files" });
      this.closeHovers();
    },
    toAccountSettings() {
      this.$router.push({ path: "/settings/profile" });
      this.closeHovers();
    },
    toGlobalSettings() {
      this.$router.push({ path: "/settings/global" });
      this.closeHovers();
    },
    help() {
      this.showHover("help");
    },
    logout: auth.logout,
  },
  watch: {
    $route: {
      handler(to) {
        if (to.path.includes("/files")) {
          this.fetchUsage();
        }
      },
      immediate: true,
    },
  },
  unmounted() {
    this.abortOngoingFetchUsage();
  },
};
</script>
