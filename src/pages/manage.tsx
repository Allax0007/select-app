import { useSelector } from "react-redux";
import AddProd from "../app/addProd";
import AddComp from "../app/addComp";
import BuildCompRel from "../app/BuildCompRel";
import PermissionDenied from "../app/403";
import React, { useState, useEffect } from "react";
import { selectUser, selectAdmin } from "../features/Auth/AuthSlice";
import { authenticate } from "../features/Auth/authenticate";
import { PlusCircleOutlined, UpOutlined, DownOutlined } from "@ant-design/icons";
import { Button, Collapse, Empty, Form, Input, message, Modal, Popconfirm, Space, Switch, Table } from "antd";

function Manage() {
  authenticate();
  const [msg, contextHolder] = message.useMessage();
  const [open, setOpen] = useState([false]);  // State to track modal
  const [data, setData] = useState<Record<string,any>>([]); // State to track data
  const [databuf, setDatabuf] = useState<Record<string,any>>([]); // State to track data
  const [columns, setColumns] = useState([]); // State to track columns
  const [tdata, setTdata] = useState([]); // State to track query 
  const [tables, setTables] = useState([]); // State to track tables
  const [currentTable, setCurrentTable] = useState(''); // State to track current table
  const disabledInsert = ['Component','Product'];

  const ListTable = async () => {
    try {
      const response = await fetch(`/get-tables`, {
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error('Failed to list table');
      }
      const data = await response.json();
      // console.log(data);
      setTables(data);
    } catch (error) {
      console.error('Error listing table:', error);
    }
  };

  useEffect(() => {
    ListTable();
  } , []);

  const onChange = async (key: string | string[]) => {
    if (key.length === 0) {
      setTdata([]);
      return;
    }
    try {
      setTdata([]);
      const response = await fetch(`/show-tables?q=${key[0]}`, {
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error(await response.text());
      }
      setCurrentTable(key[0]);
      const data = await response.json();
      setColumns(data[0]);
      setTdata(data.slice(1));
    }catch (error) {
      console.error('Error showing table:', error);
      msg.error((error as any).message);
    }
  };

  const edit = (record: any) => {
    setData(record);
    setDatabuf(record);
    setOpen([true, false]);
  };
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setData({...data, [e.target.name]: e.target.value});
  };

  const handleSwitch = (record: any) => async (checked: boolean) => {
    console.log('Switch', record, checked);
    const checkInt = checked ? 1 : 0;
    try {
      const response = await fetch(`/update-data?t=users&&col=3&&val=${checkInt}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(record),
      });
      if (!response.ok) {
        throw new Error('Failed to update data');
      }
      const text = await response.text();
      // console.log(text);
      msg.success(text);
      setOpen([false]);
    } catch (error: any) {
      console.error('Error updating data:', error);
      msg.error(error.message);
    }
    onChange(['users']);
  }

  const insert = async (table: string) => {
    console.log('Insert', data, 'on', table);
    try {
      // if data is empty, throw an error
      if (Object.keys(data).length === 0) {
        throw new Error('You must fill in the data');
      }
      const response = await fetch(`/insert-data?t=${table}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      const text = await response.text();
      if (!response.ok) {
        throw new Error(text);
      }
      // console.log(text);
      msg.success(text);
      setOpen([false]);
    } catch (error: any) {
      console.error('Error inserting data:', error);
      msg.error(error.message);
    }
    onChange([table]);
    setData([]);
    setDatabuf([]);
  };
  const update = async (table: string) => {
    console.log('Update', data, 'on', table);
    for (const key in data) {
      console.log("original",databuf);
      console.log(columns[Number(key)], data[key] === databuf[key]);
      if (data[key] === databuf[key]) {
        continue;
      }
      try {
        const response = await fetch(`/update-data?t=${table}&&col=${key}&&val=${data[key]}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(databuf),
        });
        if (!response.ok) {
          throw new Error('Failed to update data');
        }
        const text = await response.text();
        // console.log(text);
        msg.success(text);
        databuf[key] = data[key];
      } catch (error: any) {
        console.error('Error updating data:', error);
        msg.error(error.message);
      }
    }
    onChange([table]);
    setOpen([false]);
    setData([]);
    setDatabuf([]);
  };
  const delData = async (table: string,record: any) => {
    console.log('Delete', record);
    try {
      if (table === 'users' && record.id === 1) {
        throw new Error('Cannot delete admin');
      }
      var url;
      if (table === 'Product') { url = `delete-product`;}
      else if (table === 'Component') { url = `del-comp-type`;}
      else { url = `delete-data?t=${table}`;}
      const response = await fetch(`${url}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(record),
      });
      if (!response.ok) {
        throw new Error('Failed to delete data');
      }
      const text = await response.text();
      // console.log(text);
      // Remove the deleted data from state
      setTdata(tdata.filter((item) => item !== record));
      ListTable();
      msg.success(text);
    } catch (error: any) {
      console.error('Error deleting data:', error);
      msg.error(error.message);
    }
  };
  const delTable = async (table: string) => {
    console.log('Delete', table);
    try {
      const response = await fetch(`drop-table?t=${table}`, {
        method: 'DELETE',
      });
      if (!response.ok) {
        throw new Error('Failed to delete table');
      }
      const text = await response.text();
      // console.log(text);
      // Remove the deleted data from state
      setTables(tables.filter((item) => item !== table));
      ListTable();
      msg.success(text);
    } catch (error: any) {
      console.error('Error deleting table:', error);
      msg.error(error.message);
    }
  }
  const cancel = () => {
    console.log('cancel');
    setOpen([false]);
    setData([]);
    setDatabuf([]);
  };
  const chgOrder = async (record: any, direction: number) => {
    console.log('Change order', record[0]);
    console.log('Direction:', direction);
    try {
      const response = await fetch(`/chg-order?t=${currentTable}&&ord=${record[0]}&&dir=${direction}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error(await response.text());
      }
      const text = await response.text();
      msg.success(text);
      onChange([currentTable]);
    }
    catch (error) {
      console.error('Error changing order:', error);
      msg.error((error as any).message);
    }
  };

  const user = useSelector(selectUser);
  const admin = useSelector(selectAdmin);

  return (
  <>
  {contextHolder}
  {!user || !admin ? 
  <PermissionDenied /> : (
  <Space direction="vertical" style={{ display: 'flex' }}>
    <Space direction="horizontal">
      <AddComp onFinish={ListTable}/>
      <AddProd onFinish={ListTable}/>
      <BuildCompRel onFinish={ListTable}/>
    </Space>
  <Collapse accordion bordered={false} onChange={(activeKey) => onChange(activeKey)}>
    {tables.map((table) => (
      <Collapse.Panel key={table} header={table} >
        <Modal
          title={"Insert Data to " + currentTable}
          open={open[1]}
          onOk={()=>insert(currentTable)}
          onCancel={cancel}
          destroyOnClose>
          <Form labelCol={{span: 4}} wrapperCol={{span: 20}}>
            {columns.map((col: any, index: any) => (
              col !== 'manager' &&
              <Form.Item label={col} name={col} key={index}>
                <Input name={col} onChange={(e)=>handleInputChange(e)}/>
              </Form.Item>
            ))}
          </Form>
        </Modal>
        <Popconfirm title="Sure to delete table?" onConfirm={() => delTable(table)}>
          <Button danger>Delete Table</Button>
        </Popconfirm>
        {tdata === null ? <Empty description="No data" /> : (
        <Table 
          dataSource={tdata}
          columns={[...tdata.length > 0 ? Object.keys(tdata[0]).map((key: any) => (
            {title: columns[key], dataIndex: key, ellipsis: true, render: (text: any,record:any) => (
              columns[key] === 'manager' ? <Switch checked={record[key]} onChange={handleSwitch(record)}/>
               : (columns[key] === 'order' ? 
                <Space>
                  <Button icon={<UpOutlined />} onClick={() => chgOrder(record,-1)}></Button>
                  {text}
                  <Button icon={<DownOutlined />} onClick={() => chgOrder(record,1)} disabled></Button>
                </Space> : text))
            })) : [],
            {title: 'Actions', dataIndex: 'actions', render: (_,record:any) => (
              <Space>
                <Button onClick={() => edit(record)}>Edit</Button>
                <Modal
                  title={"Edit Data"}
                  open={open[0]}
                  onOk={()=>update(currentTable)}
                  onCancel={cancel}
                  destroyOnClose>
                  <Form labelCol={{span: 4}} wrapperCol={{span: 20}}>
                    {Object.keys(data).map((key: any) => (
                      columns[key] !== 'manager' &&
                      <Form.Item label={columns[key]} name={columns[key]} initialValue={data[key]} key={key}>
                        <Input name={key} value={data[key]} onChange={(e)=>handleInputChange(e)} disabled={columns[key] === 'id'} />
                      </Form.Item>
                    ))}
                  </Form>
                </Modal>
                <Popconfirm title="Sure to delete?" onConfirm={() => delData(currentTable,record)}>
                  <Button danger>Delete</Button>
                </Popconfirm>
              </Space>
            )}]}/>
        )}
        {disabledInsert.includes(currentTable) ? null : 
        <Button icon={<PlusCircleOutlined />} type="text" onClick={() => setOpen([false,true])} />}
      </Collapse.Panel>
    ))}
  </Collapse>
  </Space>
  )}
  </>
);
}

export default Manage;