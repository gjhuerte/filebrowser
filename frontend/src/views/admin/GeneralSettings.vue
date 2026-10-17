<template>
  <errors v-if="error" :errorCode="error.status" />
  <div class="row" v-else>
    <div class="column">
      <form class="card" @submit.prevent="$emit('save')">
        <div class="card-title">
          <h2>{{ t("settings.general") }}</h2>
        </div>

        <div class="card-content">
          <p class="small">{{ t("settings.generalDescription") }}</p>

          <h3>{{ t("settings.accessControl") }}</h3>
          <p class="small">{{ t("settings.accessControlDescription") }}</p>
          <div class="setting-list">
            <label class="setting-toggle">
              <span>{{ t("settings.allowSignup") }}</span>
              <input type="checkbox" class="switch" v-model="settings.signup" />
            </label>
            <label class="setting-toggle">
              <span>{{ t("settings.createUserDir") }}</span>
              <input
                type="checkbox"
                class="switch"
                v-model="settings.createUserDir"
              />
            </label>
            <label class="setting-toggle">
              <span>{{ t("settings.hideLoginButton") }}</span>
              <input
                type="checkbox"
                class="switch"
                v-model="settings.hideLoginButton"
              />
            </label>
          </div>

          <h3>{{ t("settings.userHomeBasePath") }}</h3>
          <p class="small">{{ t("settings.userHomeBasePathDescription") }}</p>
          <p>
            <input
              class="input input--block"
              type="text"
              v-model="settings.userHomeBasePath"
            />
          </p>

          <h3>{{ t("settings.minimumPasswordLength") }}</h3>
          <p class="small">
            {{ t("settings.minimumPasswordLengthDescription") }}
          </p>
          <p>
            <vue-number-input
              controls
              v-model.number="settings.minimumPasswordLength"
              :min="1"
            />
          </p>

          <h3>{{ t("settings.rules") }}</h3>
          <p class="small">{{ t("settings.globalRules") }}</p>
          <rules v-model:rules="settings.rules" />

          <div v-if="enableExec">
            <h3>{{ t("settings.executeOnShell") }}</h3>
            <p class="small">{{ t("settings.executeOnShellDescription") }}</p>
            <input
              class="input input--block"
              type="text"
              placeholder="bash -c, cmd /c, ..."
              :value="shellValue"
              @input="
                $emit(
                  'update:shellValue',
                  ($event.target as HTMLInputElement).value
                )
              "
            />
          </div>

          <h3>{{ t("settings.tusUploads") }}</h3>
          <p class="small">{{ t("settings.tusUploadsHelp") }}</p>
          <div class="tusConditionalSettings">
            <p>
              <label for="tus-chunkSize">{{
                t("settings.tusUploadsChunkSize")
              }}</label>
              <input
                class="input input--block"
                type="text"
                v-model="formattedChunkSize"
                id="tus-chunkSize"
              />
            </p>
            <p>
              <label for="tus-retryCount">{{
                t("settings.tusUploadsRetryCount")
              }}</label>
              <vue-number-input
                controls
                v-model.number="settings.tus.retryCount"
                id="tus-retryCount"
                :min="0"
              />
            </p>
          </div>
        </div>

        <div class="card-action">
          <input class="button" type="submit" :value="t('buttons.update')" />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import Rules from "@/components/settings/Rules.vue";
import Errors from "@/views/Errors.vue";
import { enableExec } from "@/utils/constants";
import { computed, onBeforeUnmount, ref } from "vue";
import { useI18n } from "vue-i18n";
import { StatusError } from "@/api/utils";

const { t } = useI18n();

const error = ref<StatusError | null>(null);
const debounceTimeout = ref<number | null>(null);
const pendingChunkSize = ref<string | null>(null);

const props = defineProps<{
  settings: ISettings;
  shellValue: string;
}>();

defineEmits<{
  save: [];
  "update:shellValue": [value: string];
}>();

const formattedChunkSize = computed({
  get() {
    return props.settings?.tus?.chunkSize
      ? formatBytes(props.settings.tus.chunkSize)
      : "";
  },
  set(value: string) {
    if (debounceTimeout.value) {
      clearTimeout(debounceTimeout.value);
    }
    pendingChunkSize.value = value;
    debounceTimeout.value = window.setTimeout(applyChunkSize, 1500);
  },
});

const applyChunkSize = () => {
  if (debounceTimeout.value) {
    clearTimeout(debounceTimeout.value);
    debounceTimeout.value = null;
  }
  if (pendingChunkSize.value === null) return;
  if (props.settings) {
    props.settings.tus.chunkSize = parseBytes(pendingChunkSize.value);
  }
  pendingChunkSize.value = null;
};

const parseBytes = (input: string) => {
  const regex = /^(\d+)(\.\d+)?(B|K|KB|M|MB|G|GB|T|TB)?$/i;
  const matches = input.match(regex);
  if (matches) {
    const size = parseFloat(matches[1].concat(matches[2] || ""));
    let unit: keyof SettingsUnit = (
      matches[3] ?? "B"
    ).toUpperCase() as keyof SettingsUnit;
    if (!unit.endsWith("B")) {
      unit += "B";
    }
    const units: SettingsUnit = {
      KB: 1024,
      MB: 1024 ** 2,
      GB: 1024 ** 3,
      TB: 1024 ** 4,
    };
    return size * (units[unit as keyof SettingsUnit] || 1);
  } else {
    return 1024 ** 2;
  }
};

const formatBytes = (bytes: number) => {
  const units = ["B", "KB", "MB", "GB", "TB"];
  let size = bytes;
  let unitIndex = 0;
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  return `${size}${units[unitIndex]}`;
};

onBeforeUnmount(() => {
  applyChunkSize();
  if (debounceTimeout.value) {
    clearTimeout(debounceTimeout.value);
  }
});
</script>
