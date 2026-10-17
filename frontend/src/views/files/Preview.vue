<template>
  <div
    id="previewer"
    :class="`previewer--${kind}`"
    @touchmove="blockTouchScroll"
    @wheel.prevent.stop
    @mousemove="toggleNavigation"
    @touchstart="toggleNavigation"
  >
    <file-header
      :name="name"
      :icon="headerIcon"
      :meta="headerMeta"
      :close-label="$t('buttons.close')"
      @close="close()"
    >
      <template #actions>
        <action
          :disabled="layoutStore.loading"
          v-if="isCsv && authStore.user?.perm.modify"
          icon="edit_note"
          :label="t('buttons.editAsText')"
          @action="editAsText"
        />
        <action
          :disabled="layoutStore.loading"
          v-if="authStore.user?.perm.rename"
          icon="drive_file_rename_outline"
          :label="$t('buttons.rename')"
          show="rename"
        />
        <action
          :disabled="layoutStore.loading"
          v-if="
            ['image', 'audio', 'video'].includes(fileStore.req?.type || '') &&
            authStore.user?.perm.download
          "
          icon="open_in_new"
          :label="t('buttons.openDirect')"
          @action="openDirect"
        />
        <action
          :disabled="layoutStore.loading"
          icon="info"
          :label="$t('buttons.info')"
          show="info"
        />
        <template v-if="authStore.user?.perm.delete">
          <span class="header__divider" aria-hidden="true"></span>
          <action
            :disabled="layoutStore.loading"
            icon="delete"
            :label="$t('buttons.delete')"
            @action="deleteFile"
            id="delete-button"
            class="action--danger"
          />
        </template>
      </template>
      <template v-if="authStore.user?.perm.download" #primary>
        <button
          type="button"
          class="button toolbar__upload file-header__primary"
          :disabled="layoutStore.loading"
          :aria-label="$t('buttons.download')"
          @click="download"
        >
          <i class="material-icons" aria-hidden="true">file_download</i>
          <span>{{ $t("buttons.download") }}</span>
        </button>
      </template>
    </file-header>

    <div class="previewer__stage">
      <MediaLoader v-if="layoutStore.loading" :label="t('files.loading')" />
      <div v-else class="preview">
        <div v-if="isEpub" class="epub-reader">
          <vue-reader
            :location="location"
            :url="previewUrl"
            :get-rendition="getRendition"
            :epubInitOptions="{
              requestCredentials: true,
            }"
            :epubOptions="{
              allowPopups: true,
            }"
            @update:location="locationChange"
          />
          <div class="size">
            <button
              @click="changeSize(Math.max(100, size - 10))"
              class="reader-button"
            >
              <i class="material-icons">remove</i>
            </button>
            <button
              @click="changeSize(Math.min(150, size + 10))"
              class="reader-button"
            >
              <i class="material-icons">add</i>
            </button>
            <span>{{ size }}%</span>
          </div>
        </div>
        <CsvViewer v-else-if="isCsv" :content="csvContent" :error="csvError" />

        <template v-else-if="kind === 'image'">
          <ExtendedImage
            v-if="image.src"
            ref="imageRef"
            :src="image.src"
            :alt="name"
            @zoom="(s: number) => (image.zoom = s)"
            @ready="onImageReady"
          />
          <MediaLoader
            v-if="image.state === 'loading'"
            :label="t('preview.loadingImage')"
            :loaded="image.loaded"
            :total="image.total"
          />
          <div v-else-if="image.state === 'error'" class="media-error">
            <i class="material-icons" aria-hidden="true">broken_image</i>
            <p class="media-error__title">{{ t("preview.imageFailed") }}</p>
            <p class="media-error__text">{{ t("preview.imageFailedText") }}</p>
            <button
              type="button"
              class="media-error__button"
              @click="loadImage"
            >
              <i class="material-icons" aria-hidden="true">refresh</i>
              {{ t("home.retry") }}
            </button>
          </div>
        </template>

        <AudioPlayer
          v-else-if="kind === 'audio'"
          :key="previewUrl"
          :src="previewUrl"
          :name="name"
          :size="fileStore.req?.size"
          :folder="folderName"
          :cover="coverUrl"
          :autoplay="autoPlay"
          :has-previous="hasPrevious"
          :has-next="hasNext"
          :can-download="!!authStore.user?.perm.download"
          @previous="prev"
          @next="next"
          @play="autoPlay = true"
          @pause="autoPlay = false"
          @download="download"
          ref="audioRef"
        />
        <div v-else-if="kind === 'video'" class="previewer__video">
          <VideoPlayer
            :key="previewUrl"
            :source="previewUrl"
            :subtitles="subtitles"
            :options="videoOptions"
            :size="fileStore.req?.size"
            :can-download="!!authStore.user?.perm.download"
            @download="download"
          />
        </div>
        <object v-else-if="isPdf" class="pdf" :data="previewUrl"></object>
        <div v-else-if="fileStore.req?.type == 'blob'" class="info">
          <span class="info__icon" aria-hidden="true">
            <i class="material-icons">description</i>
          </span>
          <p class="info__title">{{ $t("files.noPreview") }}</p>
          <p class="info__text">{{ t("preview.noPreviewHint") }}</p>
          <div class="info__actions">
            <a
              v-if="authStore.user?.perm.download"
              target="_blank"
              :href="downloadUrl"
              class="button button--flat"
            >
              <i class="material-icons" aria-hidden="true">file_download</i
              >{{ $t("buttons.download") }}
            </a>
            <a
              target="_blank"
              :href="previewUrl"
              class="button button--flat button--grey"
              v-if="!fileStore.req?.isDir"
            >
              <i class="material-icons" aria-hidden="true">open_in_new</i
              >{{ $t("buttons.openFile") }}
            </a>
          </div>
        </div>
      </div>
    </div>

    <!-- Image toolbar -->
    <div
      v-if="kind === 'image' && image.state === 'ready' && !layoutStore.loading"
      class="previewer__dock"
      :class="{ 'is-hidden': !showNav }"
      role="toolbar"
      :aria-label="t('preview.imageTools')"
      @mouseover="hoverNav = true"
      @mouseleave="hoverNav = false"
    >
      <button
        type="button"
        :aria-label="t('preview.zoomOut')"
        :title="t('preview.zoomOut') + ' (−)'"
        @click="imageRef?.zoomOut()"
      >
        <i class="material-icons" aria-hidden="true">remove</i>
      </button>
      <button
        type="button"
        class="previewer__zoom"
        :aria-label="t('preview.resetZoom')"
        :title="t('preview.resetZoom') + ' (0)'"
        @click="imageRef?.reset()"
      >
        {{ Math.round(image.zoom * 100) }}%
      </button>
      <button
        type="button"
        :aria-label="t('preview.zoomIn')"
        :title="t('preview.zoomIn') + ' (+)'"
        @click="imageRef?.zoomIn()"
      >
        <i class="material-icons" aria-hidden="true">add</i>
      </button>
      <span class="previewer__divider" aria-hidden="true"></span>
      <button
        type="button"
        :aria-label="t('preview.rotate')"
        :title="t('preview.rotate') + ' (R)'"
        @click="imageRef?.rotate()"
      >
        <i class="material-icons" aria-hidden="true">rotate_right</i>
      </button>
      <button
        v-if="isResizeEnabled"
        type="button"
        :class="{ 'is-active': fullSize }"
        :aria-pressed="fullSize"
        :aria-label="t('preview.original')"
        :title="t('preview.original')"
        @click="toggleSize"
      >
        <i class="material-icons" aria-hidden="true">hd</i>
      </button>
      <span v-if="imageMeta" class="previewer__info">
        <span class="previewer__divider" aria-hidden="true"></span>
        <span class="previewer__meta">{{ imageMeta }}</span>
      </span>
    </div>

    <button
      class="previewer__nav previewer__nav--prev"
      @click="prev"
      @mouseover="hoverNav = true"
      @mouseleave="hoverNav = false"
      :class="{ hidden: !hasPrevious || !showNav }"
      :aria-label="$t('buttons.previous')"
      :title="$t('buttons.previous')"
    >
      <i class="material-icons" aria-hidden="true">chevron_left</i>
    </button>
    <button
      class="previewer__nav previewer__nav--next"
      @click="next"
      @mouseover="hoverNav = true"
      @mouseleave="hoverNav = false"
      :class="{ hidden: !hasNext || !showNav }"
      :aria-label="$t('buttons.next')"
      :title="$t('buttons.next')"
    >
      <i class="material-icons" aria-hidden="true">chevron_right</i>
    </button>
    <link rel="prefetch" :href="previousRaw" />
    <link rel="prefetch" :href="nextRaw" />
  </div>
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import { files as api } from "@/api";
import { createURL } from "@/api/utils";
import { resizePreview } from "@/utils/constants";
import url from "@/utils/url";
import { throttle } from "lodash-es";
import FileHeader from "@/components/files/FileHeader.vue";
import Action from "@/components/header/Action.vue";
import ExtendedImage from "@/components/files/ExtendedImage.vue";
import VideoPlayer from "@/components/files/VideoPlayer.vue";
import CsvViewer from "@/components/files/CsvViewer.vue";
import AudioPlayer from "@/components/files/AudioPlayer.vue";
import MediaLoader from "@/components/files/MediaLoader.vue";
import { downloadWithProgress } from "@/utils/download";
import prettyBytes from "pretty-bytes";
import dayjs from "dayjs";
import { VueReader } from "vue-reader";
import {
  computed,
  inject,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch,
} from "vue";
import { useRoute, useRouter } from "vue-router";
import type { Rendition } from "epubjs";
import { getTheme } from "@/utils/theme";
import { useI18n } from "vue-i18n";

