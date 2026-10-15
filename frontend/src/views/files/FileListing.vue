<template>
  <div class="file-listing" @contextmenu="onContextMenu">
    <header-bar showMenu showLogo>
      <search />
      <title />
      <action
        class="search-button"
        icon="search"
        :label="t('buttons.search')"
        @action="openSearch()"
      />
    </header-bar>

    <breadcrumbs base="/files">
      <div
        class="toolbar"
        role="toolbar"
        :aria-label="t('files.folderActions')"
      >
        <template v-if="fileStore.selectedCount > 0">
          <span class="toolbar__selection">
            {{
              t(
                "files.selected",
                { count: fileStore.selectedCount },
                fileStore.selectedCount
              )
            }}
            <button
              type="button"
              class="toolbar__clear"
              :aria-label="t('files.clearSelection')"
              :title="t('files.clearSelection')"
              @click="fileStore.selected = []"
            >
              <i class="material-icons" aria-hidden="true">close</i>
            </button>
          </span>
          <button
            type="button"
            class="toolbar__btn"
            :aria-label="t('files.moreActions')"
            :title="t('files.moreActions')"
            @click.stop="openMoreMenu"
          >
            <i class="material-icons" aria-hidden="true">more_horiz</i>
          </button>
          <span class="toolbar__divider" aria-hidden="true"></span>
        </template>

        <button
          v-if="folderFavorite"
          type="button"
          class="toolbar__btn"
          :class="{ on: isFolderFavorite }"
          :aria-pressed="isFolderFavorite"
          :aria-label="
            isFolderFavorite
              ? t('sidebar.removeFavorite')
              : t('sidebar.addFavorite')
          "
          :title="
            isFolderFavorite
              ? t('sidebar.removeFavorite')
              : t('sidebar.addFavorite')
          "
          @click="toggleFolderFavorite"
        >
          <i class="material-icons" aria-hidden="true">{{
            isFolderFavorite ? "star" : "star_border"
          }}</i>
        </button>
        <button
          v-if="headerButtons.upload"
          type="button"
          class="toolbar__btn"
          :aria-label="t('sidebar.newFolder')"
          :title="t('sidebar.newFolder')"
          @click="layoutStore.showHover('newDir')"
        >
          <i class="material-icons" aria-hidden="true">create_new_folder</i>
        </button>
        <button
          v-if="headerButtons.upload"
          type="button"
          class="toolbar__btn"
          :aria-label="t('sidebar.newFile')"
          :title="t('sidebar.newFile')"
          @click="layoutStore.showHover('newFile')"
        >
          <i class="material-icons" aria-hidden="true">note_add</i>
        </button>
        <button
          v-if="headerButtons.shell"
          type="button"
          class="toolbar__btn"
          :class="{ on: layoutStore.showShell }"
          :aria-pressed="layoutStore.showShell"
          :aria-label="t('sidebar.terminal')"
          :title="t('sidebar.terminal')"
          @click="layoutStore.toggleShell"
        >
          <i class="material-icons" aria-hidden="true">terminal</i>
        </button>
        <button
          type="button"
          class="toolbar__btn"
          :class="{ on: fileStore.multiple }"
          :aria-pressed="fileStore.multiple"
          :aria-label="t('buttons.selectMultiple')"
          :title="t('buttons.selectMultiple')"
          @click="toggleMultipleSelection"
        >
          <i class="material-icons" aria-hidden="true">checklist</i>
        </button>
        <button
          type="button"
          class="toolbar__btn"
          :aria-label="t('buttons.info')"
          :title="t('buttons.info')"
          @click="layoutStore.showHover('info')"
        >
          <i class="material-icons" aria-hidden="true">info</i>
        </button>

        <div class="segmented" role="group" :aria-label="t('files.view')">
          <button
            v-for="v in viewModes"
            :key="v.mode"
            type="button"
            :class="{ on: currentViewMode === v.mode }"
            :aria-pressed="currentViewMode === v.mode"
            :aria-label="t(v.label)"
            :title="t(v.label)"
            @click="setView(v.mode)"
          >
            <i class="material-icons" aria-hidden="true">{{ v.icon }}</i>
          </button>
        </div>

        <button
          v-if="headerButtons.upload"
          type="button"
          class="button toolbar__upload"
          @click="uploadFunc"
        >
          <i class="material-icons" aria-hidden="true">file_upload</i>
          <span>{{ t("buttons.upload") }}</span>
        </button>
      </div>
    </breadcrumbs>

    <div
      v-if="isMobile && fileStore.selectedCount > 0"
      id="file-selection"
      :class="{
        'file-selection-margin-bottom': fileStore.multiple,
      }"
    >
      <span v-if="fileStore.selectedCount > 0">
        {{ t("prompts.filesSelected", fileStore.selectedCount) }}
      </span>
      <action
        v-if="headerButtons.share"
        icon="share"
        :label="t('buttons.share')"
        show="share"
      />
      <action
        v-if="headerButtons.rename"
        icon="mode_edit"
        :label="t('buttons.rename')"
        show="rename"
      />
      <action
        v-if="headerButtons.copy"
        icon="content_copy"
        :label="t('buttons.copyFile')"
        show="copy"
      />
      <action
        v-if="headerButtons.move"
        icon="forward"
        :label="t('buttons.moveFile')"
        show="move"
      />
      <action
        v-if="headerButtons.delete"
        icon="delete"
        :label="t('buttons.delete')"
        show="delete"
      />
    </div>

    <div v-if="layoutStore.loading">
      <h2 class="message delayed">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span>{{ t("files.loading") }}</span>
      </h2>
    </div>
    <template v-else>
      <div
        v-if="
          (fileStore.req?.numDirs ?? 0) + (fileStore.req?.numFiles ?? 0) == 0
        "
      >
        <empty-state
          v-if="headerButtons.upload"
          icon="create_new_folder"
          :title="
            dragCounter > 0 ? t('files.dropToUpload') : t('files.emptyTitle')
          "
          :description="t('files.emptyText')"
          dropzone
          :active="dragCounter > 0"
        >
          <button type="button" class="button" @click="uploadFunc">
            <i class="material-icons" aria-hidden="true">file_upload</i>
            {{ t("files.uploadFiles") }}
          </button>
          <button
            type="button"
            class="button button--secondary"
            @click="layoutStore.showHover('newDir')"
          >
            <i class="material-icons" aria-hidden="true">create_new_folder</i>
            {{ t("sidebar.newFolder") }}
          </button>
        </empty-state>
        <empty-state
          v-else
          icon="folder_open"
          :title="t('files.emptyTitle')"
          :description="t('files.emptyTextReadOnly')"
        />
        <input
          style="display: none"
          type="file"
          id="upload-input"
          @change="uploadInput($event)"
          multiple
        />
        <input
          style="display: none"
          type="file"
          id="upload-folder-input"
          @change="uploadInput($event)"
          webkitdirectory
          multiple
        />
      </div>
      <div
        v-else
        id="listing"
        ref="listing"
        class="file-icons"
        data-clear-on-click="true"
        :class="authStore.user?.viewMode ?? ''"
        @click="handleEmptyAreaClick"
      >
        <div>
          <div class="item header">
            <div>
              <p
                :class="{ active: nameSorted }"
                class="name"
                role="button"
                tabindex="0"
                @click="sort('name')"
                :title="t('files.sortByName')"
                :aria-label="t('files.sortByName')"
              >
                <span>{{ t("files.name") }}</span>
                <i class="material-icons">{{ nameIcon }}</i>
              </p>

              <p
                :class="{ active: sizeSorted }"
                class="size"
                role="button"
                tabindex="0"
                @click="sort('size')"
                :title="t('files.sortBySize')"
                :aria-label="t('files.sortBySize')"
              >
                <span>{{ t("files.size") }}</span>
                <i class="material-icons">{{ sizeIcon }}</i>
              </p>
              <p
                :class="{ active: modifiedSorted }"
                class="modified"
                role="button"
                tabindex="0"
                @click="sort('modified')"
                :title="t('files.sortByLastModified')"
                :aria-label="t('files.sortByLastModified')"
              >
                <span>{{ t("files.lastModified") }}</span>
                <i class="material-icons">{{ modifiedIcon }}</i>
              </p>
            </div>
          </div>
        </div>

        <h2 data-clear-on-click="true" v-if="fileStore.req?.numDirs ?? false">
          {{ t("files.folders") }}
        </h2>
        <div v-if="fileStore.req?.numDirs ?? false" data-clear-on-click="true">
          <item
            v-for="item in dirs"
            :key="base64(item.name)"
            v-bind:index="item.index"
            v-bind:name="item.name"
            v-bind:isDir="item.isDir"
            v-bind:url="item.url"
            v-bind:modified="item.modified"
            v-bind:type="item.type"
            v-bind:size="item.size"
            v-bind:path="item.path"
          >
          </item>
        </div>

        <h2 data-clear-on-click="true" v-if="fileStore.req?.numFiles ?? false">
          {{ t("files.files") }}
        </h2>
        <div v-if="fileStore.req?.numFiles ?? false" data-clear-on-click="true">
          <item
            v-for="item in files"
            :key="base64(item.name)"
            v-bind:index="item.index"
            v-bind:name="item.name"
            v-bind:isDir="item.isDir"
            v-bind:url="item.url"
            v-bind:modified="item.modified"
            v-bind:type="item.type"
            v-bind:size="item.size"
            v-bind:path="item.path"
          >
          </item>
        </div>

        <input
          style="display: none"
          type="file"
          id="upload-input"
          @change="uploadInput($event)"
          multiple
        />
        <input
          style="display: none"
          type="file"
          id="upload-folder-input"
          @change="uploadInput($event)"
          webkitdirectory
          multiple
        />

        <div :class="{ active: fileStore.multiple }" id="multiple-selection">
          <p>{{ t("files.multipleSelectionEnabled") }}</p>
          <div
            @click="() => (fileStore.multiple = false)"
            tabindex="0"
            role="button"
            :title="t('buttons.clear')"
            :aria-label="t('buttons.clear')"
            class="action"
          >
            <i class="material-icons">clear</i>
          </div>
        </div>
      </div>
    </template>

    <context-menu
      :show="isContextMenuVisible"
      :pos="contextMenuPos"
      @hide="hideContextMenu"
    >
      <template v-if="contextMenuMode === 'item'">
        <action
          v-if="headerButtons.share"
          icon="share"
          :label="t('buttons.share')"
          show="share"
        />
        <action
          v-if="headerButtons.rename"
          icon="mode_edit"
          :label="t('buttons.rename')"
          show="rename"
        />
        <action
          v-if="headerButtons.copy"
          icon="content_copy"
          :label="t('buttons.copyFile')"
          show="copy"
        />
        <action
          v-if="headerButtons.move"
          icon="drive_file_move"
          :label="t('buttons.moveFile')"
          show="move"
        />
        <action
          v-if="headerButtons.download"
          icon="file_download"
          :label="t('buttons.download')"
          @action="download"
          :counter="fileStore.selectedCount"
        />
        <action
          v-if="favoriteTarget"
          :icon="isFavoriteTarget ? 'star' : 'star_border'"
          :label="
            isFavoriteTarget
              ? t('sidebar.removeFavorite')
              : t('sidebar.addFavorite')
          "
          @action="toggleFavorite"
        />
        <action icon="info" :label="t('buttons.info')" show="info" />
        <template v-if="headerButtons.delete">
          <hr class="context-menu__divider" />
          <action
            class="action--danger"
            icon="delete"
            :label="t('buttons.delete')"
            show="delete"
          />
        </template>
      </template>

      <template v-else>
        <template v-if="headerButtons.upload">
          <action
            icon="create_new_folder"
            :label="t('sidebar.newFolder')"
            show="newDir"
          />
          <action
            icon="note_add"
            :label="t('sidebar.newFile')"
            show="newFile"
          />
          <action
            icon="file_upload"
            :label="t('buttons.upload')"
            @action="uploadFunc"
          />
          <hr class="context-menu__divider" />
        </template>
        <action
          v-if="(fileStore.req?.items.length ?? 0) > 0"
          icon="select_all"
          :label="t('files.selectAll')"
          @action="selectAll"
        />
        <action
          v-if="folderFavorite"
          :icon="isFolderFavorite ? 'star' : 'star_border'"
          :label="
            isFolderFavorite
              ? t('sidebar.removeFavorite')
              : t('sidebar.addFavorite')
          "
          @action="toggleFolderFavorite"
        />
        <action icon="info" :label="t('buttons.info')" show="info" />
      </template>
    </context-menu>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useClipboardStore } from "@/stores/clipboard";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { useFavoritesStore } from "@/stores/favorites";

