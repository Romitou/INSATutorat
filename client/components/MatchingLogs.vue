<template>
  <div class="bg-gray-900 text-gray-100 rounded-lg p-4 font-mono text-sm max-h-[500px] overflow-y-auto shadow">
    <div v-for="(log, index) in logs" :key="index" class="flex items-start gap-2 mb-1">
      <component
          :is="getIcon(log)"
          class="w-4 h-4 mt-0.5"
          :class="getIconColor(log)"
          aria-hidden="true"
      />
      <span :class="getTextColor(log)">
        {{ log }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { CheckCircleIcon, ExclamationCircleIcon, InformationCircleIcon, ArrowRightCircleIcon } from '@heroicons/vue/24/solid'

const props = defineProps({
  logs: {
    type: Array,
    required: true
  }
})

function getIcon(log) {
  if (log.toLowerCase().includes('×')) return ExclamationCircleIcon
  if (log.startsWith('✔')) return CheckCircleIcon
  if (log.startsWith('→')) return ArrowRightCircleIcon
  return InformationCircleIcon
}

function getIconColor(log) {
  if (log.toLowerCase().includes('×')) return 'text-red-400'
  if (log.startsWith('✔')) return 'text-green-400'
  if (log.startsWith('→')) return 'text-blue-400'
  return 'text-gray-400'
}

function getTextColor(log) {
  if (log.toLowerCase().includes('×')) return 'text-red-300'
  if (log.startsWith('✔')) return 'text-green-300'
  if (log.startsWith('→')) return 'text-blue-300'
  return 'text-gray-100'
}
</script>
