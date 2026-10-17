<template>
  <div class="audio-player" :class="{ 'is-playing': playing }">
    <div class="audio-player__art">
      <img
        v-if="cover && !coverFailed"
        :src="cover"
        alt=""
        @error="coverFailed = true"
      />
      <div v-else class="audio-player__art-fallback" aria-hidden="true">
        <span class="audio-player__ring"></span>
        <span class="audio-player__ring audio-player__ring--2"></span>
        <i class="material-icons">music_note</i>
      </div>

      <MediaLoader
        v-if="phase === 'loading'"
        class="audio-player__loader"
        :label="t('preview.loadingAudio')"
        :loaded="loadedBytes"
        :total="size"
        :progress="bufferedPct || null"
      />
      <MediaLoader
        v-else-if="phase === 'buffering'"
        class="audio-player__chip"
        variant="chip"
        :label="t('preview.buffering')"
        :progress="bufferedPct || null"
      />
    </div>

    <div class="audio-player__meta">
      <p class="audio-player__title" :title="name">{{ title }}</p>
      <p class="audio-player__sub">
        <span v-if="folder">{{ folder }}</span>
        <span v-if="ext" class="audio-player__badge">{{ ext }}</span>
        <span v-if="size">{{ bytes(size) }}</span>
      </p>
    </div>

    <div v-if="phase === 'error'" class="media-error media-error--inline">
      <i class="material-icons" aria-hidden="true">music_off</i>
      <p class="media-error__title">{{ t("preview.cantPlayAudio") }}</p>
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

    <template v-else>
      <div class="audio-player__seek">
        <input
          type="range"
          class="audio-player__range"
          min="0"
          :max="duration || 0"
          step="0.1"
          :value="current"
          :disabled="!duration"
          :style="{
            '--played': playedPct + '%',
            '--buffered': (bufferedPct ?? 0) + '%',
          }"
          :aria-label="t('preview.seek')"
          :aria-valuetext="`${clock(current)} / ${clock(duration)}`"
          @input="onSeek"
        />
        <div class="audio-player__times">
          <span>{{ clock(current) }}</span>
          <span v-if="bufferText" class="audio-player__buffer">{{
            bufferText
          }}</span>
          <span>{{
            duration ? "-" + clock(duration - current) : "--:--"
          }}</span>
        </div>
      </div>

      <div class="audio-player__controls">
        <button
          type="button"
          class="audio-player__btn audio-player__btn--small"
          :class="{ 'is-active': loop }"
          :aria-pressed="loop"
          :aria-label="t('preview.repeat')"
          :title="t('preview.repeat')"
          @click="loop = !loop"
        >
          <i class="material-icons" aria-hidden="true">repeat_one</i>
        </button>
        <button
          type="button"
          class="audio-player__btn"
          :disabled="!hasPrevious"
          :aria-label="t('buttons.previous')"
          :title="t('buttons.previous')"
          @click="emit('previous')"
        >
          <i class="material-icons" aria-hidden="true">skip_previous</i>
        </button>
        <button
          type="button"
          class="audio-player__btn"
          :aria-label="t('preview.back10')"
          :title="t('preview.back10')"
          @click="skip(-10)"
        >
          <i class="material-icons" aria-hidden="true">replay_10</i>
        </button>
        <button
          type="button"
          class="audio-player__play"
          :aria-label="playing ? t('preview.pause') : t('preview.play')"
          :title="playing ? t('preview.pause') : t('preview.play')"
          @click="toggle"
        >
          <i class="material-icons" aria-hidden="true">{{
            playing ? "pause" : "play_arrow"
          }}</i>
        </button>
        <button
          type="button"
          class="audio-player__btn"
          :aria-label="t('preview.forward10')"
          :title="t('preview.forward10')"
          @click="skip(10)"
        >
          <i class="material-icons" aria-hidden="true">forward_10</i>
        </button>
        <button
          type="button"
          class="audio-player__btn"
          :disabled="!hasNext"
          :aria-label="t('buttons.next')"
          :title="t('buttons.next')"
          @click="emit('next')"
        >
          <i class="material-icons" aria-hidden="true">skip_next</i>
        </button>
        <button
          type="button"
          class="audio-player__btn audio-player__btn--small audio-player__speed"
          :aria-label="t('preview.speed')"
          :title="t('preview.speed')"
          @click="cycleSpeed"
        >
          {{ rate }}×
        </button>
      </div>

      <div class="audio-player__volume">
        <button
          type="button"
          class="audio-player__btn audio-player__btn--small"
          :aria-pressed="muted"
          :aria-label="muted ? t('preview.unmute') : t('preview.mute')"
          :title="muted ? t('preview.unmute') : t('preview.mute')"
          @click="muted = !muted"
        >
          <i class="material-icons" aria-hidden="true">{{ volumeIcon }}</i>
        </button>
        <input
          type="range"
          class="audio-player__range audio-player__range--volume"
          min="0"
          max="1"
          step="0.05"
          :value="muted ? 0 : volume"
          :style="{
            '--played': (muted ? 0 : volume * 100) + '%',
            '--buffered': '0%',
          }"
          :aria-label="t('preview.volume')"
          @input="onVolume"
        />
      </div>

      <p class="audio-player__hint">{{ t("preview.audioHint") }}</p>
    </template>

    <audio
      ref="audio"
      :src="src"
      preload="auto"
      :autoplay="autoplay"
      :loop="loop"
      @loadstart="onLoadStart"
      @progress="updateBuffered"
      @durationchange="onDuration"
      @loadedmetadata="onDuration"
      @timeupdate="onTime"
      @canplay="onReady"
      @playing="onReady"
      @waiting="onWaiting"
      @play="onPlay"
      @pause="onPause"
      @ended="onEnded"
      @error="onError"
    ></audio>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useStorage } from "@vueuse/core";