import { users, files as api } from "@/api";
import { enableExec } from "@/utils/constants";
import * as upload from "@/utils/upload";
import buttons from "@/utils/buttons";
import css from "@/utils/css";
import { throttle } from "lodash-es";
import { Base64 } from "js-base64";

import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";
import Search from "@/components/Search.vue";
import Item from "@/components/files/ListingItem.vue";
import ContextMenu from "@/components/ContextMenu.vue";
import EmptyState from "@/components/EmptyState.vue";
import Breadcrumbs from "@/components/Breadcrumbs.vue";
import {
  computed,
  inject,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  watch,
} from "vue";
import { useRoute, onBeforeRouteUpdate } from "vue-router";
import { useI18n } from "vue-i18n";
import { storeToRefs } from "pinia";
import { removePrefix } from "@/api/utils";

const showLimit = ref<number>(50);
const columnWidth = ref<number>(280);
const dragCounter = ref<number>(0);
const width = ref<number>(window.innerWidth);
const itemWeight = ref<number>(0);
const isContextMenuVisible = ref<boolean>(false);
// "item": actions for the selection; "folder": actions for this folder
const contextMenuMode = ref<"item" | "folder">("item");
const contextMenuPos = ref<{ x: number; y: number }>({ x: 0, y: 0 });

