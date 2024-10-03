import './index.css'
import React from 'react'
import App from './App.tsx'
import store from './app/store'
import { Provider } from 'react-redux'
import ReactDOM from 'react-dom/client'
import { HashRouter as Router } from 'react-router-dom'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Router>
      <Provider store={store}>
        <App />
      </Provider>
    </Router>
  </React.StrictMode>
)
