<template>
  <div id="editor-container">
    <file-header
      :name="fileStore.req?.name ?? ''"
      :icon="isMarkdownFile ? 'article' : 'description'"
      :meta="headerMeta"
      :close-label="t('buttons.close')"
      @close="close()"
    >
      <template #status>
        <span v-if="readOnly" class="file-header__badge">
          <i class="material-icons" aria-hidden="true">lock</i>
          {{ t("editor.readOnly") }}
        </span>
        <span
          v-else-if="isDirty"
          class="file-header__badge file-header__badge--warn"
        >
          {{ t("editor.unsaved") }}
        </span>
      </template>
      <template v-if="authStore.user?.perm.modify && !readOnly" #primary>
        <button
          type="button"
          id="save-button"
          class="button toolbar__upload file-header__primary"
          :disabled="layoutStore.loading"
          :aria-label="t('buttons.save')"
          :title="t('buttons.save') + ' (Ctrl+S)'"
          @click="save()"
        >
          <i class="material-icons" aria-hidden="true">save</i>
          <span>{{ t("buttons.save") }}</span>
        </button>
      </template>
    </file-header>

    <!-- preview container -->
    <div class="loading delayed" v-if="layoutStore.loading">
      <div class="spinner">
        <div class="bounce1"></div>
        <div class="bounce2"></div>
        <div class="bounce3"></div>
      </div>
    </div>
    <template v-else>
      <div class="editor-bar">
        <Breadcrumbs base="/files" noLink />

        <div
          class="toolbar editor-bar__tools"
          role="toolbar"
          :aria-label="t('editor.tools')"
        >
          <div
            v-if="isMarkdownFile"
            class="segmented"
            role="group"
            :aria-label="t('editor.mode')"
          >
            <button
              type="button"
              :class="{ on: !isPreview }"
              :aria-pressed="!isPreview"
              :aria-label="t('editor.edit')"
              :title="t('editor.edit')"
              @click="isPreview = false"
            >
              <i class="material-icons" aria-hidden="true">edit</i>
            </button>
            <button
              type="button"
              :class="{ on: isPreview }"
              :aria-pressed="isPreview"
              :aria-label="t('buttons.preview')"
              :title="t('buttons.preview')"
              @click="isPreview = true"
            >
              <i class="material-icons" aria-hidden="true">visibility</i>
            </button>
          </div>
          <span
            v-if="isMarkdownFile"
            class="toolbar__divider"
            aria-hidden="true"
          ></span>

          <button
            type="button"
            class="toolbar__btn"
            :disabled="isSelectionEmpty"
            :aria-label="t('buttons.copy')"
            :title="t('buttons.copy')"
            @click="executeEditorCommand('copy')"
          >
            <i class="material-icons" aria-hidden="true">content_copy</i>
          </button>
          <button
            type="button"
            class="toolbar__btn"
            :disabled="isSelectionEmpty || readOnly"
            :aria-label="t('editor.cut')"
            :title="t('editor.cut')"
            @click="executeEditorCommand('cut')"
          >
            <i class="material-icons" aria-hidden="true">content_cut</i>
          </button>
          <button
            type="button"
            class="toolbar__btn"
            :disabled="readOnly"
            :aria-label="t('editor.paste')"
            :title="t('editor.paste')"
            @click="executeEditorCommand('paste')"
          >
            <i class="material-icons" aria-hidden="true">content_paste</i>
          </button>
          <span class="toolbar__divider" aria-hidden="true"></span>

          <div
            class="segmented"
            role="group"
            :aria-label="t('editor.fontSize')"
          >
            <button
              type="button"
              :aria-label="t('buttons.decreaseFontSize')"
              :title="t('buttons.decreaseFontSize')"
              @click="decreaseFontSize"
            >
              <i class="material-icons" aria-hidden="true">text_decrease</i>
            </button>
            <span class="editor-bar__size" aria-live="polite"
              >{{ fontSize }}px</span
            >
            <button
              type="button"
              :aria-label="t('buttons.increaseFontSize')"
              :title="t('buttons.increaseFontSize')"
              @click="increaseFontSize"
            >
              <i class="material-icons" aria-hidden="true">text_increase</i>
            </button>
          </div>

          <button
            type="button"
            class="toolbar__btn"
            :aria-label="t('editor.commands')"
            :title="t('editor.commands') + ' (F1)'"
            @click="executeEditorCommand('openCommandPalette')"
          >
            <i class="material-icons" aria-hidden="true">more_horiz</i>
          </button>
        </div>
      </div>

      <div
        v-show="isPreview && isMarkdownFile"
        id="preview-container"
        class="md_preview"
        v-html="previewContent"
      ></div>
      <form v-show="!isPreview || !isMarkdownFile" id="editor"></form>
    </template>
  </div>
