<template>
  <div class="flex flex-col min-h-screen w-full bg-zinc-50">
    <Navbar />
    <slot />
  </div>
</template>

<script setup>
import Navbar from "~/components/Navbar.vue";

const userStore = useUserStore();
onMounted(() => {
  userStore.fetchUser();
})
</script>

<style scoped>

</style>