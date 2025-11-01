<template>
  <div class="h-full bg-gray-50 flex flex-col">
    <div class="flex-1 flex flex-col justify-center items-center px-4 sm:px-6 lg:px-8">
      <div class="w-full max-w-md space-y-8">
        <div>
          <h2 class="text-center text-3xl font-bold text-gray-900">Connexion requise</h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            Cette plateforme utilise le système d'authentification central de l'INSA Rouen. Veuillez vous connecter avec via votre compte INSA pour accéder à cette plateforme.
          </p>
        </div>

        <form class="mt-8 space-y-6">
          <div>
            <a
                class="group relative flex w-full justify-center rounded-md border border-transparent bg-[#e61115] py-2 px-4 text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                href="https://cas.insa-rouen.fr/cas/login?service=https%3A%2F%2Fstpi-tutorat.insa-rouen.fr%2Fvalidate"
            >
              Se connecter avec l'INSA Rouen
            </a>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue'
import {CheckCircleIcon} from '@heroicons/vue/24/solid'
import {useToast} from "vue-toastification";

definePageMeta({
  layout: 'public'
})

</script>

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
