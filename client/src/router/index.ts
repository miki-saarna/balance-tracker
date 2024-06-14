import { createRouter, createWebHistory } from "vue-router";

// @ts-ignore
import Home from "../components/Home.vue";
// @ts-ignore
import Balances from "../components/Balances.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/balances", component: Balances },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
