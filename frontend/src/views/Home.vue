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
        <p class="home__status" aria-live="polite">
          <template v-if="refreshing">
            <span class="home__pulse" aria-hidden="true"></span>
            {{ t("home.updating") }}
          </template>
        </p>
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
            :message="t('home.scanning')"
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
            :message="t('home.scanning')"
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
            :message="t('home.scanning')"
          />
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import dayjs from "dayjs";
import prettyBytes from "pretty-bytes";
import HeaderBar from "@/components/header/HeaderBar.vue";
import SectionState from "@/components/home/SectionState.vue";
import { useAuthStore } from "@/stores/auth";
import { useHomeStore } from "@/stores/home";
import { files as filesApi } from "@/api";
import { fetchJSON, fetchURL, StatusError } from "@/api/utils";
import { encodePath } from "@/utils/url";
import {
  TYPE_ICONS,
  fileType,
  summarize,
  type ScanSummary,
} from "@/utils/fileStats";

type Status = "loading" | "ready" | "error";
interface Section<T> {
  status: Status;
  data: T | null;
}
interface Usage {
  used: number;
  total: number;
}

const { t, locale } = useI18n();
const authStore = useAuthStore();
const homeStore = useHomeStore();

const isAdmin = computed(() => authStore.user?.perm.admin ?? false);
const canShare = computed(
  () => !!(authStore.user?.perm.share && authStore.user?.perm.download)
);

const usage = reactive<Section<Usage>>({ status: "loading", data: null });
const scan = reactive<Section<ScanSummary>>({ status: "loading", data: null });
const shares = reactive<Section<number>>({ status: "loading", data: null });
const users = reactive<Section<number>>({ status: "loading", data: null });
const refreshing = ref(false);

let controller: AbortController | null = null;

async function run<T>(
  section: Section<T>,
  loader: (signal: AbortSignal) => Promise<T>,
  signal: AbortSignal,
  store: (value: T) => void
) {
  if (section.data === null) section.status = "loading";
  try {
    const value = await loader(signal);
    if (signal.aborted) return;
    section.data = value;
    section.status = "ready";
    store(value);
  } catch (e) {
    if (signal.aborted || (e instanceof StatusError && e.is_canceled)) return;
    // Keep showing stale data if we have it; otherwise show the error state
    if (section.data === null) section.status = "error";
  }
}

function load() {
  controller?.abort();
  controller = new AbortController();
  const { signal } = controller;
  const c = homeStore.cacheFor(authStore.user?.username ?? "");
  refreshing.value = true;

  // Every section loads independently and in parallel; leaving the page
  // aborts them, which also stops the server-side directory walk.
  const jobs: Promise<void>[] = [
    run(
      usage,
      (s) => filesApi.usage("/files/", s) as Promise<Usage>,
      signal,
      (v) => (c.usage = v)
    ),
    run(
      scan,
      async (s) => {
        const res = await fetchURL("/api/resources/recursive/", { signal: s });
        return summarize(await res.json());
      },
      signal,
      (v) => (c.scan = v)
    ),
  ];
  if (canShare.value) {
    jobs.push(
      run(
        shares,
        async (s) =>
          (await fetchJSON<unknown[]>("/api/shares", { signal: s })).length,
        signal,
        (v) => (c.shares = v)
      )
    );
  }
  if (isAdmin.value) {
    jobs.push(
      run(
        users,
        async (s) =>
          (await fetchJSON<unknown[]>("/api/users", { signal: s })).length,
        signal,
        (v) => (c.users = v)
      )
    );
  }

  Promise.allSettled(jobs).then(() => {
    if (!signal.aborted) refreshing.value = false;
  });
}

onMounted(() => {
  // Show the previous result instantly, then refresh in the background
  const cache = homeStore.cache;
  if (cache && cache.owner === (authStore.user?.username ?? "")) {
    if (cache.usage)
      Object.assign(usage, { status: "ready", data: cache.usage });
    if (cache.scan) Object.assign(scan, { status: "ready", data: cache.scan });
    if (cache.shares !== null)
      Object.assign(shares, { status: "ready", data: cache.shares });
    if (cache.users !== null)
      Object.assign(users, { status: "ready", data: cache.users });
  }
  load();
});

onBeforeUnmount(() => controller?.abort());

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
  const u = usage.data;
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
  const total = scan.data?.bytes ?? 0;
  if (!total) return "0%";
  const pct = (n / total) * 100;
  return pct > 0 && pct < 1 ? "<1%" : `${Math.round(pct)}%`;
};
const barWidth = (n: number) => {
  const total = scan.data?.bytes ?? 0;
  if (!total) return "0%";
  return `${Math.max((n / total) * 100, n > 0 ? 1 : 0)}%`;
};
</script>
