<template>
  <div
    class="context-menu"
    ref="contextMenu"
    v-show="show"
    :style="{
      top: `${top}px`,
      left: `${left}px`,
    }"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, nextTick, onUnmounted } from "vue";

const emit = defineEmits(["hide"]);
const props = defineProps<{ show: boolean; pos: { x: number; y: number } }>();
const contextMenu = ref<HTMLElement | null>(null);

// Measured once the menu is on screen (it is hidden until then)
const size = ref({ width: 0, height: 0 });

// Keep the menu inside the window: shift left near the right edge and
// flip up near the bottom
const left = computed(() => {
  return Math.max(
    8,
    Math.min(props.pos.x, window.innerWidth - size.value.width - 8)
  );
});

const top = computed(() => {
  const bottom = window.scrollY + window.innerHeight - 8;
  return props.pos.y + size.value.height > bottom
    ? Math.max(window.scrollY + 8, bottom - size.value.height)
    : props.pos.y;
});

const hideContextMenu = () => {
  emit("hide");
};

watch(
  () => [props.show, props.pos] as const,
  async ([show]) => {
    if (!show) return;
    await nextTick();
    const el = contextMenu.value;
    if (el) size.value = { width: el.offsetWidth, height: el.offsetHeight };
  }
);

watch(
  () => props.show,
  (val) => {
    if (val) {
      document.addEventListener("click", hideContextMenu);
    } else {
      document.removeEventListener("click", hideContextMenu);
    }
  }
);

onUnmounted(() => {
  document.removeEventListener("click", hideContextMenu);
});
</script>
