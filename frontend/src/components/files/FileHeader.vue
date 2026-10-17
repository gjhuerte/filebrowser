<template>
  <header-bar class="file-header">
    <button
      type="button"
      class="toolbar__btn file-header__back"
      :aria-label="closeLabel"
      :title="closeLabel + ' (Esc)'"
      @click="emit('close')"
    >
      <i class="material-icons" aria-hidden="true">arrow_back</i>
    </button>

    <div class="file-header__id">
      <span class="file-header__icon" aria-hidden="true">
        <i class="material-icons">{{ icon }}</i>
      </span>
      <div class="file-header__text">
        <h1 class="file-header__name" :title="name">{{ name }}</h1>
        <p v-if="meta.length || $slots.status" class="file-header__meta">
          <span v-for="(item, i) in meta" :key="i">{{ item }}</span>
          <slot name="status" />
        </p>
      </div>
    </div>

    <template v-if="$slots.actions" #actions>
      <slot name="actions" />
    </template>

    <template v-if="$slots.primary" #primary>
      <slot name="primary" />
    </template>
  </header-bar>
</template>

<script setup lang="ts">
import HeaderBar from "@/components/header/HeaderBar.vue";

/*
 * Top bar for file previews and the text editor: back button, file
 * identity (type icon, name, a short meta line), secondary actions that
 * fold into the ⋮ menu on phones, and one primary button.
 */
withDefaults(
  defineProps<{
    name: string;
    icon?: string;
    /** Short facts shown under the name, e.g. folder, size, "4 of 12". */
    meta?: string[];
    closeLabel: string;
  }>(),
  { icon: "insert_drive_file", meta: () => [] }
);

const emit = defineEmits<{ close: [] }>();
</script>