const $showError = inject<IToastError>("$showError")!;

const clipboardStore = useClipboardStore();
const authStore = useAuthStore();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const { req } = storeToRefs(fileStore);

const route = useRoute();
onBeforeRouteUpdate(() => {
  hideContextMenu();
});

const { t } = useI18n();

const listing = ref<HTMLElement | null>(null);

const nameSorted = computed(() =>
  fileStore.req ? fileStore.req.sorting.by === "name" : false
);

const sizeSorted = computed(() =>
  fileStore.req ? fileStore.req.sorting.by === "size" : false
);

const modifiedSorted = computed(() =>
  fileStore.req ? fileStore.req.sorting.by === "modified" : false
);

const ascOrdered = computed(() =>
  fileStore.req ? fileStore.req.sorting.asc : false
);

const dirs = computed(() => items.value.dirs.slice(0, showLimit.value));

const items = computed(() => {
  const dirs: any[] = [];
  const files: any[] = [];

  fileStore.req?.items.forEach((item) => {
    if (item.isDir) {
      dirs.push(item);
    } else {
      files.push(item);
    }
  });

  return { dirs, files };
});

const files = computed((): Resource[] => {
  let _showLimit = showLimit.value - items.value.dirs.length;

  if (_showLimit < 0) _showLimit = 0;

  return items.value.files.slice(0, _showLimit);
});

