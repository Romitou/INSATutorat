import {defineStore} from 'pinia';
import {useToast} from "vue-toastification";

export interface User {
    id: number;
    firstName: string;
    lastName: string;
    mail: string;
    isTutor: boolean;
    isTutee: boolean;
    isAdmin: boolean;
}

export interface UserState {
    user: User | null;
}

export const useUserStore = defineStore('user', {
    state: (): UserState => ({
        user: null,
    }),
    // getters: {
    //
    // },
    actions: {
        async fetchUser() {
            try {
                const response = await useApiFetch('/auth/self')
                if (!response.ok) {
                    void useRouter().push('/')
                    return
                }

                this.user = await response.json()
            } catch (error) {
                this.user = null
            }
        },
        async logout() {
            const toast = useToast()
            const res = await useApiFetch('/auth/logout');
            if (res.ok) {
                this.user = null
                toast.success('Déconnexion réussie')
            }
        }
    },
})
