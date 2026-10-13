<template>
  <div class="rules">
    <div v-for="(rule, index) in rules" :key="index" class="rules__row">
      <input
        @keypress.enter.prevent
        class="input"
        type="text"
        v-if="rule.regex"
        v-model="rule.regexp.raw"
        :placeholder="$t('settings.insertRegex')"
      />
      <input
        @keypress.enter.prevent
        class="input"
        type="text"
        v-else
        v-model="rule.path"
        :placeholder="$t('settings.insertPath')"
      />

      <label class="rules__option">
        <input type="checkbox" class="switch" v-model="rule.regex" />
        Regex
      </label>
      <label class="rules__option">
        <input type="checkbox" class="switch" v-model="rule.allow" />
        Allow
      </label>

      <button
        type="button"
        class="action"
        @click="remove($event, index)"
        :aria-label="$t('buttons.delete')"
        :title="$t('buttons.delete')"
      >
        <i class="material-icons">delete</i>
      </button>
    </div>

    <button
      type="button"
      class="button button--flat rules__add"
      @click="create"
    >
      <i class="material-icons" aria-hidden="true">add</i>
      {{ $t("buttons.new") }}
    </button>
  </div>
</template>

<script>
export default {
  name: "rules-textarea",
  props: ["rules"],
  methods: {
    remove(event, index) {
      event.preventDefault();
      const rules = [...this.rules];
      rules.splice(index, 1);
      this.$emit("update:rules", [...rules]);
    },
    create(event) {
      event.preventDefault();

      this.$emit("update:rules", [
        ...this.rules,
        {
          allow: true,
          path: "",
          regex: false,
          regexp: {
            raw: "",
          },
        },
      ]);
    },
  },
};
</script>
