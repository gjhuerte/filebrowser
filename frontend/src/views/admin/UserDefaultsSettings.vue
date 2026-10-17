<template>
  <div class="row">
    <div class="column">
      <form class="card" @submit.prevent="$emit('save')">
        <div class="card-title">
          <h2>{{ t("settings.userDefaults") }}</h2>
        </div>

        <div class="card-content">
          <p class="small">{{ t("settings.defaultUserDescription") }}</p>

          <user-form
            :isNew="false"
            :isDefault="true"
            v-model:user="settings.defaults"
          />
        </div>

        <div class="card-action">
          <input class="button" type="submit" :value="t('buttons.update')" />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import UserForm from "@/components/settings/UserForm.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

defineProps<{
  settings: ISettings;
}>();

defineEmits<{
  save: [];
}>();
</script>
