import { useState, useEffect } from 'react';
import { Button, Input, Form, Modal, message, Checkbox } from 'antd';

export default function AddProd({ onFinish }: { onFinish: () => void }){
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
      const response = await fetch(`create-product`, {
        method: 'PUT',
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
      console.error('Error adding product:', error);
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
    <Button onClick={()=>setOpen(true)}>Add Product</Button>
    <Modal
      title="Add a new product"
      footer={null}
      open={open}
      onCancel={cancel}
      destroyOnClose>
      <Form labelCol={{span: 4}} onFinish={handleSubmit} wrapperCol={{span: 20}}>
        <Form.Item label='Name' name='name'>
          <Input placeholder='Product name' required/>
        </Form.Item>
        <Form.Item label='Desc' name='desc'>
          <Input placeholder='Descriptions'/>
        </Form.Item>
        <Form.Item label='Comp' name='comp'
          rules={[{
            required: true,
            message: 'Please select at least one component',
          }]}>
          <Checkbox.Group 
            options={comp.map((item: any) => ({ label: item, value: item }))}>
            {comp.map((item: any) => (
              <Checkbox key={item} value={item}>{item}</Checkbox>
            ))}
          </Checkbox.Group>
        </Form.Item>
        <Form.Item>
          <Button type='primary' htmlType='submit'>Submit</Button>
        </Form.Item>
      </Form>
    </Modal>
    </>
  );
}