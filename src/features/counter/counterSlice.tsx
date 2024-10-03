import { createSlice, Slice } from '@reduxjs/toolkit'

export const counterSlice: Slice = createSlice({
  name: 'counter',
  initialState: {
    value: 0
  },
  reducers: {
    increment: state => {
      state.value += 1
    },
    decrement: state => {
      state.value -= 1
    }
  }
})

export const { increment, decrement, incrementByAmount } = counterSlice.actions

export const selectCount = (state:any) => state.counter.value

export default counterSlice.reducer