<script setup lang="ts">
import {ref, watch, computed} from 'vue'
import type {User} from "~/types/api";

const props = defineProps<{
  show: boolean
  user: User | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', payload: { id: number, isAdmin: boolean, isTutor: boolean, isTutee: boolean }): void
}>()

const userStore = useUserStore()

const isAdmin = ref(false)
const isTutor = ref(false)
const isTutee = ref(false)

const isSelf = computed(() => props.user?.id === userStore.user?.id)

watch(
    () => props.user,
    (data) => {
      isAdmin.value = data?.isAdmin ?? false
      isTutor.value = data?.isTutor ?? false
      isTutee.value = data?.isTutee ?? false
    },
    {immediate: true}
)

function submit() {
  if (!props.user) return

  emit('submit', {
    id: props.user.id,
    isAdmin: isAdmin.value,
    isTutor: isTutor.value,
    isTutee: isTutee.value,
  })

  emit('close')
}
</script>

<template>
  <div v-if="show && user" class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-lg w-full max-w-md shadow-xl p-6 relative">
      <h2 class="text-lg font-semibold text-zinc-800 mb-1">
        Modifier les rôles
      </h2>
      <p class="text-sm text-zinc-500 mb-4">{{ user.firstName }} {{ user.lastName }} — {{ user.mail }}</p>

      <div class="space-y-3">
        <label class="flex items-center gap-2 text-sm text-zinc-700">
          <input type="checkbox" v-model="isAdmin" :disabled="isSelf" class="rounded border-zinc-300"/>
          Administrateur
        </label>
        <p v-if="isSelf" class="text-xs text-zinc-400 -mt-2">
          Vous ne pouvez pas modifier vos propres droits d'administrateur.
        </p>

        <label class="flex items-center gap-2 text-sm text-zinc-700">
          <input type="checkbox" v-model="isTutor" class="rounded border-zinc-300"/>
          Tuteur
        </label>

        <label class="flex items-center gap-2 text-sm text-zinc-700">
          <input type="checkbox" v-model="isTutee" class="rounded border-zinc-300"/>
          Tutee
        </label>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <button
            @click="$emit('close')"
            class="px-4 py-2 rounded-md bg-zinc-200 hover:bg-zinc-300 text-sm text-zinc-700 transition"
        >
          {{ $t("cancel") }}
        </button>
        <button
            @click="submit"
            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm transition"
        >
          {{ $t("save") }}
        </button>
      </div>
    </div>
  </div>
</template>
