import { useState, useEffect } from 'react';
import { Button, Form, Modal, message, Select } from 'antd';

export default function BuildCompRel({ onFinish }: { onFinish: () => void }){
  const [open, setOpen] = useState(false);  // State to track modal
  const [comp, setComp] = useState([]);  // State to track components
  const [msg, contextHolder] = message.useMessage();  // State to track message

  const getComp = async () => {
    try {
      const response = await fetch(`all-components`, {
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error('Failed to fetch data');
      }
      const data = await response.json();
      console.log(data);
      setComp(data);
    } catch (error) {
      console.error(error);
    }
  };
  useEffect(() => {
    getComp();
  } , []);

  const handleSubmit = async (value: any) => {
    console.log(value);
    try {
      const response = await fetch(`create-comp-rel`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(value),
      });
      const text = await response.text();
      if (!response.ok) {
        throw new Error(text);
      }
      setOpen(false);
      msg.success(text);

    } catch (error:any) {
      console.error('Error adding components relations:', error);
      msg.error(error.message);
    }
    onFinish();
  };
  const cancel = () => {
    onFinish();
    setOpen(false);
  }

  return(
    <>
    {contextHolder}
    <Button onClick={()=>setOpen(true)}>Add Components Relations</Button>
    <Modal
      title="Add a new component relation"
      footer={null}
      open={open}
      onCancel={cancel}
      destroyOnClose>
      <Form labelCol={{span: 8}} onFinish={handleSubmit} wrapperCol={{span: 10}}>
        <Form.Item label='ComponentA' name='ComponentA'
          rules={[{
            required: true,
            message: 'Please select a component',
          }]}>
          <Select placeholder='Select a component' >
            {comp.map((item: any) => (
              <Select.Option key={item} value={item}>{item}</Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item label='ComponentB' name='ComponentB'
          rules={[{
            required: true,
            message: 'Please select a component',
            },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('ComponentA') !== value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('ComponentB must be different from ComponentA'));
              }
            })]}>
          <Select placeholder='Select a component'>
            {comp.map((item: any) => (
              <Select.Option key={item} value={item}>{item}</Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type='primary' htmlType='submit'>Submit</Button>
        </Form.Item>
      </Form>
    </Modal>
    </>
  );
}