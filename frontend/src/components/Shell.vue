<template>
  <section
    class="shell"
    :class="{ 'shell--hidden': !showShell, 'shell--dragging': dragging }"
    :style="{ height: `${height}em` }"
    :aria-hidden="!showShell"
    :aria-label="$t('shell.title')"
    dir="ltr"
    @keydown.esc.stop="hide"
  >
    <!-- Drag handle: resize the panel -->
    <div
      class="shell__grip"
      role="separator"
      aria-orientation="horizontal"
      :aria-label="$t('shell.resize')"
      @pointerdown.prevent="startDrag"
    ></div>

    <div class="shell__bar">
      <div class="shell__tabs" role="tablist" :aria-label="$t('shell.tabs')">
        <div
          v-for="(tab, i) in tabs"
          :key="tab.id"
          class="shell__tab"
          :class="{ 'shell__tab--active': tab.id === activeId }"
        >
          <button
            type="button"
            role="tab"
            class="shell__tab-label"
            :aria-selected="tab.id === activeId"
            @click="selectTab(tab.id)"
            @auxclick.middle="closeTab(tab.id)"
          >
            <span
              class="shell__status"
              :class="{ 'shell__status--busy': tab.running }"
              aria-hidden="true"
            ></span>
            {{ $t("shell.tabName", { n: i + 1 }) }}
          </button>
          <button
            type="button"
            class="shell__tab-close"
            :aria-label="$t('shell.closeTab', { n: i + 1 })"
            :title="$t('shell.closeTab', { n: i + 1 })"
            @click="closeTab(tab.id)"
          >
            <i class="material-icons" aria-hidden="true">close</i>
          </button>
        </div>
        <button
          type="button"
          class="shell__icon"
          :aria-label="$t('shell.newTab')"
          :title="$t('shell.newTab')"
          @click="newTab"
        >
          <i class="material-icons" aria-hidden="true">add</i>
        </button>
      </div>

      <div class="shell__actions">
        <span class="shell__cwd" :title="$t('shell.runsIn', { path: folder })">
          <i class="material-icons" aria-hidden="true">folder_open</i>
          <span>{{ folder }}</span>
        </span>
        <button
          type="button"
          class="shell__icon"
          :aria-label="$t('shell.clear')"
          :title="$t('shell.clear')"
          @click="clear"
        >
          <i class="material-icons" aria-hidden="true">clear_all</i>
        </button>
        <button
          type="button"
          class="shell__icon"
          :aria-label="$t('shell.hide')"
          :title="$t('shell.hide')"
          @click="hide"
        >
          <i class="material-icons" aria-hidden="true">keyboard_arrow_down</i>
        </button>
      </div>
    </div>

    <div class="shell__content" ref="scrollable" @click="focusInput">
      <p v-if="tab.content.length === 0" class="shell__intro">
        {{ $t("shell.intro") }}
      </p>

      <div v-for="(c, index) in tab.content" :key="index" class="shell__block">
        <div class="shell__line">
          <span class="shell__prompt" aria-hidden="true">❯</span>
          <span class="shell__cmd">{{ c.cmd }}</span>
          <span class="shell__where">{{ c.folder }}</span>
        </div>
        <pre v-if="c.output" class="shell__output">{{ c.output }}</pre>
        <p v-if="c.running" class="shell__running">
          <span class="shell__status shell__status--busy"></span>
          {{ $t("shell.running") }}
        </p>
      </div>

      <div
        class="shell__line shell__line--input"
        :class="{ 'shell__line--busy': tab.running }"
      >
        <span class="shell__prompt" aria-hidden="true">❯</span>
        <pre
          ref="input"
          class="shell__input"
          role="textbox"
          tabindex="0"
          :contenteditable="!tab.running"
          :aria-label="$t('shell.inputLabel')"
          :data-placeholder="
            tab.running ? $t('shell.waiting') : $t('shell.placeholder')
          "
          spellcheck="false"
          @keydown.up.prevent="historyUp"
          @keydown.down.prevent="historyDown"
          @keydown.enter.prevent="submit"
          @input="onInput"
        />
      </div>
    </div>

    <div class="shell__overlay" v-show="dragging"></div>
  </section>
</template>

<script>
import { markRaw } from "vue";
import { mapState, mapActions } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import { commands } from "@/api";