</template>

<script setup lang="ts">
import { files as api } from "@/api";
import buttons from "@/utils/buttons";
import url from "@/utils/url";
import ace, { Ace, version as ace_version } from "ace-builds";
import "ace-builds/src-noconflict/ext-language_tools";
import modelist from "ace-builds/src-noconflict/ext-modelist";
import DOMPurify from "dompurify";

import Breadcrumbs from "@/components/Breadcrumbs.vue";
import FileHeader from "@/components/files/FileHeader.vue";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { getEditorTheme } from "@/utils/theme";
import { marked } from "marked";
import markedKatex from "marked-katex-extension";
import {
  computed,
  inject,
  onBeforeUnmount,
  onMounted,
  ref,
  watchEffect,
} from "vue";
import prettyBytes from "pretty-bytes";
import { useI18n } from "vue-i18n";
import { onBeforeRouteUpdate, useRoute, useRouter } from "vue-router";
import { read, copy } from "@/utils/clipboard";

const $showError = inject<IToastError>("$showError")!;

const fileStore = useFileStore();
const authStore = useAuthStore();
const layoutStore = useLayoutStore();

const { t } = useI18n();

const route = useRoute();
const router = useRouter();

const editor = ref<Ace.Editor | null>(null);
const fontSize = ref(parseInt(localStorage.getItem("editorFontSize") || "14"));

const isPreview = ref(false);
const previewContent = ref("");
const isMarkdownFile =
  fileStore.req?.name.endsWith(".md") ||
  fileStore.req?.name.endsWith(".markdown");
const katexOptions = {
  output: "mathml" as const,
  throwOnError: false,
};
marked.use(markedKatex(katexOptions));

const isSelectionEmpty = ref(true);
const isDirty = ref(false);
const readOnly = computed(() => fileStore.req?.type === "textImmutable");

// "Folder · 4 KiB" under the file name
const headerMeta = computed(() => {
  const req = fileStore.req;
  if (!req) return [];
  const dirs = route.path.replace(/\/+$/, "").split("/");
  const folder = decodeURIComponent(dirs[dirs.length - 2] ?? "");
  return [
    folder && folder !== "files" ? folder : t("sidebar.myFiles"),
    prettyBytes(req.size, { binary: true }),
  ];
});

const executeEditorCommand = (name: string) => {
  if (name == "paste") {
    read()
      .then((data) => {
        editor.value?.execCommand("paste", {
          text: data,
        });
      })
      .catch((e) => {
        if (
          document.queryCommandSupported &&
          document.queryCommandSupported("paste")
        ) {
          document.execCommand("paste");
        } else {
          console.warn("the clipboard api is not supported", e);
        }
      });
    return;
  }
  if (name == "copy" || name == "cut") {
    const selectedText = editor.value?.getCopyText();
    copy({ text: selectedText });
  }
  editor.value?.execCommand(name);
};

