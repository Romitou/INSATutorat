<template>
</template>

<script setup>
import {onMounted} from 'vue'
import {useRouter} from 'vue-router'
import {useToast} from "vue-toastification";

definePageMeta({
  layout: 'public'
})

const router = useRouter()
const userStore = useUserStore();

onMounted(async () => {
  const route = useRoute()
  const ticket = route.query?.ticket

  if (ticket) {
    try {
      const authResult = await useApiFetch(`/auth/validate?ticket=${ticket}`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'}
      })

      if (authResult.ok) {
        await userStore.fetchUser();
        await router.push('/')
        useToast().success('Connexion réussie')
      } else {
        useToast().error('Impossible de se connecter. Veuillez réessayer.')
      }

    } catch (error) {
      useToast().error('Impossible de se connecter. Veuillez réessayer.')
    }
  } else {
    useToast().error('Jeton manquant. Veuillez vérifier le lien dans votre email.')
  }
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
