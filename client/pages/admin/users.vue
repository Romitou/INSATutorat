<script lang="ts" setup>
import {AllCommunityModule, ModuleRegistry} from 'ag-grid-community';
import {AgGridVue} from "ag-grid-vue3";
import type {User} from "~/types/api";
import UserRoleModal from '~/components/UserRoleModal.vue';

ModuleRegistry.registerModules([AllCommunityModule]);

definePageMeta({
  layout: 'loggedin'
})

const users = ref<User[]>([]);
const editingUser = ref<User | null>(null);
const showEditModal = ref(false);

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
    field: "studyYear",
    headerName: "Année d'étude",
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
  },
  {
    field: "lastLoginAt",
    headerName: "Dernière connexion",
    valueFormatter: (params) => params.value ? new Date(params.value).toLocaleDateString() : "Jamais",
  },
  {
    field: "status",
    headerName: "Statut",
    cellRenderer: (params) => {
      if (params.data.anonymizedAt) {
        return `<span class="text-zinc-400">Anonymisé le ${new Date(params.data.anonymizedAt).toLocaleDateString()}</span>`;
      }
      if (params.data.eligibleForAnonymization) {
        return `<span class="text-orange-600 font-medium">A quitté le cycle STPI</span>`;
      }
      return `<span class="text-green-600">Actif</span>`;
    },
  },
  {
    field: "edit",
    headerName: "Modifier",
    cellRenderer: (params) => {
      return `<button class="text-blue-600 hover:text-blue-700" data-edit-user="${params.data.id}">Modifier</button>`;
    },
    onCellClicked: (params) => {
      if (params.event?.target?.dataset?.editUser) {
        openEditModal(params.data);
      }
    },
  },
  {
    field: "anonymize",
    headerName: "Anonymiser",
    cellRenderer: (params) => {
      if (params.data.anonymizedAt || params.data.isAdmin) return "";
      return `<button class="text-red-600 hover:text-red-700" data-anonymize-user="${params.data.id}">Anonymiser</button>`;
    },
    onCellClicked: (params) => {
      if (params.event?.target?.dataset?.anonymizeUser) {
        handleAnonymizeUser(params.data);
      }
    },
  },
]);

const fetchUsers = async () => {
  const res = await useApiFetch(`/admin/users`);
  if (res.ok) {
    users.value = await res.json() as User[];
  } else {
    console.error("Failed to fetch user data");
  }
};

function openEditModal(user: User) {
  editingUser.value = user;
  showEditModal.value = true;
}

async function handleUpdateUser(payload: { id: number, isAdmin: boolean, isTutor: boolean, isTutee: boolean }) {
  const res = await useApiFetch(`/admin/users/${payload.id}`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({
      isAdmin: payload.isAdmin,
      isTutor: payload.isTutor,
      isTutee: payload.isTutee,
    }),
  });
  if (res.ok) {
    await fetchUsers();
  } else {
    console.error('Erreur lors de la mise à jour de l\'utilisateur:', await res.text());
  }
}

async function handleAnonymizeUser(user: User) {
  if (!confirm(
      `Anonymiser le compte de ${user.firstName} ${user.lastName} ?\n\n` +
      `Cette action est irréversible : son nom, son email et son identifiant CAS seront définitivement effacés. ` +
      `Les heures, séances et inscriptions déjà enregistrées sont conservées pour les statistiques.`
  )) return;

  const res = await useApiFetch(`/admin/users/${user.id}/anonymize`, {
    method: 'POST',
  });
  if (res.ok) {
    await fetchUsers();
  } else {
    console.error('Erreur lors de l\'anonymisation de l\'utilisateur:', await res.text());
  }
}

const eligibleCount = computed(() =>
    users.value.filter(u => u.eligibleForAnonymization).length
);

async function handleAnonymizeEligibleUsers() {
  if (eligibleCount.value === 0) return;

  if (!confirm(
      `Anonymiser les ${eligibleCount.value} compte(s) détecté(s) comme ayant quitté le cycle STPI ?\n\n` +
      `Cette action est irréversible : leur nom, leur email et leur identifiant CAS seront définitivement effacés. ` +
      `Les heures, séances et inscriptions déjà enregistrées sont conservées pour les statistiques.`
  )) return;

  const res = await useApiFetch(`/admin/users/anonymize-eligible`, {
    method: 'POST',
  });
  if (res.ok) {
    const { anonymizedCount } = await res.json();
    alert(`${anonymizedCount} compte(s) anonymisé(s).`);
    await fetchUsers();
  } else {
    console.error('Erreur lors de l\'anonymisation groupée:', await res.text());
  }
}

onMounted(() => {
  fetchUsers();
});

</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Liste des utilisateurs</h1>
        <button
            v-if="eligibleCount > 0"
            class="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
            @click="handleAnonymizeEligibleUsers"
        >
          Tout anonymiser ({{ eligibleCount }})
        </button>
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

    <UserRoleModal
        :show="showEditModal"
        :user="editingUser"
        @close="showEditModal = false"
        @submit="handleUpdateUser"
    />
  </div>
</template>

<style scoped>
</style>
