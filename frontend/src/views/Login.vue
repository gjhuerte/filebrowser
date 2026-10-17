<template>
  <div id="login" :class="{ recaptcha: recaptcha }">
    <ThemeSwitcher class="login__theme" />

    <div class="login__card">
      <div class="login__brand">
        <span class="login__logo">
          <img :src="logoURL" alt="" />
        </span>
        <h1 class="login__title">{{ name }}</h1>
        <p class="login__subtitle">
          {{ createMode ? t("login.subtitleSignup") : t("login.subtitle") }}
        </p>
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
            :placeholder="t('login.usernamePlaceholder')"
            v-model="username"
          />
        </label>

        <label class="login__field">
          <span class="login__label">{{ t("login.password") }}</span>
          <div class="password-wrap">
            <input
              class="login__input"
              :type="showPassword ? 'text' : 'password'"
              :autocomplete="createMode ? 'new-password' : 'current-password'"
              :placeholder="t('login.passwordPlaceholder')"
              v-model="password"
              @keydown="checkCapsLock"
              @keyup="checkCapsLock"
              @blur="capsLock = false"
            />
            <button
              class="password-toggle"
              type="button"
              @click="showPassword = !showPassword"
              :aria-label="
                showPassword ? t('login.hidePassword') : t('login.showPassword')
              "
              tabindex="-1"
            >
              <i class="material-icons" aria-hidden="true">{{
                showPassword ? "visibility_off" : "visibility"
              }}</i>
            </button>
          </div>
          <span v-if="capsLock" class="login__hint" role="status">
            <i class="material-icons" aria-hidden="true">keyboard_capslock</i>
            {{ t("login.capsLock") }}
          </span>
          <PasswordStrength v-if="createMode" :password="password" />
        </label>

        <label v-if="createMode" class="login__field">
          <span class="login__label">{{ t("login.passwordConfirm") }}</span>
          <div class="password-wrap">
            <input
              class="login__input"
              :type="showPasswordConfirm ? 'text' : 'password'"
              autocomplete="new-password"
              :placeholder="t('login.passwordConfirmPlaceholder')"
              v-model="passwordConfirm"
            />
            <button
              class="password-toggle"
              type="button"
              @click="showPasswordConfirm = !showPasswordConfirm"
              :aria-label="
                showPasswordConfirm
                  ? t('login.hidePassword')
                  : t('login.showPassword')
              "
              tabindex="-1"
            >
              <i class="material-icons" aria-hidden="true">{{
                showPasswordConfirm ? "visibility_off" : "visibility"
              }}</i>
            </button>
          </div>
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

      <p v-if="!createMode" class="login__note">
        {{ t("login.forgotPassword") }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { StatusError } from "@/api/utils";
import * as auth from "@/utils/auth";
import PasswordStrength from "@/components/PasswordStrength.vue";
import ThemeSwitcher from "@/components/ThemeSwitcher.vue";
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

const showPassword = ref<boolean>(false);
const showPasswordConfirm = ref<boolean>(false);
const capsLock = ref<boolean>(false);
const checkCapsLock = (event: KeyboardEvent) => {
  capsLock.value = event.getModifierState?.("CapsLock") ?? false;
};

const $showError = inject<IToastError>("$showError")!;

const reason = route.query["logout-reason"] ?? null;

const submit = async (event: Event) => {
  event.preventDefault();
  event.stopPropagation();

  const redirect = (route.query.redirect || "/home") as string;

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
  // If arriving from logout, a cover div hides the teardown. Fade it out
  // after the login card has had a moment to render.
  const cover = document.getElementById("logout-cover");
  if (cover) {
    requestAnimationFrame(() => {
      cover.style.transition = "opacity 0.25s ease-out";
      cover.style.opacity = "0";
      cover.addEventListener("transitionend", () => cover.remove());
    });
  }

  if (!recaptcha) return;

  window.grecaptcha.ready(function () {
    window.grecaptcha.render("recaptcha", {
      sitekey: recaptchaKey,
    });
  });
});
</script>
