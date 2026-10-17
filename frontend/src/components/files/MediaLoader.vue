<template>
  <div
    class="media-loader"
    :class="`media-loader--${variant}`"
    role="status"
    aria-live="polite"
  >
    <div class="media-loader__ring" aria-hidden="true">
      <svg viewBox="0 0 48 48" :class="{ 'is-spinning': percent === null }">
        <circle class="media-loader__track" cx="24" cy="24" r="20" />
        <circle
          class="media-loader__arc"
          cx="24"
          cy="24"
          r="20"
          :stroke-dasharray="CIRCUMFERENCE"
          :stroke-dashoffset="dashOffset"
        />
      </svg>
      <span v-if="percent !== null && variant === 'center'">
        {{ percent }}<small>%</small>
      </span>
    </div>
    <div class="media-loader__text">
      <span class="media-loader__label">
        {{ label
        }}<template v-if="percent !== null && variant === 'chip'">
          · {{ percent }}%</template
        >
      </span>
      <span v-if="sizeText" class="media-loader__size">{{ sizeText }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import prettyBytes from "pretty-bytes";

const props = withDefaults(
  defineProps<{
    label: string;
    /** Bytes received (or estimated). */
    loaded?: number;
    /** Total bytes; 0 = unknown. */
    total?: number;
    /** Overrides the computed percentage (0–100), e.g. buffered media. */
    progress?: number | null;
    variant?: "center" | "chip";
  }>(),
  { loaded: 0, total: 0, progress: null, variant: "center" }
);

const { t } = useI18n();
const CIRCUMFERENCE = 2 * Math.PI * 20;

const percent = computed<number | null>(() => {
  if (props.progress !== null) return Math.round(props.progress);
  if (props.total > 0)
    return Math.min(100, Math.round((props.loaded / props.total) * 100));
  return null;
});

// Indeterminate: a quarter arc that spins
const dashOffset = computed(() =>
  percent.value === null
    ? CIRCUMFERENCE * 0.72
    : CIRCUMFERENCE * (1 - percent.value / 100)
);

const bytes = (n: number) => prettyBytes(n, { binary: true });
const sizeText = computed(() => {
  if (props.total > 0 && props.loaded > 0)
    return t("preview.sizeOf", {
      loaded: bytes(props.loaded),
      total: bytes(props.total),
    });
  if (props.loaded > 0)
    return t("preview.sizeLoaded", { size: bytes(props.loaded) });
  // Nothing received yet: at least say how big the file is
  if (props.total > 0) return bytes(props.total);
  return "";
});
</script>
