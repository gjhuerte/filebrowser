import { theme } from "./constants";
import "ace-builds";
import { themesByName } from "ace-builds/src-noconflict/ext-themelist";

export const getTheme = (): UserTheme => {
  return (document.documentElement.className as UserTheme) || theme;
};

export const setTheme = (theme: UserTheme) => {
  const html = document.documentElement;
  if (!theme) {
    html.className = getMediaPreference();
  } else {
    html.className = theme;
  }
};

export const toggleTheme = (): void => {
  const activeTheme = getTheme();
  if (activeTheme === "light") {
    setTheme("dark");
  } else {
    setTheme("light");
  }
};

export const getMediaPreference = (): UserTheme => {
  const hasDarkPreference = window.matchMedia(
    "(prefers-color-scheme: dark)"
  ).matches;
  if (hasDarkPreference) {
    return "dark";
  } else {
    return "light";
  }
};

/*
 * Per-browser theme preference ("light", "dark" or "" for system).
 * A saved choice overrides the instance default from Global Settings;
 * "system" follows the OS setting live.
 */
const THEME_STORAGE_KEY = "theme";
const darkQuery = window.matchMedia("(prefers-color-scheme: dark)");

const followSystem = () => {
  document.documentElement.className = getMediaPreference();
};

export const getThemePreference = (): UserTheme => {
  try {
    const saved = localStorage.getItem(THEME_STORAGE_KEY);
    if (saved === "light" || saved === "dark") return saved;
    if (saved === "system") return "";
  } catch {
    // storage unavailable (private mode, blocked cookies): use the default
  }
  return theme;
};

export const applyThemePreference = (preference: UserTheme) => {
  darkQuery.removeEventListener("change", followSystem);
  setTheme(preference);
  if (!preference) {
    darkQuery.addEventListener("change", followSystem);
  }
};

export const setThemePreference = (preference: UserTheme) => {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, preference || "system");
  } catch {
    // not persisted, but still applied for this page view
  }
  applyThemePreference(preference);
};

export const getEditorTheme = (themeName: string) => {
  if (!themeName.startsWith("ace/theme/")) {
    themeName = `ace/theme/${themeName}`;
  }
  const themeKey = themeName.replace("ace/theme/", "");
  if (themesByName[themeKey] !== undefined) {
    return themeName;
  } else if (getTheme() === "dark") {
    return "ace/theme/twilight";
  } else {
    return "ace/theme/chrome";
  }
};
