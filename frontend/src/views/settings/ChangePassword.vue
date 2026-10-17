<template>
  <form
    class="card"
    v-if="!noAuth && !authStore.user?.lockPassword"
    @submit="updatePassword"
  >
    <div class="card-title">
      <h2>{{ t("settings.changePassword") }}</h2>
    </div>

    <div class="card-content">
      <p class="small">{{ t("settings.changePasswordDescription") }}</p>

      <div>
        <label for="newPassword">{{ t("settings.newPassword") }}</label>
        <div class="password-wrap">
          <input
            :class="passwordClass"
            :type="showNewPwd ? 'text' : 'password'"
            id="newPassword"
            v-model="password"
            name="password"
            autocomplete="new-password"
          />
          <button
            class="password-toggle"
            type="button"
            @click="showNewPwd = !showNewPwd"
            :aria-label="
              showNewPwd ? t('login.hidePassword') : t('login.showPassword')
            "
            tabindex="-1"
          >
            <i class="material-icons" aria-hidden="true">{{
              showNewPwd ? "visibility_off" : "visibility"
            }}</i>
          </button>
        </div>
        <PasswordStrength :password="password" />
      </div>
      <div>
        <label for="newPasswordConf">{{
          t("settings.newPasswordConfirm")
        }}</label>
        <div class="password-wrap">
          <input
            :class="passwordClass"
            :type="showNewPwdConf ? 'text' : 'password'"
            id="newPasswordConf"
            v-model="passwordConf"
            name="passwordConf"
            autocomplete="new-password"
          />
          <button
            class="password-toggle"
            type="button"
            @click="showNewPwdConf = !showNewPwdConf"
            :aria-label="
              showNewPwdConf ? t('login.hidePassword') : t('login.showPassword')
            "
            tabindex="-1"
          >
            <i class="material-icons" aria-hidden="true">{{
              showNewPwdConf ? "visibility_off" : "visibility"
            }}</i>
          </button>
        </div>
      </div>
      <div v-if="isCurrentPasswordRequired">
        <label for="currentPassword">{{ t("settings.currentPassword") }}</label>
        <div class="password-wrap">
          <input
            :class="passwordClass"
            :type="showCurrentPwd ? 'text' : 'password'"
            id="currentPassword"
            v-model="currentPassword"
            name="current_password"
            autocomplete="current-password"
          />
          <button
            class="password-toggle"
            type="button"
            @click="showCurrentPwd = !showCurrentPwd"
            :aria-label="
              showCurrentPwd ? t('login.hidePassword') : t('login.showPassword')
            "
            tabindex="-1"
          >
            <i class="material-icons" aria-hidden="true">{{
              showCurrentPwd ? "visibility_off" : "visibility"
            }}</i>
          </button>
        </div>
      </div>
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

  <div class="card" v-else-if="authStore.user?.lockPassword">
    <div class="card-title">
      <h2>{{ t("settings.changePassword") }}</h2>
    </div>
    <div class="card-content">
      <p class="small">{{ t("settings.passwordLockedDescription") }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import PasswordStrength from "@/components/PasswordStrength.vue";
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
const showNewPwd = ref<boolean>(false);
const showNewPwdConf = ref<boolean>(false);
const showCurrentPwd = ref<boolean>(false);

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
  isCurrentPasswordRequired.value = authMethod == "json";
  layoutStore.loading = false;
});

const updatePassword = async (event: Event) => {
  event.preventDefault();

  if (
    password.value !== passwordConf.value ||
    password.value === "" ||
    authStore.user === null
  ) {
    return;
  }

  if (isCurrentPasswordRequired.value && currentPassword.value === "") {
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
</script>
