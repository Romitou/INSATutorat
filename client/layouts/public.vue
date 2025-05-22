<template>
  <div class="flex flex-col h-screen w-full">
    <nav class="bg-[#e61115] w-full">
      <div class="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
        <div class="relative flex h-16 items-center justify-between">
          <div class="flex flex-1 items-center justify-center sm:items-stretch sm:justify-start">
            <div class="flex shrink-0 items-center">
              <img class="h-12 w-auto" src="https://cas.insa-rouen.fr/cas/images/cas-logo.png" alt="INSA Rouen Normandie" />
            </div>
            <div class="hidden h-16 sm:ml-6 sm:flex items-center justify-center">
            </div>
          </div>
          <div class="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
          </div>
        </div>
      </div>
    </nav>
    <slot />
  </div>
</template>

<script setup>
</script>

<style scoped>

</style>