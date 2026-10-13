<template>
  <div>
    <h3>{{ $t("settings.permissions") }}</h3>
    <p class="small">{{ $t("settings.permissionsHelp") }}</p>

    <div class="setting-list">
      <label class="setting-toggle">
        <span>{{ $t("settings.administrator") }}</span>
        <input type="checkbox" class="switch" v-model="admin" />
      </label>
      <label class="setting-toggle">
        <span>{{ $t("settings.perm.create") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin"
          v-model="perm.create"
        />
      </label>
      <label class="setting-toggle">
        <span>{{ $t("settings.perm.delete") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin"
          v-model="perm.delete"
        />
      </label>
      <label class="setting-toggle">
        <span>{{ $t("settings.perm.download") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin || perm.share"
          v-model="perm.download"
        />
      </label>
      <label class="setting-toggle">
        <span>{{ $t("settings.perm.modify") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin"
          v-model="perm.modify"
        />
      </label>
      <label class="setting-toggle" v-if="isExecEnabled">
        <span>{{ $t("settings.perm.execute") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin"
          v-model="perm.execute"
        />
      </label>
      <label class="setting-toggle">
        <span>{{ $t("settings.perm.rename") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin"
          v-model="perm.rename"
        />
      </label>
      <label class="setting-toggle">
        <span>{{ $t("settings.perm.share") }}</span>
        <input
          type="checkbox"
          class="switch"
          :disabled="admin"
          v-model="perm.share"
        />
      </label>
    </div>
  </div>
</template>

<script>
import { enableExec } from "@/utils/constants";
export default {
  name: "permissions",
  props: ["perm"],
  computed: {
    admin: {
      get() {
        return this.perm.admin;
      },
      set(value) {
        if (value) {
          for (const key in this.perm) {
            this.perm[key] = true;
          }
        }

        this.perm.admin = value;
      },
    },
    isExecEnabled: () => enableExec,
  },
  watch: {
    perm: {
      deep: true,
      handler() {
        if (this.perm.share === true) {
          this.perm.download = true;
        }
      },
    },
  },
};
</script>