// CSV file size limit for preview (5MB)
// Prevents browser memory issues with large files
const CSV_MAX_SIZE = 5 * 1024 * 1024;

const location = useStorage("book-progress", 0, undefined, {
  serializer: {
    read: (v) => JSON.parse(v),
    write: (v) => JSON.stringify(v),
  },
});
const size = useStorage("book-size", 120, undefined, {
  serializer: {
    read: (v) => JSON.parse(v),
    write: (v) => JSON.stringify(v),
  },
});

const locationChange = (epubcifi: number) => {
  location.value = epubcifi;
};
let rendition: Rendition | null = null;
const changeSize = (val: number) => {
  size.value = val;
  rendition?.themes.fontSize(`${val}%`);
};

const getRendition = (_rendition: Rendition) => {
  rendition = _rendition;
  switch (getTheme()) {
    case "dark": {
      rendition.themes.override("color", "rgba(255, 255, 255, 0.6)");
      break;
    }
    case "light": {
      rendition.themes.override("color", "rgb(111, 111, 111)");
      break;
    }
  }
  rendition.themes.registerRules("h2Transparent", {
    "h1,h2,h3,h4": {
      "background-color": "transparent !important",
    },
  });
  rendition?.themes.fontSize(`${size.value}%`);
  rendition.themes.select("h2Transparent");
  rendition.themes.override("background-color", "transparent", true);
};