const nameIcon = computed(() => {
  if (nameSorted.value && !ascOrdered.value) {
    return "arrow_upward";
  }

  return "arrow_downward";
});

const sizeIcon = computed(() => {
  if (sizeSorted.value && ascOrdered.value) {
    return "arrow_downward";
  }

  return "arrow_upward";
});

const modifiedIcon = computed(() => {
  if (modifiedSorted.value && ascOrdered.value) {
    return "arrow_downward";
  }

  return "arrow_upward";
});

const headerButtons = computed(() => {
  return {
    upload: authStore.user?.perm.create,
    download: authStore.user?.perm.download,
    shell: authStore.user?.perm.execute && enableExec,
    delete: fileStore.selectedCount > 0 && authStore.user?.perm.delete,
    rename: fileStore.selectedCount === 1 && authStore.user?.perm.rename,
    share:
      fileStore.selectedCount === 1 &&
      authStore.user?.perm.share &&
      authStore.user?.perm.download,
    move: fileStore.selectedCount > 0 && authStore.user?.perm.rename,
    copy: fileStore.selectedCount > 0 && authStore.user?.perm.create,
  };
});

const favoritesStore = useFavoritesStore();
favoritesStore.load();

// What the star action pins: the single selected item, or the current
// folder when nothing is selected (the root is already under Locations).
const favoriteTarget = computed(() => {
  const req = fileStore.req;
  if (!req) return null;
  if (fileStore.selectedCount === 1) {
    const item = req.items.find((i) => i.index === fileStore.selected[0]);
    if (!item) return null;
    return {
      path: item.path,
      url: item.url,
      name: item.name,
      isDir: item.isDir,
    };
  }
  if (fileStore.selectedCount === 0 && req.path !== "/") {
    return { path: req.path, url: req.url, name: req.name, isDir: true };
  }
  return null;
});

const isFavoriteTarget = computed(() => {
  return favoriteTarget.value
    ? favoritesStore.has(favoriteTarget.value.path)
    : false;
});

const toggleFavorite = () => {
  if (favoriteTarget.value) favoritesStore.toggle(favoriteTarget.value);
};

