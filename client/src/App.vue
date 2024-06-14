<template>
  <div class="flex flex-col">
    <div class="px-8">
      <nav className="flex space-x-8 py-4 border-b border-gray-300">
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/balances">Balances</RouterLink>
      </nav>
    </div>
    <main className="grow px-8 py-4">
      <RouterView />
    </main>
  </div>
</template>

<script setup></script>
