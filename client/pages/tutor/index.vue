<script lang="ts" setup>
import {onMounted, ref, type Slot} from 'vue'
import {CheckCircleIcon, ClockIcon, FolderOpenIcon, UserGroupIcon, XCircleIcon} from '@heroicons/vue/24/outline'
import type {Campaign, TuteeRegistration, TutorSubject} from "~/types/api";

interface TutorAssignment extends TutorSubject {
  tutees: TuteeRegistration[];
}

interface CampaignTutorOverview extends Campaign {
  assignments: TutorAssignment[]
  availability: Slot[];
}

definePageMeta({
  layout: 'loggedin'
})

const campaigns = ref<CampaignTutorOverview[]>([])

onMounted(async () => {
  const assignmentsResult = await useApiFetch('/assignments/tutor', {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })

  if (assignmentsResult.ok) {
    campaigns.value = (await assignmentsResult.json()) as CampaignTutorOverview[];
  }

  for (const campaign of campaigns.value) {
    const agendaResult = await useApiFetch(`/campaign/${campaign.id}/availabilities`, {
      method: 'GET',
      headers: {'Content-Type': 'application/json'}
    })

    if (agendaResult.ok) {
      campaign.availability = await agendaResult.json() as Slot[];
    }
  }
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <h1 class="text-4xl font-extrabold text-gray-800 mb-8 text-center">
      {{ $t("mySemestersTitle") }}
    </h1>

    <div class="grid gap-8 md:grid-cols-2">
      <div
          v-for="campaign in campaigns"
          :key="campaign.id"
          class="bg-white rounded-lg shadow-lg border border-gray-200 hover:shadow-xl"
      >
        <div class="p-6">
          <h2 class="text-lg text-center font-semibold mb-4">
            {{ $t("semesterTitle") }} {{ campaign.semester }} — {{ campaign.schoolYear }}
          </h2>

          <div class="mb-6">
            <div class="flex items-center space-x-2 mb-2">
              <FolderOpenIcon class="h-5 w-5 text-gray-600"/>
              <h3 class="text-md font-medium text-gray-800">{{ $t("registrationsTitle") }}</h3>
            </div>
            <div class="text-sm text-gray-600">
              <div v-if="campaign.registrationStatus === 'CLOSED'" class="flex items-start space-x-2">
                <XCircleIcon class="h-5 w-5 text-red-600"/>
                <span>{{ $t("closedSinceThe") }} <strong>{{
                    new Date(campaign.registrationEndDate).toLocaleDateString()
                  }}</strong>.</span>
              </div>
              <div v-else class="flex items-start space-x-2">
                <CheckCircleIcon class="h-5 w-5 text-green-600"/>
                <span>
                  {{ $t("openUntilThe") }} <strong>{{
                    new Date(campaign.registrationEndDate).toLocaleDateString()
                  }}</strong>.
                </span>
              </div>
              <div v-if="campaign.registrationStatus === 'OPEN'" class="mt-4">
                <NuxtLink
                    :to="`/tutor/register/${campaign.id}`"
                    class="block text-center px-4 py-2 text-sm font-medium bg-green-700 text-white rounded-lg hover:bg-green-600"
                >
                  {{ $t("registerButton") }}
                </NuxtLink>
              </div>
            </div>
          </div>

          <div class="mb-6">
            <div class="flex items-center space-x-2 mb-2">
              <ClockIcon class="h-5 w-5 text-gray-600"/>
              <h3 class="text-md font-medium text-gray-800">{{ $t("availabilitiesTitle") }}</h3>
            </div>
            <p class="text-sm text-gray-600">
              {{ $t("availabilitiesDesc") }}
            </p>
            <div class="mt-4">
              <NuxtLink
                  :to="`/campaign/${campaign.id}/availabilities`"
                  class="block text-center px-4 py-2 text-sm font-medium bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
              >
                {{ $t("fillAvailabilities") }}
              </NuxtLink>
            </div>
          </div>

          <div>
            <div class="flex items-center space-x-2 mb-4">
              <UserGroupIcon class="h-6 w-6 text-gray-600"/>
              <h3 class="text-md font-medium text-gray-800">{{ $t("myTuteesTitle") }}</h3>
            </div>

            <div v-if="!campaign.assignments || campaign.assignments.length === 0" class="text-sm text-gray-600">
              <div class="bg-gray-50 p-4 rounded-lg border">
                <ClockIcon class="h-5 w-5 text-gray-500 inline-block mr-2 align-middle"/>
                <span>{{ $t("noTuteesYet") }}</span>
              </div>
            </div>

            <div v-else class="grid gap-4 mt-4">
              <div
                  v-for="assignment in campaign.assignments"
                  :key="assignment.id"
                  class="bg-gray-50 p-4 rounded-lg shadow border hover:shadow-md transition"
              >
                <p class="text-gray-700 font-medium mb-2">
                  📘 <span class="font-semibold">{{ assignment.subject.shortName }} — {{
                    assignment.subject.name
                  }}</span>
                </p>
                <p class="text-gray-600 text-sm">
                  👤 <span class="font-semibold">{{ $t("tuteesTitle") }} :</span>
                </p>
                <ul class="list-disc list-inside text-sm text-gray-600 mt-2">
                  <li
                      v-for="tuteeReg in assignment.tutees"
                      :key="tuteeReg.id"
                  >
                    {{ tuteeReg.tutee.lastName }} {{ tuteeReg.tutee.firstName }}
                  </li>
                </ul>
                <p class="text-xs text-gray-500 mt-2">
                  ⏳ Total : {{ assignment.totalHours }} {{ $t("hour") }}{{ assignment.totalHours > 1 ? 's' : '' }}
                </p>

                <NuxtLink
                    :to="`/tutoring/${assignment.id}`"
                    class="mt-4 block text-center px-4 py-2 text-sm font-medium text-white bg-blue-500 hover:bg-blue-400 rounded-lg shadow"
                >
                  {{ $t("moreInfoButton") }}
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

