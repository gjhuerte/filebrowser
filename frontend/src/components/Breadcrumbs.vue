<template>
  <div class="breadcrumbs">
    <div
      class="breadcrumbs__trail"
      role="navigation"
      :aria-label="t('files.location')"
    >
      <ol>
        <li>
          <component
            :is="element"
            :to="base || ''"
            class="breadcrumbs__home"
            :class="{ 'breadcrumbs__home--current': items.length === 0 }"
            :aria-label="t('files.home')"
            :title="t('files.home')"
            :aria-current="items.length === 0 ? 'page' : undefined"
          >
            <i class="material-icons" aria-hidden="true">home</i>
          </component>
        </li>
        <li v-for="(link, index) in items" :key="link.url">
          <i class="material-icons chevron" aria-hidden="true">chevron_right</i>
          <span
            v-if="index === items.length - 1"
            class="breadcrumbs__current"
            aria-current="page"
            :title="link.name"
            >{{ link.name }}</span
          >
          <component
            v-else
            :is="element"
            :to="link.url"
            class="breadcrumbs__link"
            :title="link.collapsed ? t('files.showPath') : link.name"
          >
            <i v-if="link.collapsed" class="material-icons" aria-hidden="true"
              >more_horiz</i
            >
            <template v-else>{{ link.name }}</template>
          </component>
        </li>
      </ol>
    </div>
    <div v-if="$slots.default" class="breadcrumbs__actions">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";

const { t } = useI18n();

const route = useRoute();

const props = defineProps<{
  base: string;
  noLink?: boolean;
}>();

const items = computed(() => {
  const relativePath = route.path.replace(props.base, "");
  const parts = relativePath.split("/");

  if (parts[0] === "") {
    parts.shift();
  }

  if (parts[parts.length - 1] === "") {
    parts.pop();
  }

  const breadcrumbs: (BreadCrumb & { collapsed?: boolean })[] = [];

  for (let i = 0; i < parts.length; i++) {
    if (i === 0) {
      breadcrumbs.push({
        name: decodeURIComponent(parts[i]),
        url: props.base + "/" + parts[i] + "/",
      });
    } else {
      breadcrumbs.push({
        name: decodeURIComponent(parts[i]),
        url: breadcrumbs[i - 1].url + parts[i] + "/",
      });
    }
  }

  if (breadcrumbs.length > 3) {
    while (breadcrumbs.length !== 4) {
      breadcrumbs.shift();
    }

    breadcrumbs[0].collapsed = true;
  }

  return breadcrumbs;
});

const element = computed(() => {
  if (props.noLink) {
    return "span";
  }

  return "router-link";
});
</script>
