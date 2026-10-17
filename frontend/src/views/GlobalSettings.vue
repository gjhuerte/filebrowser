<template>
  <div class="dashboard">
    <header-bar showMenu showLogo />

    <div class="settings">
      <h1 class="settings__title">{{ t("settings.globalSettings") }}</h1>
      <p class="settings__subtitle">
        {{ t("settings.globalSettingsDescription") }}
      </p>

      <div
        class="tabs"
        role="navigation"
        :aria-label="t('settings.globalSettings')"
      >
        <router-link
          to="/admin/settings/general"
          class="tabs__tab"
          :class="{
            'tabs__tab--active': $route.path === '/admin/settings/general',
          }"
        >
          <i class="material-icons" aria-hidden="true">tune</i>
          <span>{{ t("settings.general") }}</span>
        </router-link>
        <router-link
          to="/admin/settings/user-defaults"
          class="tabs__tab"
          :class="{
            'tabs__tab--active':
              $route.path === '/admin/settings/user-defaults',
          }"
        >
          <i class="material-icons" aria-hidden="true">person_outline</i>
          <span>{{ t("settings.userDefaults") }}</span>
        </router-link>
        <router-link
          to="/admin/settings/branding"
          class="tabs__tab"
          :class="{
            'tabs__tab--active': $route.path === '/admin/settings/branding',
          }"
        >
          <i class="material-icons" aria-hidden="true">palette</i>
          <span>{{ t("settings.branding") }}</span>
        </router-link>
        <router-link
          v-if="enableExec"
          to="/admin/settings/commands"
          class="tabs__tab"
          :class="{
            'tabs__tab--active': $route.path === '/admin/settings/commands',
          }"
        >
          <i class="material-icons" aria-hidden="true">terminal</i>
          <span>{{ t("settings.commandRunner") }}</span>
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

      <router-view
        v-if="settings !== null"
        :settings="settings"
        :commandObject="commandObject"
        :shellValue="shellValue"
        @save="save"
        @update:shellValue="shellValue = $event"
        @update:commandObject="commandObject = $event"
      ></router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLayoutStore } from "@/stores/layout";
import HeaderBar from "@/components/header/HeaderBar.vue";
import { settings as api } from "@/api";
import { enableExec } from "@/utils/constants";
import { getTheme, setTheme } from "@/utils/theme";
import { computed, inject, onBeforeUnmount, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const layoutStore = useLayoutStore();

const $showError = inject<IToastError>("$showError")!;
const $showSuccess = inject<IToastSuccess>("$showSuccess")!;

const settings = ref<ISettings | null>(null);
const commandObject = ref<{ [key: string]: string[] | string }>({});
const shellValue = ref<string>("");
const debounceTimeout = ref<number | null>(null);

const loading = computed(() => layoutStore.loading);

onMounted(async () => {
  try {
    layoutStore.loading = true;
    const original: ISettings = await api.get();
    const newSettings: ISettings = { ...original, commands: {} };

    const keys = Object.keys(original.commands) as Array<keyof SettingsCommand>;
    for (const key of keys) {
      newSettings.commands[key] = original.commands[key];
      commandObject.value[key] = original.commands[key]!.join("\n");
    }

    settings.value = newSettings;
    shellValue.value = newSettings.shell.join(" ");
  } catch (err) {
    if (err instanceof Error) {
      $showError(err);
    }
  } finally {
    layoutStore.loading = false;
  }
});

onBeforeUnmount(() => {
  if (debounceTimeout.value) {
    clearTimeout(debounceTimeout.value);
  }
});

const save = async () => {
  if (settings.value === null) return false;

  const newSettings: ISettings = {
    ...settings.value,
    shell:
      settings.value?.shell
        .join(" ")
        .trim()
        .split(" ")
        .filter((s: string) => s !== "") ?? [],
    commands: {},
  };

  const keys = Object.keys(settings.value.commands) as Array<
    keyof SettingsCommand
  >;
  for (const key of keys) {
    const newValue = commandObject.value[key];
    if (!newValue) continue;

    if (Array.isArray(newValue)) {
      newSettings.commands[key] = newValue;
    } else if (key in commandObject.value) {
      newSettings.commands[key] = newValue
        .split("\n")
        .filter((cmd: string) => cmd !== "");
    }
  }
  newSettings.shell = shellValue.value
    .trim()
    .split(" ")
    .filter((s) => s !== "");

  if (newSettings.branding.theme !== getTheme()) {
    setTheme(newSettings.branding.theme);
  }

  try {
    await api.update(newSettings);
    $showSuccess(t("settings.settingsUpdated"));
  } catch (e: any) {
    $showError(e);
  }

  return true;
};
</script>
