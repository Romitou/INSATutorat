<script lang="ts" setup>
import {AllCommunityModule, ModuleRegistry} from 'ag-grid-community';
import {AgGridVue} from "ag-grid-vue3";
import type {Campaign} from "~/types/api";
import CampaignModal from '~/components/CampaignModal.vue';

ModuleRegistry.registerModules([AllCommunityModule]);

definePageMeta({
  layout: 'loggedin'
})

const colDefs = ref([
  {
    field: "schoolYear",
    headerName: "Année scolaire",
  },
  {
    field: "semester",
    headerName: "Semestre",
  },
  {
    field: "startDate",
    headerName: "Date de début",
    valueFormatter: (params) => new Date(params.value).toLocaleDateString(),
  },
  {
    field: "endDate",
    headerName: "Date de fin",
    valueFormatter: (params) => new Date(params.value).toLocaleDateString(),
  },
  {
    field: "registrationStatus",
    headerName: "Statut d'inscription",
  },
  {
    field: "view",
    headerName: "Voir",
    cellRenderer: (params) => {
      return `<button class="text-blue-600 hover:text-blue-700" onclick="window.location.href='/admin/campaign/${params.data.id}/overview'">Voir</button>`;
    },
  },
  {
    field: "assignments",
    headerName: "Affectations",
    cellRenderer: (params) => {
      return `<button class="text-blue-600 hover:text-blue-700" onclick="window.location.href='/admin/campaign/${params.data.id}/assignments'">Affectations</button>`;
    },
  }
]);

const campaigns = ref<Campaign[]>([]);
const showCreateModal = ref(false);

const fetchCampaigns = async () => {
  const res = await useApiFetch(`/admin/campaigns`, {
  });
  if (res.ok) {
    campaigns.value = await res.json() as Campaign[];
  } else {
    console.error("Erreur lors de la récupération des campagnes");
  }
};

async function handleCreateCampaign(newCampaign: Campaign) {
  newCampaign.semester = parseInt(newCampaign.semester.toString()); // malgré le type int, il est envoyé en string
  const res = await useApiFetch('/admin/campaigns', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(newCampaign),
  });
  if (res.ok) {
    showCreateModal.value = false;
    await fetchCampaigns();
  } else {
    console.error('Erreur lors de la création de la campagne');
  }
}

onMounted(() => {
  fetchCampaigns();
});

</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Liste des campagnes</h1>
        <button
            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm transition"
            @click="showCreateModal = true"
        >
          + Nouvelle campagne
        </button>
      </div>

      <div class="bg-white rounded-2xl shadow p-6 space-y-6">
        <ag-grid-vue
            :columnDefs="colDefs"
            :rowData="campaigns"
            class="ag-theme-alpine"
            style="height: 600px"
        />
      </div>
    </div>

    <CampaignModal
        :show="showCreateModal"
        @close="showCreateModal = false"
        @submit="handleCreateCampaign"
    />
  </div>
</template>

<style scoped>
</style>
