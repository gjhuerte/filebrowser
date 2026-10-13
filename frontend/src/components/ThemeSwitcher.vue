<template>
  <div
    class="theme-switcher"
    role="group"
    :aria-label="t('settings.themes.title')"
  >
    <button
      v-for="option in options"
      :key="option.value || 'system'"
      type="button"
      class="theme-switcher__option"
      :class="{ 'theme-switcher__option--active': preference === option.value }"
      :aria-pressed="preference === option.value"
      :aria-label="t(option.label)"
      :title="t(option.label)"
      @click="choose(option.value)"
    >
      <i class="material-icons" aria-hidden="true">{{ option.icon }}</i>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { getThemePreference, setThemePreference } from "@/utils/theme";

const { t } = useI18n();

const options: { value: UserTheme; icon: string; label: string }[] = [
  { value: "light", icon: "light_mode", label: "settings.themes.light" },
  { value: "dark", icon: "dark_mode", label: "settings.themes.dark" },
  { value: "", icon: "computer", label: "settings.themes.default" },
];

const preference = ref<UserTheme>(getThemePreference());

const choose = (value: UserTheme) => {
  preference.value = value;
  setThemePreference(value);
};
</script>