// The folder being viewed (the root can't be pinned: it's "My files")
const folderFavorite = computed(() => {
  const req = fileStore.req;
  if (!req || req.path === "/") return null;
  return { path: req.path, url: req.url, name: req.name, isDir: true };
});

const isFolderFavorite = computed(() =>
  folderFavorite.value ? favoritesStore.has(folderFavorite.value.path) : false
);

const toggleFolderFavorite = () => {
  if (folderFavorite.value) favoritesStore.toggle(folderFavorite.value);
};

const viewModes: { mode: ViewModeType; icon: string; label: string }[] = [
  { mode: "list", icon: "view_list", label: "files.viewList" },
  { mode: "mosaic", icon: "grid_view", label: "files.viewGrid" },
  { mode: "mosaic gallery", icon: "photo_library", label: "files.viewGallery" },
];

const currentViewMode = computed(() => authStore.user?.viewMode ?? "list");

const isMobile = computed(() => {
  return width.value <= 736;
});

watch(req, () => {
  // Reset the show value
  showLimit.value = 50;

  nextTick(() => {
    // Ensures that the listing is displayed
    // How much every listing item affects the window height
    setItemWeight();

    // Scroll to the item opened previously
    if (!revealPreviousItem()) {
      // Fill and fit the window with listing items
      fillWindow(true);
    }
  });
});

onMounted(() => {
  // Check the columns size for the first time.
  columnsResize();

  // How much every listing item affects the window height
  setItemWeight();

  // Scroll to the item opened previously
  if (!revealPreviousItem()) {
    // Fill and fit the window with listing items
    fillWindow(true);
  }

  // Add the needed event listeners to the window and document.
  window.addEventListener("keydown", keyEvent);
  window.addEventListener("scroll", scrollEvent);
  window.addEventListener("resize", windowsResize);

  if (!authStore.user?.perm.create) return;
  document.addEventListener("dragover", preventDefault);
  document.addEventListener("dragenter", dragEnter);
  document.addEventListener("dragleave", dragLeave);
  document.addEventListener("drop", drop);
});

onBeforeUnmount(() => {
  // Remove event listeners before destroying this page.
  window.removeEventListener("keydown", keyEvent);
  window.removeEventListener("scroll", scrollEvent);
  window.removeEventListener("resize", windowsResize);

  if (authStore.user && !authStore.user?.perm.create) return;
  document.removeEventListener("dragover", preventDefault);
  document.removeEventListener("dragenter", dragEnter);
  document.removeEventListener("dragleave", dragLeave);
  document.removeEventListener("drop", drop);
});

const base64 = (name: string) => Base64.encodeURI(name);

const keyEvent = (event: KeyboardEvent) => {
  // No prompts are shown
  if (layoutStore.currentPrompt !== null) {
    return;
  }

  if (event.key === "Escape") {
    // Reset files selection.
    fileStore.selected = [];
  }

  if (event.key === "Delete") {
    if (!authStore.user?.perm.delete || fileStore.selectedCount == 0) return;

    // Show delete prompt.
    layoutStore.showHover("delete");
  }

  if (event.key === "F2") {
    if (!authStore.user?.perm.rename || fileStore.selectedCount !== 1) return;

    // Show rename prompt.
    layoutStore.showHover("rename");
  }

  // Ctrl is pressed
  if (!event.ctrlKey && !event.metaKey) {
    return;
  }

  switch (event.key) {
    case "f":
    case "F":
      if (event.shiftKey) {
        event.preventDefault();
        layoutStore.showHover("search");
      }
      break;
    case "c":
    case "x":
      copyCut(event);
      break;
    case "v":
      paste(event);
      break;
    case "a":
      event.preventDefault();
      selectAll();
      break;
    case "s":
      event.preventDefault();
      if (headerButtons.value.download) download();
      break;
  }
};

const preventDefault = (event: Event) => {
  // Wrapper around prevent default.
  event.preventDefault();
};

const copyCut = (event: Event | KeyboardEvent): void => {
  if ((event.target as HTMLElement).tagName?.toLowerCase() === "input") return;

  if (fileStore.req === null) return;

  const items = [];

  for (const i of fileStore.selected) {
    items.push({
      from: fileStore.req.items[i].url,
      name: fileStore.req.items[i].name,
      size: fileStore.req.items[i].size,
      isDir: fileStore.req.items[i].isDir,
      modified: fileStore.req.items[i].modified,
    });
  }

  if (items.length === 0) {
    return;
  }

  clipboardStore.$patch({
    key: (event as KeyboardEvent).key,
    items,
    path: route.path,
  });
};

