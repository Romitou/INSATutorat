<template>
  <div class="flex flex-col min-h-screen w-full bg-zinc-50">
    <Navbar />
    <slot />
    <footer class="w-full text-center py-4 text-xs text-gray-500 bg-gray-50 border-t mt-auto">
      <a href="https://gitlab.insa-rouen.fr/rleprevost/stpi-tutorat" target="_blank" rel="noopener" class="text-gray-600 hover:underline">gitlab.insa-rouen.fr/rleprevost/stpi-tutorat</a>
      - Déploiement {{ buildId }}
    </footer>
  </div>
</template>

<script setup>
import Navbar from "~/components/Navbar.vue";

const buildId = useRuntimeConfig().app?.buildId;

const userStore = useUserStore();
onMounted(() => {
  userStore.fetchUser();
})
</script>

<style scoped>

</style>