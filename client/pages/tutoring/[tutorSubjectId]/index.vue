<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {useRoute} from 'vue-router'
import {AcademicCapIcon, BookOpenIcon, CalendarDaysIcon, ClockIcon, PlusIcon, UsersIcon} from '@heroicons/vue/24/solid'
import {useToast} from "vue-toastification";
import type {Subject, TutoringHour, TutoringLesson, User} from "~/types/api";

definePageMeta({layout: 'loggedin'})

interface TuteeWithHours extends User {
  hours: TutoringHour[]
}

interface TutorSubjectSummary {
  subject: Subject
  tutor: User
  lessons: TutoringLesson[]
  tutees: TuteeWithHours[]
}

const tutorSubjectId = useRoute().params.tutorSubjectId as string
const tutorSubject = ref<TutorSubjectSummary | undefined>(undefined)

const userId = computed(() => {
  const user = useUserStore().user
  if (!user) return null
  return user.id
})

onMounted(async () => {
  const response = await useApiFetch(`/tutoring/${tutorSubjectId}/summary`, {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })
  if (response.ok) {
    tutorSubject.value = (await response.json()) as TutorSubjectSummary
  }
})
const sortedLessons = computed(() => {
  return tutorSubject.value?.lessons.slice().sort((a, b) => new Date(b.startDate).getTime() - new Date(a.startDate).getTime()) || []
})

function sortedHours(hours: any[]) {
  return hours.slice().sort((a, b) => new Date(b.startDate).getTime() - new Date(a.startDate).getTime())
}

function getTotalHours(hours: any[]) {
  let totalMs = 0
  for (const h of hours) {
    const start = new Date(h.startDate).getTime()
    const end = new Date(h.endDate).getTime()
    totalMs += (end - start)
  }
  return (totalMs / 3600000).toFixed(2)
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('fr-FR', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

function getTotalHoursAll() {
  let total = 0
  for (const tutee of tutorSubject.value?.tutees || []) {
    for (const h of tutee.hours) {
      const start = new Date(h.startDate).getTime()
      const end = new Date(h.endDate).getTime()
      total += end - start
    }
  }
  return (total / 3600000).toFixed(2)
}

const lastLessonDate = computed(() => {
  if (!sortedLessons.value.length) return null
  return formatDate(sortedLessons.value[0].startDate)
});

const showModal = ref(false)
const editingLesson = ref<TutoringLesson | undefined>(undefined)

function handleAddLesson() {
  editingLesson.value = undefined
  showModal.value = true
}

function handleEditLesson(id: number) {
  if (!tutorSubject.value) return
  editingLesson.value = tutorSubject.value.lessons.find((l) => l.id === id)
  showModal.value = true
}

async function handleSubmitLesson(lesson: TutoringLesson) {
  if (!tutorSubject.value) return
  if (lesson.id) {
    const response = await useApiFetch(`/tutoring/${tutorSubjectId}/lesson/${lesson.id}`, {
      method: 'PATCH',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(lesson)
    })
    if (response.ok) {
      const updatedLesson = await response.json()
      lesson = updatedLesson as TutoringLesson
      useToast().success('Séance mise à jour avec succès')
    } else {
      useToast().error('Erreur lors de la mise à jour de la séance')
      console.error('Erreur lors de la mise à jour de la séance', response.statusText)
      return
    }

    const idx = tutorSubject.value.lessons.findIndex((l) => l.id === lesson.id)
    if (idx !== -1) tutorSubject.value.lessons[idx] = lesson
  } else {
    const response = await useApiFetch(`/tutoring/${tutorSubjectId}/lessons`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(lesson)
    })
    if (response.ok) {
      const apiLesson = await response.json()
      tutorSubject.value.lessons.push(apiLesson as TutoringLesson)
      useToast().success('Séance créée avec succès')
    } else {
      useToast().error('Erreur lors de la création de la séance')
      console.error('Erreur lors de la création de la séance', response.statusText)
      return
    }
  }
}

async function handleDeleteLesson(lessonId: number) {
  if (!tutorSubject.value) return
  if (confirm('Supprimer cette séance ?')) {
    const response = await useApiFetch(`/tutoring/${tutorSubjectId}/lesson/${lessonId}`, {
      method: 'DELETE',
      headers: {'Content-Type': 'application/json'},
    })
    if (response.ok) {
      const idx = tutorSubject.value.lessons.findIndex((l) => l.id === lessonId)
      if (idx !== -1) tutorSubject.value.lessons.splice(idx, 1)
      useToast().success('Séance supprimée avec succès')
    } else {
      useToast().error('Erreur lors de la suppression de la séance')
      console.error('Erreur lors de la suppression de la séance', response.statusText)
      return
    }
  }
}

const selectedTuteeId = ref<number | null>(null)

async function handleAddHour(hour: TutoringHour) {
  if (!tutorSubject.value) return
  const tutee = tutorSubject.value?.tutees.find(t => t.id === selectedTuteeId.value)
  if (!tutee) {
    useToast().error('Erreur : tuteur introuvable')
    return
  }

  hour.tuteeId = tutee.id
  const response = await useApiFetch(`/tutoring/${tutorSubjectId}/hours`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(hour)
  })
  if (response.ok) {
    const newHour = await response.json()
    tutee.hours.push(newHour)
    useToast().success('Heure ajoutée avec succès')
    selectedTuteeId.value = null
  } else {
    useToast().error('Erreur lors de l\'ajout de la séance')
    console.error('Erreur lors de l\'ajout de la séance', response.statusText)
    return
  }
}


