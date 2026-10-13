<template>
  <div class="dashboard">
    <header-bar showMenu showLogo />

    <div class="settings">
      <h1 class="settings__title">{{ t("sidebar.settings") }}</h1>

      <div class="tabs" role="navigation" :aria-label="t('sidebar.settings')">
        <router-link
          to="/settings/profile"
          class="tabs__tab"
          :aria-label="t('settings.profileSettings')"
          :title="t('settings.profileSettings')"
          :class="{ 'tabs__tab--active': $route.path === '/settings/profile' }"
        >
          <i class="material-icons" aria-hidden="true">person</i>
          <span>{{ t("settings.profileSettings") }}</span>
        </router-link>
        <router-link
          v-if="user?.perm.share"
          to="/settings/shares"
          class="tabs__tab"
          :aria-label="t('settings.shareManagement')"
          :title="t('settings.shareManagement')"
          :class="{ 'tabs__tab--active': $route.path === '/settings/shares' }"
        >
          <i class="material-icons" aria-hidden="true">share</i>
          <span>{{ t("settings.shareManagement") }}</span>
        </router-link>
        <router-link
          v-if="user?.perm.admin"
          to="/settings/global"
          class="tabs__tab"
          :aria-label="t('settings.globalSettings')"
          :title="t('settings.globalSettings')"
          :class="{ 'tabs__tab--active': $route.path === '/settings/global' }"
        >
          <i class="material-icons" aria-hidden="true">tune</i>
          <span>{{ t("settings.globalSettings") }}</span>
        </router-link>
        <router-link
          v-if="user?.perm.admin"
          to="/settings/users"
          class="tabs__tab"
          :aria-label="t('settings.userManagement')"
          :title="t('settings.userManagement')"
          :class="{
            'tabs__tab--active':
              $route.path === '/settings/users' || $route.name === 'User',
          }"
        >
          <i class="material-icons" aria-hidden="true">group</i>
          <span>{{ t("settings.userManagement") }}</span>
        </router-link>
      </div>

      <div v-if="loading">
        <h2 class="message delayed">
          <div class="spinner">
            <div class="bounce1"></div>
            <div class="bounce2"></div>
            <div class="bounce3"></div>
          </div>
          <span>{{ t("files.loading") }}</span>
        </h2>
      </div>

      <router-view></router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import HeaderBar from "@/components/header/HeaderBar.vue";
import { computed } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const authStore = useAuthStore();
const layoutStore = useLayoutStore();

const user = computed(() => authStore.user);
const loading = computed(() => layoutStore.loading);
</script>
