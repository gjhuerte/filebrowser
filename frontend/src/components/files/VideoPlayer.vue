<template>
  <div
    class="video-frame"
    :class="{ 'video-frame--ready': phase !== 'loading' }"
  >
    <video ref="videoPlayer" class="video-max video-js" controls preload="auto">
      <source />
      <track
        kind="subtitles"
        v-for="(sub, index) in subtitles"
        :key="index"
        :src="sub"
        :label="subLabel(sub)"
        :default="index === 0"
      />
      <p class="vjs-no-js">
        {{ t("preview.cantPlayVideo") }}
        <a :href="source">{{ t("buttons.download") }}</a>
      </p>
    </video>

    <MediaLoader
      v-if="phase === 'loading'"
      class="video-frame__loader"
      :label="t('preview.loadingVideo')"
      :loaded="loadedBytes"
      :total="size"
      :progress="bufferedPct || null"
    />
    <MediaLoader
      v-else-if="phase === 'buffering'"
      class="video-frame__chip"
      variant="chip"
      :label="t('preview.buffering')"
      :progress="bufferedPct || null"
    />

    <div v-if="phase === 'error'" class="media-error" role="alert">
      <i class="material-icons" aria-hidden="true">videocam_off</i>
      <p class="media-error__title">{{ t("preview.cantPlayVideo") }}</p>
      <p class="media-error__text">{{ errorText }}</p>
      <button
        v-if="canDownload"
        type="button"
        class="media-error__button"
        @click="emit('download')"
      >
        <i class="material-icons" aria-hidden="true">file_download</i>
        {{ t("buttons.download") }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onBeforeUnmount, nextTick } from "vue";
import { useI18n } from "vue-i18n";
import videojs from "video.js";
import type Player from "video.js/dist/types/player";
import "videojs-mobile-ui";
import "videojs-hotkeys";
import "video.js/dist/video-js.min.css";
import "videojs-mobile-ui/dist/videojs-mobile-ui.css";
import MediaLoader from "@/components/files/MediaLoader.vue";
import { bufferedFraction } from "@/utils/download";
import { FLV_TYPE, registerFlvSupport } from "@/utils/flv";

const videoPlayer = ref<HTMLElement | null>(null);
const player = ref<Player | null>(null);

const props = withDefaults(
  defineProps<{
    source: string;
    subtitles?: string[];
    options?: any;
    /** File size in bytes, used to show how much has loaded. */
    size?: number;
    canDownload?: boolean;
  }>(),
  {
    options: {},
    size: 0,
    canDownload: false,
  }
);
const emit = defineEmits<{ download: [] }>();
const { t } = useI18n();

const source = ref(props.source);
const sourceType = ref("");

/* Loading / buffering state shown on top of the player */
type Phase = "loading" | "ready" | "buffering" | "error";
const phase = ref<Phase>("loading");
const buffered = ref<number | null>(null);
const errorCode = ref(0);

const bufferedPct = computed(() =>
  buffered.value === null ? null : buffered.value * 100
);
const loadedBytes = computed(() =>
  buffered.value === null ? 0 : Math.round(buffered.value * props.size)
);
const errorText = computed(() =>
  errorCode.value === 2
    ? t("preview.networkError")
    : t("preview.unsupportedVideo")
);

nextTick(() => {
  initVideoPlayer();
});

onBeforeUnmount(() => {
  if (player.value) {
    player.value.dispose();
    player.value = null;
  }
});

const initVideoPlayer = async () => {
  try {
    registerFlvSupport(t("preview.unsupportedVideo"));

    const lang = document.documentElement.lang;
    const languagePack = await (
      languageImports[lang] || languageImports.en
    )?.();
    const code = languageImports[lang] ? lang : "en";
    videojs.addLanguage(code, languagePack.default);
    sourceType.value = getSourceType(source.value);

    const srcOpt = { sources: { src: props.source, type: sourceType.value } };
    //Supporting localized language display.
    const langOpt = { language: code };
    // support for playback at different speeds.
    const playbackRatesOpt = { playbackRates: [0.5, 1, 1.25, 1.5, 2, 3] };
    const options = getOptions(
      props.options,
      langOpt,
      srcOpt,
      playbackRatesOpt
    );
    const p = videojs(videoPlayer.value!, options, () => {});
    player.value = p;
    watchState(p);

    // @ts-expect-error no ts definition for mobileUi
    player.value!.mobileUi();
  } catch (error) {
    console.error("Error initializing video player:", error);
  }
};

