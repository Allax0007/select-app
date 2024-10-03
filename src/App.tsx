import './App.css'
import Header from './app/header';
import Sider from './app/sider';
import Home from './pages/home';
import Login from './pages/login';
import Manage from './pages/manage';
import Select from './pages/select';
import { Layout, theme } from 'antd';
import { useSelector } from "react-redux";
import { selectUser } from "./features/Auth/AuthSlice";
import { Routes, Route } from 'react-router-dom';
const { Content } = Layout;

function App() {
  const user = useSelector(selectUser);
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  return (
    <Layout>
      <Header />
      <Layout hasSider style={{minHeight: '90vh'}}>
        {user ? <Sider /> : null}
        <Layout style={{ padding: '0 24px 24px', background: '#dff7ff' }}>
          <br />
          <Content style={{ padding: '24px 16px 0', minHeight: 280, background: colorBgContainer, borderRadius: borderRadiusLG }}>
            <Routes>
              <Route index element={<Home />} />
              <Route path='database-management' element={<Manage />} />
              <Route path='select' element={<Select />} />
              <Route path='login' element={<Login />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
}

export default App