const mediaTypes: ResourceType[] = ["image", "video", "audio", "blob"];

const previousLink = ref<string>("");
const nextLink = ref<string>("");
const listing = ref<ResourceItem[] | null>(null);
const name = ref<string>("");
const fullSize = ref<boolean>(false);
const showNav = ref<boolean>(true);
const navTimeout = ref<null | number>(null);
const hoverNav = ref<boolean>(false);
const autoPlay = ref<boolean>(false);
const previousRaw = ref<string>("");
const nextRaw = ref<string>("");
const csvContent = ref<ArrayBuffer | string>("");
const csvError = ref<string>("");

const position = ref<{ n: number; total: number } | null>(null);
const coverUrl = ref<string | null>(null);
const imageRef = ref<InstanceType<typeof ExtendedImage> | null>(null);
const audioRef = ref<InstanceType<typeof AudioPlayer> | null>(null);

/* Images are downloaded with progress, then shown from a blob URL */
const image = reactive({
  state: "idle" as "idle" | "loading" | "ready" | "error",
  src: "",
  loaded: 0,
  total: 0,
  zoom: 1,
  width: 0,
  height: 0,
});
let imageAbort: AbortController | null = null;

const $showError = inject<IToastError>("$showError")!;

const authStore = useAuthStore();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const { t } = useI18n();

