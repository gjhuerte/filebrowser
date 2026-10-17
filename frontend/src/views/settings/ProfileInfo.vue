<template>
  <form class="card" @submit="updateSettings">
    <div class="card-title">
      <h2>{{ t("settings.profileInfo") }}</h2>
    </div>

    <div class="card-content">
      <p class="small">{{ t("settings.profileInfoDescription") }}</p>

      <p>
        <label for="username">{{ t("settings.username") }}</label>
        <input
          class="input input--block"
          type="text"
          id="username"
          :value="authStore.user?.username"
          disabled
        />
      </p>

      <p>
        <label for="locale">{{ t("settings.language") }}</label>
        <languages
          class="input input--block"
          id="locale"
          v-model:locale="locale"
        ></languages>
      </p>
    </div>

    <div class="card-action">
      <input
        class="button"
        type="submit"
        name="submitProfile"
        :value="t('buttons.update')"
      />
    </div>
  </form>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import Languages from "@/components/settings/Languages.vue";
import { inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const layoutStore = useLayoutStore();
const authStore = useAuthStore();
const { t } = useI18n();

const $showSuccess = inject<IToastSuccess>("$showSuccess")!;
const $showError = inject<IToastError>("$showError")!;

const locale = ref<string>("");

onMounted(async () => {
  layoutStore.loading = true;
  if (authStore.user === null) return false;
  locale.value = authStore.user.locale;
  layoutStore.loading = false;
  return true;
});

const updateSettings = async (event: Event) => {
  event.preventDefault();

  try {
    if (authStore.user === null) throw new Error("User is not set!");

    const data = {
      ...authStore.user,
      id: authStore.user.id,
      locale: locale.value,
    };

    await api.update(data, ["locale"]);
    authStore.updateUser(data);
    $showSuccess(t("settings.settingsUpdated"));
  } catch (err) {
    if (err instanceof Error) {
      $showError(err);
    }
  }
};
</script>