const MIN_HEIGHT = 10; // em
const ANSI = /\u001b\[[0-9;]*[A-Za-z]/g;

let nextId = 1;
const makeTab = () => ({
  id: nextId++,
  content: [],
  history: [],
  historyPos: 0,
  draft: "",
  running: false,
  conn: null,
});

export default {
  name: "shell",
  data: () => {
    const first = makeTab();
    return {
      tabs: [first],
      activeId: first.id,
      dragging: false,
      height: 22,
      fontsize: parseFloat(getComputedStyle(document.documentElement).fontSize),
    };
  },
  computed: {
    ...mapState(useLayoutStore, ["showShell"]),
    ...mapState(useFileStore, ["isFiles"]),
    tab() {
      return this.tabs.find((t) => t.id === this.activeId) ?? this.tabs[0];
    },
    // Commands run in the folder you're looking at
    path() {
      return this.isFiles ? this.$route.path : "";
    },
    folder() {
      let p = this.path.replace(/^\/files/, "");
      try {
        p = decodeURIComponent(p);
      } catch {
        // keep as-is
      }
      p = p.replace(/\/+$/, "");
      return p === "" ? "/" : p;
    },
  },
  watch: {
    showShell(open) {
      if (open) this.$nextTick(this.focusInput);
      this.reserveSpace();
    },
    height() {
      this.reserveSpace();
    },
  },
  mounted() {
    window.addEventListener("resize", this.clampHeight);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.clampHeight);
    this.stopDrag();
    this.tabs.forEach((t) => t.conn?.close());
    document.body.style.paddingBottom = "";
  },
  methods: {
    ...mapActions(useLayoutStore, ["toggleShell"]),

    /* Panel */
    // Leave room under the page so the last rows can scroll above the panel
    reserveSpace() {
      document.body.style.paddingBottom = this.showShell
        ? `${this.height + 1}em`
        : "";
    },
    hide() {
      if (this.showShell) this.toggleShell();
    },
    maxHeight() {
      return window.innerHeight / this.fontsize - 6;
    },
    clampHeight() {
      this.height = Math.min(
        Math.max(this.height, MIN_HEIGHT),
        Math.max(this.maxHeight(), MIN_HEIGHT)
      );
    },
    startDrag() {
      this.dragging = true;
      document.addEventListener("pointermove", this.onDrag);
      document.addEventListener("pointerup", this.stopDrag);
    },
    onDrag(event) {
      this.height = (window.innerHeight - event.clientY) / this.fontsize;
      this.clampHeight();
    },
    stopDrag() {
      this.dragging = false;
      document.removeEventListener("pointermove", this.onDrag);
      document.removeEventListener("pointerup", this.stopDrag);
    },

    /* Tabs */
    newTab() {
      const t = makeTab();
      this.tabs.push(t);
      this.selectTab(t.id);
    },
    selectTab(id) {
      this.saveDraft();
      this.activeId = id;
      this.$nextTick(() => {
        this.$refs.input.innerText = this.tab.draft;
        this.scroll();
        this.focusInput();
      });
    },
    closeTab(id) {
      const i = this.tabs.findIndex((t) => t.id === id);
      if (i === -1) return;
      // The server finishes a running command on its own; we just stop listening
      this.tabs[i].conn?.close();
      this.tabs.splice(i, 1);
      if (this.tabs.length === 0) {
        // Closing the last tab hides the panel and leaves a fresh one
        this.tabs.push(makeTab());
        this.activeId = this.tabs[0].id;
        this.$refs.input.innerText = "";
        this.hide();
        return;
      }
      if (id === this.activeId) {
        this.activeId = this.tabs[Math.max(0, i - 1)].id;
      }
      this.$nextTick(() => {
        this.$refs.input.innerText = this.tab.draft;
        this.scroll();
        this.focusInput();
      });
    },
    saveDraft() {
      if (this.$refs.input) this.tab.draft = this.$refs.input.innerText;
    },

    /* Input */
    focusInput() {
      if (window.getSelection()?.toString()) return; // let people copy output
      this.$refs.input?.focus();
    },
    onInput() {
      this.tab.draft = this.$refs.input.innerText;
    },
    scroll() {
      const el = this.$refs.scrollable;
      if (el) el.scrollTop = el.scrollHeight;
    },
    setInput(text) {
      this.$refs.input.innerText = text;
      this.tab.draft = text;
    },
    historyUp() {
      const t = this.tab;
      if (t.historyPos > 0) this.setInput(t.history[--t.historyPos]);
    },
    historyDown() {
      const t = this.tab;
      if (t.historyPos < t.history.length - 1) {
        this.setInput(t.history[++t.historyPos]);
      } else {
        t.historyPos = t.history.length;
        this.setInput("");
      }
    },
    clear() {
      this.tab.content = [];
      this.focusInput();
    },
    submit() {
      const t = this.tab;
      if (t.running) return;
      const cmd = this.$refs.input.innerText.trim();
      if (cmd === "") return;

      this.setInput("");
      if (cmd === "clear") return this.clear();
      if (cmd === "exit") return this.hide();

      t.history.push(cmd);
      t.historyPos = t.history.length;
      t.content.push({ cmd, folder: this.folder, output: "", running: true });
      const block = t.content[t.content.length - 1];
      t.running = true;
      this.$nextTick(this.scroll);

      t.conn = markRaw(
        commands(
          this.path,
          cmd,
          (event) => {
            block.output += `${String(event.data).replace(ANSI, "")}\n`;
            if (t.id === this.activeId) this.$nextTick(this.scroll);
          },
          () => {
            block.output = block.output.trimEnd();
            block.running = false;
            t.running = false;
            t.conn = null;
            if (t.id === this.activeId && this.showShell) {
              this.$nextTick(() => {
                this.scroll();
                this.focusInput();
              });
            }
          }
        )
      );
    },
  },
};
</script>
