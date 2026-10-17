<template>
  <form class="card" @submit="updateSettings">
    <div class="card-title">
      <h2>{{ t("settings.system") }}</h2>
    </div>

    <div class="card-content">
      <p class="small">{{ t("settings.systemDescription") }}</p>

      <h3>{{ t("settings.themes.title") }}</h3>
      <div class="theme-switcher">
        <button
          type="button"
          class="theme-switcher__option"
          :class="{
            'theme-switcher__option--active': currentTheme === 'light',
          }"
          :aria-label="t('settings.themes.light')"
          :title="t('settings.themes.light')"
          @click="setCurrentTheme('light')"
        >
          <i class="material-icons" aria-hidden="true">light_mode</i>
        </button>
        <button
          type="button"
          class="theme-switcher__option"
          :class="{
            'theme-switcher__option--active': currentTheme === 'dark',
          }"
          :aria-label="t('settings.themes.dark')"
          :title="t('settings.themes.dark')"
          @click="setCurrentTheme('dark')"
        >
          <i class="material-icons" aria-hidden="true">dark_mode</i>
        </button>
        <button
          type="button"
          class="theme-switcher__option"
          :class="{
            'theme-switcher__option--active': currentTheme === '',
          }"
          :aria-label="t('settings.themes.default')"
          :title="t('settings.themes.default')"
          @click="setCurrentTheme('')"
        >
          <i class="material-icons" aria-hidden="true">contrast</i>
        </button>
      </div>

      <h3>{{ t("settings.filePreferences") }}</h3>
      <div class="setting-list">
        <label class="setting-toggle">
          <span>{{ t("settings.hideDotfiles") }}</span>
          <input
            type="checkbox"
            class="switch"
            name="hideDotfiles"
            v-model="hideDotfiles"
          />
        </label>
        <label class="setting-toggle">
          <span>{{ t("settings.singleClick") }}</span>
          <input
            type="checkbox"
            class="switch"
            name="singleClick"
            v-model="singleClick"
          />
        </label>
        <label class="setting-toggle">
          <span>{{ t("settings.redirectAfterCopyMove") }}</span>
          <input
            type="checkbox"
            class="switch"
            name="redirectAfterCopyMove"
            v-model="redirectAfterCopyMove"
          />
        </label>
        <label class="setting-toggle">
          <span>{{ t("settings.setDateFormat") }}</span>
          <input
            type="checkbox"
            class="switch"
            name="dateFormat"
            v-model="dateFormat"
          />
        </label>
      </div>
    </div>

    <div class="card-action">
      <input
        class="button"
        type="submit"
        name="submitSystem"
        :value="t('buttons.update')"
      />
    </div>
  </form>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import { inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { getTheme, setTheme } from "@/utils/theme";

const layoutStore = useLayoutStore();
const authStore = useAuthStore();
const { t } = useI18n();

const $showSuccess = inject<IToastSuccess>("$showSuccess")!;
const $showError = inject<IToastError>("$showError")!;

const hideDotfiles = ref<boolean>(false);
const singleClick = ref<boolean>(false);
const redirectAfterCopyMove = ref<boolean>(false);
const dateFormat = ref<boolean>(false);
const currentTheme = ref<string>("");

onMounted(async () => {
  layoutStore.loading = true;
  if (authStore.user === null) return false;
  hideDotfiles.value = authStore.user.hideDotfiles;
  singleClick.value = authStore.user.singleClick;
  redirectAfterCopyMove.value = authStore.user.redirectAfterCopyMove;
  dateFormat.value = authStore.user.dateFormat;
  currentTheme.value = getTheme();
  layoutStore.loading = false;
  return true;
});

const setCurrentTheme = (theme: string) => {
  currentTheme.value = theme;
  setTheme(theme);
};

const updateSettings = async (event: Event) => {
  event.preventDefault();

  try {
    if (authStore.user === null) throw new Error("User is not set!");

    const data = {
      ...authStore.user,
      id: authStore.user.id,
      hideDotfiles: hideDotfiles.value,
      singleClick: singleClick.value,
      redirectAfterCopyMove: redirectAfterCopyMove.value,
      dateFormat: dateFormat.value,
    };

    await api.update(data, [
      "hideDotfiles",
      "singleClick",
      "redirectAfterCopyMove",
      "dateFormat",
    ]);
    authStore.updateUser(data);
    $showSuccess(t("settings.settingsUpdated"));
  } catch (err) {
    if (err instanceof Error) {
      $showError(err);
    }
  }
};
</script>