const route = useRoute();
const router = useRouter();

const hasPrevious = computed(() => previousLink.value !== "");

const hasNext = computed(() => nextLink.value !== "");

const downloadUrl = computed(() =>
  fileStore.req ? api.getDownloadURL(fileStore.req, false) : ""
);

const directUrl = computed(() =>
  fileStore.req ? api.getDownloadURL(fileStore.req, true) : ""
);

const previewUrl = computed(() => {
  if (!fileStore.req) {
    return "";
  }

  if (fileStore.req.type === "image" && !fullSize.value) {
    return api.getPreviewURL(fileStore.req, "big");
  }

  if (isEpub.value) {
    return createURL("api/raw" + fileStore.req.path, {});
  }

  return api.getDownloadURL(fileStore.req, true);
});

const isPdf = computed(() => fileStore.req?.extension.toLowerCase() == ".pdf");
const isEpub = computed(
  () => fileStore.req?.extension.toLowerCase() == ".epub"
);
const isCsv = computed(
  () =>
    fileStore.req?.extension.toLowerCase() == ".csv" &&
    fileStore.req.size <= CSV_MAX_SIZE
);

const isResizeEnabled = computed(() => resizePreview);

// What kind of viewer to show; drives layout classes too
const kind = computed(() => {
  if (isEpub.value) return "epub";
  if (isCsv.value) return "csv";
  if (isPdf.value) return "pdf";
  const type = fileStore.req?.type;
  if (type === "image" || type === "audio" || type === "video") return type;
  return "other";
});

const folderName = computed(() => {
  const dirs = route.path.replace(/\/+$/, "").split("/");
  const folder = decodeURIComponent(dirs[dirs.length - 2] ?? "");
  return folder === "files" ? "" : folder;
});

// Top bar: type icon and "Folder · 2.1 MiB · 4 of 12 · Modified 2 days ago"
const HEADER_ICONS: Record<string, string> = {
  image: "image",
  video: "movie",
  audio: "audiotrack",
  pdf: "picture_as_pdf",
  epub: "menu_book",
  csv: "table_chart",
};
const headerIcon = computed(
  () => HEADER_ICONS[kind.value] ?? "insert_drive_file"
);
const headerMeta = computed(() => {
  const req = fileStore.req;
  if (!req || layoutStore.loading) return [];
  const parts = [folderName.value || t("sidebar.myFiles")];
  if (!req.isDir) parts.push(prettyBytes(req.size, { binary: true }));
  // The count comes before the date so it survives truncation on phones
  if (position.value) parts.push(t("preview.position", position.value));
  if (req.modified)
    parts.push(t("preview.modified", { time: dayjs(req.modified).fromNow() }));
  return parts;
});

const imageMeta = computed(() => {
  const parts: string[] = [];
  if (image.width && image.height)
    parts.push(`${image.width} × ${image.height}`);
  if (fileStore.req?.size)
    parts.push(prettyBytes(fileStore.req.size, { binary: true }));
  return parts.join(" · ");
});

const subtitles = computed(() => {
  if (fileStore.req?.subtitles) {
    return api.getSubtitlesURL(fileStore.req);
  }
  return [];
});

