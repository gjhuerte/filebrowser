<template>
  <div id="login" :class="{ recaptcha: recaptcha }">
    <div class="login__backdrop" aria-hidden="true">
      <div class="login__orb login__orb--primary"></div>
      <div class="login__orb login__orb--secondary"></div>
    </div>

    <div class="login__card">
      <div class="login__brand">
        <span class="login__logo">
          <img :src="logoURL" alt="" />
        </span>
        <h1 class="login__title">{{ name }}</h1>
      </div>

      <form class="login__form" @submit="submit">
        <p v-if="reason != null" class="login__message" role="status">
          <i class="material-icons" aria-hidden="true">info</i>
          <span>{{ t(`login.logout_reasons.${reason}`) }}</span>
        </p>
        <p
          v-if="error !== ''"
          class="login__message login__message--error"
          role="alert"
        >
          <i class="material-icons" aria-hidden="true">error_outline</i>
          <span>{{ error }}</span>
        </p>

        <label class="login__field">
          <span class="login__label">{{ t("login.username") }}</span>
          <input
            autofocus
            class="login__input"
            type="text"
            autocapitalize="off"
            autocomplete="username"
            v-model="username"
          />
        </label>

        <label class="login__field">
          <span class="login__label">{{ t("login.password") }}</span>
          <input
            class="login__input"
            type="password"
            :autocomplete="createMode ? 'new-password' : 'current-password'"
            v-model="password"
          />
        </label>

        <label v-if="createMode" class="login__field">
          <span class="login__label">{{ t("login.passwordConfirm") }}</span>
          <input
            class="login__input"
            type="password"
            autocomplete="new-password"
            v-model="passwordConfirm"
          />
        </label>

        <div v-if="recaptcha" id="recaptcha"></div>

        <button class="login__submit" type="submit">
          <span>{{ createMode ? t("login.signup") : t("login.submit") }}</span>
          <i class="material-icons" aria-hidden="true">arrow_forward</i>
        </button>

        <button
          v-if="signup"
          class="login__switch"
          type="button"
          @click="toggleMode"
        >
          {{
            createMode ? t("login.loginInstead") : t("login.createAnAccount")
          }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { StatusError } from "@/api/utils";
import * as auth from "@/utils/auth";
import {
  name,
  logoURL,
  recaptcha,
  recaptchaKey,
  signup,
} from "@/utils/constants";
import { inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";

// Define refs
const createMode = ref<boolean>(false);
const error = ref<string>("");
const username = ref<string>("");
const password = ref<string>("");
const passwordConfirm = ref<string>("");

const route = useRoute();
const router = useRouter();
const { t } = useI18n({});
// Define functions
const toggleMode = () => (createMode.value = !createMode.value);

const $showError = inject<IToastError>("$showError")!;

const reason = route.query["logout-reason"] ?? null;

const submit = async (event: Event) => {
  event.preventDefault();
  event.stopPropagation();

  const redirect = (route.query.redirect || "/files/") as string;

  let captcha = "";
  if (recaptcha) {
    captcha = window.grecaptcha.getResponse();

    if (captcha === "") {
      error.value = t("login.wrongCredentials");
      return;
    }
  }

  if (createMode.value) {
    if (password.value !== passwordConfirm.value) {
      error.value = t("login.passwordsDontMatch");
      return;
    }
  }

  try {
    if (createMode.value) {
      await auth.signup(username.value, password.value);
    }

    await auth.login(username.value, password.value, captcha);
    router.push({ path: redirect });
  } catch (e: any) {
    // console.error(e);
    if (e instanceof StatusError) {
      if (e.status === 409) {
        error.value = t("login.usernameTaken");
      } else if (e.status === 403) {
        error.value = t("login.wrongCredentials");
      } else if (e.status === 400) {
        const match = e.message.match(/minimum length is (\d+)/);
        if (match) {
          error.value = t("login.passwordTooShort", { min: match[1] });
        } else {
          error.value = e.message;
        }
      } else {
        $showError(e);
      }
    }
  }
};

// Run hooks
onMounted(() => {
  if (!recaptcha) return;

  window.grecaptcha.ready(function () {
    window.grecaptcha.render("recaptcha", {
      sitekey: recaptchaKey,
    });
  });
});
</script>
