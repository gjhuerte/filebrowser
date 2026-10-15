<template>
  <div v-if="status === 'error'" class="section-state section-state--error">
    <p>{{ t("home.loadError") }}</p>
    <button type="button" class="button button--flat" @click="emit('retry')">
      {{ t("home.retry") }}
    </button>
  </div>
  <div
    v-else
    class="section-state"
    :class="`section-state--${skeleton}`"
    role="status"
    :aria-label="message || t('files.loading')"
  >
    <template v-if="skeleton === 'tile'">
      <span class="skeleton skeleton--value"></span>
      <span class="skeleton skeleton--meta"></span>
    </template>
    <template v-else>
      <p v-if="message" class="section-state__message">
        <span class="home__pulse" aria-hidden="true"></span>{{ message }}
      </p>
      <span v-for="n in 4" :key="n" class="skeleton skeleton--row"></span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

defineProps<{
  status: "loading" | "ready" | "error";
  skeleton: "tile" | "list";
  message?: string;
}>();
const emit = defineEmits<{ retry: [] }>();
const { t } = useI18n();
</script>
