import Timer from "./timer";
import { useState } from "react";
import { Link } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import { Layout,Menu,MenuProps,Space,Modal } from "antd";
import { logout, selectUser } from "../features/Auth/AuthSlice";

type MenuItem = Required<MenuProps>['items'][number];

function Head() {

  const user = useSelector(selectUser);
  const dispatch = useDispatch();
  const [open, setOpen] = useState(false);  // State to track modal

  const LOGOUT = () => {
    setOpen(false);
    dispatch(logout(user));
  }

  const items1: MenuItem[] = [{
    key: '1',
    label: <Link to='login'>Login</Link>,
  }]
  const items2: MenuItem[] = [{
    key: '1',
    label: `${user}`,
    children: [{
      key: '2',
      label: 'Logout',
      onClick: () => { setOpen(true); },
    }]
  }]

  return( 
    <Layout.Header 
      style={{
        display: 'flex',
        alignItems: 'center',
        position: 'sticky',
        justifyContent: 'space-between',
        top: 0,
        zIndex: 1,
        color: '#ffffff',
        minHeight: '8vh',
      }}>
      <h1>SQLite Web Server</h1>
      <Space size={'middle'}>
        <Timer />
        <Menu 
          items={user ? items2 : items1} 
          theme="dark"
          selectedKeys={[]}></Menu>
      </Space>
      <Modal 
          title=' Are you sure you want to logout?' 
          open={open} 
          onOk={LOGOUT} 
          okText='Yes'
          closable={false}
          onCancel={() => setOpen(false)}>
      </Modal>
    </Layout.Header>
  );
}

export default Head;