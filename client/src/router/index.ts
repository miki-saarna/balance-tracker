import { createRouter, createWebHistory } from "vue-router";

// @ts-ignore
import Home from "../views/Home.vue";
// @ts-ignore
import Balances from "../views/Balances.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/balances", component: Balances },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
