<script lang="ts" setup>
import {CheckCircleIcon} from "@heroicons/vue/24/solid";

definePageMeta({
  layout: 'public'
})

const userStore = useUserStore();

const logout = () => {
  void userStore.logout()
}

onMounted(logout);
</script>

<template>
  <div class="h-full bg-gray-50 flex flex-col">
    <div class="flex-1 flex flex-col justify-center items-center px-4 sm:px-6 lg:px-8">
      <div class="w-full max-w-md space-y-8">
        <div class="text-center space-y-2">
          <CheckCircleIcon class="mx-auto h-16 w-16 text-green-500"/>
          <h3 class="text-xl font-semibold text-gray-800">Déconnexion effectuée</h3>
          <p class="text-gray-600">Vous avez été déconnecté(e) de la plateforme de tutorat.</p>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
@keyframes pop {
  0% {
    transform: scale(0.8);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

svg {
  animation: pop 0.5s ease;
}
</style>
