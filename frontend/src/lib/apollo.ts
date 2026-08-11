import { ApolloClient, HttpLink, InMemoryCache } from '@apollo/client'
import { SetContextLink } from '@apollo/client/link/context'
import { getToken } from './token'

const httpLink = new HttpLink({
    uri: `${import.meta.env.VITE_API_URL ?? ''}/query`,
})

const authLink = new SetContextLink(({ headers }) => {
    const token = getToken()
    return {
        headers: {
            ...headers,
            authorization: token ? `Bearer ${token}` : '',
        },
    }
})

export const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
})
