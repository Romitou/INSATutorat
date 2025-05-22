<script setup lang="ts">
import { ref, computed } from 'vue'
import type {User} from "~/types/api";

const props = defineProps<{
  show: boolean
  subjectId: number
  assignedTutees: User[]
  allTutees: User[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', user: User): void
}>()

const selectedTuteeId = ref<number | null>(null)

const availableTutees = computed(() => {
  if (!props.show) return []
  if (!props.allTutees || !props.assignedTutees) return []
  return props.allTutees.filter(
      tutee => !props.assignedTutees.some(assigned => assigned.id === tutee.id)
  )
})

const submit = () => {
  if (!selectedTuteeId.value) return

  const tutee = availableTutees.value.find(t => t.id === selectedTuteeId.value)
  if (tutee) {
    emit('submit', tutee)
    emit('close')
  }
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-lg w-full max-w-md shadow-xl p-6 relative">
      <h2 class="text-lg font-semibold text-zinc-800 mb-4">{{ $t('addTuteeTitle') }}</h2>

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-zinc-700">{{ $t('tuteeLabel') }}</label>
          <select
              v-model="selectedTuteeId"
              class="w-full mt-1 border border-zinc-300 rounded-md p-2"
          >
            <option :value="null" disabled>{{ $t('selectTuteePlaceholder') }}</option>
            <option v-for="user in availableTutees" :key="user.id" :value="user.id">
              {{ user.lastName }} {{ user.firstName }}
            </option>
          </select>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <button
            @click="$emit('close')"
            class="px-4 py-2 rounded-md bg-zinc-200 hover:bg-zinc-300 text-sm text-zinc-700"
        >
          {{ $t("cancel") }}
        </button>
        <button
            @click="submit"
            :disabled="!selectedTuteeId"
            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm"
        >
          {{ $t("addTuteeTitle") }}
        </button>
      </div>
    </div>
  </div>
</template>
