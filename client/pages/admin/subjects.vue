<script lang="ts" setup>
import {AllCommunityModule, ModuleRegistry} from 'ag-grid-community';
import {AgGridVue} from "ag-grid-vue3";
import type {Subject} from "~/types/api";
import SubjectModal from '~/components/SubjectModal.vue';

ModuleRegistry.registerModules([AllCommunityModule]);

definePageMeta({
  layout: 'loggedin'
})

const subjects = ref<Subject[]>([]);
const showModal = ref(false);
const editingSubject = ref<Subject | null>(null);

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
  {
    field: "edit",
    headerName: "Modifier",
    cellRenderer: (params) => {
      return `<button class="text-blue-600 hover:text-blue-700" data-edit-subject="${params.data.id}">Modifier</button>`;
    },
    onCellClicked: (params) => {
      if (params.event?.target?.dataset?.editSubject) {
        openEditModal(params.data);
      }
    },
  },
  {
    field: "delete",
    headerName: "Supprimer",
    cellRenderer: (params) => {
      return `<button class="text-red-600 hover:text-red-700" data-delete-subject="${params.data.id}">Supprimer</button>`;
    },
    onCellClicked: (params) => {
      if (params.event?.target?.dataset?.deleteSubject) {
        handleDeleteSubject(params.data);
      }
    },
  },
]);

const fetchSubjects = async () => {
  const res = await useApiFetch(`/admin/subjects`);
  if (res.ok) {
    subjects.value = await res.json() as Subject[];
  } else {
    console.error("Erreur lors de la récupération des matières");
  }
};

function openCreateModal() {
  editingSubject.value = null;
  showModal.value = true;
}

function openEditModal(subject: Subject) {
  editingSubject.value = subject;
  showModal.value = true;
}

async function handleSubmitSubject(subject: Subject) {
  subject.semester = parseInt(subject.semester.toString()); // malgré le type int, il est envoyé en string

  const isEdit = !!subject.id;
  const res = await useApiFetch(isEdit ? `/admin/subjects/${subject.id}` : '/admin/subjects', {
    method: isEdit ? 'PATCH' : 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(subject),
  });
  if (res.ok) {
    await fetchSubjects();
  } else {
    console.error('Erreur lors de l\'enregistrement de la matière:', await res.text());
  }
}

async function handleDeleteSubject(subject: Subject) {
  if (!confirm(`Supprimer la matière "${subject.name}" ?`)) return;

  const res = await useApiFetch(`/admin/subjects/${subject.id}`, {
    method: 'DELETE',
  });
  if (res.ok) {
    await fetchSubjects();
  } else if (res.status === 409) {
    alert('Cette matière est utilisée dans une campagne et ne peut pas être supprimée.');
  } else {
    console.error('Erreur lors de la suppression de la matière:', await res.text());
  }
}

onMounted(() => {
  fetchSubjects();
});

</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Liste des matières</h1>
        <button
            class="px-4 py-2 rounded-md bg-blue-600 hover:bg-blue-700 text-white text-sm transition"
            @click="openCreateModal"
        >
          + Nouvelle matière
        </button>
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

    <SubjectModal
        :show="showModal"
        :initial-data="editingSubject"
        @close="showModal = false"
        @submit="handleSubmitSubject"
    />
  </div>
</template>

<style scoped>
</style>
