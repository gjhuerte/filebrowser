<template>
  <div class="card floating help" style="max-width: 700px; max-height: 80vh; overflow-y: auto;">
    <div class="card-content" style="padding: 1.5em;">
      <!-- Header -->
      <div style="display: flex; align-items: center; gap: 0.75em; margin-bottom: 1.25em;">
        <i class="material-icons" style="font-size: 1.75em; color: var(--primary);">help_outline</i>
        <div>
          <h2 style="margin: 0; font-size: 1.3em;">{{ $t("help.help") }}</h2>
          <p style="margin: 0; font-size: 0.85em; color: var(--textSecondary);">User Guide & Quick Reference</p>
        </div>
      </div>

      <!-- Tab navigation -->
      <div style="display: flex; gap: 0; border-bottom: 2px solid var(--borderPrimary); margin-bottom: 1.25em;">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="button button--flat"
          :style="{
            padding: '0.5em 1em',
            fontSize: '0.85em',
            fontWeight: activeTab === tab.id ? '600' : '400',
            borderBottom: activeTab === tab.id ? '2px solid var(--primary)' : '2px solid transparent',
            marginBottom: '-2px',
            color: activeTab === tab.id ? 'var(--primary)' : 'var(--textSecondary)',
            borderRadius: 0,
            background: 'none'
          }"
        >
          <i class="material-icons" style="font-size: 1.1em; vertical-align: middle; margin-right: 4px;">{{ tab.icon }}</i>
          {{ tab.label }}
        </button>
      </div>

      <!-- Getting Started Tab -->
      <div v-if="activeTab === 'start'">
        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">folder_open</i>
            Navigating Files
          </h3>
          <div class="help-list">
            <p><strong>Click</strong> a file or folder to select it. <strong>Double-click</strong> to open it.</p>
            <p>Use the <strong>sidebar</strong> on the left to jump to your home directory, favorites, or other locations.</p>
            <p>The <strong>breadcrumb path</strong> at the top shows where you are — click any part to navigate back.</p>
            <p>Use <strong>Ctrl + Click</strong> to select multiple items at once.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">cloud_upload</i>
            Uploading Files
          </h3>
          <div class="help-list">
            <p><strong>Drag and drop</strong> files directly onto the file list to upload them.</p>
            <p>Or use the <strong>New</strong> button in the top bar to upload files or create new files and folders.</p>
            <p>You can upload multiple files at once. Large files are supported.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">cloud_download</i>
            Downloading
          </h3>
          <div class="help-list">
            <p>Click the <strong>download</strong> button or press <strong>Ctrl + S</strong> to download the current file or folder.</p>
            <p>When downloading a folder, it will be packaged as a <strong>ZIP archive</strong>.</p>
            <p>Select multiple files and download them all at once as a ZIP.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">star</i>
            Favorites
          </h3>
          <div class="help-list">
            <p>Right-click a file or folder and choose <strong>Add to favorites</strong> to pin it to the sidebar.</p>
            <p>Favorites appear in the sidebar for quick access across all directories.</p>
            <p>You can rename or remove favorites at any time.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">tune</i>
            Customizing Locations
          </h3>
          <div class="help-list">
            <p>Click the <strong>tune</strong> button next to <strong>Locations</strong>, or press and hold a location, to edit the list.</p>
            <p>Drag a handle to move a folder, or focus it and use the <strong>arrow keys</strong>. Use the <strong>eye</strong> to hide or show a folder. Tap <strong>✓</strong> or press <strong>Esc</strong> when you're done.</p>
            <p>You can also right-click a location and choose <strong>Hide from sidebar</strong>. Hidden folders are counted under the list; click that line to bring them back.</p>
            <p>Your order and hidden folders are saved in this browser for your account. New folders appear at the end.</p>
          </div>
        </div>
      </div>

      <!-- Keyboard Shortcuts Tab -->
      <div v-if="activeTab === 'keys'">
        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">description</i>
            File Operations
          </h3>
          <table class="help-table">
            <tr><td class="key-cell"><kbd>F1</kbd></td><td>{{ $t("help.f1") }}</td></tr>
            <tr><td class="key-cell"><kbd>F2</kbd></td><td>{{ $t("help.f2") }}</td></tr>
            <tr><td class="key-cell"><kbd>DEL</kbd></td><td>{{ $t("help.del") }}</td></tr>
            <tr><td class="key-cell"><kbd>ESC</kbd></td><td>{{ $t("help.esc") }}</td></tr>
            <tr><td class="key-cell"><kbd>Ctrl + S</kbd></td><td>{{ $t("help.ctrl.s") }}</td></tr>
          </table>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">search</i>
            Search & Navigation
          </h3>
          <table class="help-table">
            <tr><td class="key-cell"><kbd>Ctrl + Shift + F</kbd></td><td>{{ $t("help.ctrl.f") }}</td></tr>
            <tr><td class="key-cell"><kbd>Click</kbd></td><td>{{ $t("help.click") }}</td></tr>
            <tr><td class="key-cell"><kbd>Double Click</kbd></td><td>{{ $t("help.doubleClick") }}</td></tr>
            <tr><td class="key-cell"><kbd>Ctrl + Click</kbd></td><td>{{ $t("help.ctrl.click") }}</td></tr>
          </table>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">visibility</i>
            Previewing Files
          </h3>
          <table class="help-table">
            <tr><td class="key-cell"><kbd>← / →</kbd></td><td>Previous / next file (seeks in videos)</td></tr>
            <tr><td class="key-cell"><kbd>Enter</kbd></td><td>Next file</td></tr>
            <tr><td class="key-cell"><kbd>+ / −</kbd></td><td>Zoom an image in or out</td></tr>
            <tr><td class="key-cell"><kbd>0</kbd></td><td>Fit the image to the screen</td></tr>
            <tr><td class="key-cell"><kbd>R</kbd></td><td>Rotate the image</td></tr>
            <tr><td class="key-cell"><kbd>Space</kbd></td><td>Play or pause audio and video</td></tr>
            <tr><td class="key-cell"><kbd>Esc</kbd></td><td>Close the preview</td></tr>
          </table>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">edit</i>
            Editor
          </h3>
          <table class="help-table">
            <tr><td class="key-cell"><kbd>Ctrl + S</kbd></td><td>Save current file</td></tr>
            <tr><td class="key-cell"><kbd>Ctrl + Z</kbd></td><td>Undo</td></tr>
            <tr><td class="key-cell"><kbd>Ctrl + Shift + Z</kbd></td><td>Redo</td></tr>
          </table>
        </div>
      </div>

      <!-- Features Tab -->
      <div v-if="activeTab === 'features'">
        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">share</i>
            Sharing Files
          </h3>
          <div class="help-list">
            <p>Select a file and click the <strong>share</strong> button to generate a shareable link.</p>
            <p>Set an <strong>expiration time</strong> for temporary links, or set a <strong>password</strong> for extra security.</p>
            <p>Manage your shared links from <strong>Settings → Shares</strong>.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">edit</i>
            File Editing
          </h3>
          <div class="help-list">
            <p>Click on any text file to open it in the <strong>built-in editor</strong>.</p>
            <p>The editor supports <strong>syntax highlighting</strong> for many programming languages.</p>
            <p>The top bar shows an <strong>Unsaved changes</strong> badge until you save. Markdown files have an <strong>Edit / Preview</strong> switch in the toolbar.</p>
            <p>Press <strong>Ctrl + S</strong> to save your changes.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">visibility</i>
            File Preview
          </h3>
          <div class="help-list">
            <p><strong>Images</strong> — Zoom, rotate, and switch to original quality from the toolbar at the bottom. Size and dimensions are shown there too.</p>
            <p><strong>Video</strong> — Built-in player with speed control, subtitles, and picture-in-picture. <strong>FLV</strong> files play too (not on iPhone).</p>
            <p><strong>Audio</strong> — Music player with album art (a <code>cover.jpg</code> or <code>folder.jpg</code> in the same folder), repeat, speed, and it moves on to the next track automatically.</p>
            <p>The top bar shows the folder, size, position in the folder and when the file was modified. <strong>Download</strong> is on the right; rename, open raw, info and delete sit beside it (in the ⋮ menu on phones).</p>
            <p>While a file downloads you'll see <strong>how much has loaded</strong> — a percentage and size where the server reports it.</p>
            <p><strong>PDF</strong> — In-browser PDF viewer.</p>
            <p><strong>Markdown</strong> — Rendered markdown preview.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">dashboard</i>
            Home Dashboard
          </h3>
          <div class="help-list">
            <p>Home shows your storage, file counts, and the largest and most recent files.</p>
            <p>Use the <strong>auto-refresh</strong> menu (top right of Home) to rescan your files every 5, 15, 30, or 60 minutes. Scans run in the background while File Browser is open and pause when the tab is hidden.</p>
            <p>Opening Home or reloading the page shows the last scan instead of starting a new one, so it loads instantly. The time of that scan is shown at the top right.</p>
            <p>The refresh button rescans right away.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">info</i>
            File Information
          </h3>
          <div class="help-list">
            <p>Select a file and click the <strong>info</strong> button to view file details.</p>
            <p>See file <strong>size, modification date, path</strong> and more.</p>
            <p>Click <strong>Technical Info</strong> to reveal permissions, checksums, and other metadata.</p>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">search</i>
            Search
          </h3>
          <div class="help-list">
            <p>Press <strong>Ctrl + Shift + F</strong> or click the search icon to search files.</p>
            <p>Search works across the current directory and all subdirectories.</p>
            <p>Results update as you type.</p>
          </div>
        </div>
      </div>

      <!-- Troubleshooting Tab -->
      <div v-if="activeTab === 'trouble'">
        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">error_outline</i>
            Common Issues
          </h3>

          <div class="help-faq">
            <button @click="toggleFaq('upload')" class="help-faq-toggle">
              <span>Files fail to upload</span>
              <i class="material-icons">{{ openFaq === 'upload' ? 'expand_less' : 'expand_more' }}</i>
            </button>
            <div v-if="openFaq === 'upload'" class="help-faq-content">
              <p>Check that you have <strong>write permissions</strong> to the current directory.</p>
              <p>For large files, verify that your reverse proxy allows the upload size (e.g., <code>client_max_body_size</code> in nginx).</p>
              <p>If using Docker, make sure the volume is mounted with the correct <strong>PUID/PGID</strong> settings.</p>
            </div>
          </div>

          <div class="help-faq">
            <button @click="toggleFaq('permission')" class="help-faq-toggle">
              <span>Permission denied errors</span>
              <i class="material-icons">{{ openFaq === 'permission' ? 'expand_less' : 'expand_more' }}</i>
            </button>
            <div v-if="openFaq === 'permission'" class="help-faq-content">
              <p>File Browser runs as a specific user. Make sure the files and directories you're accessing are owned by or readable by that user.</p>
              <p>In Docker, set <code>PUID</code> and <code>PGID</code> environment variables to match the file owner on the host.</p>
            </div>
          </div>

          <div class="help-faq">
            <button @click="toggleFaq('preview')" class="help-faq-toggle">
              <span>File preview not working</span>
              <i class="material-icons">{{ openFaq === 'preview' ? 'expand_less' : 'expand_more' }}</i>
            </button>
            <div v-if="openFaq === 'preview'" class="help-faq-content">
              <p>Not all file formats support preview. Supported formats include common image, video, audio, PDF, and text files.</p>
              <p>If using a reverse proxy, ensure <strong>WebSocket connections</strong> are allowed.</p>
              <p>Try clearing your browser cache and reloading the page.</p>
            </div>
          </div>

          <div class="help-faq">
            <button @click="toggleFaq('login')" class="help-faq-toggle">
              <span>Can't log in</span>
              <i class="material-icons">{{ openFaq === 'login' ? 'expand_less' : 'expand_more' }}</i>
            </button>
            <div v-if="openFaq === 'login'" class="help-faq-content">
              <p>Default credentials are <strong>admin / admin</strong>. If you've changed them and forgotten, you'll need to reset the database or use the CLI.</p>
              <p>Clear your browser cookies for this site and try again.</p>
              <p>If using proxy authentication, check that the proxy headers are configured correctly.</p>
              <p>Seeing <strong>"Too many failed attempts"</strong>? After 10 wrong passwords in 15 minutes, logins from your address pause for the rest of that window. Wait a few minutes, or ask your administrator to restart File Browser.</p>
            </div>
          </div>

          <div class="help-faq">
            <button @click="toggleFaq('docker')" class="help-faq-toggle">
              <span>Docker container issues</span>
              <i class="material-icons">{{ openFaq === 'docker' ? 'expand_less' : 'expand_more' }}</i>
            </button>
            <div v-if="openFaq === 'docker'" class="help-faq-content">
              <p>Make sure volumes are mounted correctly: <code>-v /host/path:/srv</code> for files, <code>-v /host/config:/config</code> for configuration.</p>
              <p>Set <code>PUID</code> and <code>PGID</code> to your host user's UID/GID to avoid permission issues.</p>
              <p>Check container logs with <code>docker logs filebrowser</code> for error messages.</p>
            </div>
          </div>
        </div>

        <div class="help-section">
          <h3 class="help-heading">
            <i class="material-icons">tips_and_updates</i>
            Tips
          </h3>
          <div class="help-list">
            <p>Use a <strong>reverse proxy</strong> (nginx, Caddy, Traefik) with HTTPS in production.</p>
            <p>Run the container as an <strong>unprivileged user</strong> for better security.</p>
            <p>Regularly <strong>back up</strong> the database file (<code>filebrowser.db</code>).</p>
            <p>Keep the <strong>command runner disabled</strong> unless you specifically need it.</p>
          </div>
        </div>
      </div>
    </div>

    <div class="card-action">
      <button
        id="focus-prompt"
        type="submit"
        @click="closeHovers"
        class="button button--flat"
        :aria-label="$t('buttons.close')"
        :title="$t('buttons.close')"
        tabindex="1"
      >
        {{ $t("buttons.close") }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions } from "pinia";
