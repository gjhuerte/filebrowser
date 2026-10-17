<template>
  <div class="row">
    <div class="column">
      <form class="card" @submit.prevent="$emit('save')">
        <div class="card-title">
          <h2>{{ t("settings.branding") }}</h2>
        </div>

        <div class="card-content">
          <i18n-t
            keypath="settings.brandingHelp"
            tag="p"
            class="small"
            scope="global"
          >
            <a
              class="link"
              target="_blank"
              href="https://github.com/filebrowser/filebrowser/blob/master/docs/customization.md#custom-branding"
              >{{ t("settings.documentation") }}</a
            >
          </i18n-t>

          <h3>{{ t("settings.appearance") }}</h3>
          <p class="small">{{ t("settings.appearanceDescription") }}</p>
          <div class="setting-list">
            <label class="setting-toggle">
              <span>{{ t("settings.disableExternalLinks") }}</span>
              <input
                type="checkbox"
                class="switch"
                v-model="settings.branding.disableExternal"
                id="branding-links"
              />
            </label>
            <label class="setting-toggle">
              <span>{{ t("settings.disableUsedDiskPercentage") }}</span>
              <input
                type="checkbox"
                class="switch"
                v-model="settings.branding.disableUsedPercentage"
                id="branding-used-disk"
              />
            </label>
          </div>

          <h3>{{ t("settings.themes.title") }}</h3>
          <p class="small">{{ t("settings.brandingThemeDescription") }}</p>
          <p>
            <themes
              class="input input--block"
              v-model:theme="settings.branding.theme"
              id="theme"
            ></themes>
          </p>

          <h3>{{ t("settings.instanceName") }}</h3>
          <p class="small">{{ t("settings.instanceNameDescription") }}</p>
          <p>
            <input
              class="input input--block"
              type="text"
              v-model="settings.branding.name"
              id="branding-name"
            />
          </p>

          <h3>{{ t("settings.brandingDirectoryPath") }}</h3>
          <p class="small">
            {{ t("settings.brandingDirectoryPathDescription") }}
          </p>
          <p>
            <input
              class="input input--block"
              type="text"
              v-model="settings.branding.files"
              id="branding-files"
            />
          </p>
        </div>

        <div class="card-action">
          <input class="button" type="submit" :value="t('buttons.update')" />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import Themes from "@/components/settings/Themes.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

defineProps<{
  settings: ISettings;
}>();

defineEmits<{
  save: [];
}>();
</script>
