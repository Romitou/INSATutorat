<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {useRoute} from 'vue-router'
import draggable from 'vuedraggable'
import {TrashIcon, UserPlusIcon} from '@heroicons/vue/24/outline'
import {useToast} from "vue-toastification";
import type {Subject, TuteeAssignment, TutorSubject, User} from "~/types/api";

interface AssignmentResponse {
  tutees: TuteeAssignment[]
  tutorSubjects: TutorSubject[]
}

interface GenerationResponse {
  affectedTutees: TuteeAssignment[]
  logs: string[]
}

definePageMeta({
  layout: 'loggedin'
})

const campaignId = useRoute().params.campaignId as string

const users = ref<User[]>([])
const subjects = ref<Subject[]>([])
const tuteeAssignments = ref<TuteeAssignment[]>([])
const tutorSubjects = ref<TutorSubject[]>([])
const generationLogs = ref<string[]>([])

const fetchUsers = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}/users`)
  if (res.ok) users.value = await res.json()
}

const fetchSubjects = async () => {
  const res = await useApiFetch(`/campaign/${campaignId}/subjects`)
  if (res.ok) subjects.value = await res.json()
}

const fetchAssignments = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}/assignments`)
  if (res.ok) {
    const result = await res.json() as AssignmentResponse
    tuteeAssignments.value = result.tutees
    tutorSubjects.value = result.tutorSubjects
  }
}