import { useLayoutStore } from "@/stores/layout";

export default {
  name: "help",
  data() {
    return {
      activeTab: "start",
      openFaq: null,
      tabs: [
        { id: "start", label: "Getting Started", icon: "play_arrow" },
        { id: "keys", label: "Shortcuts", icon: "keyboard" },
        { id: "features", label: "Features", icon: "widgets" },
        { id: "trouble", label: "Troubleshooting", icon: "build" },
      ],
    };
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    toggleFaq(id) {
      this.openFaq = this.openFaq === id ? null : id;
    },
  },
};
</script>

<style scoped>
.help-section {
  margin-bottom: 1.25em;
}

.help-heading {
  display: flex;
  align-items: center;
  gap: 0.5em;
  font-size: 0.95em;
  font-weight: 600;
  margin: 0 0 0.5em 0;
  color: var(--textPrimary);
}

.help-heading .material-icons {
  font-size: 1.2em;
  color: var(--primary);
}

.help-list p {
  margin: 0.35em 0;
  padding-left: 0.5em;
  font-size: 0.85em;
  color: var(--textSecondary);
  line-height: 1.5;
}

.help-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85em;
}

.help-table tr {
  border-bottom: 1px solid var(--borderPrimary);
}

.help-table td {
  padding: 0.4em 0.5em;
  color: var(--textSecondary);
}

.key-cell {
  width: 160px;
  white-space: nowrap;
}

kbd {
  display: inline-block;
  padding: 0.15em 0.5em;
  font-size: 0.85em;
  font-family: inherit;
  background: var(--surfaceSecondary);
  border: 1px solid var(--borderPrimary);
  border-radius: 4px;
  color: var(--textPrimary);
}

.help-faq {
  border: 1px solid var(--borderPrimary);
  border-radius: 6px;
  margin-bottom: 0.5em;
  overflow: hidden;
}

.help-faq-toggle {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.6em 0.75em;
  background: var(--surfaceSecondary);
  border: none;
  cursor: pointer;
  font-size: 0.85em;
  font-weight: 500;
  color: var(--textPrimary);
  text-align: left;
}

.help-faq-toggle:hover {
  opacity: 0.8;
}

.help-faq-content {
  padding: 0.5em 0.75em;
}

.help-faq-content p {
  margin: 0.3em 0;
  font-size: 0.83em;
  color: var(--textSecondary);
  line-height: 1.5;
}

.help-faq-content code {
  background: var(--surfaceSecondary);
  padding: 0.1em 0.4em;
  border-radius: 3px;
  font-size: 0.9em;
}
</style>
