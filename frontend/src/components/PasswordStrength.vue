<template>
  <div v-if="password" class="password-strength">
    <div class="password-strength__bar">
      <div
        class="password-strength__fill"
        :class="`password-strength__fill--${level}`"
        :style="{ width: `${percent}%` }"
      />
    </div>
    <span
      class="password-strength__label"
      :class="`password-strength__label--${level}`"
    >
      {{ t(`settings.passwordStrength.${level}`) }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const props = defineProps<{ password: string }>();
const emit = defineEmits<{ (e: "strength", level: string): void }>();

const level = computed(() => {
  const p = props.password;
  if (!p) return "weak";

  let score = 0;
  if (p.length >= 4) score++;
  if (p.length >= 8) score++;
  if (p.length >= 12) score++;
  if (/[a-z]/.test(p) && /[A-Z]/.test(p)) score++;
  if (/\d/.test(p)) score++;
  if (/[^a-zA-Z\d]/.test(p)) score++;

  if (score <= 1) return "weak";
  if (score <= 2) return "fair";
  if (score <= 4) return "good";
  return "strong";
});

const percent = computed(() => {
  switch (level.value) {
    case "weak":
      return 25;
    case "fair":
      return 50;
    case "good":
      return 75;
    case "strong":
      return 100;
    default:
      return 0;
  }
});

watch(level, (val) => emit("strength", val), { immediate: true });
</script>
