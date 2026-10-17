<template>
  <errors v-if="error" :errorCode="error.status" />
  <div class="row" v-else-if="!layoutStore.loading">
    <div class="column">
      <div class="card">
        <div class="card-title">
          <h2>{{ t("settings.users") }}</h2>
          <router-link to="/admin/users/new" class="button button--icon">
            <i class="material-icons" aria-hidden="true">add</i>
            {{ t("buttons.new") }}
          </router-link>
        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>{{ t("settings.username") }}</th>
              <th>{{ t("settings.admin") }}</th>
              <th>{{ t("settings.scope") }}</th>
              <th></th>
            </tr>

            <tr v-for="user in users" :key="user.id">
              <td>
                <router-link :to="'/admin/users/' + user.id" class="cell-user">
                  <span class="avatar" aria-hidden="true">{{
                    user.username.charAt(0).toUpperCase()
                  }}</span>
                  {{ user.username }}
                </router-link>
              </td>
              <td>
                <span v-if="user.perm.admin" class="badge">{{
                  t("settings.admin")
                }}</span>
                <span v-else class="badge badge--muted">—</span>
              </td>
              <td>
                <code>{{ user.scope }}</code>
              </td>
              <td class="cell-actions">
                <router-link
                  :to="'/admin/users/' + user.id"
                  class="action"
                  :aria-label="`${t('settings.user')} ${user.username}`"
                  :title="`${t('settings.user')} ${user.username}`"
                  ><i class="material-icons">edit</i></router-link
                >
              </td>
            </tr>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import Errors from "@/views/Errors.vue";
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { StatusError } from "@/api/utils";

const error = ref<StatusError | null>(null);
const users = ref<IUser[]>([]);

const layoutStore = useLayoutStore();
const { t } = useI18n();

onMounted(async () => {
  layoutStore.loading = true;

  try {
    users.value = await api.getAll();
  } catch (err) {
    if (err instanceof Error) {
      error.value = err;
    }
  } finally {
    layoutStore.loading = false;
  }
});
</script>
