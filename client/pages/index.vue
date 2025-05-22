<script lang="ts" setup>
import {useToast} from "vue-toastification";

definePageMeta({
  layout: 'public'
})

const router = useRouter();
const userStore = useUserStore();
onMounted(async () => {
  await userStore.fetchUser();
  if (userStore.user) {
    if (userStore.user.isTutor) {
      await router.push('/tutor');
      useToast().info('Redirection vers l\'espace tuteur');
    } else if (userStore.user.isTutee) {
      await router.push('/tutee');
      useToast().info('Redirection vers l\'espace tutoré');
    }
  } else {
    await router.push('/send-link');
  }
})
</script>

<template>

</template>

<style scoped>

</style>