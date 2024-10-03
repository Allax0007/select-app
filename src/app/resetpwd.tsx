import { useState } from 'react';
import { Button, Input, Form, Modal } from 'antd';
// import { useNavigate } from 'react-router-dom';
import { UserOutlined,LockOutlined } from '@ant-design/icons';

export default function resetpwd({onClose}: {onClose: () => void}){
  // const navigate = useNavigate();
  const [open, setOpen] = useState(false);  // State to track modal
  const [msg, setMsg] = useState('');  // State to track message
  const handleSubmit = async (value: any) => {
    // console.log(value);
    try {
      const response = await fetch(`edit-user`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(value),
      });
      const text = await response.text();
      setMsg(text);
      if (!response.ok) {
        throw new Error(text);
      }
    } catch (error) {
      console.error('Error adding user:', error);
    }
    setOpen(true);
  };

  return(
    <Form style={{padding:'8px'}} onFinish={handleSubmit} layout='vertical'>
    <Form.Item 
      name='name' 
      label='Username'
      rules={[{ 
        required: true,
        message: 'Please input your username!' }]}>
      <Input 
        prefix={<UserOutlined />}
        placeholder='Username'/>
    </Form.Item>
    <Form.Item 
      name='pwd' 
      label='New Password'
      hasFeedback
      rules={[{
        required: true, 
        message: 'Please input your new password!' 
      }, {
        validator(_, value) {
          if (!value || /^(?=.*[a-zA-Z])(?=.*[0-9]).{8,}$/.test(value)) {
            return Promise.resolve();
          }
          return Promise.reject(new Error('Password must be a combination of letters and numbers, at least 8 characters'));
        }
      }]}>
      <Input.Password
        prefix={<LockOutlined />}
        placeholder='Password'/>
    </Form.Item>
    <Form.Item 
      name='pwd2' 
      label='Confirm Password'
      dependencies={['pwd']} 
      hasFeedback
      rules={[{
        required: true,
        message: 'Please confirm your password!',
        },
        ({ getFieldValue }) => ({
          validator(_, value) {
            if (!value || getFieldValue('pwd') === value) {
              return Promise.resolve();
            }
            return Promise.reject(new Error('The two passwords that you entered do not match!'));
          },
        }),
      ]}>
      <Input.Password
        prefix={<LockOutlined />}
        placeholder='Password'/>
    </Form.Item>
    <Form.Item>
      <Button type='primary' htmlType='submit'>Submit</Button>
    </Form.Item>
    <Modal 
      title={msg === 'User updated successfully' ? 'Success' : 'Error'}
      open={open} 
      footer={<Button type='primary' onClick={()=>{
        setOpen(false); 
        if(msg === 'User updated successfully') onClose();
      }}>Ok</Button>}
      closable={false}
      destroyOnClose
      onCancel={()=>setOpen(false)}>
      <>{msg}</>
    </Modal>
    </Form>
  )
}