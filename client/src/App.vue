<template>
  <div class="flex flex-col">
    <div class="px-8">
      <nav className="flex space-x-8 py-4 border-b border-gray-300">
        <RouterLink
          v-for="route of routes"
          :to="route.path"
          :class="[
            route.path === $route.path && 'bg-blue-300',
            'py-1 px-2 rounded select-none',
          ]"
        >
          {{ route.name }}
        </RouterLink>
      </nav>
    </div>
    <main className="grow px-8 py-4">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
const routes = [
  { path: "/", name: "Home" },
  { path: "/balances", name: "Balances" },
];
</script>