onMounted(() => {
  window.addEventListener("keydown", keyEvent);
  window.addEventListener("beforeunload", handlePageChange);

  const fileContent = fileStore.req?.content || "";

  watchEffect(async () => {
    if (isMarkdownFile && isPreview.value) {
      const new_value = editor.value?.getValue() || "";
      try {
        previewContent.value = DOMPurify.sanitize(await marked(new_value));
      } catch (error) {
        console.error("Failed to convert content to HTML:", error);
        previewContent.value = "";
      }
    }
  });

  ace.config.set(
    "basePath",
    `https://cdn.jsdelivr.net/npm/ace-builds@${ace_version}/src-min-noconflict/`
  );

  if (!layoutStore.loading) {
    initEditor(fileContent);
  } else {
    const unwatch = watchEffect(() => {
      // Initialize editor when layout is loaded
      if (!layoutStore.loading) {
        setTimeout(() => {
          initEditor(fileContent);
          unwatch();
        }, 50);
      }
    });
  }
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", keyEvent);
  window.removeEventListener("beforeunload", handlePageChange);
  editor.value?.destroy();
});

onBeforeRouteUpdate((to, from, next) => {
  if (editor.value?.session.getUndoManager().isClean()) {
    next();

    return;
  }

  layoutStore.showHover({
    prompt: "discardEditorChanges",
    confirm: (event: Event) => {
      event.preventDefault();
      next();
    },
    saveAction: async () => {
      await save();
      next();
    },
  });
});

const initEditor = (fileContent: string) => {
  editor.value = ace.edit("editor", {
    value: fileContent,
    showPrintMargin: false,
    readOnly: readOnly.value,
    theme: getEditorTheme(authStore.user?.aceEditorTheme ?? ""),
    mode: modelist.getModeForPath(fileStore.req!.name).mode,
    wrap: true,
    enableBasicAutocompletion: true,
    enableLiveAutocompletion: true,
    enableSnippets: true,
  });

  editor.value.setFontSize(fontSize.value);
  editor.value.focus();

  // The undo manager updates just after "change", so check on the next tick
  editor.value.on("change", () => {
    setTimeout(() => {
      isDirty.value = !editor.value?.session.getUndoManager().isClean();
    });
  });

  const selection = editor.value?.getSelection();
  selection.on("changeSelection", function () {
    isSelectionEmpty.value = selection.isEmpty();
  });
};

const keyEvent = (event: KeyboardEvent) => {
  if (event.code === "Escape") {
    close();
  }

  if (!event.ctrlKey && !event.metaKey) {
    return;
  }

  if (event.key !== "s") {
    return;
  }

  event.preventDefault();
  save();
};

const handlePageChange = (event: BeforeUnloadEvent) => {
  if (!editor.value?.session.getUndoManager().isClean()) {
    event.preventDefault();
    // returnValue is now depecrated, though keeping in for legacy browser support
    // https://developer.mozilla.org/en-US/docs/Web/API/BeforeUnloadEvent/returnValue
    event.returnValue = true;
  }
};

const save = async (throwError?: boolean) => {
  const button = "save";
  buttons.loading("save");

  try {
    await api.put(route.path, editor.value?.getValue());
    editor.value?.session.getUndoManager().markClean();
    isDirty.value = false;
    buttons.success(button);
  } catch (e: any) {
    buttons.done(button);
    $showError(e);
    if (throwError) throw e;
  }
};

const increaseFontSize = () => {
  fontSize.value += 1;
  editor.value?.setFontSize(fontSize.value);
  localStorage.setItem("editorFontSize", fontSize.value.toString());
};

const decreaseFontSize = () => {
  if (fontSize.value > 1) {
    fontSize.value -= 1;
    editor.value?.setFontSize(fontSize.value);
    localStorage.setItem("editorFontSize", fontSize.value.toString());
  }
};

const close = () => {
  if (!editor.value?.session.getUndoManager().isClean()) {
    layoutStore.showHover({
      prompt: "discardEditorChanges",
      confirm: (event: Event) => {
        event.preventDefault();
        editor.value?.session.getUndoManager().reset();
        isDirty.value = false;
        finishClose();
      },
      saveAction: async () => {
        try {
          await save(true);
          finishClose();
        } catch {}
      },
    });
    return;
  }
  finishClose();
};

const finishClose = () => {
  const uri = url.removeLastDir(route.path) + "/";
  router.push({ path: uri });
};
</script>
