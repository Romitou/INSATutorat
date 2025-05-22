<script lang="ts" setup>
import {AllCommunityModule, ModuleRegistry} from 'ag-grid-community';
import {AgGridVue} from "ag-grid-vue3";
import type {Subject} from "~/types/api";

ModuleRegistry.registerModules([AllCommunityModule]);

definePageMeta({
  layout: 'loggedin'
})

const colDefs = ref([
  {
    field: "semester",
    headerName: "Semestre",
  },
  {
    field: "shortName",
    headerName: "Abréviation",
  },
  {
    field: "name",
    headerName: "Nom",
  },
  // {
  //   field: "view",
  //   headerName: "Voir",
  //   cellRenderer: (params) => {
  //     return `<button class="text-blue-600 hover:text-blue-700" onclick="window.location.href='/admin/subjects/${params.data.id}'">Voir</button>`;
  //   },
  // }
]);

const subjects = ref<Subject[]>([]);

const fetchSubjects = async () => {
  const res = await useApiFetch(`/admin/subjects`);
  if (res.ok) {
    subjects.value = await res.json() as Subject[];
  } else {
    console.error("Erreur lors de la récupération des matières");
  }
};

// async function handleCreateSubject(newSubject: Subject) {
//   const res = await useApiFetch('/admin/subjects', {
//     method: 'POST',
//     credentials: 'include',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify(newSubject),
//   });
//   if (res.ok) {
//     showCreateModal.value = false;
//     await fetchSubjects();
//   } else {
//     console.error('Erreur lors de la création de la matière');
//   }
// }

onMounted(() => {
  fetchSubjects();
});

</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Liste des matières</h1>
        <!--        <button-->
        <!--            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm transition"-->
        <!--            @click="showCreateModal = true"-->
        <!--        >-->
        <!--          + Nouvelle campagne-->
        <!--        </button>-->
      </div>

      <div class="bg-white rounded-2xl shadow p-6 space-y-6">
        <ag-grid-vue
            :columnDefs="colDefs"
            :rowData="subjects"
            class="ag-theme-alpine"
            style="height: 600px"
        />
      </div>
    </div>

    <!--    <subjectModal-->
    <!--        :show="showCreateModal"-->
    <!--        @close="showCreateModal = false"-->
    <!--        @submit="handleCreatesubject"-->
    <!--    />-->
  </div>
</template>

<style scoped>
</style>
