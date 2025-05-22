<template>
  <Disclosure as="nav" class="bg-[#e61115] text-white shadow-md z-50" v-slot="{ open }">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">

        <div class="flex items-center space-x-4">
          <img
              class="h-10 w-auto"
              src="https://cas.insa-rouen.fr/cas/images/cas-logo.png"
              alt="INSA Rouen Normandie"
          />
          <span class="text-lg font-semibold hidden sm:inline-block">Tutorat STPI</span>
        </div>


        <div class="hidden sm:flex space-x-4">
          <NuxtLink
              v-for="item in getAvailableNavigation()"
              :key="item.key"
              :to="localePath(item.to)"
              :class="[
              route.path.startsWith(localePath(item.to))
                ? 'bg-white text-[#e61115]'
                : 'hover:bg-white hover:text-[#e61115]',
              'px-3 py-2 rounded-md text-sm font-medium transition'
            ]"
          >
            {{ t(item.label) }}
          </NuxtLink>
        </div>


        <div class="flex items-center space-x-4">

          <Menu as="div" class="relative">
            <MenuButton class="inline-flex items-center px-2 py-1 text-sm bg-white/20 hover:bg-white/30 rounded-md">
              🌐 <ChevronDownIcon class="w-4 h-4 ml-1" />
            </MenuButton>
            <MenuItems
                class="absolute right-0 mt-2 w-32 bg-white text-gray-800 rounded-md shadow-lg ring-1 ring-black/5 focus:outline-none z-50"
            >
              <MenuItem
                  v-for="locale in locales"
                  :key="locale.code"
                  v-slot="{ active }"
              >
                <button
                    @click="changeLocale(locale.code)"
                    :class="[
                    active ? 'bg-gray-100' : '',
                    'w-full text-left px-4 py-2 text-sm'
                  ]"
                >
                  {{ locale.name }}
                </button>
              </MenuItem>
            </MenuItems>
          </Menu>


          <div v-if="userStore.user" class="hidden sm:flex items-center space-x-3">
            <span class="font-medium">
              {{ userStore.user.lastName }} {{ userStore.user.firstName }}
            </span>
            <NuxtLink to="/logout" class="hover:text-gray-200 transition" aria-label="Logout">
              <ArrowRightEndOnRectangleIcon class="w-6 h-6" />
            </NuxtLink>
          </div>


          <DisclosureButton class="sm:hidden inline-flex items-center justify-center p-2 hover:bg-white/10 rounded-md">
            <span class="sr-only">Toggle Menu</span>
            <Bars3Icon v-if="!open" class="h-6 w-6" />
            <XMarkIcon v-else class="h-6 w-6" />
          </DisclosureButton>
        </div>
      </div>
    </div>


    <DisclosurePanel class="sm:hidden px-2 pt-2 pb-4 space-y-1">
      <NuxtLink
          v-for="item in getAvailableNavigation()"
          :key="item.key"
          :to="localePath(item.to)"
          :class="[
          route.path.startsWith(localePath(item.to))
            ? 'bg-white text-[#e61115]'
            : 'hover:bg-white hover:text-[#e61115]',
          'block rounded-md px-3 py-2 text-base font-medium transition'
        ]"
      >
        {{ t(item.label) }}
      </NuxtLink>

      <div class="border-t border-white/20 pt-2">
        <div v-if="userStore.user" class="px-3 text-sm">
          {{ userStore.user.lastName }} {{ userStore.user.firstName }}
          <NuxtLink to="/logout" class="block mt-1 text-sm text-white hover:underline">
            {{ t('logout') }}
          </NuxtLink>
        </div>

        <div class="mt-3 space-y-1">
          <button
              v-for="locale in locales"
              :key="locale.code"
              @click="changeLocale(locale.code)"
              class="block w-full text-left px-3 py-1 text-sm hover:bg-white/10 rounded"
          >
            🌐 {{ locale.name }}
          </button>
        </div>
      </div>
    </DisclosurePanel>
  </Disclosure>
</template>

<script setup lang="ts">
import {
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
  Menu,
  MenuButton,
  MenuItem,
  MenuItems
} from '@headlessui/vue'
import {
  Bars3Icon,
  XMarkIcon,
  ChevronDownIcon,
  ArrowRightEndOnRectangleIcon
} from '@heroicons/vue/24/outline'

import { useRoute, useRouter, useLocalePath, useI18n } from '#imports'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const { locales, locale, setLocale, t } = useI18n()
const localePath = useLocalePath()
const userStore = useUserStore()

interface NavigationItem {
  key: string
  label: string
  to: string
  role: 'isTutee' | 'isTutor' | 'isAdmin'
}

const navigation = ref<NavigationItem[]>([
  { key: 'tutee', label: 'tuteeSpace', to: '/tutee', role: 'isTutee' },
  { key: 'tutor', label: 'tutorSpace', to: '/tutor', role: 'isTutor' },
  { key: 'admin', label: 'adminSpace', to: '/admin', role: 'isAdmin' },
])

const getAvailableNavigation = () => {
  if (!userStore.user) return []
  return navigation.value.filter(item => userStore.user[item.role])
}

function changeLocale(code: 'fr' | 'en') {
  if (code !== locale.value) {
    setLocale(code)
    router.push(localePath(route.fullPath))
  }
}
</script>
