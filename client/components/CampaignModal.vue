<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { DateTime } from 'luxon'
import type {Campaign} from "~/types/api";

const props = defineProps<{
  show: boolean
  initialData?: Campaign
}>()

const { t } = useI18n()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', campaign: Campaign): void
}>()

const schoolYear = ref('')
const semester = ref(1)
const startDate = ref('')
const endDate = ref('')
const registrationStartDate = ref('')
const registrationEndDate = ref('')
const registrationStatus = ref<'OPEN' | 'CLOSED'>('OPEN')
const errors = ref<{ [key: string]: string }>({})

watch(
    () => props.initialData,
    (data) => {
      if (data) {
        schoolYear.value = data.schoolYear
        semester.value = data.semester
        startDate.value = DateTime.fromISO(data.startDate).toISODate()
        endDate.value = DateTime.fromISO(data.endDate).toISODate()
        registrationStartDate.value = DateTime.fromISO(data.registrationStartDate).toISODate()
        registrationEndDate.value = DateTime.fromISO(data.registrationEndDate).toISODate()
        registrationStatus.value = data.registrationStatus
      } else {
        schoolYear.value = ''
        semester.value = 1
        startDate.value = ''
        endDate.value = ''
        registrationStartDate.value = ''
        registrationEndDate.value = ''
        registrationStatus.value = 'OPEN'
      }

      errors.value = {}
    },
    { immediate: true }
)

const isValid = computed(() => {
  errors.value = {}
  if (!schoolYear.value.trim()) {
    errors.value.schoolYear = 'L\'année scolaire est requise.'
  }
  if (!startDate.value) errors.value.startDate = 'La date de début est requise.'
  if (!endDate.value) errors.value.endDate = 'La date de fin est requise.'
  if (!registrationStartDate.value) errors.value.registrationStartDate = 'La date de début des inscriptions est requise.'
  if (!registrationEndDate.value) errors.value.registrationEndDate = 'La date de fin des inscriptions est requise.'

  return Object.keys(errors.value).length === 0
})

function submit() {
  if (!isValid.value) return

  emit('submit', {
    id: props.initialData?.id,
    schoolYear: schoolYear.value.trim(),
    semester: semester.value,
    startDate: DateTime.fromISO(startDate.value).toISO(),
    endDate: DateTime.fromISO(endDate.value).toISO(),
    registrationStatus: registrationStatus.value,
    registrationStartDate: DateTime.fromISO(registrationStartDate.value).toISO(),
    registrationEndDate: DateTime.fromISO(registrationEndDate.value).toISO(),
  })

  emit('close')
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-lg w-full max-w-md shadow-xl p-6 relative">
      <h2 class="text-lg font-semibold text-zinc-800 mb-4">
        {{ props.initialData?.id ? t("edit") : t("newFeminine") }} campagne
      </h2>

      <div class="space-y-4">

        <div>
          <label class="block text-sm font-medium text-zinc-700">Année scolaire</label>
          <input
              type="text"
              v-model="schoolYear"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.schoolYear" class="text-sm text-red-500 mt-1">{{ errors.schoolYear }}</p>
        </div>


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
          <label class="block text-sm font-medium text-zinc-700">Date de début</label>
          <input
              type="date"
              v-model="startDate"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.startDate" class="text-sm text-red-500 mt-1">{{ errors.startDate }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-zinc-700">Date de fin</label>
          <input
              type="date"
              v-model="endDate"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.endDate" class="text-sm text-red-500 mt-1">{{ errors.endDate }}</p>
        </div>


        <div>
          <label class="block text-sm font-medium text-zinc-700">Début des inscriptions</label>
          <input
              type="date"
              v-model="registrationStartDate"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.registrationStartDate" class="text-sm text-red-500 mt-1">{{ errors.registrationStartDate }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-zinc-700">Fin des inscriptions</label>
          <input
              type="date"
              v-model="registrationEndDate"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.registrationEndDate" class="text-sm text-red-500 mt-1">{{ errors.registrationEndDate }}</p>
        </div>


        <div>
          <label class="block text-sm font-medium text-zinc-700">Statut des inscriptions</label>
          <select
              v-model="registrationStatus"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="OPEN">Ouverte</option>
            <option value="CLOSED">Fermée</option>
          </select>
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