async function handleDeleteHour(hourId: number, tuteeId: number) {
  if (!tutorSubject.value) return
  const tutee = tutorSubject.value?.tutees.find(t => t.id === tuteeId)
  if (!tutee) {
    useToast().error('Erreur : tuteur introuvable')
    return
  }

  if (confirm('Supprimer cette heure ?')) {
    const response = await useApiFetch(`/tutoring/${tutorSubjectId}/hour/${hourId}`, {
      method: 'DELETE',
      headers: {'Content-Type': 'application/json'},
    })
    if (response.ok) {
      const idx = tutee.hours.findIndex((l) => l.id === hourId)
      if (idx !== -1) tutee.hours.splice(idx, 1)
      useToast().success('Heure supprimée avec succès')
    } else {
      useToast().error('Erreur lors de la suppression de l\'heure')
      console.error('Erreur lors de la suppression de l\'heure', response.statusText)
      return
    }
  }
}

const editingHour = ref<TutoringHour | undefined>(undefined)

function handleEditHour(hour: TutoringHour, tuteeId: number) {
  editingHour.value = {...hour, tuteeId}
  selectedTuteeId.value = tuteeId
}

async function handleSubmitHour(hour: TutoringHour) {
  const tutee = tutorSubject.value?.tutees.find(t => t.id === selectedTuteeId.value)
  if (!tutee) {
    useToast().error('Erreur : tutoré introuvable')
    return
  }

  if (hour.id) {
    hour.tuteeId = tutee.id
    const response = await useApiFetch(`/tutoring/${tutorSubjectId}/hour/${hour.id}`, {
      method: 'PATCH',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(hour)
    })
    if (response.ok) {
      const updated = await response.json()
      const idx = tutee.hours.findIndex(h => h.id === hour.id)
      if (idx !== -1) tutee.hours[idx] = updated
      useToast().success('Heure modifiée avec succès')
    } else {
      useToast().error('Erreur lors de la modification')
      return
    }
  } else {
    await handleAddHour(hour)
  }

  selectedTuteeId.value = null
  editingHour.value = undefined
}

const isTutorOrAdmin = computed(() => {
  const user = useUserStore().user
  if (!user) return false
  if (user.isAdmin) return true
  if (!tutorSubject.value) return false
  return user.id === tutorSubject.value.tutor.id
})

const sortedTutees = computed(() => {
  if (!tutorSubject.value) return []
  return tutorSubject.value.tutees.slice().sort((a, b) => {
    if (a.id === userId.value) return -1
    if (b.id === userId.value) return 1
    return 0
  })
})
</script>

