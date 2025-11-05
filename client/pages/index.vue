<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useToast } from "vue-toastification";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
import { ArrowPathIcon } from "@heroicons/vue/24/outline";

definePageMeta({
  layout: 'public'
})

const router = useRouter();
const userStore = useUserStore();
const toast = useToast();
const loading = ref(true);

onMounted(async () => {
  try {
    await userStore.fetchUser();
    if (!userStore.user) {
      toast.info("Votre session a expiré. Veuillez vous reconnecter.");
      return router.replace('/login');
    }

    if (userStore.user.isTutor) {
      toast.info("Redirection vers l'espace tuteur");
      router.replace('/tutor');
    } else if (userStore.user.isTutee) {
      toast.info("Redirection vers l'espace tutoré");
      router.replace('/tutee');
    } else if (userStore.user.isAdmin) {
      toast.info("Redirection vers l'espace admin");
      router.replace('/admin');
    } else {
      toast.warning("Vous n'avez aucun rôle attribué. Contactez un administrateur.");
    }
  } catch (error) {
    console.error("Erreur lors de la récupération de l'utilisateur :", error);
    toast.error("Impossible de récupérer vos informations. Veuillez réessayer.");
    router.replace('/login');
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div v-if="loading" class="flex justify-center items-center h-screen">
    <ArrowPathIcon class="mx-auto h-16 w-16 text-gray-500 animate-spin"/>
  </div>
</template>

<style scoped>
</style>
