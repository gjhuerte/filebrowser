<template>
  <div class="card floating" id="info-card">
    <div class="card-content">
      <!-- Multi-select summary -->
      <p v-if="selected.length > 1" class="info-summary">
        <i class="material-icons info-icon">folder_copy</i>
        <span>
          {{ $t("prompts.filesSelected", { count: selected.length }) }}
          &mdash; {{ humanSize }}
        </span>
      </p>

      <!-- Single file/folder view -->
      <template v-else>
        <!-- File header -->
        <div class="info-header">
          <i class="material-icons info-type-icon">{{ typeIcon }}</i>
          <div class="info-header-text">
            <p class="info-name break-word">{{ name }}</p>
            <p class="info-meta">{{ typeLabel }}{{ extLabel }}</p>
          </div>
        </div>

        <!-- Basic info -->
        <div class="info-grid">
          <template v-if="!dir">
            <span class="info-label">{{ $t("prompts.size") }}</span>
            <span class="info-value">{{ humanSize }} <small class="info-bytes">({{ rawSize }} bytes)</small></span>
          </template>

          <template v-if="dir">
            <span class="info-label">{{ $t("prompts.items") }}</span>
            <span class="info-value">{{ req.numFiles }} {{ $t("prompts.numberFiles").toLowerCase() }}, {{ req.numDirs }} {{ $t("prompts.numberDirs").toLowerCase() }}</span>

            <span class="info-label">{{ $t("prompts.totalSize") }}</span>
            <span class="info-value">{{ humanSize }}</span>
          </template>

          <span class="info-label">{{ $t("prompts.lastModified") }}</span>
          <span class="info-value" :title="modTimeISO">{{ modTime }}</span>

          <template v-if="resolution">
            <span class="info-label">{{ $t("prompts.resolution") }}</span>
            <span class="info-value">{{ resolution.width }} &times; {{ resolution.height }} px</span>
          </template>

          <span class="info-label">{{ $t("prompts.path") }}</span>
          <span class="info-value info-path break-word">{{ filePath }}</span>
        </div>

        <!-- Advanced toggle (files only) -->
        <template v-if="!dir">
          <button type="button" class="info-toggle" @click="showAdvanced = !showAdvanced">
            <i class="material-icons">{{ showAdvanced ? 'expand_less' : 'expand_more' }}</i>
            <span>{{ $t("prompts.technicalInfo") }}</span>
          </button>

          <div v-if="showAdvanced" class="info-advanced">
            <div class="info-grid">
              <span class="info-label">{{ $t("prompts.permissions") }}</span>
              <span class="info-value"><code class="info-code">{{ permString }}</code></span>

              <template v-if="extension">
                <span class="info-label">{{ $t("prompts.extension") }}</span>
                <span class="info-value">{{ extension }}</span>
              </template>

              <template v-if="isSymlink">
                <span class="info-label">{{ $t("prompts.symlink") }}</span>
                <span class="info-value">Yes</span>
              </template>
            </div>

            <!-- Checksums -->
            <p class="info-sub-title">{{ $t("prompts.checksums") }}</p>
            <div class="info-grid">
              <template v-for="algo in ['md5', 'sha1', 'sha256', 'sha512']" :key="algo">
                <span class="info-label">{{ algo.toUpperCase() }}</span>
                <span class="info-value">
                  <code class="info-code info-hash">
                    <a
                      v-if="!checksums[algo]"
                      @click="checksum($event, algo)"
                      @keypress.enter="checksum($event, algo)"
                      tabindex="0"
                      class="info-hash-link"
                    >{{ $t("prompts.show") }}</a>
                    <span v-else class="info-hash-value" @click="copyHash(algo)">{{ checksums[algo] }}</span>
                  </code>
                </span>
              </template>
            </div>
          </div>
        </template>
      </template>
    </div>

    <div class="card-action">
      <button
        id="focus-prompt"
        type="submit"
        @click="closeHovers"
        class="button button--flat"
        :aria-label="$t('buttons.ok')"
        :title="$t('buttons.ok')"
      >
        {{ $t("buttons.ok") }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { filesize } from "@/utils";
import dayjs from "dayjs";
import { files as api } from "@/api";

const TYPE_ICONS = {
  dir: "folder",
  video: "movie",
  audio: "audiotrack",
  image: "image",
  pdf: "picture_as_pdf",
  text: "description",
  textImmutable: "lock",
  blob: "insert_drive_file",
};

const TYPE_LABELS = {
  dir: "Folder",
  video: "Video",
  audio: "Audio",
  image: "Image",
  pdf: "PDF Document",
  text: "Text",
  textImmutable: "Text (read-only)",
  blob: "File",
};

export default {
  name: "info",
  inject: ["$showError"],
  data() {
    return {
      checksums: {},
      showAdvanced: false,
    };
  },
  computed: {
    ...mapState(useFileStore, [
      "req",
      "selected",
      "selectedCount",
      "isListing",
    ]),
    item() {
      if (this.selectedCount === 0) return this.req;
      return this.req.items[this.selected[0]];
    },
    humanSize() {
      if (this.selectedCount === 0 || !this.isListing) {
        return filesize(this.req.size);
      }
      let sum = 0;
      for (const s of this.selected) {
        sum += this.req.items[s].size;
      }
      return filesize(sum);
    },
    rawSize() {
      if (this.selectedCount === 0) return this.req.size.toLocaleString();
      if (this.selectedCount === 1) return this.item.size.toLocaleString();
      let sum = 0;
      for (const s of this.selected) sum += this.req.items[s].size;
      return sum.toLocaleString();
    },
    modTime() {
      const mod = this.selectedCount === 0 ? this.req.modified : this.item.modified;
      return new Date(Date.parse(mod)).toLocaleString();
    },
    modTimeISO() {
      const mod = this.selectedCount === 0 ? this.req.modified : this.item.modified;
      return mod;
    },
    name() {
      return this.selectedCount === 0 ? this.req.name : this.item.name;
    },
    dir() {
      return this.selectedCount > 1 || (this.selectedCount === 0 ? this.req.isDir : this.item.isDir);
    },
    filePath() {
      return this.selectedCount === 0 ? this.req.path : this.item.path;
    },
    extension() {
      const ext = this.selectedCount === 0 ? this.req.extension : this.item.extension;
      return ext || "";
    },
    isSymlink() {
      const item = this.selectedCount === 0 ? this.req : this.item;
      return item.isSymlink;
    },
    typeIcon() {
      const t = this.selectedCount === 0 ? this.req.type : this.item.type;
      return TYPE_ICONS[t] || "insert_drive_file";
    },
    typeLabel() {
      const t = this.selectedCount === 0 ? this.req.type : this.item.type;
      return TYPE_LABELS[t] || "File";
    },
    extLabel() {
      if (this.dir || !this.extension) return "";
      return " · " + this.extension.toUpperCase();
    },
    permString() {
      const mode = this.selectedCount === 0 ? this.req.mode : this.item.mode;
      return this.formatMode(mode);
    },
    resolution() {
      if (this.selectedCount === 1) {
        const selectedItem = this.req.items[this.selected[0]];
        if (selectedItem && selectedItem.type === "image") return selectedItem.resolution;
      } else if (this.req && this.req.type === "image") {
        return this.req.resolution;
      }
      return null;
    },
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    formatMode(mode) {
      const perms = mode & 0o777;
      const octal = "0" + perms.toString(8);
      const chars = ["r", "w", "x"];
      let str = "";
      for (let i = 2; i >= 0; i--) {
        const bits = (perms >> (i * 3)) & 7;
        for (let j = 2; j >= 0; j--) {
          str += bits & (1 << j) ? chars[2 - j] : "-";
        }
      }
      return octal + " (" + str + ")";
    },
    async checksum(event, algo) {
      event.preventDefault();
      let link;
      if (this.selectedCount) {
        link = this.req.items[this.selected[0]].url;
      } else {
        link = this.$route.path;
      }
      try {
        const hash = await api.checksum(link, algo);
        this.checksums = { ...this.checksums, [algo]: hash };
      } catch (e) {
        this.$showError(e);
      }
    },
    async copyHash(algo) {
      try {
        await navigator.clipboard.writeText(this.checksums[algo]);
      } catch {
        // fallback: select text
      }
    },
  },
};
</script>

<style scoped>
#info-card {
  max-width: 28em;
}