const videoOptions = computed(() => {
  return { autoplay: autoPlay.value };
});

watch(route, () => {
  updatePreview();
  toggleNavigation();
});

// Specify hooks
onMounted(async () => {
  window.addEventListener("keydown", key);
  listing.value = fileStore.oldReq?.items ?? null;
  updatePreview();
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", key);
  imageAbort?.abort();
  if (image.src) URL.revokeObjectURL(image.src);
});

const loadImage = async () => {
  imageAbort?.abort();
  imageAbort = null;
  if (kind.value !== "image" || layoutStore.loading || !fileStore.req) {
    return;
  }
  const controller = new AbortController();
  imageAbort = controller;
  Object.assign(image, { state: "loading", loaded: 0, total: 0, zoom: 1 });

  try {
    const blobUrl = await downloadWithProgress(previewUrl.value, {
      signal: controller.signal,
      // Thumbnails don't always report a size; originals match the file
      fallbackTotal: fullSize.value ? fileStore.req.size : 0,
      onProgress: (p) => {
        image.loaded = p.loaded;
        image.total = p.total;
      },
    });
    if (controller.signal.aborted) {
      URL.revokeObjectURL(blobUrl);
      return;
    }
    const old = image.src;
    image.src = blobUrl;
    image.state = "ready";
    if (old) URL.revokeObjectURL(old);
  } catch {
    if (!controller.signal.aborted) image.state = "error";
  }
};

const onImageReady = (size: { width: number; height: number }) => {
  image.width = size.width;
  image.height = size.height;
};

watch(
  () => [kind.value === "image" ? previewUrl.value : "", layoutStore.loading],
  () => loadImage(),
  { immediate: true }
);

// Stop touch scrolling of the page, but leave sliders draggable
const blockTouchScroll = (event: TouchEvent) => {
  event.stopPropagation();
  const target = event.target as HTMLElement | null;
  if (target?.closest('input[type="range"]')) return;
  event.preventDefault();
};

// Specify methods
const deleteFile = () => {
  layoutStore.showHover({
    prompt: "delete",
    confirm: () => {
      if (listing.value === null) {
        return;
      }

      const index = listing.value.findIndex((item) => item.name == name.value);
      listing.value.splice(index, 1);

      if (hasNext.value) {
        next();
      } else if (!hasPrevious.value && !hasNext.value) {
        const nearbyItem = listing.value[Math.max(0, index - 1)];
        fileStore.preselect = nearbyItem?.path;

        close();
      } else {
        prev();
      }
    },
  });
};

const prev = () => {
  hoverNav.value = false;
  router.replace({ path: previousLink.value });
};

const next = () => {
  hoverNav.value = false;
  router.replace({ path: nextLink.value });
};

