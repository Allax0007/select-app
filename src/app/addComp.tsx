import { useState } from 'react';
import { Button, Input, Form, Modal, message } from 'antd';

export default function NewComp({ onFinish }: { onFinish: () => void }){
  const [open, setOpen] = useState(false);  // State to track modal
  const [msg, contextHolder] = message.useMessage();  // State to track message

  const handleSubmit = async (value: any) => {
    console.log(value);
    try {
      const response = await fetch(`add-comp-type`, {
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
    <Button onClick={()=>setOpen(true)}>New Components Type</Button>
    <Modal
      title="Add a new component type"
      footer={null}
      open={open}
      onCancel={cancel}
      destroyOnClose>
      <Form labelCol={{span: 8}} onFinish={handleSubmit} wrapperCol={{span: 10}}>
        <Form.Item label='Type' name='type'>
          <Input required/>
        </Form.Item>
        <Form.Item>
          <Button type='primary' htmlType='submit'>Submit</Button>
        </Form.Item>
      </Form>
    </Modal>
    </>
  );
}