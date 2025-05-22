<script lang="ts" setup>
import {AcademicCapIcon, BookOpenIcon, UsersIcon} from "@heroicons/vue/24/outline";

definePageMeta({
  layout: 'loggedin'
})
</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Dashboard administrateur</h1>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">

        <NuxtLink class="block bg-white rounded-2xl shadow p-8 hover:shadow-lg transition group" to="/admin/campaigns">
          <div class="flex items-center space-x-4">
            <div class="bg-blue-100 text-blue-600 rounded-full p-3">
              <AcademicCapIcon class="h-8 w-8"/>
            </div>
            <div>
              <h2 class="text-xl font-semibold text-gray-900 group-hover:text-blue-700">Campagnes</h2>
              <p class="text-gray-500">Gérer les campagnes de tutorat</p>
            </div>
          </div>
        </NuxtLink>

        <NuxtLink class="block bg-white rounded-2xl shadow p-8 hover:shadow-lg transition group" to="/admin/users">
          <div class="flex items-center space-x-4">
            <div class="bg-green-100 text-green-600 rounded-full p-3">
              <UsersIcon class="h-8 w-8"/>
            </div>
            <div>
              <h2 class="text-xl font-semibold text-gray-900 group-hover:text-green-700">Utilisateurs</h2>
              <p class="text-gray-500">Voir et gérer les utilisateurs</p>
            </div>
          </div>
        </NuxtLink>
        <NuxtLink class="block bg-white rounded-2xl shadow p-8 hover:shadow-lg transition group" to="/admin/subjects">
          <div class="flex items-center space-x-4">
            <div class="bg-orange-100 text-orange-600 rounded-full p-3">
              <BookOpenIcon class="h-8 w-8"/>
            </div>
            <div>
              <h2 class="text-xl font-semibold text-gray-900 group-hover:text-orange-700">Matières</h2>
              <p class="text-gray-500">Voir et gérer les matières</p>
            </div>
          </div>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
