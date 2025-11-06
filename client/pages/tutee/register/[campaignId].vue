<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {FolderOpenIcon} from '@heroicons/vue/24/outline'
import draggable from 'vuedraggable'
import {useToast} from "vue-toastification";
import type {Slots, Subject} from "~/types/api";

const campaignId = useRoute().params.campaignId as string;

const availabilitiesLink = computed(() => `/campaign/${campaignId}/availabilities`)

const { t: $t } = useI18n()

definePageMeta({
  layout: 'loggedin'
})

const subjects = ref<Subject[]>([])
const selectedSubjects = ref<Subject[]>([])

const availabilities = ref<Record<string, number[]>>({});

const areAvailabilitiesFilled = computed(() => {
  if (!availabilities.value) return false;
  return Object.keys(availabilities.value).length > 0;
})

onMounted(async () => {
  const availabilitiesResult = await useApiFetch(`/campaign/${campaignId}/availabilities`, {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })

  if (availabilitiesResult.ok) {
    availabilities.value = await availabilitiesResult.json() as Slots;
  }

  const subjectsResult = await useApiFetch(`/campaign/${campaignId}/subjects`, {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })

  if (subjectsResult.ok) {
    subjects.value = await subjectsResult.json() as Subject[];
  }

  const registrationsResult = await useApiFetch(`/campaign/${campaignId}/tutee/registrations`, {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })

  if (registrationsResult.ok) {
    selectedSubjects.value = await registrationsResult.json();
    const selectedSubjectsIds = selectedSubjects.value.map(selectedSubject => selectedSubject.id);
    subjects.value = subjects.value.filter(subject => !selectedSubjectsIds.includes(subject.id));
  }
})

async function submitChoices() {
  try {
    const response = await useApiFetch(`/campaign/${campaignId}/tutee/registrations`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        subjects: selectedSubjects.value.map(subject => subject.id),
      })
    });

    if (response.ok) {
      useToast().success('Vos choix ont bien été enregistrés');
    } else {
      console.error("Erreur lors de l\'enregistrement des choix :", await response.text());
      useToast().error('Erreur lors de l\'enregistrement des choix');
    }
  } catch (err) {
    console.error("Erreur lors de l\'enregistrement des choix :", err);
    useToast().error('Erreur lors de l\'enregistrement des choix');
  }
}
</script>


<template>
  <div class="max-w-6xl w-full mx-auto px-4 md:px-8 py-8">
    <h1 class="text-3xl font-bold mb-10 text-zinc-800 text-center md:text-left">{{ $t('registerTitle') }}</h1>

    <div class="flex flex-col gap-10">


      <section
          class="flex flex-col sm:flex-row justify-between bg-white rounded-2xl shadow-md p-6 md:p-8 sm:space-x-2 space-y-4 sm:space-y-0">
        <div>
          <h2 class="text-2xl font-semibold text-red-700 mb-3">{{ $t('availabilitySectionTitle') }}</h2>
          <div v-if="areAvailabilitiesFilled">
            <p class="text-sm text-zinc-600" v-html="$t('availabilitySectionFilled')"></p>
          </div>
          <div v-else class="flex flex-row space-x-2 items-center justify-center">
            <p class="text-sm text-zinc-600" v-html="$t('availabilitySectionEmpty')"></p>
          </div>
        </div>
        <div class="flex flex-shrink-0 items-center justify-center">
          <div class="flex flex-col">
            <MiniAgenda
                v-if="availabilities"
                :slots="availabilities"
            />
            <NuxtLink
                :to="availabilitiesLink"
                class="md:hidden inline-block mt-1 px-3 py-1.5 text-sm bg-zinc-100 text-zinc-800 hover:bg-zinc-200 rounded-md text-center"
            >
              {{ $t('editAvailabilitiesMobileButton') }}
            </NuxtLink>
          </div>
        </div>
      </section>


      <section class="bg-white rounded-2xl shadow-md p-6 md:p-8">
        <div class="flex items-center gap-3 mb-6">
          <FolderOpenIcon class="w-6 h-6 text-red-600"/>
          <h2 class="text-2xl font-semibold text-red-700">{{ $t('subjectsSectionTitle') }}</h2>
        </div>

        <div v-if="!areAvailabilitiesFilled" class="mb-4 flex items-center justify-between rounded-lg border border-red-200 bg-red-50 p-4">
          <div class="text-sm text-red-800">
            {{ $t('pleaseFillAvailabilitiesFirst') }}
          </div>
          <NuxtLink
              :to="availabilitiesLink"
              class="inline-block ml-4 px-3 py-1.5 text-sm bg-red-600 text-white rounded-md hover:bg-red-700"
          >
            {{ $t('editAvailabilitiesButton') }}
          </NuxtLink>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

          <div :class="{ 'opacity-60 pointer-events-none': !areAvailabilitiesFilled }">
            <h3 class="text-lg font-medium text-zinc-700 mb-3">{{ $t('availableSubjectsTitle') }}</h3>
            <draggable
                v-model="subjects"
                :group="{ name: 'subjects', pull: true, put: true }"
                class="bg-zinc-50 border border-zinc-200 rounded-xl p-4 min-h-[200px]"
                item-key="id"
                :disabled="!areAvailabilitiesFilled"
            >
              <template #item="{ element }">
                <div class="p-3 mb-2 rounded-lg bg-white shadow hover:bg-red-50 transition cursor-grab">
                  {{ element.shortName }} - {{ element.name }}
                </div>
              </template>
              <template v-if="subjects.length === 0" #footer>
                <p class="text-center text-zinc-400 italic pb-3">{{ $t('noSubjectsAvailable') }}</p>
              </template>
            </draggable>
          </div>


          <div :class="{ 'opacity-60 pointer-events-none': !areAvailabilitiesFilled }">
            <h3 class="text-lg font-medium text-zinc-700 mb-3">{{ $t('selectedSubjectsTitle') }}</h3>
            <draggable
                v-model="selectedSubjects"
                :group="{ name: 'subjects', pull: true, put: true }"
                class="bg-zinc-50 border border-zinc-200 rounded-xl p-4 min-h-[200px]"
                item-key="id"
                :disabled="!areAvailabilitiesFilled"
            >
              <template #item="{ element, index }">
                <div
                    class="p-3 mb-2 rounded-lg bg-red-100 text-red-800 font-medium shadow-inner flex items-center gap-2">
                  <span class="text-sm text-zinc-500">{{ index + 1 }}.</span>
                  {{ element.shortName }} - {{ element.name }}
                </div>
              </template>
              <template v-if="selectedSubjects.length === 0" #footer>
                <p class="text-center text-zinc-400 italic pb-3">{{ $t('dragHereSubjects') }}</p>
              </template>
            </draggable>
          </div>
        </div>
      </section>

      <div class="flex justify-end mt-6">
        <button
            :disabled="!areAvailabilitiesFilled"
            :aria-disabled="!areAvailabilitiesFilled"
            @click="submitChoices"
            :class="( !areAvailabilitiesFilled ) ? 'px-6 py-3 rounded-xl bg-zinc-300 text-zinc-600 font-semibold shadow cursor-not-allowed' : 'px-6 py-3 rounded-xl bg-red-600 text-white font-semibold shadow hover:bg-red-700 transition'"
        >
          {{ $t('submitChoicesButton') }}
        </button>
      </div>
    </div>
  </div>
</template>




