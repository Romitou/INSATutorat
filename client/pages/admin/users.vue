<script lang="ts" setup>
import {AllCommunityModule, ModuleRegistry} from 'ag-grid-community';
import {AgGridVue} from "ag-grid-vue3";
import type {User} from "~/types/api";

ModuleRegistry.registerModules([AllCommunityModule]);

definePageMeta({
  layout: 'loggedin'
})

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
    field: "mail",
    headerName: "Email",
  },
  {
    field: "schoolYear",
    headerName: "Année scolaire",
  },
  {
    field: "groups",
    headerName: "Groupes",
    valueFormatter: (params) => params.value.join(", "),
  },
  {
    field: "isTutor",
    headerName: "Tuteur",
  },
  {
    field: "isTutee",
    headerName: "Tutee",
  },
  {
    field: "isAdmin",
    headerName: "Admin",
  }
]);

const users = ref<User[]>([]);
// const showCreateModal = ref(false);

const fetchUsers = async () => {
  const res = await useApiFetch(`/admin/users`);
  if (res.ok) {
    users.value = await res.json() as User[];
  } else {
    console.error("Failed to fetch user data");
  }
};

// async function handleCreateUser(newUser: User) {
//   const res = await useApiFetch('/admin/users', {
//     method: 'POST',
//     credentials: 'include',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify(newUser),
//   });
//   if (res.ok) {
//     showCreateModal.value = false;
//     await fetchUsers();
//   } else {
//     console.error('Erreur lors de la création de l\'utilisateur:', await res.text());
//   }
// }

onMounted(() => {
  fetchUsers();
});

</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Liste des utilisateurs</h1>
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
            :rowData="users"
            class="ag-theme-alpine"
            style="height: 600px"
        />
      </div>
    </div>

    <!--    <userModal-->
    <!--        :show="showCreateModal"-->
    <!--        @close="showCreateModal = false"-->
    <!--        @submit="handleCreateuser"-->
    <!--    />-->
  </div>
</template>

<style scoped>
</style>
