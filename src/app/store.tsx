import { configureStore } from '@reduxjs/toolkit'
import counterReducer from '../features/counter/counterSlice'
import authReducer from '../features/Auth/AuthSlice'

export default configureStore({
  reducer: {
    counter: counterReducer,
    auth: authReducer,
  },
})