#info-card .card-content {
  padding: 1.25em 1.5em 0.75em;
}

.info-summary {
  display: flex;
  align-items: center;
  gap: 0.5em;
  font-size: 1rem;
}

.info-icon {
  color: var(--primary);
  font-size: 1.4em;
}

/* Header */
.info-header {
  display: flex;
  align-items: flex-start;
  gap: 0.75em;
  margin-bottom: 1em;
}

.info-type-icon {
  font-size: 2em;
  color: var(--primary);
  flex-shrink: 0;
  margin-top: 0.05em;
}

.info-header-text {
  min-width: 0;
}

.info-name {
  font-weight: 600;
  font-size: 1.05rem;
  margin: 0;
  line-height: 1.3;
  color: var(--textSecondary);
}

.info-meta {
  margin: 0.15em 0 0;
  font-size: 0.82rem;
  color: var(--textPrimary);
}

/* Key-value grid */
.info-grid {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 0.35em 1em;
  font-size: 0.88rem;
  align-items: baseline;
}

.info-label {
  color: var(--textPrimary);
  white-space: nowrap;
}

.info-value {
  color: var(--textSecondary);
  word-break: break-word;
}

.info-path {
  font-size: 0.82em;
  font-family: var(--font-mono, monospace);
}

.info-bytes {
  color: var(--textPrimary);
  font-size: 0.85em;
}

/* Advanced toggle */
.info-toggle {
  display: flex;
  align-items: center;
  gap: 0.25em;
  margin: 1em 0 0;
  padding: 0.4em 0;
  background: none;
  border: none;
  border-top: 1px solid var(--borderPrimary);
  width: 100%;
  cursor: pointer;
  font-size: 0.78rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--textPrimary);
}

.info-toggle:hover {
  color: var(--primary);
}

.info-toggle .material-icons {
  font-size: 1.3em;
}

.info-advanced {
  padding-top: 0.5em;
}

.info-sub-title {
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--textPrimary);
  margin: 1em 0 0.4em;
}

.info-code {
  font-size: 0.9em;
  background: var(--surfaceSecondary);
  padding: 0.1em 0.4em;
  border-radius: 4px;
}

.info-hash {
  display: inline-block;
  max-width: 100%;
}

.info-hash-link {
  cursor: pointer;
  color: var(--primary);
  text-decoration: none;
}

.info-hash-link:hover {
  text-decoration: underline;
}

.info-hash-value {
  cursor: copy;
  word-break: break-all;
  font-size: 0.85em;
  user-select: all;
}
</style>