const key = (event: KeyboardEvent) => {
  if (layoutStore.currentPrompt !== null) {
    return;
  }
  // When previewing a video, let arrow keys fall through to video.js for
  // seeking instead of switching to the prev/next file. Enter still advances.
  const isVideo = fileStore.req?.type === "video";
  const target = event.target as HTMLElement | null;
  const onControl = !!target?.closest("button, input, select, textarea, a");

  // Image tools: + − 0 R
  if (kind.value === "image" && image.state === "ready" && !onControl) {
    const actions: Record<string, (() => void) | undefined> = {
      "+": imageRef.value?.zoomIn,
      "=": imageRef.value?.zoomIn,
      "-": imageRef.value?.zoomOut,
      "0": imageRef.value?.reset,
      r: imageRef.value?.rotate,
      R: imageRef.value?.rotate,
    };
    const action = actions[event.key];
    if (action && !event.ctrlKey && !event.metaKey) {
      event.preventDefault();
      action();
      return;
    }
  }

  // Audio: Space / K play or pause
  if (kind.value === "audio" && !onControl) {
    if (event.key === " " || event.key === "k") {
      event.preventDefault();
      audioRef.value?.toggle();
      return;
    }
  }

  if (event.which === 13) {
    // enter
    if (hasNext.value) next();
  } else if (event.which === 39) {
    // right arrow
    if (isVideo) return;
    if (hasNext.value) next();
  } else if (event.which === 37) {
    // left arrow
    if (isVideo) return;
    if (hasPrevious.value) prev();
  } else if (event.which === 27) {
    // esc
    close();
  }
};
const updatePreview = async () => {
  const dirs = route.fullPath.split("/");
  name.value = decodeURIComponent(dirs[dirs.length - 1]);

  // Load CSV content if it's a CSV file
  if (isCsv.value && fileStore.req) {
    csvContent.value = "";
    csvError.value = "";

    if (fileStore.req.size > CSV_MAX_SIZE) {
      csvError.value = t("files.csvTooLarge");
    } else {
      if (fileStore.req.rawContent != null) {
        csvContent.value = fileStore.req.rawContent;
      } else {
        csvContent.value = fileStore.req.content ?? "";
      }
    }
  }

  if (!listing.value) {
    try {
      const path = url.removeLastDir(route.path);
      const res = await api.fetch(path);
      listing.value = res.items;
    } catch (e: any) {
      $showError(e);
    }
  }

  previousLink.value = "";
  nextLink.value = "";
  position.value = null;
  coverUrl.value = null;
  if (listing.value) {
    const media = listing.value.filter((item) =>
      mediaTypes.includes(item.type)
    );
    const index = media.findIndex((item) => item.name === name.value);
    if (index !== -1 && media.length > 1) {
      position.value = { n: index + 1, total: media.length };
    }
    if (fileStore.req?.type === "audio") {
      coverUrl.value = findCover(listing.value, name.value);
    }

    for (let i = 0; i < listing.value.length; i++) {
      if (listing.value[i].name !== name.value) {
        continue;
      }

      for (let j = i - 1; j >= 0; j--) {
        if (mediaTypes.includes(listing.value[j].type)) {
          previousLink.value = listing.value[j].url;
          previousRaw.value = prefetchUrl(listing.value[j]);
          break;
        }
      }
      for (let j = i + 1; j < listing.value.length; j++) {
        if (mediaTypes.includes(listing.value[j].type)) {
          nextLink.value = listing.value[j].url;
          nextRaw.value = prefetchUrl(listing.value[j]);
          break;
        }
      }

      return;
    }
  }
};

// Album art: "cover.jpg", "folder.png", … or an image named like the track
const COVER_NAMES =
  /^(cover|folder|front|album|albumart|artwork)\.(jpe?g|png|webp)$/i;
const findCover = (items: ResourceItem[], track: string) => {
  const base = track.replace(/\.[^/.]+$/, "").toLowerCase();
  const match =
    items.find(
      (i) =>
        i.type === "image" &&
        i.name.replace(/\.[^/.]+$/, "").toLowerCase() === base
    ) ?? items.find((i) => i.type === "image" && COVER_NAMES.test(i.name));
  return match ? api.getPreviewURL(match, "thumb") : null;
};

const prefetchUrl = (item: ResourceItem) => {
  if (item.type !== "image") {
    return "";
  }

  return fullSize.value
    ? api.getDownloadURL(item, true)
    : api.getPreviewURL(item, "big");
};

const toggleSize = () => (fullSize.value = !fullSize.value);

const toggleNavigation = throttle(function () {
  showNav.value = true;

  if (navTimeout.value) {
    clearTimeout(navTimeout.value);
  }

  navTimeout.value = window.setTimeout(() => {
    showNav.value = false || hoverNav.value;
    navTimeout.value = null;
  }, 1500);
}, 500);

const close = () => {
  const uri = url.removeLastDir(route.path) + "/";
  router.push({ path: uri });
};

const download = () => window.open(downloadUrl.value);
const openDirect = () => window.open(directUrl.value);

const editAsText = () => {
  router.push({ path: route.path, query: { edit: "true" } });
};
</script>
