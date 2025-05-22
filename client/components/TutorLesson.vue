<script setup lang="ts">
import { CalendarDaysIcon, PencilSquareIcon, TrashIcon, ChatBubbleLeftEllipsisIcon } from '@heroicons/vue/24/solid'
import type {TutoringLesson} from "~/types/api";

defineProps<{
  canEdit: boolean,
  lesson: TutoringLesson,
  onEdit: (lessonId: number) => void
  onDelete: (lessonId: number) => void
}>()

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('fr-FR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="relative bg-zinc-50 p-4 rounded-md border border-zinc-200 group hover:shadow transition">

    <div v-if="canEdit" class="absolute top-2 right-2 flex gap-2">
      <button @click="onEdit(lesson.id)" class="p-1 rounded hover:bg-blue-100 text-blue-600" title="Modifier">
        <PencilSquareIcon class="w-5 h-5" />
      </button>
      <button @click="onDelete(lesson.id)" class="p-1 rounded hover:bg-red-100 text-red-600" title="Supprimer">
        <TrashIcon class="w-5 h-5" />
      </button>
    </div>


    <div class="flex items-center gap-2 text-sm text-zinc-500 mb-2">
      <CalendarDaysIcon class="w-5 h-5 text-zinc-400" />
      <span>{{ formatDate(lesson.startDate) }} → {{ formatDate(lesson.endDate) }}</span>
    </div>


    <div class="flex items-start gap-2 text-zinc-700">
      <ChatBubbleLeftEllipsisIcon class="w-5 h-5 text-zinc-500 mt-0.5" />
      <p class="whitespace-pre-line text-sm leading-relaxed">{{ lesson.content }}</p>
    </div>
  </div>
</template>