const paste = async (event: Event) => {
  if ((event.target as HTMLElement).tagName?.toLowerCase() === "input") return;

  // TODO router location should it be
  const items: any[] = [];

  for (const item of clipboardStore.items) {
    const from = item.from.endsWith("/") ? item.from.slice(0, -1) : item.from;
    const to = route.path + encodeURIComponent(item.name);
    items.push({
      from,
      to,
      name: item.name,
      size: item.size,
      isDir: item.isDir,
      modified: item.modified,
      overwrite: false,
      rename: clipboardStore.path == route.path,
    });
  }

  if (items.length === 0) {
    return;
  }

  const preselect = removePrefix(route.path) + items[0].name;

  let action = (overwrite?: boolean, rename?: boolean) => {
    api
      .copy(items, overwrite, rename)
      .then(() => {
        fileStore.preselect = preselect;
        fileStore.reload = true;
      })
      .catch($showError);
  };

  if (clipboardStore.key === "x") {
    action = (overwrite, rename) => {
      api
        .move(items, overwrite, rename)
        .then(() => {
          clipboardStore.resetClipboard();
          fileStore.preselect = preselect;
          fileStore.reload = true;
        })
        .catch($showError);
    };
  }

  const path = route.path.endsWith("/") ? route.path : route.path + "/";
  const conflict = await upload.checkConflict(items, path, true);

  if (conflict.length > 0) {
    layoutStore.showHover({
      prompt: "resolve-conflict",
      props: {
        conflict: conflict,
      },
      confirm: (event: Event, result: Array<ConflictingResource>) => {
        event.preventDefault();
        layoutStore.closeHovers();
        for (let i = result.length - 1; i >= 0; i--) {
          const item = result[i];
          if (item.checked.length == 2) {
            items[item.index].rename = true;
          } else if (item.checked.length == 1 && item.checked[0] == "origin") {
            items[item.index].overwrite = true;
          } else {
            items.splice(item.index, 1);
          }
        }
        if (items.length > 0) {
          action();
        }
      },
    });

    return;
  }

  action(false, false);
};

const columnsResize = () => {
  // Update the columns size based on the window width.
  const items_ = css(["#listing.mosaic .item", ".mosaic#listing .item"]);
  if (items_ === null) return;

  let columns = Math.floor(
    (document.querySelector("main")?.offsetWidth ?? 0) / columnWidth.value
  );
  if (columns === 0) columns = 1;
  items_.style.width = `calc(${100 / columns}% - 1em)`;
};

const scrollEvent = throttle(() => {
  const totalItems =
    (fileStore.req?.numDirs ?? 0) + (fileStore.req?.numFiles ?? 0);

  // All items are displayed
  if (showLimit.value >= totalItems) return;

  const currentPos = window.innerHeight + window.scrollY;

  // Trigger at the 75% of the window height
  const triggerPos = document.body.offsetHeight - window.innerHeight * 0.25;

  if (currentPos > triggerPos) {
    // Quantity of items needed to fill 2x of the window height
    const showQuantity = Math.ceil((window.innerHeight * 2) / itemWeight.value);

    // Increase the number of displayed items
    showLimit.value += showQuantity;
  }
}, 100);

const dragEnter = () => {
  dragCounter.value++;

  // When the user starts dragging an item, put every
  // file on the listing with 50% opacity.
  const items = document.getElementsByClassName("item");

  Array.from(items).forEach((file: Element) => {
    (file as HTMLElement).style.opacity = "0.5";
  });
};

const dragLeave = () => {
  dragCounter.value--;

  if (dragCounter.value == 0) {
    resetOpacity();
  }
};

