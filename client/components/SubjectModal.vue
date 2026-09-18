<script setup lang="ts">
import {ref, watch, computed} from 'vue'
import type {Subject} from "~/types/api";

const props = defineProps<{
  show: boolean
  initialData?: Subject | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', subject: Subject): void
}>()

const semester = ref(1)
const shortName = ref('')
const name = ref('')
const errors = ref<{ [key: string]: string }>({})

watch(
    () => props.initialData,
    (data) => {
      if (data) {
        semester.value = data.semester
        shortName.value = data.shortName
        name.value = data.name
      } else {
        semester.value = 1
        shortName.value = ''
        name.value = ''
      }

      errors.value = {}
    },
    {immediate: true}
)

const isValid = computed(() => {
  errors.value = {}
  if (!shortName.value.trim()) {
    errors.value.shortName = 'L\'abréviation est requise.'
  }
  if (!name.value.trim()) {
    errors.value.name = 'Le nom est requis.'
  }

  return Object.keys(errors.value).length === 0
})

function submit() {
  if (!isValid.value) return

  emit('submit', {
    id: props.initialData?.id,
    semester: semester.value,
    shortName: shortName.value.trim(),
    name: name.value.trim(),
  } as Subject)

  emit('close')
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-lg w-full max-w-md shadow-xl p-6 relative">
      <h2 class="text-lg font-semibold text-zinc-800 mb-4">
        {{ props.initialData?.id ? $t("edit") : $t("newFeminine") }} matière
      </h2>

      <div class="space-y-4">

        <div>
          <label class="block text-sm font-medium text-zinc-700">Semestre</label>
          <select
              v-model="semester"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="1">1</option>
            <option value="2">2</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-zinc-700">Abréviation</label>
          <input
              type="text"
              v-model="shortName"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.shortName" class="text-sm text-red-500 mt-1">{{ errors.shortName }}</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-zinc-700">Nom</label>
          <input
              type="text"
              v-model="name"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.name" class="text-sm text-red-500 mt-1">{{ errors.name }}</p>
        </div>
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
            :disabled="!isValid"
            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ $t("save") }}
        </button>
      </div>
    </div>
  </div>
</template>