import prettyBytes from "pretty-bytes";
import MediaLoader from "@/components/files/MediaLoader.vue";
import { bufferedFraction } from "@/utils/download";

const props = withDefaults(
  defineProps<{
    src: string;
    name: string;
    size?: number;
    folder?: string;
    cover?: string | null;
    autoplay?: boolean;
    hasPrevious?: boolean;
    hasNext?: boolean;
    canDownload?: boolean;
  }>(),
  {
    size: 0,
    folder: "",
    cover: null,
    autoplay: false,
    hasPrevious: false,
    hasNext: false,
    canDownload: false,
  }
);
const emit = defineEmits<{
  previous: [];
  next: [];
  download: [];
  play: [];
  pause: [];
}>();
const { t } = useI18n();

const audio = ref<HTMLAudioElement | null>(null);
type Phase = "loading" | "ready" | "buffering" | "error";
const phase = ref<Phase>("loading");
const playing = ref(false);
const current = ref(0);
const duration = ref(0);
const buffered = ref<number | null>(null);
const errorCode = ref(0);
const coverFailed = ref(false);

// Remembered between tracks and visits
const volume = useStorage("audio-volume", 1);
const muted = useStorage("audio-muted", false);
const rate = useStorage("audio-rate", 1);
const loop = ref(false);
const RATES = [0.75, 1, 1.25, 1.5, 2];

const title = computed(() => props.name.replace(/\.[^/.]+$/, ""));
const ext = computed(() => {
  const dot = props.name.lastIndexOf(".");
  return dot > 0 ? props.name.slice(dot + 1).toUpperCase() : "";
});
const playedPct = computed(() =>
  duration.value ? (current.value / duration.value) * 100 : 0
);
const bufferedPct = computed(() =>
  buffered.value === null ? null : buffered.value * 100
);
const loadedBytes = computed(() =>
  buffered.value === null ? 0 : Math.round(buffered.value * props.size)
);
// "Loaded 40% · 1.2 MiB of 3 MiB" until the whole file is buffered
const bufferText = computed(() => {
  if (buffered.value === null || buffered.value >= 0.999) return "";
  const pct = Math.round(buffered.value * 100);
  return props.size
    ? t("preview.loadedOf", {
        percent: pct,
        loaded: bytes(loadedBytes.value),
        total: bytes(props.size),
      })
    : t("preview.loadedPercent", { percent: pct });
});
const volumeIcon = computed(() =>
  muted.value || volume.value === 0
    ? "volume_off"
    : volume.value < 0.5
      ? "volume_down"
      : "volume_up"
);
const errorText = computed(() =>
  errorCode.value === 2
    ? t("preview.networkError")
    : t("preview.unsupportedAudio")
);