const drop = async (event: DragEvent) => {
  event.preventDefault();
  dragCounter.value = 0;
  resetOpacity();

  const dt = event.dataTransfer;
  let el: HTMLElement | null = event.target as HTMLElement;

  if (fileStore.req === null || dt === null || dt.files.length <= 0) return;

  for (let i = 0; i < 5; i++) {
    if (el !== null && !el.classList.contains("item")) {
      el = el.parentElement;
    }
  }

  const files: UploadList = (await upload.scanFiles(dt)) as UploadList;
  let path = route.path.endsWith("/") ? route.path : route.path + "/";

  if (
    el !== null &&
    el.classList.contains("item") &&
    el.dataset.dir === "true"
  ) {
    // Get url from ListingItem instance
    // TODO: Don't know what is happening here
    path = el.__vue__.url;

    try {
      (await api.fetch(path)).items;
    } catch (error: any) {
      $showError(error);
      return;
    }
  }

  // Checking the destination hits the server, so show it is working rather
  // than leaving the action looking inert until the upload starts.
  buttons.loading("upload");
  const conflict = await upload.checkConflict(files, path);

  const preselect = removePrefix(path) + (files[0].fullPath || files[0].name);

  if (conflict.length > 0) {
    buttons.done("upload");
    layoutStore.showHover({
      prompt: "resolve-conflict",
      props: {
        conflict: conflict,
        isUploadAction: true,
      },
      confirm: (event: Event, result: Array<ConflictingResource>) => {
        event.preventDefault();
        layoutStore.closeHovers();
        for (let i = result.length - 1; i >= 0; i--) {
          const item = result[i];
          if (item.checked.length == 2) {
            continue;
          } else if (item.checked.length == 1 && item.checked[0] == "origin") {
            files[item.index].overwrite = true;
          } else {
            files.splice(item.index, 1);
          }
        }
        if (files.length > 0) {
          upload.handleFiles(files, path);
          fileStore.preselect = preselect;
        }
      },
    });

    return;
  }

  upload.handleFiles(files, path);
  fileStore.preselect = preselect;
};

const uploadInput = async (event: Event) => {
  const files = (event.currentTarget as HTMLInputElement)?.files;
  if (files === null) return;

  const folder_upload = !!files[0].webkitRelativePath;

  const uploadFiles: UploadList = [];
  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    const fullPath = folder_upload ? file.webkitRelativePath : undefined;
    uploadFiles.push({
      file,
      name: file.name,
      size: file.size,
      isDir: false,
      fullPath,
    });
  }

  const path = route.path.endsWith("/") ? route.path : route.path + "/";

  // Checking the destination hits the server, so show it is working rather
  // than leaving the action looking inert until the upload starts.
  buttons.loading("upload");
  const conflict = await upload.checkConflict(uploadFiles, path);

  if (conflict.length > 0) {
    buttons.done("upload");
    layoutStore.showHover({
      prompt: "resolve-conflict",
      props: {
        conflict: conflict,
        isUploadAction: true,
      },
      confirm: (event: Event, result: Array<ConflictingResource>) => {
        event.preventDefault();
        layoutStore.closeHovers();
        for (let i = result.length - 1; i >= 0; i--) {
          const item = result[i];
          if (item.checked.length == 2) {
            continue;
          } else if (item.checked.length == 1 && item.checked[0] == "origin") {
            uploadFiles[item.index].overwrite = true;
          } else {
            uploadFiles.splice(item.index, 1);
          }
        }
        if (uploadFiles.length > 0) {
          upload.handleFiles(uploadFiles, path);
        }
      },
    });

    return;
  }

  upload.handleFiles(uploadFiles, path);
};

const resetOpacity = () => {
  const items = document.getElementsByClassName("item");

  Array.from(items).forEach((file: Element) => {
    (file as HTMLElement).style.opacity = "1";
  });
};

const sort = async (by: string) => {
  let asc = false;

  if (by === "name") {
    if (nameIcon.value === "arrow_upward") {
      asc = true;
    }
  } else if (by === "size") {
    if (sizeIcon.value === "arrow_upward") {
      asc = true;
    }
  } else if (by === "modified") {
    if (modifiedIcon.value === "arrow_upward") {
      asc = true;
    }
  }

  try {
    if (authStore.user?.id) {
      await users.update({ id: authStore.user?.id, sorting: { by, asc } }, [
        "sorting",
      ]);
    }
  } catch (e: any) {
    $showError(e);
  }

  fileStore.reload = true;
};

const openSearch = () => {
  layoutStore.showHover("search");
};

const toggleMultipleSelection = () => {
  fileStore.toggleMultiple();
  layoutStore.closeHovers();
};

const windowsResize = throttle(() => {
  columnsResize();
  width.value = window.innerWidth;

  // Listing element is not displayed
  if (listing.value == null) return;

  // How much every listing item affects the window height
  setItemWeight();

  // Fill but not fit the window
  fillWindow();
}, 100);