const fetchGeneratedAssignments = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}/generate-assignments`)

  if (!res.ok) {
    useToast().error('Erreur lors de la génération automatique des affectations');
    return
  }

  const result: GenerationResponse = await res.json()
  const {affectedTutees, logs} = result;

  generationLogs.value = logs

  tuteeAssignments.value = tuteeAssignments.value.map((assignment: TuteeAssignment) => {
    const affected = affectedTutees.find((tutee: TuteeAssignment) => tutee.id === assignment.id)
    if (affected) {
      return {
        ...assignment,
        tutorSubjectId: affected.tutorSubjectId
      }
    }
    return assignment
  })

  useToast().success('Affectations générées automatiquement');
}

const saveAssignments = async () => {
  const tutees = tuteeAssignments.value.map(assignment => ({
    id: assignment.id,
    tutee: {
      id: assignment.tutee.id,
      firstName: assignment.tutee.firstName,
      lastName: assignment.tutee.lastName
    },
    tuteeId: assignment.tuteeId,
    subjectId: assignment.subjectId,
    tutorSubjectId: assignment.tutorSubjectId
  }));

  const apiTutorSubjects = tutorSubjects.value.map(tutorSubject => ({
    id: tutorSubject.id,
    subjectId: tutorSubject.subjectId,
    tutor: {
      id: tutorSubject.tutor.id,
      firstName: tutorSubject.tutor.firstName,
      lastName: tutorSubject.tutor.lastName
    },
    maxTutees: tutorSubject.maxTutees
  }));

  const payload = {
    tutees,
    tutorSubjects: apiTutorSubjects
  };

  try {
    const res = await useApiFetch(`/admin/campaign/${campaignId}/assignments`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    });

    if (res.ok) {
      useToast().success('Les modifications ont été enregistrées avec succès');
    } else {
      useToast().error('Erreur lors de l\'enregistrement des modifications');
    }
  } catch (error) {
    useToast().error('Erreur lors de l\'enregistrement des modifications');
    console.error(error);
  }

  await fetchAssignments()
}

const showAddTutorModal = ref(false)
const selectedSubjectId = ref<number | null>(null)

const openAddTutorModal = (subjectId: number) => {
  selectedSubjectId.value = subjectId
  showAddTutorModal.value = true
}

const closeAddTutorModal = () => {
  showAddTutorModal.value = false
  selectedSubjectId.value = null
}

const handleAddTutor = (user: User, maxTutees: number) => {
  if (selectedSubjectId.value === null) return

  tutorSubjects.value.push({
    subjectId: selectedSubjectId.value,
    tutor: {
      id: user.id,
      firstName: user.firstName,
      lastName: user.lastName,
    },
    maxTutees,
  } as TutorSubject)
}

const showAddTuteeModal = ref(false)

const openAddTuteeModal = (subjectId: number) => {
  selectedSubjectId.value = subjectId
  showAddTuteeModal.value = true
}

const closeAddTuteeModal = () => {
  showAddTuteeModal.value = false
  selectedSubjectId.value = null
}

const handleAddTutee = (user: User) => {
  if (selectedSubjectId.value === null) return

  tuteeAssignments.value.push({
    tuteeId: user.id,
    subjectId: selectedSubjectId.value,
    tutorSubjectId: null,
    tutee: {
      id: user.id,
      firstName: user.firstName,
      lastName: user.lastName,
    }
  } as TuteeAssignment)
}

const removeTutor = async (tutorSubjectId: number) => {
  if (!confirm("Êtes-vous sûr de vouloir supprimer ce tuteur ?")) return
  const res = await useApiFetch(`/admin/campaign/${campaignId}/assignments/tutor`, {
    method: 'DELETE',
    body: JSON.stringify({
      id: tutorSubjectId
    })
  })
  if (res.ok) {
    const index = tutorSubjects.value.findIndex(ts => ts.id === tutorSubjectId)
    if (index !== -1) {
      tutorSubjects.value = tutorSubjects.value.filter(ts => ts.id !== tutorSubjectId)
    } else {
      useToast().error('Erreur lors de la suppression du tuteur')
    }
  }
}

const removeTutee = async (subjectId: number, tuteeId: number) => {
  const tuteeReg = tuteeAssignments.value.find(ts => ts.tuteeId === tuteeId && ts.subjectId === subjectId);
  if (!tuteeReg) {
    useToast().error('Erreur lors de la suppression du tutoré')
    return
  }
  if (!confirm("Êtes-vous sûr de vouloir supprimer ce tutoré ?")) return
  const res = await useApiFetch(`/admin/campaign/${campaignId}/assignments/tutee`, {
    method: 'DELETE',
    body: JSON.stringify({
      id: tuteeReg.id
    })
  })
  if (res.ok) {
    const index = tuteeAssignments.value.findIndex(ts => ts.id === tuteeReg.id)
    if (index !== -1) {
      tuteeAssignments.value = tuteeAssignments.value.filter(ts => ts.id !== tuteeReg.id)
    } else {
      useToast().error('Erreur lors de la suppression du tutoré')
    }
  }
}

const getTutorsBySubject = (subjectId: number) => {
  return tutorSubjects.value.filter(ts => ts.subjectId === subjectId)
}

const getTuteesBySubject = (subjectId: number): TuteeAssignment[] => {
  return tuteeAssignments.value.filter(a => a.subjectId === subjectId)
}

const getUnassignedTuteesBySubject = (subjectId: number) => {
  return tuteeAssignments.value
      .filter(a => a.subjectId === subjectId && a.tutorSubjectId === null)
      .map(a => a.tutee)
}

const getTuteesByTutorSubject = (tutorSubjectId: number) => {
  return tuteeAssignments.value
      .filter(a => a.tutorSubjectId === tutorSubjectId)
      .map(a => a.tutee)
}

const updateAssignments = (subjectId: number, tutorSubjectId: number | null, e: any) => {
  const movedTutees: User[] = [];

  if (Array.isArray(e)) {
    movedTutees.push(...e);
  } else if (e?.added?.element) {
    movedTutees.push(e.added.element);
  } else {
    console.warn("Aucun tutoré détecté dans updateAssignments", e);
    return;
  }

  movedTutees.forEach((tutee) => {
    const matchingAssignments = tuteeAssignments.value.filter(
        (assignment) =>
            assignment.subjectId === subjectId &&
            (assignment.tuteeId === tutee.id || assignment.tutee?.id === tutee.id)
    );

    if (matchingAssignments.length === 0) {
      console.warn(
          `Aucuen ligne d\'affectation trouvée pour le tutoré ID ${tutee.id} dans la matière ${subjectId}`
      );
      return;
    }

    matchingAssignments.forEach((assignment) => {
      assignment.tutorSubjectId = tutorSubjectId;
    });
  });
};

onMounted(() => {
  fetchUsers()
  fetchSubjects()
  fetchAssignments()
})
</script>