<template>
  <div class="max-w-6xl w-full mx-auto p-6 space-y-10">
    <div class="mt-10 bg-white rounded-lg shadow p-6 border space-y-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-semibold text-zinc-800 flex items-center gap-2">
            <AcademicCapIcon class="w-6 h-6 text-blue-600"/>
            {{ $t('tutoring') }} : {{ tutorSubject?.subject.name }}
          </h2>
          <p class="text-zinc-600 text-sm mt-1">
            {{ $t('semester') }} : <span class="font-medium">{{ tutorSubject?.subject.semester }}</span> –
            {{ $t('subject') }} : <span class="font-medium">{{ tutorSubject?.subject.shortName }}</span>
          </p>
        </div>
        <div class="text-right">
          <p class="text-sm text-zinc-500">{{ $t('tutor') }} :</p>
          <p class="font-medium text-zinc-800">{{ tutorSubject?.tutor.firstName }} {{
              tutorSubject?.tutor.lastName
            }}</p>
        </div>
      </div>


      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-center text-sm">
        <div class="flex flex-col items-center justify-center">
          <BookOpenIcon class="w-6 h-6 text-indigo-500 mb-1"/>
          <p class="text-zinc-600">{{ $t('sessions') }}</p>
          <p class="font-bold text-zinc-800 text-lg">{{ tutorSubject?.lessons.length }}</p>
        </div>
        <div class="flex flex-col items-center justify-center">
          <UsersIcon class="w-6 h-6 text-green-500 mb-1"/>
          <p class="text-zinc-600">{{ $t('tutees') }}</p>
          <p class="font-bold text-zinc-800 text-lg">{{ tutorSubject?.tutees.length }}</p>
        </div>
        <div v-if="isTutorOrAdmin" class="flex flex-col items-center justify-center">
          <ClockIcon class="w-6 h-6 text-yellow-500 mb-1"/>
          <p class="text-zinc-600">{{ $t('totalHours') }}</p>
          <p class="font-bold text-zinc-800 text-lg">{{ getTotalHoursAll() }} h</p>
        </div>
        <div class="flex flex-col items-center justify-center">
          <CalendarDaysIcon class="w-6 h-6 text-purple-500 mb-1"/>
          <p class="text-zinc-600">{{ $t('lastSession') }}</p>
          <p class="font-bold text-zinc-800 text-sm">
            <span v-if="lastLessonDate">{{ lastLessonDate }}</span>
            <span v-else class="italic text-zinc-500">—</span>
          </p>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow p-6 border">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold text-zinc-800">{{ $t('latestSessions') }}</h2>
        <button
            v-if="isTutorOrAdmin"
            class="inline-flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition text-sm"
            @click="handleAddLesson"
        >
          <PlusIcon class="w-5 h-5"/>
          {{ $t('addSession') }}
        </button>
      </div>

      <div v-if="sortedLessons.length" class="grid sm:grid-cols-2 gap-4">
        <TutorLesson
            v-for="lesson in sortedLessons"
            :key="lesson.id"
            :can-edit="isTutorOrAdmin"
            :lesson="lesson"
            :on-delete="handleDeleteLesson"
            :on-edit="handleEditLesson"
        />
      </div>
      <div v-else class="text-zinc-500 text-sm italic">
        {{ $t('noSessionsYet') }}
      </div>
    </div>

    <div class="bg-white rounded-lg shadow p-6 border">
      <h2 class="text-xl font-semibold text-zinc-800 mb-6">{{ $t('tutees') }}</h2>
      <div class="space-y-6">
        <div v-for="tutee in sortedTutees" :key="tutee.id" class="p-4 border rounded-md bg-zinc-50">
          <div class="flex justify-between items-center mb-2">
            <div>
              <p class="text-lg font-medium text-zinc-800">{{ tutee.firstName }} {{ tutee.lastName }}</p>
              <p v-if="isTutorOrAdmin || tutee.id === userId" class="text-sm text-zinc-500">
                {{ $t('entryCount', {count: tutee.hours.length}) }}
              </p>
              <p v-else>
                <span class="text-sm text-zinc-500">{{ $t('noAccessToHours') }}</span>
              </p>
            </div>
            <div v-if="isTutorOrAdmin || tutee.id === userId" class="flex items-center gap-3">
              <span class="text-sm bg-blue-100 text-blue-800 font-semibold px-3 py-1 rounded-full">
                {{ $t('total') }} : {{ getTotalHours(tutee.hours) }} h
              </span>
              <button
                  class="inline-flex items-center gap-1 bg-green-600 text-white px-3 py-1.5 rounded-md hover:bg-green-700 transition text-sm"
                  @click="selectedTuteeId = tutee.id"
              >
                <ClockIcon class="w-4 h-4"/>
                {{ $t('addHours') }}
              </button>
            </div>
          </div>
          <ul class="space-y-2 mt-2 text-sm text-zinc-600">
            <li
                v-for="hour in sortedHours(tutee.hours)"
                :key="hour.id"
                class="flex items-center justify-between p-2 bg-white rounded border hover:shadow-sm transition"
            >
              <div class="flex items-center gap-2">
                <ClockIcon class="w-4 h-4 text-zinc-400"/>
                <span>{{ formatDate(hour.startDate) }} → {{ formatDate(hour.endDate) }}</span>
              </div>
              <div class="flex items-center gap-1">
                <button class="text-blue-600 hover:underline text-xs" @click="handleEditHour(hour, tutee.id)">
                  {{ $t('edit') }}
                </button>
                <button class="text-red-500 hover:underline text-xs" @click="handleDeleteHour(hour.id, tutee.id)">
                  {{ $t('delete') }}
                </button>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
  <TutorLessonModal
      :initial-data="editingLesson"
      :show="showModal"
      @close="showModal = false"
      @submit="handleSubmitLesson"
  />
  <TutorHourModal
      :initial-data="editingHour"
      :show="selectedTuteeId !== null"
      :tutee-id="selectedTuteeId"
      @close="selectedTuteeId = null"
      @submit="handleSubmitHour"
  />
</template>

<style scoped>
</style>
