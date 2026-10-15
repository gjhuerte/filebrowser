<template>
  <div class="card floating confirm-logout">
    <div class="card-content">
      <span class="confirm-logout__icon" aria-hidden="true">
        <i class="material-icons">logout</i>
      </span>
      <h2 class="confirm-logout__title">{{ $t("prompts.logoutTitle") }}</h2>
      <p class="confirm-logout__text">
        {{
          $t("prompts.logoutMessage", { name: authStore.user?.username ?? "" })
        }}
      </p>
      <p v-if="uploading" class="confirm-logout__warning" role="alert">
        <i class="material-icons" aria-hidden="true">warning_amber</i>
        <span>{{ $t("prompts.logoutUploading") }}</span>
      </p>
    </div>
    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="layoutStore.closeHovers()"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
        tabindex="2"
      >
        {{ $t("buttons.cancel") }}
      </button>
      <button
        id="focus-prompt"
        class="button button--flat button--red"
        @click="logOut"
        :aria-label="$t('prompts.logoutConfirm')"
        :title="$t('prompts.logoutConfirm')"
        tabindex="1"
      >
        {{ $t("prompts.logoutConfirm") }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { useUploadStore } from "@/stores/upload";
import { logout } from "@/utils/auth";

const authStore = useAuthStore();
const layoutStore = useLayoutStore();
const uploadStore = useUploadStore();

// Leaving the page stops any upload that's still sending
const uploading = computed(() => uploadStore.totalBytes > 0);

const logOut = () => {
  layoutStore.closeHovers();
  logout();
};
</script>
