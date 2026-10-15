<template>
  <div
    id="search"
    @click="onRootClick"
    :class="{ active, ongoing }"
    :role="active ? 'dialog' : undefined"
    :aria-modal="active ? 'true' : undefined"
    :aria-label="active ? $t('search.title') : undefined"
  >
    <div class="search__panel">
      <div id="input">
        <button
          v-if="active"
          class="action"
          @click="close"
          :aria-label="closeButtonTitle"
          :title="closeButtonTitle"
        >
          <i v-if="ongoing" class="material-icons">stop_circle</i>
          <i v-else class="material-icons">arrow_back</i>
        </button>
        <i v-else class="material-icons" aria-hidden="true">search</i>
        <input
          type="search"
          @keyup.exact="keyup"
          @keyup.enter="submit"
          ref="input"
          :autofocus="active"
          v-model.trim="prompt"
          enterkeyhint="search"
          autocomplete="off"
          spellcheck="false"
          :aria-label="$t('search.search')"
          :placeholder="$t('search.placeholder')"
        />
        <button
          v-if="active && prompt.length > 0"
          class="action search__clear"
          @click.stop="clear"
          :aria-label="$t('buttons.clear')"
          :title="$t('buttons.clear')"
        >
          <i class="material-icons">close</i>
        </button>
        <kbd v-if="!active" class="search__kbd" aria-hidden="true">{{
          shortcut
        }}</kbd>
      </div>

      <!-- Thin indeterminate bar while the server is walking folders -->
      <div v-if="active" class="search__progress" :class="{ on: ongoing }" />

      <div id="result" ref="result">
        <div class="search__body">
          <!-- Where we search + type filters (always available) -->
          <div class="search__context">
            <p class="search__scope">
              <i class="material-icons" aria-hidden="true">folder_open</i>
              <span>
                {{ $t("search.searchingIn") }}
                <strong>{{ scopeLabel }}</strong>
              </span>
            </p>
            <p class="search__scope-hint">{{ $t("search.scopeHint") }}</p>
          </div>

          <div
            class="search__chips"
            role="group"
            :aria-label="$t('search.filterByType')"
          >
            <span class="search__chips-label">{{
              $t("search.filterByType")
            }}</span>
            <button
              v-for="f in filters"
              :key="f.token"
              type="button"
              class="search__chip"
              :class="{ on: hasToken(f.token) }"
              :aria-pressed="hasToken(f.token)"
              @click="toggleToken(f.token)"
            >
              <i class="material-icons" aria-hidden="true">{{
                hasToken(f.token) ? "check" : f.icon
              }}</i>
              {{ $t("search." + f.label) }}
            </button>
          </div>

          <!-- Nothing typed yet: examples + notes -->
          <template v-if="prompt.length === 0 && !ongoing">
            <section class="search__section">
              <div class="search__section-head">
                <h3 class="search__heading">{{ $t("search.tips") }}</h3>
                <span class="search__section-hint">{{
                  $t("search.tipsHint")
                }}</span>
              </div>
              <ul class="search__tips">
                <li v-for="tip in tips" :key="tip.example">
                  <button type="button" @click="init(tip.example)">
                    <code>{{ tip.example }}</code>
                    <span>{{ $t(tip.label) }}</span>
                  </button>
                </li>
              </ul>
            </section>

            <section class="search__section search__section--notes">
              <h3 class="search__heading">
                <i class="material-icons" aria-hidden="true">info</i>
                {{ $t("search.notes") }}
              </h3>
              <ul class="search__notes">
                <li v-for="n in notes" :key="n">{{ $t(n) }}</li>
              </ul>
            </section>
          </template>

          <!-- Typed, not yet submitted: one clear next step -->
          <button
            v-else-if="!ongoing && !searched && results.length === 0"
            type="button"
            class="search__go"
            @click="submit"
          >
            <i class="material-icons" aria-hidden="true">search</i>
            <span class="search__go-text">
              <span class="search__go-title">{{
                $t("search.searchFor", { query: prompt })
              }}</span>
              <span class="search__go-sub">{{
                $t("search.inScope", { scope: scopeLabel })
              }}</span>
            </span>
            <kbd>Enter ↵</kbd>
          </button>

          <!-- Status line while / after searching -->
          <p
            v-if="ongoing || results.length > 0"
            class="search__status"
            aria-live="polite"
          >
            <span v-if="ongoing && results.length === 0">{{
              $t("search.searching")
            }}</span>
            <span v-else>{{
              $t("search.results", { count: results.length }, results.length)
            }}</span>
            <span v-if="ongoing" class="search__status-hint">{{
              $t("search.liveHint")
            }}</span>
          </p>

          <ul v-if="results.length > 0" class="search__results">
            <li v-for="(s, k) in filteredResults" :key="k">
              <router-link @click="close" :to="s.url" class="search__result">
                <span
                  class="search__icon"
                  :class="{ 'search__icon--dir': s.dir }"
                  aria-hidden="true"
                >
                  <i class="material-icons">{{ iconOf(s) }}</i>
                </span>
                <span class="search__text">
                  <span class="search__name">{{ nameOf(s.path) }}</span>
                  <span class="search__path">{{ folderOf(s.path) }}</span>
                </span>
                <i class="material-icons search__open" aria-hidden="true"
                  >arrow_forward</i
                >
              </router-link>
            </li>
          </ul>

          <!-- Finished with nothing found -->
          <div
            v-if="searched && !ongoing && results.length === 0"
            class="search__empty"
            aria-live="polite"
          >
            <span class="search__empty-icon" aria-hidden="true">
              <i class="material-icons">search_off</i>
            </span>
            <p class="search__empty-title">
              {{ $t("search.noResults", { query: searched }) }}
            </p>
            <p>{{ $t("search.noResultsHint") }}</p>
          </div>
        </div>
      </div>

      <div v-if="active" class="search__footer">
        <span><kbd>Enter</kbd> {{ $t("search.toSearch") }}</span>
        <span><kbd>Esc</kbd> {{ $t("search.toClose") }}</span>
        <span class="search__footer-note">{{ $t("search.namesOnly") }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import url from "@/utils/url";
import { search } from "@/api";
import { computed, inject, onMounted, ref, watch, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";
import { storeToRefs } from "pinia";
import { StatusError } from "@/api/utils";
import { TYPE_ICONS, fileType } from "@/utils/fileStats";

// Type filters understood by the server (several = match any of them)
const filters = [
  { token: "type:image", label: "images", icon: "image" },
  { token: "type:audio", label: "music", icon: "audiotrack" },
  { token: "type:video", label: "video", icon: "movie" },
  { token: "type:pdf", label: "pdf", icon: "picture_as_pdf" },
];

const notes = [
  "search.noteNames",
  "search.noteSubfolders",
  "search.noteCombine",
  "search.noteLive",
];

const tips = [
  { example: "report", label: "search.tipWord" },
  { example: "invoice receipt", label: "search.tipAny" },
  { example: '"annual report"', label: "search.tipPhrase" },
  { example: "type:docx", label: "search.tipType" },
  { example: "case:sensitive README", label: "search.tipCase" },
];

// Shortcut handled by the file listing (Ctrl/⌘ + Shift + F)
const shortcut = /Mac|iPhone|iPad/.test(navigator.userAgent)
  ? "⌘ ⇧ F"
  : "Ctrl ⇧ F";

const layoutStore = useLayoutStore();
const fileStore = useFileStore();
let searchAbortController = new AbortController();

const { currentPromptName } = storeToRefs(layoutStore);

const prompt = ref<string>("");
const active = ref<boolean>(false);
const ongoing = ref<boolean>(false);
const results = ref<any[]>([]);
const reload = ref<boolean>(false);
const resultsCount = ref<number>(50);
// The query of the last finished search ("" when none)
const searched = ref<string>("");

const $showError = inject<IToastError>("$showError")!;

const input = ref<HTMLInputElement | null>(null);
const result = ref<HTMLElement | null>(null);

const { t } = useI18n();

const route = useRoute();

watch(currentPromptName, (newVal, oldVal) => {
  active.value = newVal === "search";

  if (oldVal === "search" && !active.value) {
    if (reload.value) {
      fileStore.reload = true;
    }

    document.body.style.overflow = "auto";
    reset();
    prompt.value = "";
    active.value = false;
    input.value?.blur();
  } else if (active.value) {
    reload.value = false;
    input.value?.focus();
    document.body.style.overflow = "hidden";
  }
});

watch(prompt, () => {
  reset();
});

// ...mapState(useFileStore, ["isListing"]),
// ...mapState(useLayoutStore, ["show"]),
// ...mapWritableState(useFileStore, { sReload: "reload" }),

const filteredResults = computed(() => {
  return results.value.slice(0, resultsCount.value);
});

// Folder the search runs in, as the user sees it
const searchPath = () =>
  fileStore.isListing ? route.path : url.removeLastDir(route.path) + "/";

const scopeLabel = computed(() => {
  let path = searchPath().replace(/^\/files/, "");
  try {
    path = decodeURIComponent(path);
  } catch {
    // keep as-is
  }
  path = path.replace(/\/+$/, "");
  return path === "" ? t("sidebar.myFiles") : path;
});

const nameOf = (path: string) => {
  const p = path.replace(/\/+$/, "");
  return p.slice(p.lastIndexOf("/") + 1);
};
const folderOf = (path: string) => {
  const p = path.replace(/\/+$/, "");
  const i = p.lastIndexOf("/");
  return i === -1 ? scopeLabel.value : `${scopeLabel.value}/${p.slice(0, i)}`;
};

const iconOf = (s: { dir: boolean; path: string }) =>
  s.dir ? "folder" : TYPE_ICONS[fileType(nameOf(s.path))];

// Filter chips add or remove their token in the query
const tokensOf = (q: string) => q.split(/\s+/).filter(Boolean);
const hasToken = (token: string) => tokensOf(prompt.value).includes(token);
const toggleToken = (token: string) => {
  const parts = tokensOf(prompt.value);
  const next = hasToken(token)
    ? parts.filter((p) => p !== token)
    : [token, ...parts];
  prompt.value = next.length ? next.join(" ") + " " : "";
  input.value?.focus();
};

const closeButtonTitle = computed(() => {
  return ongoing.value ? t("buttons.stopSearch") : t("buttons.close");
});

onMounted(() => {
  if (result.value === null) {
    return;
  }
  result.value.addEventListener("scroll", (event: Event) => {
    if (
      (event.target as HTMLElement).offsetHeight +
        (event.target as HTMLElement).scrollTop >=
      (event.target as HTMLElement).scrollHeight - 100
    ) {
      resultsCount.value += 50;
    }
  });
});

onUnmounted(() => {
  abortLastSearch();
});

const open = () => {
  !active.value && layoutStore.showHover("search");
};

// Closed: any click opens. Open: a click on the backdrop closes.
const onRootClick = (event: MouseEvent) => {
  if (!active.value) {
    open();
  } else if (event.target === event.currentTarget && !ongoing.value) {
    layoutStore.closeHovers();
  }
};

const close = (event: Event) => {
  if (ongoing.value) {
    abortLastSearch();
    ongoing.value = false;
  } else {
    event.stopPropagation();
    event.preventDefault();
    layoutStore.closeHovers();
  }
};

const keyup = (event: KeyboardEvent) => {
  if (event.key === "Escape") {
    close(event);
    return;
  }
  results.value.length = 0;
};

const init = (string: string) => {
  prompt.value = `${string} `;
  input.value !== null ? input.value.focus() : "";
};

const clear = () => {
  prompt.value = "";
  input.value?.focus();
};

const reset = () => {
  searched.value = "";
  abortLastSearch();
  ongoing.value = false;
  resultsCount.value = 50;
  results.value = [];
};

const abortLastSearch = () => {
  searchAbortController.abort();
};

const submit = async (event: Event) => {
  event.preventDefault();

  if (prompt.value === "") {
    return;
  }

  const path = searchPath();
  const query = prompt.value;

  ongoing.value = true;
  searched.value = "";

  try {
    abortLastSearch();
    searchAbortController = new AbortController();
    results.value = [];
    await search(path, prompt.value, searchAbortController.signal, (item) =>
      results.value.push(item)
    );
  } catch (error: any) {
    if (error instanceof StatusError && error.is_canceled) {
      return;
    }
    $showError(error);
  }

  ongoing.value = false;
  searched.value = query;
};
</script>
