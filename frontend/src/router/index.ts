import type { RouteLocation } from "vue-router";
import { createRouter, createWebHistory } from "vue-router";
import Login from "@/views/Login.vue";
import Layout from "@/views/Layout.vue";
import Files from "@/views/Files.vue";
import Home from "@/views/Home.vue";
import Share from "@/views/Share.vue";
import Users from "@/views/settings/Users.vue";
import User from "@/views/settings/User.vue";
import Settings from "@/views/Settings.vue";
import ProfileInfo from "@/views/settings/ProfileInfo.vue";
import ChangePassword from "@/views/settings/ChangePassword.vue";
import SystemSettings from "@/views/settings/System.vue";
import Shares from "@/views/settings/Shares.vue";
import GlobalSettingsWrapper from "@/views/GlobalSettings.vue";
import GeneralSettings from "@/views/admin/GeneralSettings.vue";
import UserDefaultsSettings from "@/views/admin/UserDefaultsSettings.vue";
import BrandingSettings from "@/views/admin/BrandingSettings.vue";
import CommandsSettings from "@/views/admin/CommandsSettings.vue";
import UserManagement from "@/views/admin/UserManagement.vue";
import Errors from "@/views/Errors.vue";
import { useAuthStore } from "@/stores/auth";
import { baseURL, name } from "@/utils/constants";
import i18n from "@/i18n";
import { recaptcha, loginPage } from "@/utils/constants";
import { login, validateLogin } from "@/utils/auth";

const titles = {
  Login: "sidebar.login",
  Share: "buttons.share",
  Home: "sidebar.home",
  Files: "files.files",
  Settings: "sidebar.settings",
  ProfileInfo: "settings.profileInfo",
  ChangePassword: "settings.changePassword",
  SystemSettings: "settings.system",
  Shares: "settings.shareManagement",
  GlobalSettings: "settings.globalSettings",
  GeneralSettings: "settings.general",
  UserDefaultsSettings: "settings.userDefaults",
  BrandingSettings: "settings.branding",
  CommandsSettings: "settings.commandRunner",
  UserManagement: "settings.userManagement",
  Users: "settings.users",
  UserNew: "settings.newUser",
  User: "settings.user",
  Forbidden: "errors.forbidden",
  NotFound: "errors.notFound",
  InternalServerError: "errors.internal",
};

const routes = [
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/share",
    component: Layout,
    children: [
      {
        path: ":path*",
        name: "Share",
        component: Share,
      },
    ],
  },
  {
    path: "/home",
    component: Layout,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        name: "Home",
        component: Home,
      },
    ],
  },
  {
    path: "/files",
    component: Layout,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: ":path*",
        name: "Files",
        component: Files,
      },
    ],
  },
  {
    path: "/settings",
    component: Layout,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        name: "Settings",
        component: Settings,
        redirect: {
          path: "/settings/profile",
        },
        children: [
          {
            path: "profile",
            name: "ProfileInfo",
            component: ProfileInfo,
          },
          {
            path: "password",
            name: "ChangePassword",
            component: ChangePassword,
          },
          {
            path: "system",
            name: "SystemSettings",
            component: SystemSettings,
          },
          {
            path: "shares",
            name: "Shares",
            component: Shares,
          },
        ],
      },
    ],
  },
  {
    path: "/admin/settings",
    component: Layout,
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
    },
    children: [
      {
        path: "",
        name: "GlobalSettings",
        component: GlobalSettingsWrapper,
        redirect: {
          path: "/admin/settings/general",
        },
        children: [
          {
            path: "general",
            name: "GeneralSettings",
            component: GeneralSettings,
          },
          {
            path: "user-defaults",
            name: "UserDefaultsSettings",
            component: UserDefaultsSettings,
          },
          {
            path: "branding",
            name: "BrandingSettings",
            component: BrandingSettings,
          },
          {
            path: "commands",
            name: "CommandsSettings",
            component: CommandsSettings,
          },
        ],
      },
    ],
  },
  {
    path: "/admin/users",
    component: Layout,
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
    },
    children: [
      {
        path: "",
        name: "UserManagement",
        component: UserManagement,
        children: [
          {
            path: "",
            name: "Users",
            component: Users,
          },
          {
            path: "new",
            name: "UserNew",
            component: User,
          },
          {
            path: ":id",
            name: "User",
            component: User,
          },
        ],
      },
    ],
  },
  {
    path: "/403",
    name: "Forbidden",
    component: Errors,
    props: {
      errorCode: 403,
      showHeader: true,
    },
  },
  {
    path: "/404",
    name: "NotFound",
    component: Errors,
    props: {
      errorCode: 404,
      showHeader: true,
    },
  },
  {
    path: "/500",
    name: "InternalServerError",
    component: Errors,
    props: {
      errorCode: 500,
      showHeader: true,
    },
  },
  {
    // Legacy redirects for old settings routes
    path: "/settings/global",
    redirect: "/admin/settings/general",
  },
  {
    path: "/settings/users",
    redirect: "/admin/users",
  },
  {
    path: "/:catchAll(.*)*",
    redirect: (to: RouteLocation) => {
      const catchAll = to.params.catchAll;
      if (!catchAll) return "/home";
      return `/files/${Array.isArray(catchAll) ? catchAll.join("/") : catchAll}`;
    },
  },
];

async function initAuth() {
  if (loginPage) {
    await validateLogin();
  } else {
    await login("", "", "");
  }

  if (recaptcha) {
    await new Promise<void>((resolve) => {
      const check = () => {
        if (typeof window.grecaptcha === "undefined") {
          setTimeout(check, 100);
        } else {
          resolve();
        }
      };

      check();
    });
  }
}

const router = createRouter({
  history: createWebHistory(baseURL),
  routes,
});

router.beforeResolve(async (to, from) => {
  const title = i18n.global.t(titles[to.name as keyof typeof titles]);
  document.title = title + " - " + name;

  const authStore = useAuthStore();

  // this will only be null on first route
  if (from.name == null) {
    try {
      await initAuth();
    } catch (error) {
      console.error(error);
    }
  }

  if (to.path.endsWith("/login") && authStore.isLoggedIn) {
    return { path: "/home" };
  }

  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!authStore.isLoggedIn) {
      return {
        path: "/login",
        query: { redirect: to.fullPath },
      };
    }

    if (to.matched.some((record) => record.meta.requiresAdmin)) {
      if (authStore.user === null || !authStore.user.perm.admin) {
        return { path: "/403" };
      }
    }
  }

  return;
});

export { router, router as default };
