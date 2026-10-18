<template>
  <div class="home">
    <header-bar showMenu showLogo />

    <div class="home__inner">
      <div class="home__heading">
        <div>
          <h1 class="home__title">
            {{ t("home.greeting", { name: authStore.user?.username }) }}
          </h1>
          <p class="home__subtitle">{{ t("home.subtitle") }}</p>
        </div>
        <div class="home__sync">
          <p class="home__status" aria-live="polite">
            <template v-if="refreshing">
              <span class="home__pulse" aria-hidden="true"></span>
              {{ progressLabel }}
            </template>
            <template v-else-if="updatedLabel">
              {{ updatedLabel
              }}<span v-if="nextLabel" class="home__next">
                · {{ nextLabel }}</span
              >
            </template>
          </p>
          <div class="home__controls">
            <label class="home__interval">
              <i class="material-icons" aria-hidden="true">schedule</i>
              <select
                v-model="intervalModel"
                :title="t('home.autoRefreshHint')"
                :aria-label="t('home.autoRefresh')"
              >
                <option v-for="m in REFRESH_CHOICES" :key="m" :value="m">
                  {{ intervalLabel(m) }}
                </option>
              </select>
            </label>
            <button
              type="button"
              class="home__refresh"
              :class="{ 'home__refresh--spinning': refreshing }"
              :disabled="refreshing"
              :aria-label="t('home.refreshNow')"
              :title="t('home.refreshNow')"
              @click="load"
            >
              <i class="material-icons" aria-hidden="true">refresh</i>
            </button>
          </div>
        </div>
      </div>

      <!-- Headline numbers -->
      <section class="home__tiles" :aria-label="t('home.overview')">
        <div class="tile">
          <p class="tile__label">{{ t("home.storage") }}</p>
          <template v-if="usage.data">
            <p class="tile__value">{{ bytes(usage.data.used) }}</p>
            <p class="tile__meta">
              {{ t("home.storageOf", { total: bytes(usage.data.total) }) }}
              · {{ t("home.percentUsed", { percent: usedPercent }) }}
            </p>
            <div
              class="meter"
              :class="meterLevel"
              role="meter"
              :aria-valuenow="usedPercent"
              aria-valuemin="0"
              aria-valuemax="100"
              :aria-label="t('home.storage')"
            >
              <div class="meter__fill" :style="{ width: usedPercent + '%' }" />
            </div>
          </template>
          <SectionState
            v-else
            :status="usage.status"
            @retry="load"
            skeleton="tile"
          />
        </div>

        <div class="tile">
          <p class="tile__label">{{ t("home.files") }}</p>
          <template v-if="scan.data">
            <p class="tile__value">{{ count(scan.data.files) }}</p>
            <p class="tile__meta">
              {{ t("home.totalSize", { size: bytes(scan.data.bytes) }) }}
            </p>
          </template>
          <SectionState
            v-else
            :status="scan.status"
            @retry="load"
            skeleton="tile"
          />
        </div>

        <div class="tile">
          <p class="tile__label">{{ t("home.folders") }}</p>
          <template v-if="scan.data">
            <p class="tile__value">{{ count(scan.data.folders) }}</p>
            <p class="tile__meta">{{ t("home.inYourFiles") }}</p>
          </template>
          <SectionState
            v-else
            :status="scan.status"
            @retry="load"
            skeleton="tile"
          />
        </div>

        <router-link
          v-if="canShare"
          to="/settings/shares"
          class="tile tile--link"
        >
          <p class="tile__label">{{ t("home.sharedLinks") }}</p>
          <template v-if="shares.data !== null">
            <p class="tile__value">{{ count(shares.data) }}</p>
            <p class="tile__meta tile__meta--link">
              {{ t("home.manage") }}
              <i class="material-icons" aria-hidden="true">arrow_forward</i>
            </p>
          </template>
          <SectionState
            v-else
            :status="shares.status"
            @retry="load"
            skeleton="tile"
          />
        </router-link>

        <router-link
          v-if="isAdmin"
          to="/settings/users"
          class="tile tile--link"
        >
          <p class="tile__label">{{ t("home.users") }}</p>
          <template v-if="users.data !== null">
            <p class="tile__value">{{ count(users.data) }}</p>
            <p class="tile__meta tile__meta--link">
              {{ t("home.manage") }}
              <i class="material-icons" aria-hidden="true">arrow_forward</i>
            </p>
          </template>
          <SectionState
            v-else
            :status="users.status"
            @retry="load"
            skeleton="tile"
          />
        </router-link>
      </section>

      <div class="home__grid">
        <!-- Storage by file type -->
        <section class="panel">
          <h2 class="panel__title">{{ t("home.byType") }}</h2>
          <ul v-if="scan.data && scan.data.byType.length" class="types">
            <li
              v-for="row in scan.data.byType"
              :key="row.key"
              class="types__row"
            >
              <i class="material-icons types__icon" aria-hidden="true">{{
                TYPE_ICONS[row.key]
              }}</i>
              <div class="types__body">
                <div class="types__line">
                  <span class="types__name">{{
                    t(`home.types.${row.key}`)
                  }}</span>
                  <span class="types__value"
                    >{{ bytes(row.bytes) }}
                    <span class="types__share">{{
                      share(row.bytes)
                    }}</span></span
                  >
                </div>
                <div class="types__track">
                  <div
                    class="types__bar"
                    :style="{ width: barWidth(row.bytes) }"
                  ></div>
                </div>
                <span class="types__count">{{
                  t("home.fileCount", { count: count(row.count) }, row.count)
                }}</span>
              </div>
            </li>
          </ul>
          <p v-else-if="scan.data" class="panel__empty">
            {{ t("home.noFiles") }}
          </p>
          <SectionState
            v-else
            :status="scan.status"
            @retry="load"
            skeleton="list"
            :message="progressLabel"
          />
        </section>

        <!-- Recently modified -->
        <section class="panel">
          <h2 class="panel__title">{{ t("home.recent") }}</h2>
          <ul v-if="scan.data && scan.data.recent.length" class="entries">
            <li v-for="f in scan.data.recent" :key="f.path">
              <router-link :to="fileLink(f.path)" class="entries__row">
                <i class="material-icons entries__icon" aria-hidden="true">{{
                  TYPE_ICONS[fileType(f.name)]
                }}</i>
                <span class="entries__text">
                  <span class="entries__name">{{ f.name }}</span>
                  <span class="entries__path">{{ folderOf(f.path) }}</span>
                </span>
                <time class="entries__meta" :datetime="f.modified">{{
                  ago(f.modified)
                }}</time>
              </router-link>
            </li>
          </ul>
          <p v-else-if="scan.data" class="panel__empty">
            {{ t("home.noFiles") }}
          </p>
          <SectionState
            v-else
            :status="scan.status"
            @retry="load"
            skeleton="list"
            :message="progressLabel"
          />
        </section>

        <!-- Largest files -->
        <section class="panel">
          <h2 class="panel__title">{{ t("home.largest") }}</h2>
          <ul v-if="scan.data && scan.data.largest.length" class="entries">
            <li v-for="f in scan.data.largest" :key="f.path">
              <router-link :to="fileLink(f.path)" class="entries__row">
                <i class="material-icons entries__icon" aria-hidden="true">{{
                  TYPE_ICONS[fileType(f.name)]
                }}</i>
                <span class="entries__text">
                  <span class="entries__name">{{ f.name }}</span>
                  <span class="entries__path">{{ folderOf(f.path) }}</span>
                </span>
                <span class="entries__meta">{{ bytes(f.size) }}</span>
              </router-link>
            </li>
          </ul>
          <p v-else-if="scan.data" class="panel__empty">
            {{ t("home.noFiles") }}
          </p>
          <SectionState
            v-else
            :status="scan.status"
            @retry="load"
            skeleton="list"
            :message="progressLabel"
          />
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useI18n } from "vue-i18n";
import dayjs from "dayjs";
import prettyBytes from "pretty-bytes";
import HeaderBar from "@/components/header/HeaderBar.vue";
import SectionState from "@/components/home/SectionState.vue";
import { useAuthStore } from "@/stores/auth";
import { REFRESH_CHOICES, useHomeStore } from "@/stores/home";
import { encodePath } from "@/utils/url";
import { TYPE_ICONS, fileType } from "@/utils/fileStats";

