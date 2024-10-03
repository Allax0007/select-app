import { Layout, Menu, theme } from 'antd';
import type { MenuProps } from 'antd';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { selectAdmin } from '../features/Auth/AuthSlice';

type MenuItem = Required<MenuProps>['items'][number];

// const admin = useSelector(selectAdmin);
const adminItems: MenuItem[] = [
  {
    key: '1',
    label: <Link to='/'>Home</Link>,
  },
  {
    key: '2',
    label: <Link to='database-management'>Database Management</Link>,
  }
];

const userItems: MenuItem[] = [
  {
    key: '1',
    label: <Link to='/'>Home</Link>,
  },
  {
    key: '2',
    label: <Link to='select'>Select</Link>,
  },
];


export default function SiderMenu() {
  const admin = useSelector(selectAdmin)
  const {token: { colorBgContainer }} = theme.useToken();
  
  return (
    <Layout.Sider 
      breakpoint="lg"
      collapsedWidth="0"
      style={{ background: colorBgContainer }}>
      <Menu 
        items={admin ? adminItems : userItems}
        theme="dark" 
        mode="vertical" 
        defaultSelectedKeys={["1"]}
        style={{ height: '100%', borderRight: 0 }}>
      </Menu>
    </Layout.Sider>
  );
}