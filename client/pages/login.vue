<template>
  <div class="h-full bg-gray-50 flex flex-col">
    <div class="flex-1 flex flex-col justify-center items-center px-4 sm:px-6 lg:px-8">
      <div v-if="!loading" class="w-full max-w-md space-y-8">
        <div>
          <h2 class="text-center text-3xl font-bold text-gray-900">
            {{ isMagicLink ? 'Connexion à la plateforme' : 'Connexion requise' }}
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            {{
              isMagicLink
                  ? 'Saisissez votre adresse mail INSA afin de recevoir un lien de connexion.'
                  : "Cette plateforme utilise le système d'authentification central de l'INSA Rouen. Veuillez vous connecter avec via votre compte INSA pour accéder à cette plateforme."
            }}
          </p>
        </div>

        <div v-if="isMagicLink">
          <form v-if="!success" class="mt-8 space-y-6" @submit.prevent="submitMagicLink">
            <div class="rounded-md shadow-sm -space-y-px">
              <div>
                <label class="sr-only" for="email-address">Adresse email INSA</label>
                <input
                    id="email-address"
                    v-model="email"
                    autocomplete="email"
                    class="relative block w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-red-500 focus:outline-none focus:ring-red-500 sm:text-sm"
                    name="email"
                    placeholder="Votre adresse email INSA"
                    required
                    type="email"
                />
              </div>
            </div>

            <div>
              <button
                  class="group relative flex w-full justify-center rounded-md border border-transparent bg-[#e61115] py-2 px-4 text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                  type="submit"
              >
                Envoyer le lien
              </button>
            </div>
          </form>

          <div v-else class="text-center space-y-4">
            <CheckCircleIcon class="mx-auto h-16 w-16 text-green-500"/>
            <h3 class="text-xl font-semibold text-gray-800">Lien envoyé !</h3>
            <p class="text-gray-600">Vérifiez votre boîte mail pour vous connecter.</p>
          </div>
        </div>

        <div v-else class="mt-8 space-y-6">
          <div>
            <a
                class="group relative flex w-full justify-center rounded-md border border-transparent bg-[#e61115] py-2 px-4 text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                :href="casLoginUrl"
            >
              Se connecter avec l'INSA Rouen
            </a>
          </div>
        </div>
      </div>
      <div v-else class="text-center">
        <p class="text-gray-600">Chargement...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref, computed} from 'vue'
import {CheckCircleIcon} from '@heroicons/vue/24/solid'
import {useToast} from "vue-toastification";
import {useRouter, useRoute} from 'vue-router'

definePageMeta({
  layout: 'public'
})

const router = useRouter()
const route = useRoute()
const userStore = useUserStore();
const toast = useToast()

const loading = ref(true)
const authMethod = ref('CAS')
const casUrl = ref('')
const serviceUrl = ref('')

// Magic Link state
const email = ref('')
const success = ref(false)

const isMagicLink = computed(() => authMethod.value === 'MAGIC_LINK')
const casLoginUrl = computed(() => {
  return `${casUrl.value}/login?service=${encodeURIComponent(serviceUrl.value)}`
})

async function fetchConfig() {
  try {
    const res = await useApiFetch('/auth/config')
    if (res.ok) {
      const data = await res.json()
      authMethod.value = data.authMethod
      casUrl.value = data.casUrl
      serviceUrl.value = data.serviceUrl || 'https://stpi-tutorat.insa-rouen.fr/validate'
    } else {
      console.error("Config fetch failed", res.status)
    }
  } catch (e) {
    console.error("Failed to fetch auth config", e)
  } finally {
    loading.value = false
  }
}

async function handleMagicLinkLogin(token) {
  try {
    const authResult = await useApiFetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify({token}),
      headers: {'Content-Type': 'application/json'}
    })

    if (!authResult.ok) {
      const err = await authResult.json().catch(() => ({}))
      throw new Error(err.error || 'Login failed')
    }

    await userStore.fetchUser();
    await router.push('/')
    toast.success('Connexion réussie')

  } catch (error) {
    toast.error('Impossible de se connecter. Veuillez réessayer.')
    await fetchConfig()
  }
}

async function submitMagicLink() {
  try {
    const res = await useApiFetch('/auth/send-link', {
      method: 'POST',
      body: JSON.stringify({mail: email.value}),
      headers: {'Content-Type': 'application/json'}
    })

    if (!res.ok) {
      throw new Error('Send link failed')
    }

    success.value = true
  } catch (error) {
    toast.error('Erreur lors de l\'envoi du lien, veuillez réessayer.')
  }
}

onMounted(async () => {
  const token = route.query?.token

  if (token) {
    // Magic Link Callback
    await handleMagicLinkLogin(token)
  } else {
    // Standard Load
    await fetchConfig()
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