const { t, locale } = useI18n();
const authStore = useAuthStore();
const homeStore = useHomeStore();
const { usage, scan, shares, users, refreshing, progress } =
  storeToRefs(homeStore);

const isAdmin = computed(() => authStore.user?.perm.admin ?? false);
const canShare = computed(
  () => !!(authStore.user?.perm.share && authStore.user?.perm.download)
);

const load = () => homeStore.refresh();

onMounted(() => {
  // Show the saved result. Rescanning happens on the auto-refresh schedule
  // or when the user presses Refresh; only a first visit loads on its own.
  homeStore.ensureOwner(authStore.user?.username ?? "");
  if (homeStore.isEmpty() && !homeStore.refreshing) load();
});

/* Auto-refresh */
const intervalModel = computed({
  get: () => homeStore.interval,
  set: (v: number) => homeStore.setInterval(Number(v)),
});
const intervalLabel = (m: number) =>
  m === 0
    ? t("home.autoRefreshOff")
    : m >= 60
      ? t("home.everyHours", { n: m / 60 }, m / 60)
      : t("home.everyMinutes", { n: m });
const updatedLabel = computed(() => {
  if (!homeStore.lastUpdated) return "";
  // Read `now` so the label re-renders on every scheduler tick
  const diff = homeStore.now - homeStore.lastUpdated;
  return diff < 60_000
    ? t("home.updatedJustNow")
    : t("home.updatedAgo", { time: dayjs(homeStore.lastUpdated).fromNow() });
});
const nextLabel = computed(() => {
  const at = homeStore.nextRefreshAt;
  if (!at) return "";
  const mins = Math.max(1, Math.ceil((at - homeStore.now) / 60_000));
  return t("home.nextScan", { n: mins });
});