const bytes = (n: number) => prettyBytes(n, { binary: true });
const clock = (s: number) => {
  if (!s || !isFinite(s)) return "0:00";
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = Math.floor(s % 60)
    .toString()
    .padStart(2, "0");
  return h ? `${h}:${m.toString().padStart(2, "0")}:${sec}` : `${m}:${sec}`;
};

const applySettings = () => {
  if (!audio.value) return;
  audio.value.volume = volume.value;
  audio.value.muted = muted.value;
  audio.value.playbackRate = rate.value;
};
watch([volume, muted, rate], applySettings);

const onLoadStart = () => {
  phase.value = "loading";
  current.value = 0;
  buffered.value = null;
  coverFailed.value = false;
  applySettings();
};
const updateBuffered = () => {
  if (audio.value) buffered.value = bufferedFraction(audio.value);
};
const onDuration = () => {
  duration.value = audio.value?.duration || 0;
  updateBuffered();
};
const onTime = () => {
  current.value = audio.value?.currentTime || 0;
  updateBuffered();
};
const onReady = () => {
  if (phase.value !== "error") phase.value = "ready";
  updateBuffered();
};
const onWaiting = () => {
  if (phase.value === "ready" && playing.value) phase.value = "buffering";
};
const onPlay = () => {
  playing.value = true;
  emit("play");
};
const onPause = () => {
  playing.value = false;
  // Pausing at the end of a track isn't a "stop"; keep autoplay going
  if (!audio.value?.ended) emit("pause");
};
const onEnded = () => {
  playing.value = false;
  // Continue with the next track in the folder
  if (!loop.value && props.hasNext) emit("next");
};
const onError = () => {
  errorCode.value = audio.value?.error?.code ?? 0;
  phase.value = "error";
  playing.value = false;
};

const toggle = () => {
  const el = audio.value;
  if (!el) return;
  if (el.paused) el.play().catch(() => {});
  else el.pause();
};
const skip = (s: number) => {
  const el = audio.value;
  if (!el || !duration.value) return;
  el.currentTime = Math.min(Math.max(0, el.currentTime + s), duration.value);
};
const onSeek = (e: Event) => {
  const v = Number((e.target as HTMLInputElement).value);
  if (audio.value) audio.value.currentTime = v;
  current.value = v;
};
const onVolume = (e: Event) => {
  volume.value = Number((e.target as HTMLInputElement).value);
  muted.value = volume.value === 0;
};
const cycleSpeed = () => {
  const i = RATES.indexOf(rate.value);
  rate.value = RATES[(i + 1) % RATES.length];
};

/* Lock-screen / media-key controls where the browser supports them */
const setupMediaSession = () => {
  if (!("mediaSession" in navigator)) return;
  navigator.mediaSession.metadata = new MediaMetadata({
    title: title.value,
    album: props.folder,
    artwork: props.cover ? [{ src: props.cover }] : [],
  });
  const handlers: [MediaSessionAction, (() => void) | null][] = [
    ["play", toggle],
    ["pause", toggle],
    ["previoustrack", props.hasPrevious ? () => emit("previous") : null],
    ["nexttrack", props.hasNext ? () => emit("next") : null],
    ["seekbackward", () => skip(-10)],
    ["seekforward", () => skip(10)],
  ];
  for (const [action, fn] of handlers) {
    try {
      navigator.mediaSession.setActionHandler(action, fn);
    } catch {
      /* action not supported */
    }
  }
};

onMounted(setupMediaSession);
watch(() => [props.src, props.hasNext, props.hasPrevious], setupMediaSession);
onBeforeUnmount(() => {
  if ("mediaSession" in navigator) navigator.mediaSession.metadata = null;
});

defineExpose({ toggle, skip });
</script>
