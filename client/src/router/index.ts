import { createRouter, createWebHistory } from "vue-router";

import Home from "../components/Home.vue";
import Balances from "../components/Balances.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/balances", component: Balances },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
