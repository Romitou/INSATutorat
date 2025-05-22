<script setup lang="ts">
import { ref, computed } from 'vue'
import type {User} from "~/types/api";

const props = defineProps<{
  show: boolean
  subjectId: number
  assignedTutors: User[]
  allTutors: User[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', user: User, maxTutees: number): void
}>()

const selectedTutorId = ref<number | null>(null)
const maxTutees = ref<number>(3)

const availableTutors = computed(() => {
  if (!props.show) return []
  if (!props.allTutors || !props.assignedTutors) return []
  return props.allTutors.filter(
      tutor => !props.assignedTutors.some(assigned => assigned.id === tutor.id)
  )
})

const submit = () => {
  if (!selectedTutorId.value) return

  const tutor = availableTutors.value.find(t => t.id === selectedTutorId.value)
  if (tutor) {
    emit('submit', tutor, maxTutees.value)
    emit('close')
  }
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-lg w-full max-w-md shadow-xl p-6 relative">
      <h2 class="text-lg font-semibold text-zinc-800 mb-4">{{ $t('addTutorTitle') }}</h2>

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-zinc-700">{{ $t('tutorLabel') }}</label>
          <select
              v-model="selectedTutorId"
              class="w-full mt-1 border border-zinc-300 rounded-md p-2"
          >
            <option :value="null" disabled>{{ $t('selectTutorPlaceholder') }}</option>
            <option v-for="user in availableTutors" :key="user.id" :value="user.id">
              {{ user.lastName }} {{ user.firstName }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-zinc-700">{{ $t('maxTuteesLabel') }}</label>
          <input
              type="number"
              v-model.number="maxTutees"
              class="w-full mt-1 border border-zinc-300 rounded-md p-2"
              min="1"
          />
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
            :disabled="!selectedTutorId"
            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm"
        >
          {{ $t("addTutorTitle") }}
        </button>
      </div>
    </div>
  </div>
</template>
