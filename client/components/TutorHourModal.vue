<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { DateTime } from 'luxon'
import type {TutoringHour} from "~/types/api";
const { t } = useI18n()

const props = defineProps<{
  show: boolean
  initialData?: TutoringHour
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', hour: TutoringHour): void
}>()

const date = ref('')
const time = ref('')
const duration = ref(60)
const errors = ref<{ date?: string; time?: string; duration?: string }>({})

watch(
    () => props.initialData,
    (data) => {
      if (data) {
        const start = DateTime.fromISO(data.startDate).setZone('local')
        const end = DateTime.fromISO(data.endDate).setZone('local')

        date.value = start.toISODate()
        time.value = start.toFormat('HH:mm')

        const diff = end.diff(start, 'minutes').minutes
        duration.value = Math.min(diff, 90)
      } else {
        date.value = ''
        time.value = ''
        duration.value = 60
      }

      errors.value = {}
    },
    { immediate: true }
)

const isValid = computed(() => {
  errors.value = {}
  if (!date.value) errors.value.date = t("dateRequired")
  if (!time.value) errors.value.time = t("timeRequired")
  if (!duration.value || duration.value > 90) {
    errors.value.duration = t("durationError")
  }

  return Object.keys(errors.value).length === 0
})

function submit() {
  if (!isValid.value) return

  const start = DateTime.fromISO(`${date.value}T${time.value}`, { zone: 'local' })
  const end = start.plus({ minutes: duration.value })

  emit('submit', {
    id: props.initialData?.id,
    startDate: start.toISO(),
    endDate: end.toISO()
  })

  emit('close')
}

// ChatGPT + https://stackoverflow.com/questions/66906787/calculate-duration-by-start-and-end-hour
function formatDurationWithEnd(mins: number): string {
  const h = Math.floor(mins / 60)
  const m = mins % 60
  const label = `${h > 0 ? `${h}h` : ''}${m > 0 ? ` ${m}min` : ''}`.trim()

  if (!date.value || !time.value) return label

  const start = DateTime.fromISO(`${date.value}T${time.value}`, { zone: 'local' })
  const end = start.plus({ minutes: mins })
  return `${label} — ${t("endAt")} ${end.toFormat('HH:mm')}`
}

const formattedDuration = computed(() => {
  const h = Math.floor(duration.value / 60)
  const m = duration.value % 60
  return `${h > 0 ? `${h}h` : ''}${m > 0 ? ` ${m}min` : ''}`.trim()
})

const calculatedEndTime = computed(() => {
  if (!date.value || !time.value) return null
  const start = DateTime.fromISO(`${date.value}T${time.value}`, { zone: 'local' })
  return start.plus({ minutes: duration.value }).toFormat('HH:mm')
})
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-lg w-full max-w-md shadow-xl p-6 relative">
      <h2 class="text-lg font-semibold text-zinc-800 mb-4">
        {{ props.initialData?.id ? $t("editHoursTitle") : $t("addHoursTitle") }}
      </h2>

      <div class="space-y-4">

        <div>
          <label class="block text-sm font-medium text-zinc-700">{{ $t("dateLabel") }}</label>
          <input
              type="date"
              v-model="date"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.date" class="text-sm text-red-500 mt-1">{{ errors.date }}</p>
        </div>


        <div>
          <label class="block text-sm font-medium text-zinc-700">{{ $t("startHour") }}</label>
          <input
              type="time"
              v-model="time"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p v-if="errors.time" class="text-sm text-red-500 mt-1">{{ errors.time }}</p>
        </div>


        <div>
          <label class="block text-sm font-medium text-zinc-700">{{ $t("duration") }}</label>
          <select
              v-model="duration"
              class="mt-1 w-full border border-zinc-300 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option v-for="d in [15, 30, 45, 60, 75, 90]" :key="d" :value="d">
              {{ formatDurationWithEnd(d) }}
            </option>
          </select>
          <p v-if="errors.duration" class="text-sm text-red-500 mt-1">{{ errors.duration }}</p>
          <p v-if="duration" class="text-sm text-zinc-600 mt-1">
            {{ $t("duration") }} : {{ formattedDuration }}<span v-if="calculatedEndTime"> — {{ $t("estimatedEndTime") }} : {{ calculatedEndTime }}</span>
          </p>
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
