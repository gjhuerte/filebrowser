<template>
  <div class="dashboard">
    <header-bar showMenu showLogo />

    <div class="settings">
      <h1 class="settings__title">{{ t("sidebar.sharedLinks") }}</h1>
      <p class="settings__subtitle">
        {{ t("settings.sharesSubtitle") }}
      </p>

      <div v-if="layoutStore.loading" class="row">
        <div class="column">
          <div class="card">
            <div class="card-content">
              <h2 class="message delayed">
                <div class="spinner">
                  <div class="bounce1"></div>
                  <div class="bounce2"></div>
                  <div class="bounce3"></div>
                </div>
                <span>{{ t("files.loading") }}</span>
              </h2>
            </div>
          </div>
        </div>
      </div>

      <errors v-else-if="error" :errorCode="error.status" />

      <div class="row" v-else>
        <div class="column">
          <div class="card">
            <div class="card-content full" v-if="links.length > 0">
              <table>
                <tr>
                  <th>{{ t("settings.path") }}</th>
                  <th>{{ t("settings.shareDuration") }}</th>
                  <th v-if="authStore.user?.perm.admin">
                    {{ t("settings.username") }}
                  </th>
                  <th></th>
                  <th></th>
                </tr>

                <tr v-for="link in links" :key="link.hash">
                  <td>
                    <a :href="buildLink(link)" target="_blank">{{
                      link.path
                    }}</a>
                  </td>
                  <td>
                    <span
                      v-if="link.expire !== 0"
                      class="badge badge--muted"
                      >{{ humanTime(link.expire) }}</span
                    >
                    <span v-else class="badge">{{ t("permanent") }}</span>
                  </td>
                  <td v-if="authStore.user?.perm.admin">
                    {{ link.username }}
                  </td>
                  <td class="small">
                    <button
                      class="action"
                      @click="deleteLink($event, link)"
                      :aria-label="t('buttons.delete')"
                      :title="t('buttons.delete')"
                    >
                      <i class="material-icons">delete</i>
                    </button>
                  </td>
                  <td class="small">
                    <button
                      class="action copy-clipboard"
                      :aria-label="t('buttons.copyToClipboard')"
                      :title="t('buttons.copyToClipboard')"
                      @click="copyToClipboard(buildLink(link))"
                    >
                      <i class="material-icons">content_paste</i>
                    </button>
                  </td>
                </tr>
              </table>
            </div>
            <empty-state
              v-else
              compact
              icon="link"
              :title="t('settings.noShares')"
              :description="t('settings.noSharesText')"
            >
              <router-link to="/files/" class="button button--secondary">
                <i class="material-icons" aria-hidden="true">folder</i>
                {{ t("sidebar.myFiles") }}
              </router-link>
            </empty-state>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useLayoutStore } from "@/stores/layout";
import { share as api, users } from "@/api";
import dayjs from "dayjs";
import HeaderBar from "@/components/header/HeaderBar.vue";
import Errors from "@/views/Errors.vue";
import EmptyState from "@/components/EmptyState.vue";
import { inject, ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { StatusError } from "@/api/utils";
import { copy } from "@/utils/clipboard";

const $showError = inject<IToastError>("$showError")!;
const $showSuccess = inject<IToastSuccess>("$showSuccess")!;
const { t } = useI18n();

const layoutStore = useLayoutStore();
const authStore = useAuthStore();

const error = ref<StatusError | null>(null);
const links = ref<Share[]>([]);

onMounted(async () => {
  layoutStore.loading = true;

  try {
    const newLinks = await api.list();
    if (authStore.user?.perm.admin) {
      const userMap = new Map<number, string>();
      for (const user of await users.getAll())
        userMap.set(user.id, user.username);
      for (const link of newLinks) {
        if (link.userID && userMap.has(link.userID))
          link.username = userMap.get(link.userID);
      }
    }
    links.value = newLinks;
  } catch (err) {
    if (err instanceof Error) {
      error.value = err;
    }
  } finally {
    layoutStore.loading = false;
  }
});

const copyToClipboard = (text: string) => {
  copy({ text }).then(
    () => {
      $showSuccess(t("success.linkCopied"));
    },
    () => {
      copy({ text }, { permission: true }).then(
        () => {
          $showSuccess(t("success.linkCopied"));
        },
        (e) => {
          $showError(e);
        }
      );
    }
  );
};

const deleteLink = async (event: Event, link: any) => {
  event.preventDefault();

  layoutStore.showHover({
    prompt: "share-delete",
    confirm: () => {
      layoutStore.closeHovers();

      try {
        api.remove(link.hash);
        links.value = links.value.filter((item) => item.hash !== link.hash);
        $showSuccess(t("settings.shareDeleted"));
      } catch (err) {
        if (err instanceof Error) {
          $showError(err);
        }
      }
    },
  });
};

const humanTime = (time: number) => {
  return dayjs(time * 1000).fromNow();
};

const buildLink = (share: Share) => {
  return api.getShareURL(share);
};
</script>
