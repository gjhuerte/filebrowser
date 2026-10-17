<template>
  <div>
    <h3 v-if="!isDefault">{{ t("settings.basicInfo") }}</h3>
    <p v-if="!isDefault" class="small">
      {{ t("settings.basicInfoDescription") }}
    </p>

    <p v-if="!isDefault && props.user !== null">
      <label for="username">{{ t("settings.username") }}</label>
      <input
        class="input input--block"
        type="text"
        v-model="user.username"
        id="username"
      />
    </p>

    <div v-if="!isDefault">
      <label for="password">{{ t("settings.password") }}</label>
      <div class="password-wrap">
        <input
          class="input input--block"
          :type="showPassword ? 'text' : 'password'"
          :placeholder="passwordPlaceholder"
          v-model="user.password"
          id="password"
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
      <PasswordStrength v-if="user.password" :password="user.password" />
    </div>

    <div>
      <label for="scope">{{ t("settings.scope") }}</label>
      <span class="small">{{ t("settings.scopeDescription") }}</span>
      <input
        :disabled="createUserDirData ?? false"
        :placeholder="scopePlaceholder"
        class="input input--block"
        type="text"
        v-model="user.scope"
        id="scope"
      />
    </div>
    <div class="setting-list" v-if="displayHomeDirectoryCheckbox">
      <label class="setting-toggle">
        <span>{{ t("settings.createUserHomeDirectory") }}</span>
        <input type="checkbox" class="switch" v-model="createUserDirData" />
      </label>
    </div>

    <p>
      <label for="locale">{{ t("settings.language") }}</label>
      <languages
        class="input input--block"
        id="locale"
        v-model:locale="user.locale"
      ></languages>
    </p>

    <div v-if="!isDefault" class="collapsible user-form__advanced">
      <input id="advanced-toggle" type="checkbox" />
      <label for="advanced-toggle">
        <p>{{ t("settings.advanced") }}</p>
        <i class="material-icons">arrow_drop_down</i>
      </label>
      <div class="collapse">
        <p class="small">{{ t("settings.advancedDescription") }}</p>

        <div class="setting-list" v-if="user.perm">
          <label class="setting-toggle">
            <span>{{ t("settings.lockPassword") }}</span>
            <input
              type="checkbox"
              class="switch"
              :disabled="user.perm.admin"
              v-model="user.lockPassword"
            />
          </label>
        </div>

        <permissions v-model:perm="user.perm" />
        <commands v-if="enableExec" v-model:commands="user.commands" />

        <h3>{{ t("settings.rules") }}</h3>
        <p class="small">{{ t("settings.rulesHelp") }}</p>
        <rules v-model:rules="user.rules" />
      </div>
    </div>

    <div v-if="isDefault">
      <div class="setting-list" v-if="user.perm">
        <label class="setting-toggle">
          <span>{{ t("settings.lockPassword") }}</span>
          <input
            type="checkbox"
            class="switch"
            :disabled="user.perm.admin"
            v-model="user.lockPassword"
          />
        </label>
      </div>

      <permissions v-model:perm="user.perm" />
      <commands v-if="enableExec" v-model:commands="user.commands" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Languages from "./Languages.vue";
import Rules from "./Rules.vue";
import Permissions from "./Permissions.vue";
import Commands from "./Commands.vue";
import PasswordStrength from "@/components/PasswordStrength.vue";
import { enableExec } from "@/utils/constants";
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const showPassword = ref<boolean>(false);
const createUserDirData = ref<boolean | null>(null);
const originalUserScope = ref<string | null>(null);

const props = defineProps<{
  user: IUserForm;
  isNew: boolean;
  isDefault: boolean;
  createUserDir?: boolean;
}>();

onMounted(() => {
  if (props.user.scope) {
    originalUserScope.value = props.user.scope;
    createUserDirData.value = props.createUserDir;
  }
});

const passwordPlaceholder = computed(() =>
  props.isNew ? "" : t("settings.avoidChanges")
);
const scopePlaceholder = computed(() =>
  createUserDirData.value ? t("settings.userScopeGenerationPlaceholder") : ""
);
const displayHomeDirectoryCheckbox = computed(
  () => props.isNew && createUserDirData.value
);

watch(
  () => props.user,
  () => {
    if (!props.user?.perm?.admin) return;
    props.user.lockPassword = false;
  }
);

watch(createUserDirData, () => {
  if (props.user?.scope) {
    props.user.scope = createUserDirData.value
      ? ""
      : (originalUserScope.value ?? "");
  }
});
</script>
