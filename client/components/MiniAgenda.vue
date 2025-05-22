<script setup lang="ts">
const props = defineProps<{
  slots: Record<string, number[]> // Format: { "1": [-1, 0, 13, ...], ... }
}>()

const dayOrder = [1, 2, 3, 4, 5]
const maxSlots = 7

const getColorClass = (val: number | undefined): string => {
  if (val === -1) return 'bg-orange-500'
  if (val === 0) return 'bg-green-500'
  return 'bg-gray-300'
}
</script>

<template>
  <div class="inline-grid grid-cols-5 gap-1 p-2 rounded-lg border border-zinc-200 w-fit">
    <div
        v-for="day in dayOrder"
        :key="day"
        class="grid grid-rows-7 gap-1"
    >
      <div
          v-for="period in maxSlots"
          :key="`${day}-${period}`"
          class="w-4 h-4 rounded-sm flex items-center justify-center"
          :class="getColorClass(slots[day]?.[period - 1])"
      >
        <div class=" text-xs text-gray-500" v-if="slots[day]?.[period - 1] !== undefined && slots[day][period - 1] > 0">
            {{ slots[day][period - 1] }}
        </div>
      </div>
    </div>
  </div>
</template>
