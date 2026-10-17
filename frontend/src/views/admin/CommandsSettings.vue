<template>
  <div class="row">
    <div class="column">
      <form class="card" @submit.prevent="$emit('save')">
        <div class="card-title">
          <h2>{{ t("settings.commandRunner") }}</h2>
        </div>

        <div class="card-content">
          <i18n-t
            keypath="settings.commandRunnerHelp"
            tag="p"
            class="small"
            scope="global"
          >
            <code>FILE</code>
            <code>SCOPE</code>
            <a
              class="link"
              target="_blank"
              href="/docs/command-execution.md#hook-runner"
              >{{ t("settings.documentation") }}</a
            >
          </i18n-t>

          <p class="small">{{ t("settings.commandRunnerDescription") }}</p>

          <div
            v-for="(command, key) in settings.commands"
            :key="key"
            class="collapsible"
          >
            <input :id="key" type="checkbox" />
            <label :for="key">
              <p>{{ capitalize(key) }}</p>
              <i class="material-icons">arrow_drop_down</i>
            </label>
            <div class="collapse">
              <textarea
                class="input input--block input--textarea"
                v-model.trim="commandObject[key]"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="card-action">
          <input class="button" type="submit" :value="t('buttons.update')" />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t } = useI18n();

defineProps<{
  settings: ISettings;
  commandObject: { [key: string]: string[] | string };
}>();

defineEmits<{
  save: [];
}>();

const capitalize = (name: string, where: string | RegExp = "_") => {
  if (where === "caps") where = /(?=[A-Z])/;
  const split = name.split(where);
  name = "";

  for (let i = 0; i < split.length; i++) {
    name += split[i].charAt(0).toUpperCase() + split[i].slice(1) + " ";
  }

  return name.slice(0, -1);
};
</script>