<template>
  <div class="p-6 space-y-12">
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">Affectations tuteurs/tutorés</h1>
      <div class="flex gap-4">
        <button
            class="bg-blue-600 text-white px-4 py-2 rounded-xl hover:bg-blue-700"
            @click="fetchGeneratedAssignments"
        >
          <UserPlusIcon class="w-5 h-5 inline mr-2"/>
          Générer les affectations
        </button>
        <button
            class="bg-green-600 text-white px-4 py-2 rounded-xl hover:bg-green-700"
            @click="saveAssignments"
        >
          💾 Sauvegarder
        </button>
      </div>
    </div>

    <MatchingLogs
        v-if="tuteeAssignments.length > 0"
        :logs="generationLogs"
    />

    <div
        v-for="subject in subjects"
        :key="subject.id"
        class="bg-white rounded-2xl shadow-md p-6 space-y-6"
    >
      <div class="flex justify-between items-center">
        <h2 class="text-2xl font-semibold text-gray-800">{{ subject.name }}</h2>
        <div class="flex flex-row space-x-2">
          <button
              class="bg-indigo-600 text-white px-3 py-1 rounded-lg hover:bg-indigo-700"
              @click="openAddTutorModal(subject.id)"
          >
            ➕ Ajouter un tuteur
          </button>
          <button
              class="bg-indigo-600 text-white px-3 py-1 rounded-lg hover:bg-indigo-700"
              @click="openAddTuteeModal(subject.id)"
          >
            ➕ Ajouter un tutoré
          </button>
        </div>
      </div>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">

        <div
            v-for="tutorSubject in getTutorsBySubject(subject.id)"
            :key="tutorSubject.id"
            class="bg-white rounded-xl border border-gray-200 p-4 shadow-sm"
        >
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-bold text-gray-700">
              {{ tutorSubject.tutor.firstName }} {{ tutorSubject.tutor.lastName }}
            </h3>
            <button
                :class="`${getTuteesByTutorSubject(tutorSubject.id).length > 0 ? 'text-gray-600' : 'text-red-600 hover:text-red-800'} p-1 rounded-md`"
                :disabled="getTuteesByTutorSubject(tutorSubject.id).length > 0"
                title="Supprimer le tuteur"
                @click="removeTutor(tutorSubject.id)"
            >
              <TrashIcon class="w-5 h-5"/>
            </button>
          </div>
          <div class="flex flex-row items-center space-x-3 mb-2">
            <NuxtLink
                :to="`/tutoring/${tutorSubject.id}`"
                class="text-sm text-blue-500 hover:underline"
            >
              Voir le profil
            </NuxtLink>
            <div class="text-sm text-gray-500">
              {{ tutorSubject.totalHours }} h
            </div>
            <div class="text-sm text-gray-500">Capacité max :
              <input
                  v-model.number="tutorSubject.maxTutees"
                  class="ml-2 w-10 border rounded px-2 py-0.5 text-sm"
                  min="0"
                  type="number"
              />
            </div>
          </div>

          <draggable
              :group="'tutees-' + subject.id"
              :list="getTuteesByTutorSubject(tutorSubject.id)"
              class="space-y-2 min-h-[3rem]"
              item-key="id"
              @change="e => e.added && updateAssignments(subject.id, tutorSubject.id, e)"
          >
            <template #item="{ element }">
              <div
                  class="bg-gray-100 px-3 py-1 rounded-lg text-sm text-gray-800 flex items-center justify-between"
              >
                {{ element.firstName }} {{ element.lastName }}
                <button
                    class="ml-2 text-red-500 hover:text-red-700"
                    @click="updateAssignments(subject.id, null, [element])"
                >
                  ✕
                </button>
              </div>
            </template>
          </draggable>

          <div v-if="getTuteesByTutorSubject(tutorSubject.id).length === 0" class="text-sm text-gray-400 italic mt-2">
            Aucun tutoré assigné.
          </div>
        </div>


        <div class="bg-white rounded-xl border border-dashed border-gray-300 p-4 shadow-sm">
          <h3 class="text-lg font-semibold text-gray-600 mb-2">Tutorés non assignés</h3>
          <draggable
              :group="'tutees-' + subject.id"
              :list="getUnassignedTuteesBySubject(subject.id)"
              class="space-y-2 min-h-[3rem]"
              item-key="id"
              @change="e => e.added && updateAssignments(subject.id, null, e)"
          >
            <template #item="{ element }">
              <div
                  class="bg-gray-100 px-3 py-1 rounded-lg text-sm text-gray-800 flex items-center justify-between"
              >
                {{ element.firstName }} {{ element.lastName }}

                <button
                    class="ml-2 text-red-500 hover:text-red-700"
                    @click="removeTutee(subject.id, element.id)"
                >
                  <TrashIcon class="w-4 h-4"/>
                </button>
              </div>
            </template>
          </draggable>

          <div v-if="getUnassignedTuteesBySubject(subject.id).length === 0" class="text-sm text-gray-400 italic mt-2">
            Aucun tutoré non assigné.
          </div>
        </div>
      </div>
    </div>
  </div>
  <AddTutorModal
      v-if="selectedSubjectId !== null"
      :all-tutors="users.filter(user => user.isTutor)"
      :assigned-tutors="getTutorsBySubject(selectedSubjectId).map(ts => ts.tutor)"
      :show="showAddTutorModal"
      :subject-id="selectedSubjectId"
      @close="closeAddTutorModal"
      @submit="handleAddTutor"
  />
  <AddTuteeModal
      v-if="selectedSubjectId !== null"
      :all-tutees="users.filter(user => user.isTutee)"
      :assigned-tutees="getTuteesBySubject(selectedSubjectId).map(ts => ts.tutee)"
      :show="showAddTuteeModal"
      :subject-id="selectedSubjectId"
      @close="closeAddTuteeModal"
      @submit="handleAddTutee"
  />
</template>
