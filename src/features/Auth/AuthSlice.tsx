import { createSlice, Slice } from '@reduxjs/toolkit'

export const authSlice: Slice = createSlice({
    name: 'auth',
    initialState: {
        user: null,
        exp: null,
        token: null,
        admin: false,
    },
    reducers: {
        setCredentials: (state, action) => {
            const { token } = action.payload;
            state.token = token;
            const arrayToken = token.split('.');
            const parsed = JSON.parse(atob(arrayToken[1]));
            state.user = parsed.usr;
            state.exp = parsed.exp;
            state.admin = parsed.adm;
        },
        logout: (state) => {
            state.user = null;
            state.exp = null;
            state.token = null;
            state.admin = false;
            console.log('Logged out');
        },
    }
})

export const { setCredentials, logout } = authSlice.actions
export const selectUser = (state: any) => state.auth.user
export const selectExp = (state: any) => state.auth.exp
export const selectToken = (state: any) => state.auth.token
export const selectAdmin = (state: any) => state.auth.admin

export default authSlice.reducer