const watchState = (p: Player) => {
  const updateBuffered = () => {
    const el = p.el()?.querySelector("video");
    if (el) buffered.value = bufferedFraction(el);
  };
  p.on("loadstart", () => {
    if (phase.value !== "error") phase.value = "loading";
  });
  p.on(["progress", "durationchange", "timeupdate"], updateBuffered);
  p.on(["canplay", "playing", "seeked"], () => {
    if (phase.value !== "error") phase.value = "ready";
    updateBuffered();
  });
  // Paused videos don't buffer forever, so a paused 'waiting' isn't shown
  p.on("waiting", () => {
    if (phase.value === "ready" && !p.paused()) phase.value = "buffering";
  });
  p.on("error", () => {
    errorCode.value = p.error()?.code ?? 0;
    phase.value = "error";
  });
};

const getOptions = (...srcOpt: any[]) => {
  const options = {
    controlBar: {
      skipButtons: {
        forward: 10,
        backward: 10,
      },
    },
    html5: {
      nativeTextTracks: false,
    },
    plugins: {
      hotkeys: {
        volumeStep: 0.1,
        seekStep: 10,
        enableModifiersForNumbers: false,
      },
    },
  };

  return videojs.obj.merge(options, ...srcOpt);
};

const getSourceType = (source: string) => {
  const ext = (source ? source.split("?")[0].split(".").pop() : "") ?? "";
  switch (ext.toLowerCase()) {
    // Most MKV files hold H.264/AAC which browsers decode fine as MP4
    case "mkv":
      return "video/mp4";
    // Flash video is remuxed in the browser (see utils/flv.ts)
    case "flv":
      return FLV_TYPE;
    default:
      return "";
  }
};

const subLabel = (subUrl: string) => {
  let url: URL;
  try {
    url = new URL(subUrl);
  } catch {
    // treat it as a relative url
    // we only need this for filename
    url = new URL(subUrl, window.location.origin);
  }

  const label = decodeURIComponent(
    url.pathname
      .split("/")
      .pop()!
      .replace(/\.[^/.]+$/, "")
  );

  return label;
};

interface LanguageImports {
  [key: string]: () => Promise<any>;
}

const languageImports: LanguageImports = {
  ar: () => import("video.js/dist/lang/ar.json"),
  bg: () => import("video.js/dist/lang/bg.json"),
  cs: () => import("video.js/dist/lang/cs.json"),
  de: () => import("video.js/dist/lang/de.json"),
  el: () => import("video.js/dist/lang/el.json"),
  en: () => import("video.js/dist/lang/en.json"),
  es: () => import("video.js/dist/lang/es.json"),
  fr: () => import("video.js/dist/lang/fr.json"),
  he: () => import("video.js/dist/lang/he.json"),
  hr: () => import("video.js/dist/lang/hr.json"),
  hu: () => import("video.js/dist/lang/hu.json"),
  it: () => import("video.js/dist/lang/it.json"),
  ja: () => import("video.js/dist/lang/ja.json"),
  ko: () => import("video.js/dist/lang/ko.json"),
  lv: () => import("video.js/dist/lang/lv.json"),
  nb: () => import("video.js/dist/lang/nb.json"),
  nl: () => import("video.js/dist/lang/nl.json"),
  "nl-be": () => import("video.js/dist/lang/nl.json"),
  pl: () => import("video.js/dist/lang/pl.json"),
  "pt-br": () => import("video.js/dist/lang/pt-BR.json"),
  "pt-pt": () => import("video.js/dist/lang/pt-PT.json"),
  ro: () => import("video.js/dist/lang/ro.json"),
  ru: () => import("video.js/dist/lang/ru.json"),
  sk: () => import("video.js/dist/lang/sk.json"),
  tr: () => import("video.js/dist/lang/tr.json"),
  uk: () => import("video.js/dist/lang/uk.json"),
  vi: () => import("video.js/dist/lang/vi.json"),
  "zh-cn": () => import("video.js/dist/lang/zh-CN.json"),
  "zh-tw": () => import("video.js/dist/lang/zh-TW.json"),
};
</script>
