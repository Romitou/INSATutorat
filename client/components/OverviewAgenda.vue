<script setup lang="ts">
import { defineProps } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps<{
  days: string[] // exemple : ['MONDAY', 'TUESDAY', ...]
  timeSlots: { label: string; period: number }[]
  agenda: { day: string; periods: { period: number; items: Record<string, number> }[] }[]
  slotStatuses: Record<string, string>
  toggleStatus: (day: string, period: number) => void
  isLoading: boolean
}>()

const getCoursesFor = (day: string, period: number): Record<string, number> | null => {
  const dayEntry = props.agenda.find(d => d.day === day)
  if (!dayEntry) return null

  const periodEntry = dayEntry.periods.find(p => p.period === period)
  return periodEntry ? periodEntry.items : null
}

const isOccupied = (day: string, period: number): boolean => {
  const courses = getCoursesFor(day, period)
  return courses !== null && Object.keys(courses).length > 0
}

const totalCourses = (courses: Record<string, number> | null): number => {
  if (!courses) return 0
  return Object.values(courses).reduce((sum, count) => sum + count, 0)
}

const getTranslatedStatus = (status: string): string => {
  switch (status) {
    case 'OCCUPIED':
      return t('statusOccupied')
    case 'AVAILABLE':
      return t('statusAvailable')
    default:
      return t('statusUnknown')
  }
}
</script>

<template>
  <div class="overflow-x-auto border border-zinc-200 rounded-2xl">
    <div v-if="isLoading" class="p-6 space-y-4 animate-pulse">
      <div class="h-6 bg-gray-200 rounded w-1/3"></div>
      <div class="grid grid-cols-6 gap-4">
        <div v-for="n in 6" :key="n" class="h-12 bg-gray-200 rounded"></div>
        <div v-for="n in 6" :key="'row-' + n" class="h-12 bg-gray-100 rounded"></div>
        <div v-for="n in 6" :key="'row2-' + n" class="h-12 bg-gray-200 rounded"></div>
      </div>
    </div>

    <table v-else class="min-w-full text-sm rounded-2xl overflow-hidden">
      <thead class="bg-gray-100 text-zinc-600 font-semibold">
      <tr>
        <th class="text-left p-3">{{ t('scheduleHeader') }}</th>
        <th v-for="day in days" :key="day" class="text-center p-3">{{ t(`dayNames.${day}`) }}</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(slot, i) in timeSlots" :key="i" class="even:bg-gray-50">
        <td class="p-3 font-medium text-zinc-700">{{ slot.label }}</td>
        <td
            v-for="day in days"
            :key="day + i"
            @click="!isOccupied(day, i) && toggleStatus(day, i)"
            :class="[
              'relative group text-center p-4 transition',
              isOccupied(day, i)
                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                : 'hover:bg-red-50 cursor-pointer'
            ]"
        >

          <template v-if="isOccupied(day, i)">
            <div class="font-semibold">{{ t('statusOccupied') }}</div>
            <div class="text-xs text-gray-500">
              {{ t('courseCount', { count: totalCourses(getCoursesFor(day, i)) }) }}
            </div>
            <div class="absolute left-1/2 top-full w-44 mt-2 p-2 bg-white border rounded-lg shadow text-xs text-gray-700 transform -translate-x-1/2 hidden group-hover:block z-20">
              <ul>
                <li v-for="(count, label) in getCoursesFor(day, i)" :key="label">
                  {{ label }} ({{ count }})
                </li>
              </ul>
            </div>
          </template>


          <template v-else>
            <div class="text-xs font-semibold mb-1"
                 :class="{
                     'text-gray-400': slotStatuses[`${day}-${i}`] === 'UNKNOWN',
                     'text-green-600': slotStatuses[`${day}-${i}`] === 'AVAILABLE',
                     'text-orange-500': slotStatuses[`${day}-${i}`] === 'OCCUPIED'
                   }">
              {{ getTranslatedStatus(slotStatuses[`${day}-${i}`]) }}
            </div>
            <div class="w-3 h-3 mx-auto rounded-full"
                 :class="{
                     'bg-gray-300': slotStatuses[`${day}-${i}`] === 'UNKNOWN',
                     'bg-green-500': slotStatuses[`${day}-${i}`] === 'AVAILABLE',
                     'bg-orange-400': slotStatuses[`${day}-${i}`] === 'OCCUPIED'
                   }"></div>
          </template>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
