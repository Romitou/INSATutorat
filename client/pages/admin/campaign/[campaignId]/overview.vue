<script lang="ts" setup>
import {AllCommunityModule, ModuleRegistry} from 'ag-grid-community';
import {AgGridVue} from "ag-grid-vue3";
import {useRoute} from "#vue-router";
import {onMounted, ref} from "vue";
import type {Campaign, TuteeAssignment, TutorSubject, User} from "~/types/api";

ModuleRegistry.registerModules([AllCommunityModule]);

const colDefs = ref([
  {
    field: "firstName",
    headerName: "Prénom",
  },
  {
    field: "lastName",
    headerName: "Nom",
  },
  {
    field: "totalHours",
    headerName: "Heures totales",
  },
]);

definePageMeta({
  layout: "loggedin",
});

interface AssignmentResponse {
  tutees: TuteeAssignment[]
  tutorSubjects: TutorSubject[]
}

const campaign = ref<Campaign | undefined>(undefined);
const isEditing = ref(false);
const tuteeAssignments = ref<TuteeAssignment[]>([])
const tutorSubjects = ref<TutorSubject[]>([])

const fetchAssignments = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}/assignments`)
  if (res.ok) {
    const result = await res.json() as AssignmentResponse
    tuteeAssignments.value = result.tutees
    tutorSubjects.value = result.tutorSubjects
  }
}

const campaignId = useRoute().params.campaignId as string;

const users = ref<User[]>([]);

const fetchCampaign = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}/overview`);
  if (res.ok) {
    campaign.value = await res.json();
  } else {
    console.error("Failed to fetch campaign data");
  }
};

const fetchUsers = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}/users`);
  if (res.ok) users.value = await res.json();
};

const updateCampaign = async () => {
  const res = await useApiFetch(`/admin/campaign/${campaignId}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(campaign.value),
  });
  if (res.ok) {
    campaign.value = await res.json();
    isEditing.value = false;
  } else {
    console.error("Erreur lors de la mise à jour de la campagne");
  }
};

function handleUpdate(updatedCampaign: any) {
  campaign.value = updatedCampaign;
  isEditing.value = false;
  updateCampaign();
}

onMounted(() => {
  fetchCampaign();
  fetchUsers();
  fetchAssignments();
});

const tutorSubjectsByTutor = (tutorId: number): TutorSubject[] => {
  return tutorSubjects.value.filter((subject) => subject.tutor.id === tutorId);
}

const tuteeRegByTutee = (tuteeId: number): TuteeAssignment[] => {
  return tuteeAssignments.value.filter((assignment) => assignment.tuteeId === tuteeId);
}

const tutors = computed(() => {
  return users.value.filter((user) => user.isTutor);
})

const tutees = computed(() => {
  return users.value.filter((user) => user.isTutee);
})

const tutorsWithTotalHours = computed(() => {
  return tutors.value.map((tutor) => {
    const tutorSubjects = tutorSubjectsByTutor(tutor.id);
    let totalHours = 0;
    for (const subject of tutorSubjects) {
      totalHours += subject.totalHours;
    }
    return {
      ...tutor,
      totalHours,
    };
  });
});

const tuteesWithTotalHours = computed(() => {
  return tutees.value.map((tutee) => {
    const tuteeRegs = tuteeRegByTutee(tutee.id);
    let totalHours = 0;
    for (const reg of tuteeRegs) {
      totalHours += reg.totalHours;
    }
    return {
      ...tutee,
      totalHours,
    };
  });
});
</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Vue d'ensemble de la campagne</h1>
        <span
            :class="{
            'bg-green-100 text-green-800': campaign?.registrationStatus === 'OPEN',
            'bg-red-100 text-red-800': campaign?.registrationStatus === 'CLOSED',
          }"
            class="flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium"
        >
          <span
              :class="{
              'bg-green-500': campaign?.registrationStatus === 'OPEN',
              'bg-red-500': campaign?.registrationStatus === 'CLOSED',
            }"
              class="w-2 h-2 rounded-full"
          ></span>
          {{ campaign?.registrationStatus === 'OPEN' ? 'Ouverte' : 'Fermée' }}
        </span>
      </div>

      <div class="bg-white rounded-2xl shadow p-6 space-y-6">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-800">📅 Campagne {{ campaign?.schoolYear }} - Semestre
              {{ campaign?.semester }}</h2>
            <div class="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500">Dates de la campagne</p>
                <p v-if="campaign" class="text-lg font-medium text-gray-800">
                  Du {{ new Date(campaign.startDate).toLocaleDateString() }} au
                  {{ new Date(campaign.endDate).toLocaleDateString() }}
                </p>
              </div>
              <div>
                <p class="text-sm text-gray-500">Période d'inscription</p>
                <p v-if="campaign" class="text-lg font-medium text-gray-800">
                  Du {{ new Date(campaign.registrationStartDate).toLocaleDateString() }} au
                  {{ new Date(campaign.registrationEndDate).toLocaleDateString() }}
                </p>
              </div>
            </div>
          </div>
          <button
              class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium bg-yellow-500 text-white rounded-xl hover:bg-yellow-600 transition"
              @click="isEditing = true"
          >
            ✏️ Modifier
          </button>
        </div>
      </div>


      <div class="bg-white rounded-2xl shadow p-6">
        <h2 class="text-xl font-semibold text-gray-800 mb-6">👥 Étudiants</h2>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div>
            <h3 class="text-lg font-semibold text-gray-700">🧑‍🏫 Tuteurs</h3>
            <p class="text-sm text-gray-500 mb-4">
              {{ tutorsWithTotalHours.length }} tuteurs référencés
            </p>
            <ag-grid-vue
                :columnDefs="colDefs"
                :rowData="tutorsWithTotalHours"
                class="ag-theme-alpine"
                style="height: 600px"
            />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-700">🎓 Tutorés</h3>
            <p class="text-sm text-gray-500 mb-4">
              {{ tuteesWithTotalHours.length }} tutorés référencés
            </p>
            <ag-grid-vue
                :columnDefs="colDefs"
                :rowData="tuteesWithTotalHours"
                class="ag-theme-alpine"
                style="height: 600px"
            />
          </div>
        </div>
      </div>
    </div>

    <CampaignModal
        :initial-data="campaign"
        :show="isEditing"
        @close="isEditing = false"
        @submit="handleUpdate"
    />
  </div>
</template>


<style scoped>
</style>