/* "Scanning… 12,345 files found so far" while the server scan runs */
const progressLabel = computed(() =>
  progress.value && progress.value.files > 0
    ? t(
        "home.scanProgress",
        { count: count(progress.value.files) },
        progress.value.files
      )
    : refreshing.value
      ? t("home.updating")
      : t("home.scanning")
);

/* Formatting */
const bytes = (n: number) => prettyBytes(n, { binary: true });
const count = (n: number) =>
  new Intl.NumberFormat(locale.value, {
    notation: n >= 10000 ? "compact" : "standard",
    maximumFractionDigits: 1,
  }).format(n);
const ago = (iso: string) => dayjs(iso).fromNow();
const folderOf = (path: string) => {
  const dir = path.slice(0, path.lastIndexOf("/"));
  return dir || "/";
};
const fileLink = (path: string) => "/files" + encodePath(path);

const usedPercent = computed(() => {
  const u = usage.value.data;
  return u && u.total > 0 ? Math.round((u.used / u.total) * 100) : 0;
});
const meterLevel = computed(() =>
  usedPercent.value >= 90
    ? "meter--danger"
    : usedPercent.value >= 80
      ? "meter--warning"
      : ""
);

const share = (n: number) => {
  const total = scan.value.data?.bytes ?? 0;
  if (!total) return "0%";
  const pct = (n / total) * 100;
  return pct > 0 && pct < 1 ? "<1%" : `${Math.round(pct)}%`;
};
const barWidth = (n: number) => {
  const total = scan.value.data?.bytes ?? 0;
  if (!total) return "0%";
  return `${Math.max((n / total) * 100, n > 0 ? 1 : 0)}%`;
};
</script>
