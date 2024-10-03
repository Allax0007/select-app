import { useState } from 'react';
import { Button, Input, Form, Modal, Space } from 'antd';
import { useNavigate } from "react-router-dom";
import { useDispatch } from 'react-redux';
import { UserOutlined,LockOutlined } from '@ant-design/icons';
import { setCredentials } from '../features/Auth/AuthSlice';
import Register from '../app/register';
import ResetPwd from '../app/resetpwd';


const Login = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [open, setOpen] = useState([false]);  // State to track modal
  const [msg, setMsg] = useState('');  // State to track message
  const loginSubmit = async (value: any) => {
    try {
      const response = await fetch('/loginAuth', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(value),
      });
      if (!response.ok) {
        throw new Error('Failed to login');
      }
      const token = await response.json();
      console.log(token);
      dispatch(setCredentials({ token }));
      navigate('/');
    } catch (error: any) {
      console.error('Error logging in:', error);
      setMsg('Invalid username or password');
      setOpen([true, false, false]);
    }
  // console.log('localStorage token:',getAuthToken()); // Debug
  };

  return (
    <Space direction='vertical' style={{
      width:'300px',  
      borderRadius: '8px', 
      boxShadow: '0 0 10px #C8DEE5'}}>
    <h2>Log in</h2>
    <Form style={{padding:'0 8px 8px'}} onFinish={loginSubmit}>
      <Form.Item 
        name='name' 
        rules={[{ 
          required: true, 
          message: 'Please input your username!' }]}>
        <Input 
          prefix={<UserOutlined />}
          placeholder='Username'/>
      </Form.Item>
      <Form.Item 
      name='pwd' 
      rules={[{ 
        required: true, 
        message: 'Please input your password!' }]}>
        <Input.Password
          prefix={<LockOutlined />}
          placeholder='Password'/>
      </Form.Item>
      <Form.Item>
        <Button type='primary' htmlType='submit'>Log in</Button>
        <Modal 
          title='Error' 
          open={open[0]} 
          footer={<Button type='primary' onClick={()=>setOpen([false])}>Ok</Button>}
          closable={false}
          onCancel={()=>setOpen([false])}>
          <>{msg}</>
        </Modal>
      </Form.Item>
      <Button type='link' onClick={()=>setOpen([false, true, false])}>Forgot Password?</Button>
      <Modal 
        title='Reset Password' 
        open={open[1]} 
        footer={null}
        onCancel={() => setOpen([false])}>
        <ResetPwd onClose={() => setOpen([false])} />
      </Modal>
      <p className="read-the-docs">
        Do not have an account? <Button type='link' onClick={()=>setOpen([false, false, true])}>Register</Button>
      </p>
      <Modal 
        title='Register' 
        open={open[2]} 
        footer={null}
        onCancel={()=>setOpen([false])}>
        <Register onClose={() => setOpen([false])} />
      </Modal>
    </Form>
    </Space>
    );
};

export default Login;