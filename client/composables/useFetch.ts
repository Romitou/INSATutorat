export const useApiFetch = async (url: string, options = {}) => {
    const config = useRuntimeConfig()
    try {
        return await fetch(`${config.public.BASE_URL}${url}`, {
            credentials: 'include',
            ...options,
        })
    } catch (error) {
        console.error('Fetch error:', error)
        throw error
    }
}
