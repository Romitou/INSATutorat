<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import OverviewAgenda from '~/components/OverviewAgenda.vue'
import {useRoute} from 'vue-router'
import {useToast} from "vue-toastification";
import {type SlotStatus, statuses} from "~/types/api";

definePageMeta({layout: 'loggedin'})

const campaignId = useRoute().params.campaignId as string

const slotStatuses = ref<Record<string, SlotStatus>>({})
const isLoading = ref(true);

const agenda = ref<{ day: string, periods: { period: number, items: Record<string, number> }[] }[]>([])

const days = ['MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY']
const dayMap = {
  MONDAY: 1,
  TUESDAY: 2,
  WEDNESDAY: 3,
  THURSDAY: 4,
  FRIDAY: 5,
}
const reverseDayMap = Object.fromEntries(Object.entries(dayMap).map(([k, v]) => [v, k]))

const timeSlots = [
  {
    label: '08:00 - 09:30',
    period: 0
  },
  {
    label: '09:45 - 11:15',
    period: 1
  },
  {
    label: '11:30 - 13:00',
    period: 2
  },
  {
    label: '13:15 - 14:45',
    period: 3
  },
  {
    label: '15:00 - 16:30',
    period: 4
  },
  {
    label: '16:45 - 18:15',
    period: 5
  },
  {
    label: '18:30 - 20:00',
    period: 6
  }
]

const toggleStatus = (day: string, period: number) => {
  const key = `${day}-${period}`
  const currentIndex = statuses.indexOf(slotStatuses.value[key])
  slotStatuses.value[key] = statuses[(currentIndex + 1) % statuses.length]
}

const occupiedSlots = new Set<string>()

const fetchAgendaAndAvailabilities = async () => {
  try {
    const [agendaRes, availabilitiesRes] = await Promise.all([
      useApiFetch(`/campaign/${campaignId}/agenda`),
      useApiFetch(`/campaign/${campaignId}/availabilities`)
    ])

    if (agendaRes.ok) {
      agenda.value = await agendaRes.json()
      for (const {day, periods} of agenda.value) {
        for (const {period} of periods) {
          occupiedSlots.add(`${day}-${period}`)
        }
      }
    }

    for (const day of days) {
      for (const {period} of timeSlots) {
        const key = `${day}-${period}`
        if (occupiedSlots.has(key)) continue
        slotStatuses.value[key] = 'AVAILABLE'
      }
    }

    if (availabilitiesRes.ok) {
      const data: Record<string, number[]> = await availabilitiesRes.json()

      for (const [dayNumStr, values] of Object.entries(data)) {
        const day = reverseDayMap[+dayNumStr]
        if (!day) continue

        values.forEach((val, period) => {
          const key = `${day}-${period}`
          if (occupiedSlots.has(key)) return

          slotStatuses.value[key] = val === -1 ? 'OCCUPIED' : 'AVAILABLE'
        })
      }
    }

    isLoading.value = false
  } catch (err) {
    console.error('Erreur lors de la récupération des données :', err)
    useToast().error('Erreur lors de la récupération des données')
  }
}

onMounted(fetchAgendaAndAvailabilities)

const mapStatusesToApiFormat = (): Record<string, number[]> => {
  return Object.fromEntries(
      days.map(day => {
        const apiDay = dayMap[day].toString()
        const values = timeSlots.map(({period}) => {
          const status = slotStatuses.value[`${day}-${period}`]
          return status === 'OCCUPIED' ? -1 : 0
        })
        return [apiDay, values]
      })
  )
}

const submitSlots = async () => {
  const payload = mapStatusesToApiFormat()

  const toastId = useToast().info('Enregistrement des créneaux...', {
    timeout: false,
  })
  try {
    const res = await useApiFetch(`/campaign/${campaignId}/availabilities`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    })

    if (res.ok) {
      useToast().update(toastId, {
        content: 'Créneaux enregistrés avec succès !',
        options: {
          type: "success",
          timeout: 5000,
        }
      })
    } else {
      console.error(await res.text())
      useToast().update(toastId, {
        content: 'Erreur lors de l\'enregistrement des créneaux',
        options: {
          type: "error",
          timeout: 5000,
        }
      })
    }
  } catch (err) {
    console.error('Erreur lors de l\'enregistrement des créneaux :', err)
    useToast().update(toastId, {
      content: 'Erreur lors de l\'enregistrement des créneaux',
      options: {
        type: "error",
        timeout: 5000,
      }
    })
  }
}
</script>

<template>
  <div class="max-w-6xl w-full mx-auto px-4 md:px-8 py-8">
    <h1 class="text-3xl font-bold text-red-700 text-center mb-10">
      {{ $t('slotSelectionTitle') }}
    </h1>

    <div class="bg-white rounded-2xl shadow p-6 md:p-10">
      <OverviewAgenda
          :agenda="agenda"
          :days="days"
          :is-loading="isLoading"
          :slotStatuses="slotStatuses"
          :timeSlots="timeSlots"
          :toggleStatus="toggleStatus"
      />
    </div>

    <div class="flex justify-end mt-6">
      <button
          class="px-6 py-3 rounded-xl bg-red-600 text-white font-semibold shadow hover:bg-red-700 transition"
          @click="submitSlots"
      >
        {{ $t('saveSlotsButton') }}
      </button>
    </div>
  </div>
</template>
