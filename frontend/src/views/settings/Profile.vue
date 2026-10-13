<template>
  <div class="row">
    <div class="column">
      <form class="card" @submit="updateSettings">
        <div class="card-title">
          <h2>{{ t("settings.profileSettings") }}</h2>
        </div>

        <div class="card-content">
          <div class="setting-list">
            <label class="setting-toggle">
              <span>{{ t("settings.hideDotfiles") }}</span>
              <input
                type="checkbox"
                class="switch"
                name="hideDotfiles"
                v-model="hideDotfiles"
              />
            </label>
            <label class="setting-toggle">
              <span>{{ t("settings.singleClick") }}</span>
              <input
                type="checkbox"
                class="switch"
                name="singleClick"
                v-model="singleClick"
              />
            </label>
            <label class="setting-toggle">
              <span>{{ t("settings.redirectAfterCopyMove") }}</span>
              <input
                type="checkbox"
                class="switch"
                name="redirectAfterCopyMove"
                v-model="redirectAfterCopyMove"
              />
            </label>
            <label class="setting-toggle">
              <span>{{ t("settings.setDateFormat") }}</span>
              <input
                type="checkbox"
                class="switch"
                name="dateFormat"
                v-model="dateFormat"
              />
            </label>
          </div>

          <p>
            <label for="locale">{{ t("settings.language") }}</label>
            <languages
              class="input input--block"
              id="locale"
              v-model:locale="locale"
            ></languages>
          </p>

          <p>
            <label for="aceTheme">{{ t("settings.aceEditorTheme") }}</label>
            <AceEditorTheme
              class="input input--block"
              v-model:aceEditorTheme="aceEditorTheme"
              id="aceTheme"
            ></AceEditorTheme>
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
    </div>

    <div v-if="!noAuth" class="column">
      <form
        class="card"
        v-if="!authStore.user?.lockPassword"
        @submit="updatePassword"
      >
        <div class="card-title">
          <h2>{{ t("settings.changePassword") }}</h2>
        </div>

        <div class="card-content">
          <p>
            <label for="newPassword">{{ t("settings.newPassword") }}</label>
            <input
              :class="passwordClass"
              type="password"
              id="newPassword"
              v-model="password"
              name="password"
              autocomplete="new-password"
            />
          </p>
          <p>
            <label for="newPasswordConf">{{
              t("settings.newPasswordConfirm")
            }}</label>
            <input
              :class="passwordClass"
              type="password"
              id="newPasswordConf"
              v-model="passwordConf"
              name="passwordConf"
              autocomplete="new-password"
            />
          </p>
          <p v-if="isCurrentPasswordRequired">
            <label for="currentPassword">{{
              t("settings.currentPassword")
            }}</label>
            <input
              :class="passwordClass"
              type="password"
              id="currentPassword"
              v-model="currentPassword"
              name="current_password"
              autocomplete="current-password"
            />
          </p>
        </div>

        <div class="card-action">
          <input
            class="button"
            type="submit"
            name="submitPassword"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import AceEditorTheme from "@/components/settings/AceEditorTheme.vue";
import Languages from "@/components/settings/Languages.vue";
import { computed, inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { authMethod, noAuth } from "@/utils/constants";

const layoutStore = useLayoutStore();
const authStore = useAuthStore();
const { t } = useI18n();

const $showSuccess = inject<IToastSuccess>("$showSuccess")!;
const $showError = inject<IToastError>("$showError")!;

const password = ref<string>("");
const passwordConf = ref<string>("");
const currentPassword = ref<string>("");
const isCurrentPasswordRequired = ref<boolean>(false);
const hideDotfiles = ref<boolean>(false);
const singleClick = ref<boolean>(false);
const redirectAfterCopyMove = ref<boolean>(false);
const dateFormat = ref<boolean>(false);
const locale = ref<string>("");
const aceEditorTheme = ref<string>("");

const passwordClass = computed(() => {
  const baseClass = "input input--block";

  if (password.value === "" && passwordConf.value === "") {
    return baseClass;
  }

  if (password.value === passwordConf.value) {
    return `${baseClass} input--green`;
  }

  return `${baseClass} input--red`;
});

onMounted(async () => {
  layoutStore.loading = true;
  if (authStore.user === null) return false;
  locale.value = authStore.user.locale;
  hideDotfiles.value = authStore.user.hideDotfiles;
  singleClick.value = authStore.user.singleClick;
  redirectAfterCopyMove.value = authStore.user.redirectAfterCopyMove;
  dateFormat.value = authStore.user.dateFormat;
  aceEditorTheme.value = authStore.user.aceEditorTheme;
  layoutStore.loading = false;
  isCurrentPasswordRequired.value = authMethod == "json";

  return true;
});

const updatePassword = async (event: Event) => {
  event.preventDefault();

  if (
    password.value !== passwordConf.value ||
    password.value === "" ||
    currentPassword.value === "" ||
    authStore.user === null
  ) {
    return;
  }

  try {
    const data = {
      ...authStore.user,
      id: authStore.user.id,
      password: password.value,
    };
    await api.update(data, ["password"], currentPassword.value);
    authStore.updateUser(data);
    $showSuccess(t("settings.passwordUpdated"));
  } catch (e: any) {
    $showError(e);
  } finally {
    password.value = passwordConf.value = "";
  }
};
const updateSettings = async (event: Event) => {
  event.preventDefault();

  try {
    if (authStore.user === null) throw new Error("User is not set!");

    const data = {
      ...authStore.user,
      id: authStore.user.id,
      locale: locale.value,
      hideDotfiles: hideDotfiles.value,
      singleClick: singleClick.value,
      redirectAfterCopyMove: redirectAfterCopyMove.value,
      dateFormat: dateFormat.value,
      aceEditorTheme: aceEditorTheme.value,
    };

    await api.update(data, [
      "locale",
      "hideDotfiles",
      "singleClick",
      "redirectAfterCopyMove",
      "dateFormat",
      "aceEditorTheme",
    ]);
    authStore.updateUser(data);
    $showSuccess(t("settings.settingsUpdated"));
  } catch (err) {
    if (err instanceof Error) {
      $showError(err);
    }
  }
};
</script>
