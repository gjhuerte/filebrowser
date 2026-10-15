<template>
  <div
    class="empty-state"
    :class="{
      'empty-state--drop': dropzone,
      'empty-state--active': active,
      'empty-state--compact': compact,
    }"
  >
    <span class="empty-state__icon" aria-hidden="true">
      <i class="material-icons">{{ active ? "file_download" : icon }}</i>
    </span>
    <h2 class="empty-state__title">{{ title }}</h2>
    <p v-if="description" class="empty-state__text">{{ description }}</p>
    <div v-if="$slots.default" class="empty-state__actions">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  icon: string;
  title: string;
  description?: string;
  /** Draw a dashed outline: the area accepts dropped files */
  dropzone?: boolean;
  /** Something is being dragged over the page right now */
  active?: boolean;
  /** Smaller variant for cards and dropdowns */
  compact?: boolean;
}>();
</script>