const download = () => {
  if (fileStore.req === null) return;

  if (
    fileStore.selectedCount === 1 &&
    !fileStore.req.items[fileStore.selected[0]].isDir
  ) {
    api.download(null, fileStore.req.items[fileStore.selected[0]].url);
    return;
  }

  layoutStore.showHover({
    prompt: "download",
    confirm: (format: any) => {
      layoutStore.closeHovers();

      const files = [];

      if (fileStore.selectedCount > 0 && fileStore.req !== null) {
        for (const i of fileStore.selected) {
          files.push(fileStore.req.items[i].url);
        }
      } else {
        files.push(route.path);
      }

      api.download(format, ...files);
    },
  });
};

const selectAll = () => {
  for (const item of [...items.value.dirs, ...items.value.files]) {
    if (fileStore.selected.indexOf(item.index) === -1) {
      fileStore.selected.push(item.index);
    }
  }
};

const setView = (mode: ViewModeType) => {
  layoutStore.closeHovers();
  if (authStore.user?.viewMode === mode) return;

  const data = { id: authStore.user?.id, viewMode: mode };

  users.update(data, ["viewMode"]).catch($showError);

  authStore.updateUser(data);

  setItemWeight();
  fillWindow();
};

const uploadFunc = () => {
  if (
    typeof window.DataTransferItem !== "undefined" &&
    typeof DataTransferItem.prototype.webkitGetAsEntry !== "undefined"
  ) {
    layoutStore.showHover("upload");
  } else {
    document.getElementById("upload-input")?.click();
  }
};

const setItemWeight = () => {
  // Listing element is not displayed
  if (listing.value === null || fileStore.req === null) return;

  let itemQuantity = fileStore.req.numDirs + fileStore.req.numFiles;
  if (itemQuantity > showLimit.value) itemQuantity = showLimit.value;

  // How much every listing item affects the window height
  itemWeight.value = listing.value.offsetHeight / itemQuantity;
};

const fillWindow = (fit = false) => {
  if (fileStore.req === null) return;

  const totalItems = fileStore.req.numDirs + fileStore.req.numFiles;

  // More items are displayed than the total
  if (showLimit.value >= totalItems && !fit) return;

  const windowHeight = window.innerHeight;

  // Quantity of items needed to fill 2x of the window height
  const showQuantity = Math.ceil(
    (windowHeight + windowHeight * 2) / itemWeight.value
  );

  // Less items to display than current
  if (showLimit.value > showQuantity && !fit) return;

  // Set the number of displayed items
  showLimit.value = showQuantity > totalItems ? totalItems : showQuantity;
};

const revealPreviousItem = () => {
  if (!fileStore.req || !fileStore.oldReq) return;

  const index = fileStore.selected[0];
  if (index === undefined) return;

  showLimit.value =
    index + Math.ceil((window.innerHeight * 2) / itemWeight.value);

  nextTick(() => {
    const items = document.querySelectorAll("#listing .item");
    items[index].scrollIntoView({ block: "center" });
  });

  return true;
};

// Right click: on an item, act on the selection; elsewhere, on the folder
const onContextMenu = (event: MouseEvent) => {
  const target = event.target as HTMLElement | null;
  // Keep the browser's own menu for the top bar, path and text fields
  if (target?.closest("header, .breadcrumbs, input, textarea, .shell")) {
    return;
  }
  const onItem = !!target?.closest(".item:not(.header)");
  if (!onItem) fileStore.selected = [];
  contextMenuMode.value =
    onItem && fileStore.selectedCount > 0 ? "item" : "folder";
  showContextMenu(event);
};

// "More" button in the toolbar: the selection menu, under the button
const openMoreMenu = (event: MouseEvent) => {
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  contextMenuMode.value = "item";
  isContextMenuVisible.value = true;
  contextMenuPos.value = {
    x: rect.left,
    y: rect.bottom + 6 + Math.floor(window.scrollY),
  };
};

const showContextMenu = (event: MouseEvent) => {
  event.preventDefault();
  isContextMenuVisible.value = true;
  contextMenuPos.value = {
    x: event.clientX + 8,
    y: event.clientY + Math.floor(window.scrollY),
  };
};

const hideContextMenu = () => {
  isContextMenuVisible.value = false;
};

const handleEmptyAreaClick = (e: MouseEvent) => {
  const target = e.target;
  if (!(target instanceof HTMLElement)) return;

  if (target.dataset.clearOnClick === "true") {
    fileStore.selected = [];
  }
};
</script>
<style scoped>
/* The whole page area below the toolbar takes right clicks and drops */
.file-listing {
  min-height: calc(100vh - 5.5em);
}

.file-selection-margin-bottom {
  margin-bottom: 3.5rem;
}
</style>
